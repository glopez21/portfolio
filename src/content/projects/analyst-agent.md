---
title: Analyst Agent
pillar: ai-ml
tagline: SOC analyst co-pilot — NL queries over security data, RAG over threat intel,
  and automated playbook triggering.
status: planned
order: 1
featured: true
path: Analyst-Agent
stack:
- Python
- ollama
- MCP
- RAG
- Tool Calling
- Streamlit
highlights:
- Natural language queries over Augur/EventFlow data
- RAG index over threat intel feeds + ShadowSim scenario data
- 'Tool-calling: invoke ThreatPulse playbooks, Augur triage, TheHive cases'
- MCP server exposing all analyst tools to Claude/Cursor
- Streamlit chat UI with live event correlation
---

The Analyst Agent is the AI flagship of the SOC ecosystem — a natural
language co-pilot that lets a security analyst query their infrastructure
using plain English and trigger complex response workflows without writing
detection queries manually.

The agent architecture is built from first principles: a ReAct loop with
structured tool calling over the Augur and ThreatPulse APIs, a RAG index
over the threat intelligence corpus and ShadowSim scenario library, and
an MCP server exposing all capabilities to Claude, Cursor, or any
MCP-compatible client.

The Streamlit chat UI sits at the top: you type "show me all SSH
brute-force attempts from the last hour" and the agent translates that
into Augur API calls, formats the results, and presents them with
contextual threat intel from the RAG index. You can also trigger
playbooks directly: "block 192.168.1.42 on all gateways" invokes
ThreatPulse through the agent's tool registry.