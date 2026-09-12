---
title: ThreatPulse
pillar: cybersecurity
tagline: Production-grade SOAR platform — HMAC-signed playbooks, 16 real enforcement
  actions, multi-tenant alerting.
status: live
order: 1
featured: true
path: ThreatPulse
repo: https://git.4rch3.io/61-72-6b-68-e9/threat-pulse
stack:
- Python
- FastAPI
- PostgreSQL
- Webhooks
- REST API
highlights:
- HMAC-signed webhook payloads for secure agent-to-platform delivery
- 16 production enforcement actions (iptables, Slack, AbuseIPDB, TheHive, etc.)
- Multi-tenant isolation with per-tenant trigger rules and action throttling
- Built-in playbook step executor with retry, timeout, and rollback
- Prometheus metrics and structured JSON logging throughout
---

ThreatPulse is the enforcement layer of the SOC ecosystem — a standalone
SOAR (Security Orchestration, Automation & Response) platform that can be
triggered by any external system via a signed REST webhook.

The platform is purpose-built to sit downstream of Augur: when the detection
engine identifies a high-severity event, it issues an HMAC-signed POST to
ThreatPulse's webhook receiver. The trigger evaluator matches the event
against a set of rule definitions, and the step executor fans out to
whatever enforcement actions are configured — blocking IPs at the firewall,
pushing alerts to Slack, creating case files, or filing automated reports.

Unlike cloud SOAR tools, ThreatPulse runs entirely on-prem and owns every
action end-to-end. There are no third-party integrations to audit — the
actions are first-party Python modules with full test coverage, structured
outputs, and deterministic retry semantics.