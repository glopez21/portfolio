---
title: ShadowSim
pillar: cybersecurity
tagline: SOC traffic simulator — 43 attack modules covering kernel, firmware, cloud,
  supply-chain, and LLM attacks.
status: live
order: 3
featured: false
path: shadowsim
repo: https://git.4rch3.io/61-72-6b-68-e9/shadowsim
stack:
- Python
- Modular
- MITRE ATT&CK
- Docker
highlights:
- 43 attack simulation modules across 5 tiers of sophistication
- 'Kernel-level evasion: DKOM, direct syscalls, SSDT hooking, callback manipulation'
- Hypervisor escape and hardware/firmware attack simulation
- SCADA/ICS, IoT/Embedded, and blockchain attack scenarios
- LLM memory poisoning and AI agent exploitation modules
---

ShadowSim is the SOC ecosystem's threat generator — a modular
simulation engine that produces realistic security events for testing
detection pipelines, training analysts, and validating response playbooks.

The module library spans five tiers: basic log spoofing through to
hypervisor escape chains and LLM memory poisoning. Each module produces
structured output compatible with Augur's ingestion API, meaning you can
pipe ShadowSim directly into the detection → correlation → triage → response
pipeline and watch the full lifecycle execute.

The most advanced modules simulate cutting-edge attack surface:
container escapes (CVE-2025-31133), GPU DMA Rowhammer, supply chain worms
across npm/PyPI/Docker, and identity chain attacks with dead-man switches.
These are not toy demonstrations — they produce structurally accurate
artifacts that a well-tuned detection engine should flag.