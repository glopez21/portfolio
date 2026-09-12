---
title: Bl4ck1c3
pillar: rust
tagline: Rust smart firewall — libpcap capture, DPI, anomaly ML, nftables enforcement,
  live dashboard.
status: live
order: 1
featured: true
path: Bl4ck1c3
repo: https://git.4rch3.io/61-72-6b-68-e9/Bl4ck1c3
stack:
- Rust
- axum
- tokio
- sqlx
- libpcap
- nftables
- pnet
- Chart.js
highlights:
- Full-flow DPI engine with per-protocol parsing
- 19-dimensional flow feature extraction for ML anomaly scoring
- 'Live dashboard: top blocked IPs, drop rates, node health, rules CRUD'
- 'Agent mesh: per-PID socket-inode connection attribution across hosts'
- XDP eBPF fast-path program (bl4ck1c3_kern.c) for kernel-level capture
---

Bl4ck1c3 is a Rust-native smart firewall deployed across the homelab
fleet — a controller on clu5t3r, edge gateways on multiple hosts, and
agents reporting host telemetry over a self-signed mesh.

The architecture is built around the libpcap capture loop: packets are
parsed via `pnet`, assembled into IP-pair flows, and handed through a
19-dimensional feature extraction step before being classified by a
flow-level anomaly detector. DPI results and ML scores are fed into
nftables rule enforcement, with drops and alerts submitted back to the
central controller.

The controller exposes a REST API and a live browser dashboard (Chart.js
powered) showing real-time telemetry, top-IPs, event streams, and
a full rules management interface. The agent runs per-host, attributing
connections to individual processes via socket-inode matching — a
capability that surfaces in the dashboard as per-PID network breakdowns.

Built entirely in Rust with async Tokio internals, a hardened XDP
capture path in C, and zero third-party cloud dependencies.