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
- `admin` (contact): the site's contact form posts urlencoded data to
  `POST /api/contact/`, which appends one JSON line per message to
  `CONTACT_FILE` (`/data/contact_messages.jsonl`, named `admin_data`
  volume) and returns the literal body `OK` (the frontend's
  `php-email-form` validate.js expects exactly `OK`). Messages are readable
  at `/admin/contact/` (auth, newest first). The `admin_data` volume is
  owned by uid 1000 — re-`chown` it if the volume ever gets recreated as
  root and store starts failing.
- `admin` (contact email notify): when `SMTP_HOST` + `CONTACT_EMAIL` are
  set (see `.env`), each contact message also emails the owner via
  net/smtp (STARTTLS, plain auth) — fire-and-forget goroutine, so the
  visitor still gets `OK` even if SMTP fails (errors only logged; the
  JSONL inbox is the durable record). `Reply-To` is set to the visitor's
  address. Visitor-controlled fields are CR/LF-sanitized against header
  injection. Verified with an in-network SMTP sink container (Gmail
  app-password goes in `.env` `SMTP_PASSWORD`; empty = notify attempts
  auth and fails with a startup-log warning).

## Development

No dev server. Edit files directly and rebuild the `web` image:

```
make sync        # docker compose build web && up -d web
docker compose up -d   # full stack
```

Verify locally on `:8095` (site + admin/API proxy) and `:8096` (direct
admin). Admin defaults: `admin` / `admin123`.

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
  .portfolio-wrap`). The `portfolio-links` hover icons are: a GitHub/git mark (`.portfolio-repo`,
  `href="#"` placeholder for now — will point at the project `repo`
  when it goes live) and the unchanged details-link venobox
  (`portfolio-details.html?slug=…`).
- `assets/js/project-details.js` — fetches a single project by slug,
  populates the detail page. Unknown slug shows a "not found" fallback.
  The hero image (`#portfolio-hero`, shown at 30% width of its column) is
  set **inline in the HTML from `?slug=`** before any JS/API runs (no race,
  no fallback trap), and the detail JS re-sets it as well. Default `src` is
  threatpulse.svg so the page never paints a broken image.
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
  `genghis.lopez@gmail.com`; Social Profiles = GitHub/LinkedIn/Medium,
  matching the hero links) — the template's fake address/phone were removed.
- The form (`class="php-email-form"`, `action="/api/contact/"`) is validated
  by `assets/vendor/php-email-form/validate.js`, which AJAX-POSTs its
  serialize() output and requires the response body to be exactly `OK`.
  The Go endpoint handles that contract (returns `OK`, or a clear error
  message for empty/invalid input; 405 for non-POST).

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

- **Repo links not live yet.** The portfolio hover GitHub/git icon is a
  placeholder (`href="#"`); when the user approves going live, set its href
  from each project's `repo` frontmatter (`<a class="portfolio-repo">`).
- Resume: the home `#resume` section and the empty `resume.html` stub were
  **removed** (nav "Resume" link dropped on all pages). A polished,
  restyled CV/skills section may be added back later.
- Dark-cyber restyle of the Bootstrap template is pending.