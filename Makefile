.PHONY: up down build logs ps sync clean

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

clean:
	docker compose down -v --rmi local