package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

const quoteCols = `id, text, author, category, tags, is_favorite, is_archived, views, upvotes, downvotes, created_at, updated_at`

type Quote struct {
	ID         int    `json:"id"`
	Text       string `json:"text"`
	Author     string `json:"author"`
	Category   string `json:"category"`
	Tags       string `json:"tags"`
	IsFavorite bool   `json:"is_favorite"`
	IsArchived bool   `json:"is_archived"`
	Views      int    `json:"views"`
	Upvotes    int    `json:"upvotes"`
	Downvotes  int    `json:"downvotes"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type QuoteStore struct {
	db *sql.DB
}

func NewQuoteStore(path string) (*QuoteStore, error) {
	dsn := path + "?_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return &QuoteStore{db: db}, nil
}

func scanQuote(scan func(dest ...any) error) (Quote, error) {
	var q Quote
	var fav, arch int
	var createdAt, updatedAt sql.NullString
	err := scan(&q.ID, &q.Text, &q.Author, &q.Category, &q.Tags, &fav, &arch,
		&q.Views, &q.Upvotes, &q.Downvotes, &createdAt, &updatedAt)
	q.IsFavorite = fav != 0
	q.IsArchived = arch != 0
	q.CreatedAt = createdAt.String
	q.UpdatedAt = updatedAt.String
	return q, err
}

func (s *QuoteStore) List(search string, includeArchived bool) ([]Quote, error) {
	conds := []string{}
	args := []any{}
	if search != "" {
		like := "%" + search + "%"
		conds = append(conds, "(text LIKE ? OR author LIKE ?)")
		args = append(args, like, like)
	}
	if !includeArchived {
		conds = append(conds, "is_archived = 0")
	}
	q := "SELECT " + quoteCols + " FROM quote"
	if len(conds) > 0 {
		q += " WHERE " + strings.Join(conds, " AND ")
	}
	q += " ORDER BY is_archived ASC, created_at DESC, id DESC"
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Quote
	for rows.Next() {
		x, err := scanQuote(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (s *QuoteStore) Get(id int) (Quote, error) {
	row := s.db.QueryRow("SELECT "+quoteCols+" FROM quote WHERE id = ?", id)
	return scanQuote(row.Scan)
}

func (s *QuoteStore) Count() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM quote").Scan(&n)
	return n, err
}

func (s *QuoteStore) CountArchived() (int, error) {
	var n int
	err := s.db.QueryRow("SELECT COUNT(*) FROM quote WHERE is_archived = 1").Scan(&n)
	return n, err
}

func (s *QuoteStore) Create(q Quote) (int64, error) {
	res, err := s.db.Exec(`INSERT INTO quote (text, author, category, tags, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		strings.TrimSpace(q.Text), strings.TrimSpace(q.Author), strings.TrimSpace(q.Category), strings.TrimSpace(q.Tags))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *QuoteStore) Update(q Quote) error {
	_, err := s.db.Exec(`UPDATE quote SET text = ?, author = ?, category = ?, tags = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		strings.TrimSpace(q.Text), strings.TrimSpace(q.Author), strings.TrimSpace(q.Category), strings.TrimSpace(q.Tags), q.ID)
	return err
}

func (s *QuoteStore) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM quote WHERE id = ?", id)
	return err
}

func (s *QuoteStore) SetFavorite(id int, v bool) error {
	_, err := s.db.Exec("UPDATE quote SET is_favorite = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", v, id)
	return err
}

func (s *QuoteStore) SetArchived(id int, v bool) error {
	_, err := s.db.Exec("UPDATE quote SET is_archived = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", v, id)
	return err
}

func (s *QuoteStore) Random() (Quote, error) {
	row := s.db.QueryRow("SELECT " + quoteCols + " FROM quote WHERE is_archived = 0 ORDER BY RANDOM() LIMIT 1")
	return scanQuote(row.Scan)
}

// --- HTTP handlers ---

func parseQuoteID(s string) (int, bool) {
	id, err := strconv.Atoi(strings.TrimSuffix(s, "/"))
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func validateQuote(q Quote) string {
	if strings.TrimSpace(q.Text) == "" {
		return "Quote text is required."
	}
	if strings.TrimSpace(q.Author) == "" {
		return "Author is required."
	}
	return ""
}

type quotePage struct {
	Quotes          []Quote
	Search          string
	IncludeArchived bool
	Total           int
	Archived        int
}

type quoteFormPage struct {
	Quote Quote
	IsNew bool
	Error string
}

func (a *App) handleQuotes(w http.ResponseWriter, r *http.Request) {
	if a.quotes == nil {
		http.NotFound(w, r)
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/admin/quotes/")
	rest = strings.TrimSuffix(rest, "/")
	switch {
	case rest == "":
		a.handleQuoteList(w, r)
	case rest == "add":
		a.handleQuoteAddForm(w, r)
	case rest == "new":
		a.handleQuoteCreate(w, r)
	case strings.HasSuffix(rest, "/save"):
		a.handleQuoteSave(w, r, strings.TrimSuffix(rest, "/save"))
	case strings.HasSuffix(rest, "/del"):
		a.handleQuoteDelete(w, r, strings.TrimSuffix(rest, "/del"))
	case strings.HasSuffix(rest, "/delete"):
		a.handleQuoteConfirm(w, r, strings.TrimSuffix(rest, "/delete"))
	case strings.HasSuffix(rest, "/archive"):
		a.handleQuoteArchive(w, r, strings.TrimSuffix(rest, "/archive"))
	case strings.HasSuffix(rest, "/favorite"):
		a.handleQuoteFavorite(w, r, strings.TrimSuffix(rest, "/favorite"))
	default:
		a.handleQuoteEditForm(w, r, rest)
	}
}

func (a *App) handleQuoteList(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")
	includeArchived := r.URL.Query().Get("archived") == "1"
	list, err := a.quotes.List(search, includeArchived)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	total, _ := a.quotes.Count()
	archived, _ := a.quotes.CountArchived()
	a.render(w, "quotes_list.html", quotePage{
		Quotes:          list,
		Search:          search,
		IncludeArchived: includeArchived,
		Total:           total,
		Archived:        archived,
	})
}

func (a *App) handleQuoteAddForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, "quote_edit.html", quoteFormPage{IsNew: true})
}

func (a *App) handleQuoteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/quotes/", http.StatusFound)
		return
	}
	r.ParseForm()
	q := Quote{
		Text:     r.FormValue("text"),
		Author:   r.FormValue("author"),
		Category: r.FormValue("category"),
		Tags:     r.FormValue("tags"),
	}
	if msg := validateQuote(q); msg != "" {
		a.render(w, "quote_edit.html", quoteFormPage{Quote: q, IsNew: true, Error: msg})
		return
	}
	if _, err := a.quotes.Create(q); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/quotes/", http.StatusSeeOther)
}

func (a *App) handleQuoteEditForm(w http.ResponseWriter, r *http.Request, idStr string) {
	id, ok := parseQuoteID(idStr)
	if !ok {
		http.NotFound(w, r)
		return
	}
	q, err := a.quotes.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	a.render(w, "quote_edit.html", quoteFormPage{Quote: q})
}

func (a *App) handleQuoteSave(w http.ResponseWriter, r *http.Request, idStr string) {
	id, ok := parseQuoteID(idStr)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/quotes/", http.StatusFound)
		return
	}
	if _, err := a.quotes.Get(id); err != nil {
		http.NotFound(w, r)
		return
	}
	r.ParseForm()
	q := Quote{
		ID:       id,
		Text:     r.FormValue("text"),
		Author:   r.FormValue("author"),
		Category: r.FormValue("category"),
		Tags:     r.FormValue("tags"),
	}
	if msg := validateQuote(q); msg != "" {
		a.render(w, "quote_edit.html", quoteFormPage{Quote: q, Error: msg})
		return
	}
	if err := a.quotes.Update(q); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/quotes/", http.StatusSeeOther)
}

func (a *App) handleQuoteConfirm(w http.ResponseWriter, r *http.Request, idStr string) {
	id, ok := parseQuoteID(idStr)
	if !ok {
		http.NotFound(w, r)
		return
	}
	q, err := a.quotes.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	a.render(w, "quote_confirm_delete.html", q)
}

func (a *App) handleQuoteDelete(w http.ResponseWriter, r *http.Request, idStr string) {
	id, ok := parseQuoteID(idStr)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/quotes/", http.StatusFound)
		return
	}
	if err := a.quotes.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/quotes/", http.StatusSeeOther)
}

func (a *App) handleQuoteArchive(w http.ResponseWriter, r *http.Request, idStr string) {
	id, ok := parseQuoteID(idStr)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/quotes/", http.StatusFound)
		return
	}
	q, err := a.quotes.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := a.quotes.SetArchived(id, !q.IsArchived); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/quotes/", http.StatusSeeOther)
}

func (a *App) handleQuoteFavorite(w http.ResponseWriter, r *http.Request, idStr string) {
	id, ok := parseQuoteID(idStr)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/quotes/", http.StatusFound)
		return
	}
	q, err := a.quotes.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := a.quotes.SetFavorite(id, !q.IsFavorite); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/quotes/", http.StatusSeeOther)
}

func (a *App) handleQuotesAPI(w http.ResponseWriter, r *http.Request) {
	if a.quotes == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	rest := strings.TrimPrefix(r.URL.Path, "/api/quotes")
	switch rest {
	case "/random":
		q, err := a.quotes.Random()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				json.NewEncoder(w).Encode(map[string]string{"error": "No quotes available"})
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(q)
	default:
		http.NotFound(w, r)
	}
}