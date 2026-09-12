---
title: m0rpheus-rs
pillar: rust
tagline: TUI Postgres workbench with embedded MCP server — vim-mode, schema tools,
  Claude/Cursor integration.
status: dev
order: 2
featured: false
path: m0rpheus-rs
stack:
- Rust
- ratatui
- tokio
- sqlx
- MCP
- PostgreSQL
highlights:
- Vim-mode TUI client for Postgres with schema browser
- Built-in MCP server (stdio) for LLM-driven database access
- Full schema introspection with table/column/type display
- Single static binary — zero runtime dependencies
---

m0rpheus-rs is a from-scratch Postgres workbench in Rust, designed as
a single static binary that brings a vim-mode TUI, schema tools, and
a built-in MCP server together into one workspace.

The TUI is built with ratatui and offers a split-pane view: a query
editor with vim keybindings, a results table, and a schema sidebar.
The schema sidebar does live introspection — tables, columns, types,
indexes — and lets you jump into any table with a keystroke.

The MCP subcommand (`m0rpheus-rs mcp`) exposes the same schema and
query capabilities over stdio, making it immediately usable as an MCP
server in Claude, Cursor, or any MCP-compatible client. This means an
LLM can query your lab databases with full schema context, without
any wrapper scripts or HTTP bridges.