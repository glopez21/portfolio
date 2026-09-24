# Augur

SOC hub — log ingestion, detection, correlation, triage, and worker orchestration in one coherent system.

**Intent**

Augur is the central nervous system of the SOC platform: a hub that ingests security telemetry from multiple agents (sentryd, LogSentry, ShadowSim), runs detection and correlation pipelines, and triages alerts by severity before forwarding high-fidelity events to ThreatPulse.

**Tech stack**

- Python
- FastAPI
- SQLite
- JWT
- Prometheus
- Workers

*Status: live*
