# 4rch3.io — Project Portfolio

> The gemstone of the homelab: a static portfolio site showcasing every
> engineering project across four pillars — Cybersecurity, AI/ML, Python,
> and Rust. Built with **Astro**, styled in a **dark-cyber** aesthetic,
> deployed to **portfolio.4rch3.io** on the clu5t3r homelab.

## Quickstart

```bash
npm install        # install dependencies
npm run dev        # dev server at http://localhost:4321
npm run build      # static build to dist/
npm run preview    # preview the production build locally
```

## Containerized deployment (Docker Compose)

The whole stack ships as two containers — the static site behind nginx
plus a Go admin CMS — and is reachable from any machine on the LAN:

```bash
cp .env.example .env        # review secrets first
make up                     # build + start both containers
```

| Service | URL | Notes |
|---|---|---|
| Site | `http://<host>:8095/` | nginx serves the static build |
| CMS | `http://<host>:8095/admin/` | Go admin, Django-admin-style UI, behind the same port |
| CMS (direct) | `http://<host>:8096/` | optional direct port `:8096` |

Both ports bind `0.0.0.0`, so any LAN device can reach the site and the
admin UI. Log in with `ADMIN_USERNAME` / `ADMIN_PASSWORD` from `.env`.

### Managing content

The **admin UI is the CMS** — edit projects there (title, pillar, status,
features, featured/order flags, body copy). It writes the Markdown source
files in `src/content/projects/` **in place**, so the content of truth is
always the files the static site is built from (no database, no export
step).

```bash
make sync   # rebuild + restart the site container so edits go live
```

Equivalent one-off commands:

```bash
make logs   # tail both containers
make down   # stop the stack (content files are on the host, nothing lost)
```

Ports and credentials live in `.env` (see `.env.example`).

## Project structure

```
portfolio/
├── src/
│   ├── content/
│   │   └── projects/          # ← project entries (Markdown + frontmatter)
│   ├── content.config.ts      # collection schema + pillar taxonomy
│   ├── components/            # Header, Footer, ProjectCard
│   ├── layouts/Base.astro     # global shell
│   ├── lib/meta.ts            # pillar colors, status labels
│   ├── pages/
│   │   ├── index.astro        # hero + featured + per-pillar sections
│   │   ├── about.astro        # lab infrastructure + philosophy
│   │   └── projects/
│   │       ├── index.astro    # filterable grid (?pillar= python)
│   │       └── [slug].astro   # deep project page
│   └── styles/global.css      # dark-cyber theme tokens
├── astro.config.mjs           # site URL + sitemap integration
└── public/                    # favicon
```

## Adding a project

The recommended path is the admin CMS (`make up` → `:8095/admin/` → add
project → `make sync`). To add one directly as Markdown, create a new file
in `src/content/projects/<slug>.md`:

```md
---
title: My Project
pillar: rust            # cybersecurity | ai-ml | python | rust
tagline: 'One-line description.'
status: dev             # live | dev | planned | research
order: 3               # lower = higher in pillar section
featured: true         # appears in the hero spotlight
path: my-project
repo: https://github.com/user/repo
live: https://app.4rch3.io
stack: ['Rust', 'axum', 'tokio']
highlights:
  - 'Feature one'
  - 'Feature with colon: still a string'
---
Narrative — problem → approach → result. Supports full Markdown.
```

The card, pillar section, `/projects` grid, and detail page all update
automatically. Highlights are YAML strings — the CMS export path handles
quoting automatically, so write them naturally in the admin UI.

## Admin/CMS — architecture

```
┌────────────┐   :8095      ┌──────────────────────────────┐   :80
│ any LAN    │ ───────────► │ web (nginx)                  │
│ device     │              │  /            → static site  │
└────────────┘              │  /admin/  ⇢ proxy ──────────┼──► admin (Go)
   ┌────────────┐  :8096     └──────────────────────────────┘      │ :8000
   │ direct     │ ──────►    (same admin container)                ▼
   └────────────┘                        src/content/projects/*.md
                                          (writes in place — the
                                           static build runs from these)
```

- `web` — multi-stage build (Node build → nginx runtime), reverse-proxies
  `/admin/` and `/admin-static/` to the admin container.
- `admin` — a single static Go binary (`admin-go/`): session-based login,
  project list/add/edit/delete UI matching the schema of the content
  collection, JSON API at `:8096/api/projects/`. No database — it reads
  and writes the Markdown files directly.
- The site keeps serving even if the admin container is down (nginx only
  proxies on demand).

## Pillars & status

- Pillars: `cybersecurity`, `ai-ml`, `python`, `rust` (defined in `content.config.ts`)
- Statuses: `live` (green), `dev` (cyan), `planned` (amber), `research` (muted)

## Deployment (pending)

When the site reaches 80% completeness, it deploys to
`portfolio.4rch3.io` via the standard clu5t3r pipeline (Traefik + systemd +
`./deploy.sh`, matching the convention used by the rest of the fleet). The
`site` URL in `astro.config.mjs` already points there for sitemap/SEO.