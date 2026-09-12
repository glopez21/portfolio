# 4rch3.io — Project Portfolio

> The gemstone of the homelab: the personal portfolio site, showcasing
> engineering projects across four pillars — Cybersecurity, AI/ML, Python,
> and Rust (+ The Lab). Dark-cyber aesthetic, deployed on the homelab.

## Stack

Static frontend + Go admin CMS, containerized (Docker Compose):

- **`web`** (nginx): serves `frontend/` (BootstrapMade "Personal"-based
  template) as static files. Reverse-proxies `/api/`, `/admin/`, and
  `/admin-static/` to the `admin` container.
- **`admin`** (Go): Django-admin-style CMS. Exposes `/api/projects/`
  (JSON, read from the Markdown source) and the edit UI at `/admin/`.
  Writes `src/content/projects/*.md` **in place** — no database, no
  export step. Also manages the `ai-quotes` SQLite DB (hero tagline
  quotes) and stores contact-form messages.

| Service | URL | Notes |
|---|---|---|
| Site | `http://<host>:8095/` | one-pager + `research.html` + `portfolio-details.html?slug=` |
| CMS | `http://<host>:8095/admin/` | behind the site port (auth required) |
| CMS (direct) | `http://<host>:8096/` | optional direct admin port |

## Quickstart

```bash
cp .env.example .env     # set ADMIN_USERNAME / ADMIN_PASSWORD (required!)
make up                  # build + start both containers
```

Credentials and ports live in `.env` (read automatically by compose).
Both ports bind `0.0.0.0` (LAN-reachable).

## Managing content

Projects live in `src/content/projects/*.md` — the Markdown frontmatter
is the **source of truth** (schema: title, pillar, tagline, status,
order, featured, path, repo, live, stack, highlights, body; pillars:
`cybersecurity`, `ai-ml`, `python`, `rust`, `homelab`).

- Edit via the admin UI (`:8095/admin/`) — the JSON API reflects edits
  immediately, **no rebuild needed** for content changes.
- Frontend file changes require `make sync` (rebuilds the `web` image).
- Project logos are generated: `make logos` (after adding/changing
  projects) → `frontend/assets/img/project-logos/<slug>.svg`.

```bash
make up      # build + start the stack
make sync    # rebuild web after frontend edits
make logs    # tail both containers
make down    # stop (content is on the host — nothing lost)
```

## Contact form

The contact section posts to `POST /api/contact/` (Go admin), which
appends each message to `/data/contact_messages.jsonl` on the
`admin_data` volume. Read messages at **`:8095/admin/contact/`**
(auth, newest first).

## Repo layout

```
├── frontend/            # static site (nginx docroot) — the live pages
│   ├── index.html       # one-pager: hero, about, skills, portfolio, contact
│   ├── research.html    # research topics + reading log (client-side)
│   ├── portfolio-details.html  # deep pages via ?slug=
│   └── assets/js/       # projects.js, project-details.js, research.js,
│                        # hero-quote.js — client-side data rendering
├── admin-go/            # Go admin CMS (single binary, no DB)
├── src/content/projects/  # project entries (Markdown) — content source of truth
├── scripts/gen_project_logos.py
├── nginx.conf           # static serving + /api,/admin proxying + no-store caching
└── docker-compose.yml
```

Note: this repo was originally scaffolded as an Astro site; the live
frontend is the static `frontend/` template (the Astro `src/` content
collection survives only as the Markdown content contract + logo
generator input). See `AGENTS.md` / `dev_notes.md` for deep detail.

## Deployment (pending — M7)

When the site reaches 80% completeness it deploys to
`portfolio.4rch3.io` via the clu5t3r homelab pipeline (Traefik + DNS +
router forwarding). The agreed routing plan lives in `dev_notes.md`.
