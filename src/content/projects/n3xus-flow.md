---
title: n3xus-flow
pillar: rust
tagline: Local-first CI/CD runner with TUI — DAG pipelines, cron triggers, and artifact
  collection.
status: dev
order: 3
featured: false
path: n3xus-flow
stack:
- Rust
- ratatui
- tokio
- YAML
- TUI
highlights:
- DAG pipeline execution with parallel stage scheduling
- Cron-triggered runs with persistent state
- Artifact collection and step-level output streaming
- ratatui TUI with live pipeline status, log viewer, and step tree
- Headless mode for agent-driven or CI-driven execution
---

n3xus-flow is a local-first CI/CD pipeline runner with a ratatui TUI —
designed to replace manual deployment scripts with a structured, auditable
pipeline that runs entirely on your local machine or homelab server.

Pipelines are defined in YAML with DAG-style stage dependencies: stages
that don't depend on each other run in parallel, with the runner
scheduling work across available threads. Cron triggers let you schedule
recurring pipelines (nightly builds, health checks, backup jobs) with
full run history and artifact collection.

The TUI provides a live view of running pipelines: stage status tree,
per-step log streaming, and a pipeline history browser. The headless
mode runs the same pipelines without a terminal, suitable for agent-driven
or CI-triggered execution where you just need the exit code and logs.