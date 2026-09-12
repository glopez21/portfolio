.PHONY: up down build logs ps sync clean screenshots diagrams og-image push a11y

# ---------------------------------------------------------------------- #
# 4rch3.io portfolio — containerized site + Go admin CMS
# ---------------------------------------------------------------------- #

## Bring up the whole stack (builds images on first run)
up:
	docker compose up -d --build

down:
	docker compose down

build:
	docker compose build

## Rebuild the web image without touching the admin container
build-web:
	docker compose build web
	docker compose up -d web

logs:
	docker compose logs -f --tail=200

ps:
	docker compose ps

## Rebuild the static site so edits made in the admin go live.
## The Go admin edits the Markdown source in place, so no export step exists.
sync:
	docker compose build web
	docker compose up -d web

## Regenerate per-project SVG logos from src/content/projects/*.md
logos:
	python3 scripts/gen_project_logos.py

## Headless-capture the live homelab dashboards (.venv: see AGENTS.md)
dashboards:
	.venv/bin/python scripts/capture_dashboards.py

## Re-render real terminal output as screenshot PNGs (commands in
## scripts/terminal_captures.json — read-only commands only)
terminal:
	.venv/bin/python scripts/capture_terminal.py

## Render src/content/diagrams/*.mmd to frontend/assets/img/diagrams/
diagrams:
	.venv/bin/python scripts/render_diagrams.py

## All three screenshot pipelines
screenshots: dashboards terminal diagrams

## Regenerate the social/OG image
og-image:
	python3 scripts/gen_og_image.py

## Accessibility audit (axe-core) against the live site — exit 1 on violations
a11y:
	.venv/bin/python scripts/a11y_audit.py

## Push to Forgejo origin (GitHub mirror follows automatically)
push:
	git push origin main

clean:
	docker compose down -v --rmi local