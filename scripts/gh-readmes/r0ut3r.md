# r0ut3r

Private AI inference router — load balancing, fallback, and per-model metrics behind one endpoint.

**Intent**

r0ut3r is a private inference router — a LiteLLM-style control plane that routes requests across multiple local model backends behind a single OpenAI-compatible API endpoint.

**Tech stack**

- Python
- LiteLLM
- ollama
- Prometheus
- OpenAI API

*Status: dev*
