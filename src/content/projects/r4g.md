---
title: r4g
pillar: ai-ml
tagline: Agentic RAG + MCP from first principles — no LangChain, no LlamaIndex, pure
  Rust/Python engineering.
status: dev
order: 1
featured: true
path: r4g
stack:
- Python
- Rust
- ollama
- BM25
- MCP
- Hybrid Retrieval
highlights:
- Hybrid BM25 + vector retrieval without LangChain or LlamaIndex
- ReAct agent loop with tool registry and structured output
- MCP server for Claude/Cursor integration
- Per-query evaluation scoring in a JSONL trace log
- From-scratch embedding store with no external vector DB
---

r4g is a from-scratch agentic RAG system built to demonstrate depth in
retrieval, tool use, and agent design — without relying on LangChain,
LlamaIndex, or any orchestration framework.

The retrieval stack is hybrid: BM25 via a hand-built inverted index
combined with a local vector store over sentence embeddings, with
reciprocal rank fusion to merge the two result streams. The agent uses
a ReAct loop with a structured tool registry — the LLM generates a
plan, executes tools (including retrieval and code execution), observes
results, and decides the next step.

Every query produces a JSONL trace entry with per-step reasoning, tool
calls, and timing — making the system evaluable and debuggable. An MCP
server exposes the same capabilities over stdio, so the RAG engine
can be plugged directly into Claude or Cursor as a knowledge source.