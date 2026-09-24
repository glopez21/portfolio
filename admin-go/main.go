package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// ---------------------------------------------------------------------------
// Configuration (env-driven, mirroring the old Django settings)
// ---------------------------------------------------------------------------

type config struct {
	Addr        string
	ProjectsDir string
	DataDir     string
	Username    string
	Password    string
	SessionTTL  time.Duration
	QuotesDB    string
	ContactFile string

	// Contact-form email notification (disabled unless SMTPHost and
	// ContactEmail are set). JSONL storage always stays on.
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
	SMTPFrom     string
	ContactEmail string

	// Public contact endpoint abuse limits (per client IP, fixed window).
	ContactRateLimit   int
	ContactRateWindow  time.Duration
}

func loadConfig() config {
	rateLimit := envInt("CONTACT_RATE_LIMIT", 5)
	window, err := time.ParseDuration(env("CONTACT_RATE_WINDOW", "1h"))
	if err != nil || window <= 0 {
		window = time.Hour
	}
	return config{
		Addr:        env("ADMIN_ADDR", "0.0.0.0:8000"),
		ProjectsDir: env("PROJECTS_DIR", "/content/projects"),
		DataDir:     env("DATA_DIR", "/data"),
		Username:    env("ADMIN_USERNAME", "admin"),
		Password:    env("ADMIN_PASSWORD", "admin123"),
		SessionTTL:  12 * time.Hour,
		QuotesDB:    env("AI_QUOTES_DB", ""),
		ContactFile: env("CONTACT_FILE", "/data/contact_messages.jsonl"),

		SMTPHost:     env("SMTP_HOST", ""),
		SMTPPort:     env("SMTP_PORT", "587"),
		SMTPUser:     env("SMTP_USER", ""),
		SMTPPass:     env("SMTP_PASSWORD", ""),
		SMTPFrom:     env("SMTP_FROM", ""),
		ContactEmail: env("CONTACT_EMAIL", ""),

		ContactRateLimit:  rateLimit,
		ContactRateWindow: window,
	}
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ---------------------------------------------------------------------------
// Project model (mirrors content collection frontmatter schema)
// ---------------------------------------------------------------------------

var pillars = map[string]bool{
	"cybersecurity": true,
	"ai-ml":         true,
	"python":        true,
	"rust":          true,
	"homelab":       true,
}

var statuses = map[string]bool{
	"live":     true,
	"dev":      true,
	"planned":  true,
	"research": true,
}

type Project struct {
	Slug       string   `json:"slug" yaml:"-"`
	Title      string   `json:"title" yaml:"title"`
	Pillar     string   `json:"pillar" yaml:"pillar"`
	Tagline    string   `json:"tagline" yaml:"tagline"`
	Status     string   `json:"status" yaml:"status"`
	Order      int      `json:"order" yaml:"order"`
	Featured   bool     `json:"featured" yaml:"featured"`
	Path       string   `json:"path" yaml:"path"`
	Repo       string   `json:"repo" yaml:"repo"`
	Live       string   `json:"live" yaml:"live"`
	Stack      []string `json:"stack" yaml:"stack"`
	Highlights []string `json:"highlights" yaml:"highlights"`
	Body       string   `json:"body" yaml:"-"`
}

func (p Project) PillarLabel() string {
	switch p.Pillar {
	case "ai-ml":
		return "AI / ML"
	case "cybersecurity":
		return "Cybersecurity"
	case "python":
		return "Python"
	case "rust":
		return "Rust"
	case "homelab":
		return "Homelab"
	}
	return p.Pillar
}

func (p Project) StatusLabel() string {
	switch p.Status {
	case "live":
		return "Live"
	case "dev":
		return "In development"
	case "planned":
		return "Planned"
	case "research":
		return "Research"
	}
	return p.Status
}

// ---------------------------------------------------------------------------
// Content store: read/write Markdown files the Astro site consumes directly.
// The files under ProjectsDir ARE the source of truth for the static site,
// so editing via the admin is live for the next `make sync` build.
// ---------------------------------------------------------------------------

type Store struct {
	mu      sync.RWMutex
	dir     string
	digest  map[string]string
	changed time.Time
}

func NewStore(dir string) *Store {
	s := &Store{dir: dir}
	return s
}

func (s *Store) slugToFile(slug string) (string, bool) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return "", false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			if strings.TrimSuffix(e.Name(), ".md") == slug {
				return filepath.Join(s.dir, e.Name()), true
			}
		}
	}
	return "", false
}

func hasFrontmatter(raw []byte) bool {
	t := bytes.TrimLeft(raw, " \t\r\n")
	return bytes.HasPrefix(t, []byte("---"))
}

func parseProject(raw []byte, filename string) (Project, error) {
	p := Project{Slug: strings.TrimSuffix(filename, ".md")}
	text := string(raw)
	if !hasFrontmatter([]byte(text)) {
		p.Body = text
		return p, nil
	}
	parts := strings.SplitN(text, "---", 3)
	if len(parts) < 3 {
		p.Body = text
		return p, nil
	}
	front := parts[1]
	body := parts[2]
	// strip leading newline from body
	body = strings.TrimPrefix(body, "\n")
	if err := yaml.Unmarshal([]byte(front), &p); err != nil {
		return p, err
	}
	if err := yaml.Unmarshal([]byte(front), &struct {
		Tagline string   `yaml:"tagline"`
		Stack   []string `yaml:"stack"`
	}{}); err != nil {
		_ = err
	}
	p.Body = body
	return p, nil
}

func (s *Store) List() []Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil
	}
	var out []Project
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(s.dir, e.Name()))
		if err != nil {
			continue
		}
		p, err := parseProject(raw, e.Name())
		if err != nil {
			continue
		}
		out = append(out, p)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Pillar != out[j].Pillar {
			return pillarRank(out[i].Pillar) < pillarRank(out[j].Pillar)
		}
		return out[i].Order < out[j].Order
	})
	return out
}

func pillarRank(p string) int {
	switch p {
	case "cybersecurity":
		return 0
	case "ai-ml":
		return 1
	case "python":
		return 2
	case "rust":
		return 3
	case "homelab":
		return 4
	}
	return 9
}

// Marshal writes a project back to Markdown with proper YAML quoting.
func (p Project) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("---\n")
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	fm := struct {
		Title      string   `yaml:"title"`
		Pillar     string   `yaml:"pillar"`
		Tagline    string   `yaml:"tagline"`
		Status     string   `yaml:"status"`
		Order      int      `yaml:"order"`
		Featured   bool     `yaml:"featured"`
		Path       string   `yaml:"path,omitempty"`
		Repo       string   `yaml:"repo,omitempty"`
		Live       string   `yaml:"live,omitempty"`
		Stack      []string `yaml:"stack,omitempty"`
		Highlights []string `yaml:"highlights,omitempty"`
	}{
		Title:      p.Title,
		Pillar:     p.Pillar,
		Tagline:    p.Tagline,
		Status:     p.Status,
		Order:      p.Order,
		Featured:   p.Featured,
		Path:       p.Path,
		Repo:       p.Repo,
		Live:       p.Live,
		Stack:      p.Stack,
		Highlights: p.Highlights,
	}
	if len(fm.Stack) == 0 {
		fm.Stack = nil
	}
	if len(fm.Highlights) == 0 {
		fm.Highlights = nil
	}
	if err := enc.Encode(fm); err != nil {
		return nil, err
	}
	enc.Close()
	buf.WriteString("---\n\n")
	buf.WriteString(p.Body)
	return buf.Bytes(), nil
}

func (s *Store) Get(slug string) (Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	path, ok := s.slugToFile(slug)
	if !ok {
		return Project{}, false
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return Project{}, false
	}
	p, err := parseProject(raw, filepath.Base(path))
	if err != nil {
		return Project{}, false
	}
	return p, true
}

func (s *Store) Save(p Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	filename := p.Slug + ".md"
	full := filepath.Join(s.dir, filename)
	raw, err := p.Marshal()
	if err != nil {
		return err
	}
	return os.WriteFile(full, raw, 0644)
}

func (s *Store) Delete(slug string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	path, ok := s.slugToFile(slug)
	if !ok {
		return os.ErrNotExist
	}
	return os.Remove(path)
}

// ---------------------------------------------------------------------------
// Session store (in-memory token -> expiry), persisted checksums not required.
// ---------------------------------------------------------------------------

type SessionStore struct {
	mu    sync.Mutex
	token string
	exp   time.Time
	ttl   time.Duration
	user  string
}

func NewSessionStore(ttl time.Duration, user string) *SessionStore {
	return &SessionStore{ttl: ttl, user: user}
}

func (s *SessionStore) Issue() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = hex.EncodeToString(b)
	s.exp = time.Now().Add(s.ttl)
	return s.token
}

func (s *SessionStore) Valid(token string) bool {
	if token == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return token == s.token && time.Now().Before(s.exp)
}

// ---------------------------------------------------------------------------
// HTTP handlers
// ---------------------------------------------------------------------------

type App struct {
	cfg      config
	store    *Store
	quotes   *QuoteStore
	sessions *SessionStore
	tmpl     *template.Template
	mu       sync.Mutex

	// Per-IP fixed-window limiter for the public contact endpoint.
	limits   map[string]*ipWindow
	limitsMu sync.Mutex
}

// ipWindow counts contact-form submissions from one address inside a
// fixed window. Zero-allocation cleanup: expired entries are dropped
// lazily on each check.
type ipWindow struct {
	count int
	reset time.Time
}

// allowContact reports whether ip may submit another contact message,
// enforcing a fixed window of at most cfg.ContactRateLimit submissions
// (used for real logging + 429s; bots that fill the honeypot are dropped
// silently before this is even consulted).
func (a *App) allowContact(ip string) (ok bool, retryAfter time.Duration) {
	a.limitsMu.Lock()
	defer a.limitsMu.Unlock()
	now := time.Now()
	if a.limits == nil {
		a.limits = map[string]*ipWindow{}
	}
	// Lazy GC so the map cannot grow without bound under address churn.
	if len(a.limits) > 4096 {
		for k, w := range a.limits {
			if now.After(w.reset) {
				delete(a.limits, k)
			}
		}
	}
	w := a.limits[ip]
	if w == nil || now.After(w.reset) {
		w = &ipWindow{count: 0, reset: now.Add(a.cfg.ContactRateWindow)}
		a.limits[ip] = w
	}
	if w.count >= a.cfg.ContactRateLimit {
		return false, time.Until(w.reset)
	}
	w.count++
	return true, 0
}

func newApp(cfg config) (*App, error) {
	a := &App{
		cfg:      cfg,
		store:    NewStore(cfg.ProjectsDir),
		sessions: NewSessionStore(cfg.SessionTTL, cfg.Username),
	}
	if cfg.QuotesDB != "" {
		qs, err := NewQuoteStore(cfg.QuotesDB)
		if err != nil {
			return nil, fmt.Errorf("quotes db %q: %w", cfg.QuotesDB, err)
		}
		a.quotes = qs
	}
	funcs := template.FuncMap{
		"quotesEnabled": func() bool { return a.quotes != nil },
		"pillarLabel": pillarLabel,
		"statusLabel": func(s string) string {
			switch s {
			case "live":
				return "Live"
			case "dev":
				return "In development"
			case "planned":
				return "Planned"
			case "research":
				return "Research"
			}
			return s
		},
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict: odd number of args")
			}
			m := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict: key is not string")
				}
				m[key] = values[i+1]
			}
			return m, nil
		},
	}
	tmpl, err := template.New("").Funcs(funcs).ParseFS(contentFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	a.tmpl = tmpl
	return a, nil
}

func (a *App) hashPassword() string {
	h := sha256.Sum256([]byte(a.cfg.Password))
	return hex.EncodeToString(h[:])
}

func (a *App) authenticated(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	return a.sessions.Valid(c.Value)
}

func (a *App) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.authenticated(r) {
			http.Redirect(w, r, "/admin/login/", http.StatusFound)
			return
		}
		next(w, r)
	}
}

func (a *App) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template error: %v", err)
	}
}

// --- routes ---

func (a *App) handleStatic(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/admin-static/")
	data, err := contentFS.ReadFile("static/" + name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	if strings.HasSuffix(name, ".css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}
	w.Write(data)
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/", http.StatusFound)
}

// group is one pillar's bucket of projects. It is the SINGLE shared shape
// used by the dashboard landing (handleAdmin), the per-pillar pages
// (handlePillar), and the groupByPillar helper — so the ordering, labels,
// and bucket shape never drift between them.
type group struct {
	Key      string
	Projects []Project
}

// pillarOrder is the canonical display order of pillars. Both the dashboard
// landing and the per-pillar pages iterate it, so a change here reorders
// both surfaces at once.
var pillarOrder = []string{"cybersecurity", "ai-ml", "python", "rust", "homelab"}

// pillarLabel returns the human label for a pillar key. It is the single
// source of truth used by the dashboard landing, the per-pillar pages, and
// the template FuncMap — so the label text is guaranteed identical on every
// surface that shows a pillar.
func pillarLabel(key string) string {
	switch key {
	case "cybersecurity":
		return "Cybersecurity"
	case "ai-ml":
		return "AI / ML"
	case "python":
		return "Python"
	case "rust":
		return "Rust"
	case "homelab":
		return "Homelab"
	}
	return key
}

func (a *App) handleAdmin(w http.ResponseWriter, r *http.Request) {
	if !a.authenticated(r) {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
		return
	}
	projects := a.store.List()
	groups := groupByPillar(projects)
	type summary struct {
		Total        int
		Groups       []group
		QuoteCount   int
		ContactCount int
		Featured     int
	}
	s := summary{Groups: groups}
	for _, p := range projects {
		s.Total++
		if p.Featured {
			s.Featured++
		}
	}
	if a.quotes != nil {
		if n, err := a.quotes.Count(); err == nil {
			s.QuoteCount = n
		}
	}
	s.ContactCount = a.contactCount()
	a.render(w, "index.html", s)
}

// groupByPillar buckets projects by pillar and returns the buckets in the
// canonical homepage order (used by both the dashboard landing and the
// per-pillar pages so the ordering and labels never drift).
func groupByPillar(projects []Project) []group {
	byPillar := map[string][]Project{}
	for _, p := range projects {
		byPillar[p.Pillar] = append(byPillar[p.Pillar], p)
	}
	var out []group
	for _, k := range pillarOrder {
		if len(byPillar[k]) > 0 {
			out = append(out, group{Key: k, Projects: byPillar[k]})
		}
	}
	return out
}

// handlePillar renders one pillar's projects as an evenly-arranged card
// grid (one card per project, no vertical scroll marathon).
func (a *App) handlePillar(w http.ResponseWriter, r *http.Request) {
	if !a.authenticated(r) {
		http.Redirect(w, r, "/admin/login/", http.StatusFound)
		return
	}
	key := strings.TrimPrefix(r.URL.Path, "/admin/pillars/")
	key = strings.TrimSuffix(key, "/")
	type summary struct {
		Key      string
		Label    string
		Total    int
		Projects []Project
	}
	s := summary{Key: key, Label: pillarLabel(key)}
	for _, p := range a.store.List() {
		if p.Pillar == key {
			s.Projects = append(s.Projects, p)
			s.Total++
		}
	}
	// unknown pillar → 404 rather than a confusing empty page
	if s.Total == 0 {
		http.NotFound(w, r)
		return
	}
	a.render(w, "pillar.html", s)
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		u := r.FormValue("username")
		p := r.FormValue("password")
		if u == a.cfg.Username && p == a.cfg.Password {
			token := a.sessions.Issue()
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   int(a.cfg.SessionTTL.Seconds()),
			})
			http.Redirect(w, r, "/admin/", http.StatusSeeOther)
			return
		}
		a.render(w, "login.html", map[string]bool{"Error": true})
		return
	}
	a.render(w, "login.html", nil)
}

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/admin/login/", http.StatusSeeOther)
}

type listPage struct {
	Projects []Project
	Pillars  map[string]bool
}

func (a *App) handleList(w http.ResponseWriter, r *http.Request) {
	projects := a.store.List()
	a.render(w, "list.html", listPage{Projects: projects, Pillars: pillars})
}

func (a *App) handleAddForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, "edit.html", struct {
		Project     Project
		IsNew       bool
		Pillars     map[string]bool
		Statuses    map[string]bool
		StackLines  string
		HighlightIn string
	}{
		Project:     Project{Status: "dev", Order: 50, Stack: []string{}, Highlights: []string{}},
		IsNew:       true,
		Pillars:     pillars,
		Statuses:    statuses,
		StackLines:  "",
		HighlightIn: "",
	})
}

func (a *App) handleSave(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	p, ok := a.store.Get(slug)
	if !ok {
		p = Project{}
	}
	r.ParseForm()
	slug = r.FormValue("slug")
	if slug == "" {
		slug = r.PathValue("slug")
	}
	slug = sanitizeSlug(slug)
	p.Slug = slug
	p.Title = r.FormValue("title")
	p.Pillar = r.FormValue("pillar")
	p.Tagline = r.FormValue("tagline")
	p.Status = r.FormValue("status")
	fmt.Sscanf(r.FormValue("order"), "%d", &p.Order)
	p.Featured = r.FormValue("featured") == "on" || r.FormValue("featured") == "1" || r.FormValue("featured") == "true"
	p.Path = r.FormValue("path")
	p.Repo = r.FormValue("repo")
	p.Live = r.FormValue("live")
	p.Stack = splitLines(r.FormValue("stack"))
	p.Highlights = splitLines(r.FormValue("highlights"))
	p.Body = r.FormValue("body")

	if !pillars[p.Pillar] {
		p.Pillar = "python"
	}
	if !statuses[p.Status] {
		p.Status = "dev"
	}
	if p.Title == "" {
		p.Title = slug
	}
	// rename existing file if slug changed
	if oldSlug := r.PathValue("slug"); oldSlug != "" && oldSlug != slug {
		if path, ok := a.store.slugToFile(oldSlug); ok {
			os.Rename(path, filepath.Join(a.store.dir, slug+".md"))
		}
	}
	if err := a.store.Save(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/projects/", http.StatusSeeOther)
}

func (a *App) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"count":    len(a.store.List()),
		"projects": a.store.List(),
	})
}

// contactMessage is one contact-form submission, persisted as one JSON line.
type contactMessage struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Time    string `json:"time"`
	IP      string `json:"ip"`
}

// handleContact accepts the site's contact form (PHP-email-form sends
// urlencoded form data and expects the success body to be exactly "OK").
func (a *App) handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// Honeypot: real visitors never see or fill the hidden "website"
	// field; bots that do are dropped silently (OK body, nothing stored)
	// so they don't learn they were caught.
	if strings.TrimSpace(r.FormValue("website")) != "" {
		log.Printf("contact: honeypot tripped (ip %s) — dropped", clientIP(r))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
		return
	}
	msg := contactMessage{
		Name:    strings.TrimSpace(r.FormValue("name")),
		Email:   strings.TrimSpace(r.FormValue("email")),
		Subject: strings.TrimSpace(r.FormValue("subject")),
		Message: strings.TrimSpace(r.FormValue("message")),
		Time:    time.Now().UTC().Format(time.RFC3339),
		IP:      clientIP(r),
	}
	if ok, retry := a.allowContact(msg.IP); !ok {
		log.Printf("contact: rate limited ip %s (retry in %s)", msg.IP, retry.Round(time.Second))
		w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())+1))
		http.Error(w, "too many messages — please try again later", http.StatusTooManyRequests)
		return
	}
	if msg.Name == "" || msg.Email == "" || msg.Message == "" {
		http.Error(w, "name, email and message are required", http.StatusBadRequest)
		return
	}
	if !strings.Contains(msg.Email, "@") {
		http.Error(w, "invalid email address", http.StatusBadRequest)
		return
	}
	line, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "could not encode message", http.StatusInternalServerError)
		return
	}
	a.mu.Lock()
	err = appendLine(a.cfg.ContactFile, string(line))
	a.mu.Unlock()
	if err != nil {
		log.Printf("contact: store: %v", err)
		http.Error(w, "could not store message", http.StatusInternalServerError)
		return
	}
	a.notifyContact(msg)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}

// notifyContact emails the site owner about a new contact message. It is
// fire-and-forget: the visitor already got their OK (the JSONL copy is the
// durable record), and any SMTP failure is only logged.
func (a *App) notifyContact(msg contactMessage) {
	c := a.cfg
	if c.SMTPHost == "" || c.ContactEmail == "" {
		return
	}
	go func() {
		if err := a.sendContactMail(msg); err != nil {
			log.Printf("contact: notify email failed: %v", err)
		}
	}()
}

// sanitizeHeader collapses CR/LF into spaces — name/subject come from the
// visitor and must never smuggle extra email headers (header injection).
func sanitizeHeader(s string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(strings.TrimSpace(s))
}

func (a *App) sendContactMail(msg contactMessage) error {
	c := a.cfg
	from := c.SMTPFrom
	if from == "" {
		from = c.SMTPUser
	}
	subject := "[4rch3.io contact] " + sanitizeHeader(msg.Subject)
	if subject == "[4rch3.io contact] " {
		subject = "[4rch3.io contact] (no subject)"
	}
	body := fmt.Sprintf(
		"New contact-form message on 4rch3.io\n\n"+
			"Name:    %s\n"+
			"Email:   %s\n"+
			"Subject: %s\n"+
			"Time:    %s\n"+
			"IP:      %s\n\n"+
			"--- message ---\n%s\n--- end ---\n"+
			"Reply directly to this email to answer the visitor (Reply-To is set).\n"+
			"All messages are also archived at /admin/contact/.\n",
		sanitizeHeader(msg.Name), sanitizeHeader(msg.Email), sanitizeHeader(msg.Subject),
		msg.Time, msg.IP, msg.Message)
	headers := map[string]string{
		"From":         from,
		"To":           c.ContactEmail,
		"Subject":      subject,
		"Reply-To":     sanitizeHeader(msg.Email),
		"Date":         time.Now().Format(time.RFC1123Z),
		"MIME-Version": "1.0",
		"Content-Type": `text/plain; charset="utf-8"`,
	}
	var b bytes.Buffer
	for _, k := range []string{"From", "To", "Subject", "Reply-To", "Date", "MIME-Version", "Content-Type"} {
		fmt.Fprintf(&b, "%s: %s\r\n", k, headers[k])
	}
	b.WriteString("\r\n" + body)

	addr := net.JoinHostPort(c.SMTPHost, c.SMTPPort)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(20 * time.Second))
	cl, err := smtp.NewClient(conn, c.SMTPHost)
	if err != nil {
		return fmt.Errorf("smtp greeting: %w", err)
	}
	defer cl.Close()
	if ok, _ := cl.Extension("STARTTLS"); ok {
		if err = cl.StartTLS(&tls.Config{ServerName: c.SMTPHost}); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}
	if c.SMTPUser != "" {
		if err = cl.Auth(smtp.PlainAuth("", c.SMTPUser, c.SMTPPass, c.SMTPHost)); err != nil {
			return fmt.Errorf("auth: %w", err)
		}
	}
	if err = cl.Mail(from); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	if err = cl.Rcpt(c.ContactEmail); err != nil {
		return fmt.Errorf("rcpt to: %w", err)
	}
	w, err := cl.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	if _, err = w.Write(b.Bytes()); err != nil {
		return fmt.Errorf("write body: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("close body: %w", err)
	}
	return cl.Quit()
}

// handleContactAdmin lists stored contact submissions, newest first.
func (a *App) contactCount() int {
	a.mu.Lock()
	data, err := os.ReadFile(a.cfg.ContactFile)
	a.mu.Unlock()
	if err != nil {
		return 0
	}
	n := 0
	for _, ln := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(ln) != "" {
			n++
		}
	}
	return n
}

func (a *App) handleContactAdmin(w http.ResponseWriter, r *http.Request) {
	var msgs []contactMessage
	a.mu.Lock()
	data, err := os.ReadFile(a.cfg.ContactFile)
	a.mu.Unlock()
	if err == nil {
		for _, ln := range strings.Split(string(data), "\n") {
			ln = strings.TrimSpace(ln)
			if ln == "" {
				continue
			}
			var m contactMessage
			if json.Unmarshal([]byte(ln), &m) == nil {
				msgs = append(msgs, m)
			}
		}
	}
	sort.Slice(msgs, func(i, j int) bool { return msgs[i].Time > msgs[j].Time })
	a.render(w, "contact.html", struct {
		Count    int
		Messages []contactMessage
	}{Count: len(msgs), Messages: msgs})
}

func appendLine(path, line string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, line)
	return err
}

// clientIP resolves the visitor's address. X-Real-IP is set by our own
// nginx proxy ($remote_addr — not client-controllable), so it wins over
// X-Forwarded-For, whose first entry a client can spoof by sending its
// own header before the proxy appends to it.
func clientIP(r *http.Request) string {
	if rip := strings.TrimSpace(r.Header.Get("X-Real-Ip")); rip != "" {
		return rip
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if parts := strings.Split(fwd, ","); len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func sanitizeSlug(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r - 'A' + 'a')
		default:
			b.WriteRune('-')
		}
	}
	res := strings.Trim(b.String(), "-")
	if res == "" {
		res = "project"
	}
	return res
}

func splitLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func main() {
	cfg := loadConfig()
	if _, err := os.Stat(cfg.ProjectsDir); os.IsNotExist(err) {
		log.Fatalf("projects dir %q does not exist", cfg.ProjectsDir)
	}
	a, err := newApp(cfg)
	if err != nil {
		log.Fatalf("init: %v", err)
	}
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatalf("data dir: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", a.handleIndex)
	mux.HandleFunc("/admin-static/", a.handleStatic)
	mux.HandleFunc("/admin/login/", a.handleLogin)
	mux.HandleFunc("/admin/logout/", a.handleLogout)
	mux.HandleFunc("/admin/", a.requireAuth(a.handleAdmin))
	mux.HandleFunc("/admin/projects/", a.requireAuth(a.handleProjects))
	mux.HandleFunc("/admin/pillars/", a.requireAuth(a.handlePillar))
	mux.HandleFunc("/admin/quotes/", a.requireAuth(a.handleQuotes))
	mux.HandleFunc("/api/projects/", a.handleAPI)
	mux.HandleFunc("/api/quotes/", a.handleQuotesAPI)
	mux.HandleFunc("/api/contact/", a.handleContact)
	mux.HandleFunc("/admin/contact/", a.requireAuth(a.handleContactAdmin))

	notify := "off"
	if cfg.SMTPHost != "" && cfg.ContactEmail != "" {
		notify = "on (" + cfg.ContactEmail + ")"
		if cfg.SMTPUser != "" && cfg.SMTPPass == "" {
			notify += " — WARNING: SMTP_PASSWORD is empty; set it in .env or mails will fail auth"
		}
	}
	log.Printf("4rch3 portfolio-admin listening on %s (projects: %s, quotes: %s, contact email notify: %s)", cfg.Addr, cfg.ProjectsDir, cfg.QuotesDB, notify)
	log.Fatal(http.ListenAndServe(cfg.Addr, mux))
}

// handleProjects dispatches on the sub-path under /admin/projects/ to avoid
// ServeMux wildcard/literal conflicts.
func (a *App) handleProjects(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/admin/projects/")
	rest = strings.TrimSuffix(rest, "/")
	switch {
	case rest == "":
		a.handleList(w, r)
	case rest == "add":
		a.handleAddForm(w, r)
	case rest == "new":
		a.handleSave(w, r)
	case strings.HasSuffix(rest, "/save"):
		a.handleSaveSlug(w, r, strings.TrimSuffix(rest, "/save"))
	case strings.HasSuffix(rest, "/delete"):
		a.handleConfirmDelete(w, r, strings.TrimSuffix(rest, "/delete"))
	case strings.HasSuffix(rest, "/del"):
		a.handleDelete(w, r, strings.TrimSuffix(rest, "/del"))
	default:
		a.handleEditForm(w, r, rest)
	}
}

func (a *App) handleSaveSlug(w http.ResponseWriter, r *http.Request, slug string) {
	r.SetPathValue("slug", slug)
	a.handleSave(w, r)
}

func (a *App) handleConfirmDelete(w http.ResponseWriter, r *http.Request, slug string) {
	r.SetPathValue("slug", slug)
	p, ok := a.store.Get(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}
	a.render(w, "confirm_delete.html", p)
}

func (a *App) handleDelete(w http.ResponseWriter, r *http.Request, slug string) {
	r.SetPathValue("slug", slug)
	if err := a.store.Delete(slug); err != nil && err != os.ErrNotExist {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/projects/", http.StatusSeeOther)
}

func (a *App) handleEditForm(w http.ResponseWriter, r *http.Request, slug string) {
	r.SetPathValue("slug", slug)
	p, ok := a.store.Get(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}
	a.render(w, "edit.html", struct {
		Project     Project
		IsNew       bool
		Pillars     map[string]bool
		Statuses    map[string]bool
		StackLines  string
		HighlightIn string
	}{
		Project:     p,
		IsNew:       false,
		Pillars:     pillars,
		Statuses:    statuses,
		StackLines:  strings.Join(p.Stack, "\n"),
		HighlightIn: strings.Join(p.Highlights, "\n"),
	})
}
