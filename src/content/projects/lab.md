---
title: The Lab
pillar: homelab
tagline: Self-hosted home cluster — nginx reverse proxy, Docker Compose services, and this very portfolio running from one box on the LAN.
status: live
order: 27
featured: false
path: TheLab
stack:
- Docker Compose
- nginx
- Go
- SQLite
- Caddy
- Tailscale
highlights:
- Runs this portfolio + Go admin CMS behind a single nginx reverse proxy
- Quotes database (SQLite) managed live through the admin UI
- Canary services (LLM proxy, vault, monitoring) on separate compose stacks
- Encrypted remote access via Tailscale; no router port exposure
---

# The Lab

A single, quiet machine on the 192.168.1.0/24 network runs the whole
house: this portfolio, the Go admin, the ai-quotes database, and a few
side services. One reverse proxy, one compose file, one box.

Services running today:

- Portfolio site + admin CMS (nginx → Go)
- Quotes API / SQLite store
- Local LLM proxy (ollama) for the analyst agent
- Vault, log shipping, and uptime probes

The next phase moves the public services behind Traefik on the LAN and
mirrors the SQLite quotes store to the clu5t3r node for redundancy.