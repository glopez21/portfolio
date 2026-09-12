---
title: termvault
pillar: python
tagline: Terminal secrets vault — encrypted credential storage with audit reporting
  and CLI-first UX.
status: dev
order: 3
featured: false
path: termvault
repo: https://git.4rch3.io/61-72-6b-68-e9/termvault
stack:
- Python
- Click
- Cryptography
- SQLite
- CLI
highlights:
- AES-256-GCM encrypted credential storage with master password
- CLI-first UX with search, copy-to-clipboard, and tree view
- Audit report feature with timestamped access logs
- No network dependencies — fully offline credential management
---

termvault is a terminal-native secrets manager — a lightweight,
encrypted credential store for managing passwords, API keys, and
certificates from the command line.

Credentials are encrypted at rest using AES-256-GCM with a master
password that never leaves the process. The CLI provides a familiar
tree-view of your vault, with fuzzy search, one-copy-to-clipboard,
and structured output for scripting. The audit report feature logs
every access and modification with timestamps, making it suitable
for shared environments where credential access needs to be tracked.