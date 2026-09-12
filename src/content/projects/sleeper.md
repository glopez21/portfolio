---
title: SLEEPER
pillar: cybersecurity
tagline: Evasion and persistence research — sleep evasion, HWBP bypass, and ghost
  frame spoofing.
status: dev
order: 8
featured: false
path: SLEEPER
stack:
- Python
- EDR Testing
- Sleep Evasion
- HWBP
- ETW
highlights:
- 'Sleep evasion simulation: FOLIAGE, Ekko, WaitForSingleObject'
- PHANTOMPULSE HWBP AMSI/WLDP/ETW bypass with PEB hash walking
- LACUNA Chain ghost frame BYOUD-Gap spoofing
- EDR detection testkit for validating defensive tooling
---

SLEEPER is an evasion research project focused on the cat-and-mouse
game between endpoint detection and attacker persistence — testing how
well EDR products detect advanced evasion techniques.

The sleep evasion module simulates three known evasion patterns:
FOLIAGE (thread stack manipulation), Ekko (timer-based wake), and
WaitForSingleObject abuse — each producing structurally accurate
behavior for EDR testing. The PHANTOMPULSE module implements HWBP
hardware breakpoint-based bypass of AMSI, WLDP, and ETW using PEB→Ldr
hash walking and private syscall stubs.

The LACUNA Chain module introduces ghost frame BYOUD-Gap spoofing —
exploiting the gap between `.pdata` exception records and `win32u` NOP
gaps to inject code that falls between detection windows. The EDR
testkit provides a structured validation suite for testing which
techniques your defensive tooling actually catches.