---
title: s3arc
pillar: python
tagline: From-scratch search engine — BM25 ranking, full-text indexing, and relevance
  scoring over a Markdown corpus.
status: dev
order: 5
featured: false
path: s3arc
repo: https://git.4rch3.io/61-72-6b-68-e9/s3arc
stack:
- Python
- BM25
- Inverted Index
- TF-IDF
- CLI
highlights:
- Hand-built inverted index with TF-IDF and BM25 ranking
- Full-text indexing over Markdown corpus (ortex, books, docs)
- Query autocomplete with frequency-ranked suggestions
- CLI and HTTP API interfaces
---

s3arc is a from-scratch information retrieval engine — building the
full indexing and ranking pipeline without relying on Elasticsearch,
Meilisearch, or any external search infrastructure.

The inverted index is built from scratch in Python: documents are
tokenized, normalized, and stored in a term-frequency structure that
supports both TF-IDF and BM25 ranking. The query engine parses search
expressions, scores them against the index, and returns ranked results
with relevance scores.

The initial corpus is the ortex notes, books, and documentation from
the homelab — giving you a working search engine over your own
knowledge base. The query autocomplete builds a frequency model from
the index, providing relevant suggestions as you type.