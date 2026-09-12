---
title: transitflow-nyc
pillar: python
tagline: Real-time NYC transit data pipeline — ETL, Postgres, ML, and live dashboards.
status: dev
order: 2
featured: true
path: transitflow-nyc
repo: https://git.4rch3.io/61-72-6b-68-e9/transitflow-nyc
stack:
- Python
- FastAPI
- PostgreSQL
- ML
- Chart.js
- Docker
highlights:
- Live MTA feed ingestion with automatic deduplication
- Postgres backend with time-series optimized schema
- 'ML layer: arrival delay prediction + corridor congestion clustering'
- Interactive dashboard with per-line and per-station views
- Scheduled scraper for historical data enrichment
---

transitflow-nyc is a full-stack data engineering project: ingesting
real-time MTA transit feeds, storing them in a time-series Postgres
schema, running ML prediction and clustering models, and rendering
the results in an interactive dashboard.

The ETL pipeline pulls from the MTA's GTFS-RT feeds, deduplicates
arrivals against the PostgreSQL store, and computes rolling features
(average delay, corridor congestion, line utilization) that feed into
two ML models: a delay predictor (scikit-learn gradient boosting) and
a corridor clustering model (K-Means over congestion vectors).

The FastAPI backend exposes a REST API consumed by a Chart.js dashboard
that shows live per-line status, per-station delay heatmaps, and
historical trend charts. The project demonstrates complete data
engineering literacy — from raw API ingestion through to live,
interactive visualization.