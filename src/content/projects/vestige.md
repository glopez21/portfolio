---
title: Vestige
pillar: cybersecurity
tagline: DFIR memory forensics toolkit — memory diffing, YARA generation, and LLM-assisted
  forensic reports.
status: dev
order: 6
featured: false
path: Vestige
stack:
- Python
- YARA
- LLM
- Forensics
- CLI
highlights:
- Automated memory diffing engine across two capture snapshots
- YARA rule generation from memory patterns
- Forensic timeline reconstruction with ETW patch detection
- TPM PCR forgery and Secure Boot chain simulation
- LLM-assisted forensic report generation
---

Vestige is a digital forensics and incident response toolkit focused on
memory analysis — the practice of extracting forensic artifacts from
running system memory dumps.

The memory diff engine takes two snapshots and identifies structural
differences: new processes, changed memory regions, injected code, and
modified kernel structures. The YARA generator produces detection rules
automatically from diff results, giving you targeted rules without manual
signature writing.

The most distinctive capability is LLM-assisted forensics: given a
memory dump and a diff report, the system generates a structured forensic
narrative — timeline, indicators, and recommended next steps — that a
human analyst can act on directly.