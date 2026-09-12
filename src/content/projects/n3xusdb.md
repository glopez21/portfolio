---
title: n3xusDB
pillar: rust
tagline: From-scratch embedded time-series engine — WAL, segment store, compaction,
  and soak benchmarks.
status: planned
order: 4
featured: false
path: n3xusDB
repo: https://git.4rch3.io/61-72-6b-68-e9/n3xusDB
stack:
- Rust
- WAL
- Segment Store
- Compaction
- Benchmarks
highlights:
- Append-only segment store with write-ahead log
- WAL with fsync batching for crash safety
- Segment compaction and retention policy support
- 1M-point soak benchmark with latency/capacity reporting
- Prometheus remote-write ingestion endpoint
---

n3xusDB is a from-scratch embedded time-series database built in Rust —
demonstrating the core systems concepts behind production databases
like Prometheus and InfluxDB, without the abstraction layer.

The storage layer is an append-only segment store backed by a write-ahead
log with fsync batching: writes go to the WAL first, then get flushed to
segments during compaction. The segment store uses a columnar layout for
time-series data, with timestamps and values stored separately for
efficient compression and range queries.

The benchmark suite writes 1 million data points and measures write
throughput, query latency, and storage efficiency — producing a report
that demonstrates both the engine's performance and the systems
reasoning behind its design decisions.