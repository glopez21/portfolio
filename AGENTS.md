## Stack

Static frontend + Go admin CMS, containerized:

- `web` (nginx): serves `frontend/` as static files. Proxies `/api/` and
  `/admin/` + `/admin-static/` to the `admin` container.
- `admin` (Go): Django-admin-style CMS. Exposes `/api/projects/` (JSON read
  from the Markdown source) and the edit UI. Writes
  `src/content/projects/*.md` **in place** — no database, no export step.
- `admin` (quotes): also manages the `ai-quotes` project's SQLite DB directly
  (pure-Go `modernc.org/sqlite`; enabled only when `AI_QUOTES_DB` is set —
  compose mounts `../ai-quotes/instance/` → `/quotes-instance`). Exposes
  `/admin/quotes/` (CRUD + favorite/archive toggles) and a public
  `/api/quotes/random` endpoint.
- `admin` (dashboard): `/admin/` is a dashboard landing — a stats row plus a
  responsive card grid, one card per pillar (project/featured/quote/contact
  counts). Each pillar card's title links to its per-pillar page.
- `admin` (per-pillar pages): `/admin/pillars/<pillar>/` (auth) shows that
  pillar's projects as an evenly-arranged card grid — one equal card per
  project, tiles across the viewport instead of one vertical scroll marathon.
- `admin` (list pages are card grids too): `/admin/projects/`,
  `/admin/quotes/`, and `/admin/contact/` all use the same
  `auto-fill minmax(300px, 1fr)` tile layout as the dashboard and pillar
  pages — one card per project/quote/message, no `<table>`, no vertical
  marathon. The shared `.card-grid` / `.card-inner` rules live in
  `admin-go/static/admin.css` (a card is a `.module` with a title, meta
  badges, and a pinned actions row). The quotes list keeps its search +
  show-archived tools above the grid.
  Pillar order + human labels come from one shared source (`pillarOrder` +
  `pillarLabel` in `admin-go/main.go`, fed to templates via the `pillarLabel`
  FuncMap) so the dashboard and pillar pages never drift.
- `admin` (contact): the site's contact form posts urlencoded data to
  `POST /api/contact/`, which appends one JSON line per message to
  `CONTACT_FILE` (`/data/contact_messages.jsonl`, named `admin_data`
  volume) and returns the literal body `OK` (the frontend's
  `php-email-form` validate.js expects exactly `OK`). Messages are readable
  at `/admin/contact/` (auth, newest first). The `admin_data` volume is
  owned by uid 1000 — re-`chown` it if the volume ever gets recreated as
  root and store starts failing.
- `admin` (contact owner-notify): **SMTP-relay-free by design.** The JSONL
  inbox above is the durable owner-facing record, and replies go out of the
  owner's own mailbox `CONTACT_EMAIL` (a `@pm.me` Proton address,
  independent of this site's MX). A Gmail relay was wired briefly but
  removed — Gmail rewrites `From` to its authenticated user, so a `@pm.me`
  `From` could never be sent through it. The Go notify code is still
  gated on empty `SMTP_HOST` (documented in `admin-go/main.go`) and stays
  off; no SMTP creds live in `.env` or the compose files.

## Development

No dev server. Edit files directly and rebuild the `web` image:

```
make sync        # docker compose build web && up -d web
docker compose up -d   # full stack
```

Verify locally on `:8095` (site + admin/API proxy) and `:8096` (direct
admin). Admin credentials live in `.env` (real password, not the old
default).

## Git

Origin is Forgejo on the homelab: `git.4rch3.io/61-72-6b-68-e9/portfolio`
(LAN-only). A **push mirror** on the Forgejo side auto-syncs every push to
the public GitHub mirror `github.com/glopez21/portfolio` — so `make push`
(git push origin main) is all that's needed; never push to GitHub directly.
`.env` is git-ignored; tokens live in the host credential store.

## Project repositories

Every portfolio-listed project is backed up on Forgejo, **all private** —
`repo:` frontmatter on each `src/content/projects/*.md` points at
`https://git.4rch3.io/61-72-6b-68-e9/<name>`. Repo names follow the source
directory, not the slug (e.g. slug `threatpulse` → repo `threat-pulse`,
`n3xusdb` → `n3xusDB`, `lab` → `homelab`). Exceptions: `sentryd` points at
the `logsentry` repo (its code lives in `logsentry/sentryd/`), and
`analyst-agent` points at a repo still to be created from the clu5t3r copy.
`n3xusDB` and `transitflow-nyc` also keep their GitHub `origin`; their
Forgejo remote is named `forgejo` (same convention as all project repos).
README polishing and going public is deferred; the GitHub side is mirrored
manually by the user. Note: `logsentry` on clu5t3r may lag its Forgejo
main (m41n holds a merge of the web-UI + health-server lines).

## Production (clu5t3r)

`docker-compose.traefik.yml` deploys the same stack onto the clu5t3r
Traefik fleet: compose project `portfolio`, checked out at
`~/projects/portfolio` on clu5t3r (cloned from Forgejo — no credential
store on that host; the clone URL embeds a token, same pattern as
threat-pulse on m41n). Differences from the dev compose:

- **No published ports** — Traefik fronts both services on `proxy-net`:
  `web` → `Host(4rch3.io) || Host(www.4rch3.io)` with
  security-headers + rate-limit; `admin` → `Host(admin.4rch3.io)` plus
  `cf-access-auth@file`.
- **Quotes DB is mirrored to prod** — m41n's `../ai-quotes/instance/ai_quotes.db`
  is the canonical copy; it's scp'd to clu5t3r `./quotes/` and mounted as
  `/quotes-instance` (`AI_QUOTES_DB`), so `/api/quotes/random` works on prod.
  After editing quotes on m41n, run `make deploy-quotes` (mirror + chown) and
  recreate the admin container. `quotes/` is gitignored.
- Prod `.env` was scp'd from m41n (same admin password, SMTP empty),
  chmod 600. Deploy/update: `docker compose -f docker-compose.traefik.yml
  up -d --build` on clu5t3r. Content edits still flow through the admin
  UI or git pull + recreate.
- The **fresh `admin_data` volume arrives root-owned** — same uid-1000
  chown gotcha as dev: `docker exec -u 0 portfolio-admin chown -R
  1000:1000 /data` on first deploy. Prod inbox was seeded with the one
  real contact message from m41n.
- Portal context: the homelab portal moved to `homelab.4rch3.io`
  (+ `portal.4rch3.io` alias, "Homelab Portal", `external_url` updated)
  when the portfolio took the apex — see dev_notes M7b.



## Frontend layout

- `frontend/index.html` — one-page site (hero, about, skills, portfolio,
  contact). Header brand is **4rch3.io** (`.cyber-brand`,
  Share Tech Mono, green CRT glow + flicker; `.brand-tag` Greek
  "αρχή · η πρώτη αρχή" sub-tagline in IBM Plex Mono; quote tagline below).
  When scrolled, `#header.header-top` becomes a fixed bar and hides the
  tagline (`h2`, `.social-links`, `.brand-tag`); the Greek tagline only
  shows in the hero under the title.
- Inner pages (`research.html`, `portfolio-details.html`) use a compact
  sticky `#inner-header` navbar: **4rch3.io** brand + Home/About/Portfolio/
  Research/Contact menu, no hero tagline/quote.
- `frontend/research.html` — research topics page, fully client-side:
  `assets/js/research.js` holds the `TOPICS` and `LOG` data and renders the
  cards, hidden note modals (`#research-notes`), and the Reading Log
  (`#research-log`). Each card opens a venobox **inline** modal with the
  note body (paragraphs + `term` chips + a `From research → build` line
  tying the topic to projects + optional `note-links`). Status per topic:
  `.status-reading` (green) / `.status-experimenting` (turquoise). The
  whole card is clickable (triggers its hidden `.read-note` venobox link).
  `.vbox-inline` is restyled as a dark card (`min(640px, 94vw)`,
  `max-height:78vh`, overflow-y auto); note links are added by filling the
  `links` arrays in `research.js` (kept empty since real URLs are up to
  the user — don't fabricate).
- `frontend/portfolio-details.html` — detail page, opened as
  `portfolio-details.html?slug=<slug>`.
- `assets/js/projects.js` — fetches `/api/projects/`, renders the grid +
  pillar filters, then inits isotope/venobox (client-side rendering).
  `homelab`-pillar items render as a research-style `.card-panel`
  (icon + title + tagline) inside the standard `.portfolio-wrap`, keeping
  the same hover/lightbox effects.
  Every other project renders a generated SVG **project logo**
  (`assets/img/project-logos/<slug>.svg`) instead of a photo. Non-lab
  cards are shown at 82% of their column width (`.portfolio-item:not(.lab-item)
  .portfolio-wrap`). The `portfolio-links` hover icons are: a GitHub/git mark
  (`.portfolio-repo`, wired to the project's `repo` frontmatter — Forgejo,
  private repos, see "Project repositories"), an optional live-site mark
  (`.portfolio-live`, `bx-link-external`, rendered only when `live` is set —
  currently just The Lab → `https://homelab.4rch3.io`), and the unchanged
  details-link venobox (`portfolio-details.html?slug=…`).
- `assets/js/project-details.js` — fetches a single project by slug,
  populates the detail page. Unknown slug shows a "not found" fallback.
  The hero image (`#portfolio-hero`, styled at 24%/230px as a logo tile)
  prefers a **real screenshot** (`assets/img/screenshots/<slug>.png`);
  on load it gets `.is-shot` (larger presentation, `min(560px,100%)`).
  Missing screenshot → `onerror` swaps back to the generated
  `assets/img/project-logos/<slug>.svg` tile. The inline pre-JS setter in
  `portfolio-details.html` implements the same contract (`data-shots` +
  `data-base`), so there is no race and no broken-image paint. The About
  column appends an **Architecture diagram** img
  (`assets/img/diagrams/<slug>.png`) + heading, removed via `onerror`
  when no diagram exists for the slug.
- `portfolio-details.html` detects when it runs inside the venobox iframe
  (`window.self !== window.top`) → hides the `#inner-header` navbar and adds
  `body.in-modal`, which lays the content out as a centered card: title +
  tagline across the top, then the logo (centered, `24%`/`max-width:230px`)
  on its own `col-lg-12` row, then two columns — left `col-lg-7` has
  "Project information" + "Highlights" (text left-justified) and right
  `col-lg-5` has "About". The 7/5 split is forced side-by-side via
  `body.in-modal .col-lg-7/.col-lg-5` flex widths at `>=576px`, because the
  iframe viewport is < Bootstrap's 992px `lg` breakpoint (and < 576px the
  modal intentionally stacks to a single column).
- The modal iframe size is capped with
  `.vbox-overlay .venoframe { width: min(940px, 96vw); height: min(74vh, 660px); }`
  — note the overlay has **class only** (`class="vbox-overlay"`, no id), so
  an id-based selector won't match.
- `assets/js/hero-quote.js` — fetches `/api/quotes/random` and repaints the
  hero tagline (quote + `— author`); 1–3 random words keep the underline.
  Failure leaves the static `h2` untouched.

The grid is rendered client-side, so data wiring lives in the two
`assets/js/*.js` files, not in the HTML. The Go admin and the frontend only
talk through `/api/projects/` — the Markdown contract in
`src/content/projects/` is the source of truth.

## Contact

- Contact section (`index.html#contact`) has two info boxes (Email Me =
  `admin.4rch3.io@pm.me`; Social Profiles = GitHub/LinkedIn/Medium,
  matching the hero links) — the template's fake address/phone were removed.
- The form (`class="php-email-form"`, `action="/api/contact/"`) is validated
  by `assets/vendor/php-email-form/validate.js`, which AJAX-POSTs its
  serialize() output and requires the response body to be exactly `OK`.
  The Go endpoint handles that contract (returns `OK`, or a clear error
  message for empty/invalid input; 405 for non-POST).
- **Abuse controls**: per-IP fixed-window rate limit (`CONTACT_RATE_LIMIT`,
  default 5/hour, `CONTACT_RATE_WINDOW`) → 429 + `Retry-After`; a hidden
  honeypot field (`name="website"`, positioned off-screen) — bots that
  fill it get a fake `OK` and are dropped silently (logged). nginx now
  sets `X-Real-IP`/`X-Forwarded-For` on `/api/` too, and the admin's
  `clientIP` prefers `X-Real-IP` (nginx's own view — not client-spoofable,
  unlike the first XFF entry). Without that, every visitor appeared as
  the web container's IP.

## Content

Projects live in `src/content/projects/*.md` (frontmatter schema: title,
pillar, tagline, status, order, featured, path, repo, live, stack,
highlights, body). Pillars: `cybersecurity`, `ai-ml`, `python`, `rust`,
`homelab` (`lab.md` = "The Lab").
Edit them via the admin UI at `:8095/admin/`; the JSON API reflects edits
immediately (no rebuild needed for content — only for frontend file changes).
Project logos are generated from frontmatter (`make logos` → runs
`scripts/gen_project_logos.py`); run it after adding/changing projects.
Each logo is a dark Matrix-style tile with a pillar-colored **concept mark
from the `LOGOS` dict** (svgrepo.com blocks automated downloads — 429 — so
marks are drawn as custom monochrome `currentColor` glyphs; swap/pick
different ones by editing a `LOGOS` entry).

Quotes live in the `ai-quotes` SQLite DB (`../ai-quotes/instance/ai_quotes.db`),
managed at `:8095/admin/quotes/`. The compose mount is a **directory** mount
(not a file mount) — SQLite must be able to create its `-journal` sibling.
Don't use a file bind for it.

## Screenshots & diagrams (M6)

Real visuals only — never stock images or fabricated dashboards.
Three pipelines, all writing into `frontend/assets/img/`:

- **Dashboards** — `make dashboards` (`scripts/capture_dashboards.py`,
  Playwright/Chromium headless): captures the live homelab services listed
  in `DASHBOARDS` (threatpulse, augur, eventflow) → `screenshots/<slug>.png`.
  Prints a verdict per page (login wall / data-rich); pages behind login
  need `DASH_CREDS='{"slug":["user","pass"]}'`. augur currently requires
  credentials (registration-based, none stored).
- **Terminal captures** — `make terminal` (`scripts/capture_terminal.py`):
  executes the commands in `scripts/terminal_captures.json` (read-only
  commands only — output must be real) with `FORCE_COLOR=1`, parses the
  ANSI, and re-renders it as a branded terminal-window PNG →
  `screenshots/<slug>.png`. Extend the JSON to add tools; ~80 lines max.
- **Architecture diagrams** — `make diagrams` (`scripts/render_diagrams.py`):
  renders `src/content/diagrams/<slug>.mmd` (mermaid, hand-written from each
  project's real narrative) via vendored `scripts/vendor/mermaid.min.js`
  (no network) → `diagrams/<slug>.png`, dark theme in brand colors.

The screenshot tooling lives in the project `.venv` (`python3 -m venv .venv
&& .venv/bin/pip install playwright pillow && .venv/bin/playwright install
chromium`); it's git-ignored. Detail pages pick these up automatically
(hero prefers screenshot; About appends diagram) — run `make sync` after
regenerating.

## Accessibility

`make a11y` (`scripts/a11y_audit.py`, vendored axe-core 4.10.2) audits
the live pages and **exits non-zero on any violation** — keep it at zero.
Current state: 0 violations across home, research, and both detail-page
modes. Hard-won gotchas baked into the code:

- `body` must keep a **solid** `background: #040404` — the template's
  `background: transparent` shorthand reset the background-color, and
  auditors/browsers then resolve text against **white** (bg.jpg comes
  from `body::before`, which axe cannot attribute). This single line
  caused most of the original contrast failures.
- CRT glow uses `filter: drop-shadow(...)`, **not** `text-shadow` —
  a dim green text-shadow is treated as worst-case foreground color by
  contrast audits. `text-shadow` also inherits to child spans.
- `brand-flicker` keyframes must not dip opacity below ~0.88 (4.5:1 at
  mid-flicker).
- Heading order: home = h1 brand → h2 section titles → h3 grid cards;
  research = h1 section → h2 card titles → h3 note-modal titles; detail
  pages = h1 project title → h2 (Project information / Highlights /
  About / Architecture). Changing a generated heading level means
  updating both the JS and the CSS selectors (`.portfolio-info h3` etc.).
- `#portfolio-details` (fixed, `overflow-y: auto`) carries
  `tabindex="0"` + `role="region"` for keyboard scroll access.
- `prefers-reduced-motion` disables all animations/transitions.
- Icon-only links (hero socials, venobox `.read-note` triggers) carry
  `aria-label`; contact inputs carry `aria-label`.

## Verification

- `docker compose build web` exits 0.
- `curl localhost:8095/` returns 200 and non-empty content; `/api/projects/`
  returns the project count.
- `/admin/quotes/` (mounted DB) lists the ai-quotes rows; `/api/quotes/random`
  returns a JSON quote; a CRUD round-trip leaves the real DB unchanged.
- Data-rendering logic is tested with jsdom (`/tmp/opencode/headtest`).

## Browser caching

nginx serves **every** asset + HTML page with
`Cache-Control: no-store, no-cache, must-revalidate, max-age=0` (never
rely on revalidation-style `no-cache` alone — Brave/other browsers have
shown stale pages despite it; `no-store` fixed those reports). Asset
`?v=` query-bumps are still used as a belt-and-suspenders on `style.css`.
If a user reports stale styles, tell them a normal reload suffices now.

## SEO / social metadata

All three pages carry proper description/keywords + canonical, Open
Graph + Twitter `summary_large_image` tags, and (index only) JSON-LD
`Person` + `WebSite` structured data (real `sameAs`: GitHub glopez21,
LinkedIn). The shared 1200x630 `og-image.png` is generated by
`scripts/gen_og_image.py` (Pillow, matches the #18d26e CRT brand) —
regenerate with `make og-image`-style: `python3 scripts/gen_og_image.py`.
Canonical/OG URLs use `https://4rch3.io` (the M7 deploy domain; the site
is LAN-only until then, so previews render the image but links won't
resolve publicly yet).

## Deferred (not bugs)

- Resume: the home `#resume` section and the empty `resume.html` stub were
  **removed** (nav "Resume" link dropped on all pages). A polished,
  restyled CV/skills section may be added back later.
- Dark-cyber restyle of the Bootstrap template is pending.