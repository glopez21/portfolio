---
title: SENTINEL
pillar: cybersecurity
tagline: Honeypot with AI-driven defense — intent classification cache and adaptive
  worm emulator.
status: dev
order: 5
featured: false
path: SENTINEL
repo: https://git.4rch3.io/61-72-6b-68-e9/SENTINEL
stack:
- Python
- FastAPI
- LLM
- Intent Cache
- CLI
highlights:
- HoneyGPT hybrid intent classifier with command caching
- AI-driven adaptive worm emulator with LLM reasoning
- Multi-protocol honeypot with realistic service emulation
- Structured attack logging with MITRE ATT&CK tagging
---

SENTINEL is a honeypot system that goes beyond static service emulation —
it uses an AI-powered intent classifier to respond to attacker actions
with realistic, adaptive behavior.

The HoneyGPT intent cache classifies incoming commands and requests into
attack categories (reconnaissance, lateral movement, credential dumping,
exfiltration) and maintains a stateful cache so repeated behaviors are
handled without redundant LLM calls. The adaptive worm emulator takes it
further: it reasons about the attacker's objectives and adapts its own
emulated responses to draw the engagement out, producing richer telemetry
for detection pipeline testing.