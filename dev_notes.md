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
      (done in M5f). Remaining: project screenshots/diagrams,
      accessibility pass. Still needed before the 80% deploy threshold.
- [ ] **M7 — Deploy (external)**: expose via the proper clu5t3r + Traefik
      pipeline on `portfolio.4rch3.io` (public DNS + router 80/443 forwarding).
      Deferred by user — Traefik already running in the homelab; nothing is
      publicly reachable today (even 4rch3.io resolves only on the LAN).
      NOT using router port-forwards on m41n / Tailscale Funnel as interim.

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
- [ ] Deploy portfolio stack to clu5t3r (fleet convention — a compose
      project on clu5t3r with Traefik labels; NOT a file-route back to m41n).
      See `docker-compose.traefik.yml`/labels planned for `web`.
- [ ] `web` Traefik router: `Host(4rch3.io) || Host(www.4rch3.io)`,
      entrypoint websecure, cloudflare certresolver, middlewares
      `security-headers@file,rate-limit@file`.
- [ ] `admin` Traefik router: `Host(admin.4rch3.io)`, + `cf-access-auth@file`
      (keep the admin off the plain web).
- [ ] Portal: router rule → `Host(homelab.4rch3.io)`; keep `portal.4rch3.io`
      as a redirect alias; portal `config.yaml` `external_url` →
      `https://homelab.4rch3.io` (keeps login/recorded links correct).
- [ ] DNS (Cloudflare): `4rch3.io`, `www.4rch3.io`, `homelab.4rch3.io`,
      `admin.4rch3.io`, plus **wildcard `*.4rch3.io`** (portal serves
      `<service>.4rch3.io` dynamically) → public IP.
- [ ] Router: forward TCP 80/443 → clu5t3r.
- [ ] Service discovery labels (`4rch3.io.*` on threat-pulse, augur, …):
      untouched — their routes are unaffected by the apex/portal moves.
- [ ] Verify: `4rch3.io` loads portfolio; `homelab.4rch3.io` loads portal;
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