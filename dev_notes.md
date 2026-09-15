# Portfolio Site — dev notes

> The gemstone site — 4rch3.io portfolio. Static Bootstrap "Fluid Mind"
> template served by nginx; real project data rendered client-side from the
> Go admin API. Deployed when ≥80% complete.

## Milestones

- [x] **M1 — Scaffold & theme**: Astro 7 minimal + strict TS; global dark-cyber
      theme tokens (surfaces, glow accents, pillar/status colors); favicon.
- [x] **M2 — Content model**: `glob` content collection + zod schema
      (title, pillar, tagline, status, order, featured, stack, highlights).
- [x] **M3 — Pages**: home hero + featured + per-pillar sections; `/projects`
      filterable grid (vanilla TS, `?pillar=` deep-link); `/projects/[slug]`
      deep pages; `/about`.
- [x] **M4 — Content seed**: 26 project entries across the four pillars with
      tagline, stack, highlights, and full narratives.
- [x] **M5 — Verify**: clean `astro build` (29 pages), sitemap generated,
      preview + routes smoke-tested.
- [x] **M5b — Containerize + CMS**: nginx web container (multi-stage) +
      Go admin CMS container (single static binary, Django-admin-style UI);
      `docker compose up` serves site + `/admin/` on one port (`:8095`),
      optional direct admin `:8096`. The admin reads/writes the Markdown
      source in place — no database, no export step; `make sync` rebuilds
      the site container from the edited content. Verified end-to-end.
- [x] **M5c — Frontend rewrite**: replaced the Astro site with the
      `frontend/` Bootstrap template (BootstrapMade "Personal" v2.3.0 base).
      `web` image now just copies `frontend/` into nginx (no node build).
      nginx proxies `/api/` → admin. `assets/js/projects.js` renders the
      26-project grid + pillar filters from `/api/projects/`;
      `assets/js/project-details.js` renders `portfolio-details.html?slug=`
      deep pages. jsdom-tested (16/16, incl. error fallbacks). Branding →
      "4rch3".
- [x] **M5d — Quotes admin**: the Go admin manages the `ai-quotes` project's
      SQLite DB directly (pure-Go `modernc.org/sqlite`, no CGO). New
      `/admin/quotes/` section (list/search, show archived, add, edit,
      save, delete, favorite/archive toggles) + public `/api/quotes/random`
      endpoint. Enabled only when `AI_QUOTES_DB` is set (compose mounts
      `../ai-quotes/instance/` → `/quotes-instance`). Verified: full local
      CRUD cycle + live round-trip through :8095 on the real DB (add→show
      in list→edit→favorite→archive→delete), DB left exactly as found.
- [x] **M5e — Hero tagline quotes**: `hero-quote.js` fetches
      `/api/quotes/random` on page load and repaints `#header h2` with the
      quote + `— author`. 1–3 random words get the template's green
      underline (`span`), the author span is styled to NOT be underlined,
      and failures keep the static tagline. Rotation is refresh-only (as
      agreed). Reviewed with the user before building; jsdom-tested (8 new
      checks: text, author, underline range, fail + empty-body fallbacks).
- [x] **M5f — Contact email notify + SEO metadata (2026-09-12)**:
      (a) contact-form SMTP notify in the Go admin (net/smtp, STARTTLS,
      fire-and-forget goroutine, header-injection-safe, Reply-To =
      visitor; JSONL inbox stays the durable record). Env-gated via
      SMTP_HOST/CONTACT_EMAIL in `.env`; Gmail app password pending from
      user. Verified with an in-network python SMTP sink (full
      transcript: From/To/Subject/Reply-To/body correct, visitor `OK`
      immediate). (b) OG + Twitter `summary_large_image` + JSON-LD
      (Person + WebSite) on all three pages; fixed the template's broken
      `descriptison` meta tag; generated 1200x630 `og-image.png` via
      `scripts/gen_og_image.py` (Pillow CRT-brand render). Also: `.env`
      created with real admin creds (default admin123 retired), repo
      git-initialized, README rewritten to match the real stack.
- [ ] **M6 — Visual polish**: ~~OG images, open-graph + JSON-LD metadata~~
      (done in M5f). ~~Project screenshots/diagrams~~ (done 2026-09-12 —
      three real pipelines: `make dashboards` Playwright captures of the
      live homelab services, `make terminal` real-ANSI terminal renders
      from `scripts/terminal_captures.json`, `make diagrams` mermaid
      renders from `src/content/diagrams/*.mmd`; detail pages prefer
      screenshots + append diagrams automatically). Remaining: accessibility
      pass; authenticated augur capture (login wall is what we have — no
      creds stored); extend terminal commands + diagrams to more projects.
      Still needed before the 80% deploy threshold.
- [x] **M6b — Accessibility + contact abuse controls (2026-09-12)**:
      (a) axe-core audit pipeline (`make a11y`, vendored axe 4.10.2,
      `scripts/a11y_audit.py`, exit-1 gate) — fixed all 9 baseline
      violation types to **0** across home/research/detail pages: body
      solid bg (the `background: transparent` reset made auditors see
      white — root cause of most contrast flags), CRT glow via
      filter/drop-shadow instead of text-shadow, flicker opacity floor
      0.88, `<main>` landmark on home, named icon-only links
      (hero socials, venobox triggers), form input aria-labels, heading
      order fixes (cards h2/h3, detail sections h2, research h1),
      scrollable `#portfolio-details` tabindex+region,
      prefers-reduced-motion support. (b) Contact hardening: per-IP
      fixed-window rate limit (5/h default, env-tunable, 429+Retry-After),
      hidden honeypot field (silent drop + log), nginx X-Real-IP/XFF on
      `/api/` + `clientIP` prefers X-Real-IP (was recording the web
      container's IP for every visitor). All verified live: 5×200→429,
      honeypot dropped, real IPs recorded.
- [ ] **M7 — Deploy (external)**: expose via the proper clu5t3r + Traefik
      pipeline — `4rch3.io` apex + `www` → portfolio, `homelab.4rch3.io` →
      portal (see routing plan below); public DNS + router 80/443 forwarding.
      Deferred by user — Traefik already running in the homelab; nothing is
      publicly reachable today (even 4rch3.io resolves only on the LAN).
      NOT using router port-forwards on m41n / Tailscale Funnel as interim.
- [x] **M7a — Project repos on Forgejo + repo/live wiring**: all 27
      portfolio projects backed up on Forgejo, **all private**. 20 new
      repos (16 fresh `git init` with tailored .gitignore — Rust `target/`,
      Python `.venv/__pycache__/*.db` — plus hermes/rust-infer which had
      zero commits, plus n3xusDB/transitflow-nyc which keep GitHub `origin`
      and gained a `forgejo` remote); augur/shadowsim/threat-pulse/homelab
      snapshot-committed + pushed; logsentry had diverged 3v3 — merged
      forgejo/main (web-UI serve mode + health server) into local main,
      docker-compose.yml conflict resolved keeping local port-removal +
      remote web-ui labels (8081 + /health), pushed 05ac94e — the deployed
      copy on clu5t3r is now behind Forgejo main. Pre-existing repos
      flipped private (incl. 4rch3.io, alertflow, m0rpheus, misp,
      portfolio). Frontmatter: `repo:` on all 27 (names follow source
      dirs: threat-pulse, n3xusDB, homelab; sentryd → logsentry repo since
      its code lives in `logsentry/sentryd/`; analyst-agent → repo pending,
      source lives on clu5t3r) and `live: https://homelab.4rch3.io` on
      lab.md — the first live link; it 404s at Traefik until the routing
      move lands. Wiring: `.portfolio-repo` href from `p.repo`, new
      `.portfolio-live` card icon (bx-link-external, rendered only when
      `live` is set); details pages already render Repository/Live rows
      from frontmatter. a11y gate caught the new repo row
      (`link-in-text-block`, serious) → details-list links underlined,
      audit's stale `tr4c3` slug replaced with `threatpulse`; 0 violations
      across 4 pages. Cache-busters: projects.js v3, style.css v13.
- [x] **M7b — Routing move executed (LAN)**: portal on clu5t3r re-labeled
      `Host(homelab.4rch3.io) || Host(portal.4rch3.io)`, renamed "Homelab
      Portal", `config.yaml` `external_url` + CORS origin →
      `https://homelab.4rch3.io`. Portfolio deployed to clu5t3r via
      `docker-compose.traefik.yml` (web: apex + www, security-headers +
      rate-limit; admin: `admin.4rch3.io` + cf-access-auth; no published
      ports; no ai-quotes mount — `/api/quotes/random` 404s on prod by
      design, hero falls back to the static tagline; quote management
      stays on m41n). Verified: 4rch3.io/www → portfolio 200,
      homelab/portal → portal 200, threatpulse/eventflow untouched 200,
      admin login 200. Prod inbox volume hit the uid-1000 chown gotcha
      (fresh volume root-owned) — fixed, and seeded with the one real
      contact message from m41n. clu5t3r's logsentry checkout
      fast-forwarded to the merged Forgejo main (05ac94e) — all three
      machines agree now. **analyst-agent source not found** on
      m41n/clu5t3r/n4n0 (4l13n offline) — repo still pending, location
      to be confirmed with user. Remaining M7: public DNS + router
      port-forward (item 4).
- [ ] **M7c — Public exposure (item 4, in progress — paused for the
      night)**: Facts: Cloudflare zone had a PUBLIC wildcard
      `*.4rch3.io → 192.168.1.33` (private IP published; currently
      load-bearing — the router is a pure forwarder, so ALL LAN
      *.4rch3.io resolution rides on it; m41n has /etc/hosts overrides
      for apex/portal only, clu5t3r has none). Public IP 66.108.23.33
      (not CGNAT; staticness unknown). Done so far: split-horizon
      **dnsmasq deployed on clu5t3r 192.168.1.33:53**
      (`address=/4rch3.io/192.168.1.33`, upstreams 1.1.1.1/9.9.9.9,
      in homelab repo infrastructure compose — published on the
      specific LAN IP because dockerd bridge DNS + systemd-resolved
      hold other :53 binds). infrastructure/.env had gone EMPTY
      (values only in running containers) — rebuilt from container
      state (DOMAIN, CLOUDFLARE_EMAIL/API/TUNNEL_TOKEN). Homelab repo
      unified 3-way (clu5t3r running tree + m41n snapshot merged at
      b9e690a; conflicts: running tree won for traefik.yml [metrics
      entrypoint] + composes [dnsmasq, siem.dev wazuh rule], m41n
      docs won for INGRESS.md). Decisions: tunnel path chosen
      (cloudflared already runs on proxy-net; ingress →
      http://infrastructure-traefik-1:80 so Traefik routing +
      middlewares stay in the loop; zero open ports, no DDNS needed);
      dev.4rch3.io public records kept per user. CF_ACCESS_ENFORCE is
      false (log-only) — needs Zero Trust Access app + TEAM/AUD/true
      before admin.4rch3.io is meaningfully protected. Next steps
      (awaiting user): (1) add 4 Public Hostnames in Zero Trust →
      tunnel (4rch3.io, www, homelab, admin → HTTP
      infrastructure-traefik-1:80 — auto-creates proxied CNAMEs), or
      hand over a Zero Trust-scoped API token to do it via API;
      (2) SSL/TLS → Full (strict); (3) router DHCP DNS → 192.168.1.33
      (unblocks deleting the public wildcard; until then LAN keeps
      riding the wildcard, tunnel names hairpin harmlessly);
      (4) Access app for admin.4rch3.io → then set CF_ACCESS_TEAM/
      AUD/ENFORCE on the verifier. Wildcard deletion + final public
      verification = mine, after (1)+(3).
- [ ] **M7d — 2026-09-15: prod bug triage + portal auth gate (in
      progress)**: User reported three issues after the M7b/M7c deploy —
      all diagnosed and fixed in prod today.
      (1) **Portfolio hero losing its background + "refresh lands on the
      Contact page"**. Root cause: `rate-limit@file` (avg 100, burst 50,
      1m) was applied to the whole `web` router; a single page load bursts
      ~50 requests (27 SVG logos + assets), so on refresh assets (incl.
      `bg.jpg` + logos) got **429s** → no background + broken grid.
      The "Contact page" weirdness was browser scroll-restoration
      restoring a viewport into the hidden-section overlay (template shows
      sections one-at-a-time behind a 100vh hero; no JS writes hashes, but
      it reads `location.hash`). Fixes: Traefik router split — static
      assets stay on the bare-host router (security-headers only), new
      `portfolio-api` router (`Host && PathPrefix(/api,/admin,/admin-static)`)
      keeps rate-limit (Traefik rule-length priority picks it);
      `main.js` sets `history.scrollRestoration=manual` + scrolls to top
      on pageshow when no hash. (2) **Projects grid collapsed/overlapping
      until first tab click**: isotope `fitRows` computed layout before the
      800×600 SVG logos (no intrinsic size) loaded → rows collapsed +29px.
      Fixes: `.project-logo` reserves `aspect-ratio: 4/3` in CSS +
      `projects.js` re-layouts isotope on img load and window load.
      (3) **Portal cards forever-"loading"**: dashboard rendered to guests
      but every card action (`/api/service-panel`, `/api/services/status`,
      `/api/deploy-panel`) is `RequireAuth` → 401 on the htmx request →
      modal stuck on its loading placeholder. `status_page_public: false`
      existed but was **never enforced** (dead config). Decided + done:
      enforce it — `Index` behind login in
      `internal/handlers/dashboard.go` (publicStatus field + ctor param,
      guests → 303 `/login`), `cmd/server/main.go` passes
      `cfg.Monitor.StatusPagePublic`. Committed "dashboard: gate index
      behind auth when status_page_public=false", pushed Forgejo 4rch3.io
      main. **Deploy note on portal**: clu5t3r `~/4rch3.io` is NOT a git
      repo (plain build dir — `git reset --hard` → "not a git
      repository"); deploy = rsync/scp changed sources over + `docker
      compose build portal && docker compose up -d portal` (multi-stage
      Dockerfile compiles in a golang builder container, no Go toolchain
      needed on clu5t3r). Verified on prod: guest `/` → 303 → `/login`
      (local + via hosts), `/api/health` still public, wrong-creds login
      round-trips properly (CSRF+htmx work). Metrics panels render once an
      admin is logged in (services have docker_ids). Portfolio prod
      re-verified with Playwright: 3 clean grid rows (y=314/553/792),
      logo slots 279×209 (4:3), hero bg painting (25,680 unique colors),
      0 asset 4xx/429s (only `/api/quotes/random` 404 — by design, no
      ai-quotes mount on prod), refresh stays at scrollY=0 on the hero,
      `make a11y` still 0 violations (4 pages). Portal admin password
      still unknown → user should log in to homelab.4rch3.io and confirm
      cards/panels/charts render.

## Routing plan (agreed with user — execute at M7)

Current state on **clu5t3r** (`192.168.1.33`, Traefik v3.6.6 on :80/:443,
cloudflare certresolver, docker provider + file provider at
`/home/x5pyd3rx/projects/homelab/configs/traefik`):
- `Host(4rch3.io) || Host(portal.4rch3.io)` → **4rch3-portal** (Go homelab
  ops hub, backend :8081). Middlewares: security-headers, rate-limit,
  cf-access-auth (all `@file`).
- Custom services discovered via `4rch3.io.service=` docker labels
  (threatpulse, augur, eventflow, forgejo, thehive, cortex, wazuh, …) and
  reverse-proxied as `<service>.4rch3.io` by the portal — **unchanged**.

Requested changes (3):
1. **`4rch3.io` → the portfolio site** (web container), not the portal.
2. **`homelab.4rch3.io` → the portal/homelab hub** (portal moves off apex).
3. **Custom project subdomains keep their current routes.**

Execution checklist:
- [x] Deploy portfolio stack to clu5t3r (fleet convention — a compose
      project on clu5t3r with Traefik labels; NOT a file-route back to m41n).
      `docker-compose.traefik.yml` — compose project `portfolio`, checked out
      at `~/projects/portfolio` on clu5t3r (cloned from Forgejo with a
      one-time token URL; no credential store on that host).
- [x] `web` Traefik router: `Host(4rch3.io) || Host(www.4rch3.io)`,
      entrypoint websecure, cloudflare certresolver, middlewares
      `security-headers@file,rate-limit@file`.
- [x] `admin` Traefik router: `Host(admin.4rch3.io)`, + `cf-access-auth@file`
      (keep the admin off the plain web).
- [x] Portal: router rule → `Host(homelab.4rch3.io)`; keep `portal.4rch3.io`
      as a redirect alias; portal `config.yaml` `external_url` →
      `https://homelab.4rch3.io` (keeps login/recorded links correct).
- [ ] DNS (Cloudflare): `4rch3.io`, `www.4rch3.io`, `homelab.4rch3.io`,
      `admin.4rch3.io`, plus **wildcard `*.4rch3.io`** (portal serves
      `<service>.4rch3.io` dynamically) → public IP.
- [ ] Router: forward TCP 80/443 → clu5t3r.
- [x] Service discovery labels (`4rch3.io.*` on threat-pulse, augur, …):
      untouched — their routes are unaffected by the apex/portal moves.
- [x] Verify: `4rch3.io` loads portfolio; `homelab.4rch3.io` loads portal;
      `threatpulse.4rch3.io` / `eventflow.4rch3.io` still resolve as today.

### Frontend rewrite (#4) — DONE (frontend is the `frontend/` template)
- The Go admin and the frontend only talk through the **content contract**,
  now via the API: the Go admin serves `/api/projects/` (schema: slug, title,
  pillar, tagline, status, order, featured, path, repo, live, stack,
  highlights, body) from the Markdown frontmatter files in
  `src/content/projects/`. The admin writes these in place.
- Frontend = plain static files in `frontend/` (nginx html root). nginx
  proxies `/api/` → admin:8000; `projects.js` and `project-details.js`
  fetch and render client-side. The admin edit UI keeps working unchanged.
- Deferred (recorded): `research.html` will list 2-3 research topics
  (Quantum Computing, Liquid Neural Networks, …); dark-cyber restyle of the
  template is still pending. (Hero tagline quotes are now LIVE — see M5e.)
- Keep the admin container bind-mount to the same `src/content/projects`
  path when relocating to clu5t3r (M7).

## Admin/CMS — known quirks

- Go template namespaces: every page template must be executed as its own
  file template (`ExecuteTemplate(w, "list.html", ...)`). Shared `{{define
  "title"}}` blocks across files collide — pass titles via data instead.
- ServeMux wildcard routes (`{slug}`) need `go >= 1.22` in `go.mod`; below
  that they're treated as literal path segments. We route `/admin/projects/`
  through one dispatcher to avoid wildcard/literal registration conflicts.
- The admin runs as uid 1000 to match the host user and write the
  bind-mounted `src/content/projects` Markdown in place.
- nginx proxies `/admin-static/` (Go admin's embedded CSS) — the old Django
  path was `/static/`. The `/api/` prefix is also proxied → admin.
- `:8095` (site + admin/API proxy) and `:8096` (direct admin) both bind
  `0.0.0.0`. Works on the LAN as-is; change the ports via `.env`.

### Quotes (ai-quotes SQLite) — gotchas
- **Directory mount, not file mount.** SQLite needs to create its `-journal`
  sibling next to the DB, so the container bind-mounts the whole `instance/`
  dir (`../ai-quotes/instance` → `/quotes-instance:rw`). A bare file mount
  puts the journal next to the file — in `/` of the container, which uid
  1000 cannot write — and SQLite then silently opens the DB **read-only**
  ("attempt to write a readonly database" on any write, reads still work).
- **NULL timestamps break scans.** The Flask DB schema leaves
  `created_at`/`updated_at` unset for direct SQL inserts (SQLAlchemy sets
  them at app level). Scanning those into a Go `string` fails the row read;
  `scanQuote` uses `sql.NullString`, and `CREATE` sets both columns to
  `CURRENT_TIMESTAMP` so new rows are never NULL.
- **`modernc.org/sqlite` needs `go >= 1.25`** (its deps bumped the toolchain
  requirement). The admin `Dockerfile` now builds from `golang:1.25-alpine`.
- Read-only "quotes" features (favorite/archive state, views, votes) follow
  the ai-quotes public API contract; display-only fields are not edited here.

## Frontend — known quirks

- The portfolio grid is rendered client-side, so isotope + venobox are
  initialized in `projects.js` after data loads (main.js no longer inits
  them on page load).
- Deep links are `portfolio-details.html?slug=<slug>`; an unknown/absent
  slug shows a "Project not found" fallback (does not silently fall back to
  the first project).
- If `/api/projects/` is unreachable (admin down) the grid section shows an
  error banner instead of an empty grid.
- `resume.html` and `research.html` are empty stubs for now; the nav
  "Research" item is a placeholder until the research page is built.
- The hero tagline is rendered client-side by `hero-quote.js` from
  `/api/quotes/random`; the static `h2` in `index.html` is the no-JS /
  fallback content. 1–3 random words get the underline (`.hero-author` span
  is explicitly un-underlined).

## Definitions of done (frontend)

- `docker compose build web` exits 0; the stack comes up clean.
- `frontend/` pages return 200 with real content from `/api/projects/`
  (grid shows 4 pillars + every project; detail pages resolve by slug).
- `/admin/` + `/admin-static/` proxied on `:8095`; admin CRUD still edits
  the Markdown source in place.
- Every project entry has a ≥100-word narrative and ≥3 highlights (content
  contract unchanged).

## Known quirks

- YAML (hand-written files only): a highlight item containing `: ` must be
  single-quoted or it parses as a key-value map and fails the string schema.
  The CMS export path auto-quotes, so content edited in the admin is safe.
- Add a new tagline/`live` link only when a project is actually deployed —
  the site must stay honest.

## Commands

```bash
docker compose up -d --build  # full stack: site :8095 + admin :8096
docker compose build web      # rebuild site only (frontend/ -> nginx)
docker compose build admin    # rebuild the Go admin only

make up                   # docker compose up -d --build
make sync                 # rebuild web after admin edits Markdown in place
make logs / make ps       # logs / status
make clean                # down + remove volumes and local images
```