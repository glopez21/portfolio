---
title: ph4nt0m
pillar: cybersecurity
tagline: Anonymity and OPSEC toolkit — DoH rotation, covert channels, and blockchain-based
  C2.
status: dev
order: 7
featured: false
path: ph4nt0m
stack:
- Python
- DoH
- Ethereum
- Covert Channels
- OPSEC
highlights:
- 8-provider DoH rotation with automatic failover
- 'Covert channel suite: DPPM, TCP timestamps, QUIC timing, history morph'
- Blockchain C2 over Ethereum/Base/Optimism smart contracts
- Production hardening with metadata scrubbing and key rotation
---

ph4nt0m is an anonymity and operational security toolkit built for
network operations that require plausible deniability and covert
communication channels.

The DNS-over-HTTPS rotation engine cycles through 8 provider endpoints
with automatic failover and metadata scrubbing — DNS queries become
difficult to correlate across sessions. The covert channel library
implements four distinct embedding methods: differential pulse-position
modulation, TCP timestamp manipulation, QUIC timing morphing, and
connection history steganography.

The blockchain C2 channel is the most novel component: commands are
encoded in Ethereum/Base/Optimism smart contract state, making them
indistinguishable from ordinary blockchain transactions. The operator
writes a command to a contract's storage slot; the implant reads it
and executes — no traditional C2 infrastructure to detect.