---
title: Augur
pillar: cybersecurity
tagline: SOC hub — log ingestion, detection, correlation, triage, and worker orchestration
  in one coherent system.
status: live
order: 2
featured: true
path: Augur
repo: https://git.4rch3.io/61-72-6b-68-e9/augur
stack:
- Python
- FastAPI
- SQLite
- JWT
- Prometheus
- Workers
highlights:
- 'Modular pipeline: ingestion → detection → correlation → triage'
- DB-backed async worker queue with replay and failure recovery
- JWT-authenticated agent enrollment with per-agent secret rotation
- DetectionRuleModel pluggable engine with hot-swap rule loading
- Prometheus metrics + structured JSON logging + Sentry error tracking
---

Augur is the central nervous system of the SOC platform: a hub that
ingests security telemetry from multiple agents (sentryd, LogSentry,
ShadowSim), runs detection and correlation pipelines, and triages
alerts by severity before forwarding high-fidelity events to ThreatPulse.

The architecture is intentionally simple: a single FastAPI process runs
the API surface and a DB-backed worker loop. Detection rules are
evaluated as Python callables loaded from a rules directory — no
external DSL, no query engine, just typed Python with full testability.
Workers pick up pending events from a SQLite queue, evaluate them against
the loaded detection rules, and write structured results back to the
database for the triage layer.

Augur is designed to run headless on the homelab, with its API surface
consumed by the web dashboard, by the CLI, and by any downstream
automation that can speak HTTP + JWT.