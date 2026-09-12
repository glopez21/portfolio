---
title: 4sk
pillar: python
tagline: NL→SQL analyst — ask questions in English, get answers from a database with
  chart output.
status: dev
order: 1
featured: false
path: 4sk
repo: https://git.4rch3.io/61-72-6b-68-e9/4sk
stack:
- Python
- FastAPI
- SQLAlchemy
- Chart.js
- SQLite
highlights:
- Natural language queries translated to safe read-only SQL
- Matplotlib / Chart.js rendering of result sets
- Source database introspection for schema-aware prompting
- Full audit trail of every generated query
---

4sk (SQL Ask) is a natural-language-to-SQL analyst that lets you ask
questions in plain English and receive rendered charts backed by
real database queries.

The system is designed around safety: the translation layer emits only
read-only SQL, with schema introspection feeding into the prompt so the
LLM has full context of available tables and columns. Every generated
query is logged for auditability — you can always see exactly what was
run and why.

The FastAPI backend serves the API surface and renders static chart
images or JSON data for frontend consumption. The project is aimed at
the "Chat with your database" pattern, demonstrated over a homelab
dataset so it can be deployed as a live, interactive showcase.