---
title: r0ut3r
pillar: ai-ml
tagline: Private AI inference router — load balancing, fallback, and per-model metrics
  behind one endpoint.
status: dev
order: 3
featured: false
path: r0ut3r
repo: https://git.4rch3.io/61-72-6b-68-e9/r0ut3r
stack:
- Python
- LiteLLM
- ollama
- Prometheus
- OpenAI API
highlights:
- Single OpenAI-compatible endpoint routing multiple backends
- Load balancing across ollama and vLLM instances
- Fallback chains with automatic health probing
- Per-model Prometheus metrics for latency and throughput
- 'Model aliasing: map task-specific aliases to concrete models and backends'
---

r0ut3r is a private inference router — a LiteLLM-style control plane
that routes requests across multiple local model backends behind a
single OpenAI-compatible API endpoint.

The router handles load balancing across ollama and vLLM instances,
with automatic health probing to avoid routing to degraded backends.
Fallback chains ensure availability: if the primary model is overloaded,
requests cascade to the next backend in priority order. Model aliases
let you map task-specific names ("fast-chat", "deep-reasoning") to
concrete models and backends, so callers don't need to know the
underlying infrastructure.

Per-model Prometheus metrics track latency percentiles and throughput,
feeding into the existing Grafana dashboard. The router is designed
to run as a lightweight sidecar on clu5t3r, with zero cloud dependencies
and full privacy for prompt/response data.