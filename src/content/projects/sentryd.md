---
title: sentryd
pillar: cybersecurity
tagline: 7-crate Rust SIEM — ingest, store, detect, agent, with per-PID connection
  attribution.
status: dev
order: 4
featured: false
path: sentryd
stack:
- Rust
- SQLite
- REST API
- Agent
- Docker
highlights:
- '7-crate workspace with clean separation: ingest → store → detect → agent'
- Per-PID connection attribution via socket-inode matching
- SSH brute-force, suspicious IP, and privilege-escalation detectors
- Agent pushes enriched events to Augur for full SOC pipeline integration
---

sentryd is a Rust-native SIEM agent and detection engine, organized as
a 7-crate workspace with clear separation of concerns: ingestion, storage,
detection, and a host agent that ties them together.

The agent reads from `/proc/net/tcp` and `/proc/net/6/tcp`, matching
socket inodes to running processes to attribute network connections to
specific PIDs — a capability that most host-based SIEMs lack entirely.
Detectors evaluate incoming events against known patterns: SSH brute-force
sequences, suspicious outbound IPs, and privilege-escalation indicators.

sentryd is designed to sit between the raw host telemetry and Augur's
detection pipeline, providing enriched, PID-attributed events that
Augur can correlate across multiple hosts.