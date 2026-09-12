---
title: logsentry
pillar: cybersecurity
tagline: Security log parsing toolkit — multi-format ingestion with MITRE ATT&CK heuristic
  detection.
status: dev
order: 10
featured: false
path: logsentry
stack:
- Python
- MITRE ATT&CK
- CLI
- YARA
- MISP
highlights:
- 'Multi-format log parsing: syslog, SSH, PAM, CloudTrail'
- MITRE ATT&CK detection heuristics with Navigator export
- Sigma rule engine for custom detection logic
- Real-time file watching + syslog listening mode
---

LogSentry is a security log parsing and enrichment toolkit designed for
SOC analysts who need to work across multiple log formats without a
full SIEM deployment.

The parser handles syslog, SSH auth, PAM, and CloudTrail logs out of
the box, normalizing them into a structured format for downstream
detection. The MITRE ATT&CK heuristic mapper tags each event with
relevant technique IDs and produces a Navigator-exportable layer
showing coverage gaps.

The Sigma rule engine lets you write detection logic in a portable
format — rules work identically whether they're running in real-time
file watching mode, batch analysis, or syslog listening mode. The
project bridges the gap between lightweight CLI log analysis and
full SIEM deployment.