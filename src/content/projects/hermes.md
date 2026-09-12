---
title: hermes
pillar: rust
tagline: TUI terminal multiplexer — workspace management, pane splitting, and session
  persistence.
status: planned
order: 5
featured: false
path: hermes
stack:
- Rust
- ratatui
- tokio
- TUI
highlights:
- Terminal multiplexer with split panes and workspace management
- Session persistence across disconnects
- ratatui-based TUI with vim-style navigation
- Seed for m0rpheus-rs TUI patterns
---

hermes is a terminal multiplexer and workspace manager in Rust — a
personal tool for managing multiple terminal sessions with pane splitting,
workspace layouts, and session persistence across SSH disconnects.

Built with ratatui and tokio, hermes provides vim-style navigation
between panes, named workspace layouts that persist to disk, and
session state that survives network interruptions — reconnect and
pick up exactly where you left off.

The project serves as the design seed for m0rpheus-rs: the TUI
patterns, input handling, and layout engine developed here directly
informed the ratatui work in the Postgres workbench. hermes is
intentionally smaller in scope, focused on getting the core TUI
abstractions right before building a more complex application on top.