---
title: rust-infer
pillar: ai-ml
tagline: Production model inference server in Rust — model registry, /predict, and
  Prometheus metrics.
status: research
order: 3
featured: false
path: rust-infer
repo: https://git.4rch3.io/61-72-6b-68-e9/rust-infer
stack:
- Rust
- axum
- ONNX Runtime
- Prometheus
- REST
highlights:
- Model registry with versioning and metadata management
- /predict endpoint with batch support and latency tracking
- 'Per-model Prometheus metrics: latency p50/p95/p99, throughput'
- Designed to serve KRAKEN classifier and neural-sim proxy models
---

rust-infer is a production-grade model inference server in Rust, designed
to serve ML models behind a well-instrumented REST API with full
Prometheus metrics integration.

The model registry manages uploaded models with version metadata, format
validation, and hot-reload capability. The /predict endpoint accepts
both single and batch inference requests, with per-request latency
tracking exposed to Prometheus.

The design is deliberately minimal — no training, no fine-tuning,
just reliable, measurable serving of pre-trained models. The initial
deployment target is serving the KRAKEN classifier (DGA family
attribution) and a neural-sim proxy model, with the Prometheus
metrics feeding into the existing Grafana dashboard.