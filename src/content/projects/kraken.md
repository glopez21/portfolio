---
title: KRAKEN
pillar: cybersecurity
tagline: DGA research toolkit — transformer-based generation, ML classification, and
  adversarial evasion.
status: dev
order: 9
featured: false
path: KRAKEN
stack:
- Python
- Transformer
- ML
- PyTorch
- CLI
highlights:
- Transformer-attention DGA with 3 distinct generation algorithms
- ML-based DGA family classifier for attribution
- Adversarial DGA variants that evade statistical detection
- Batch generation and real-time DGA monitoring modes
---

KRAKEN is a Domain Generation Algorithm research toolkit — both an
offensive tool for generating realistic DGA domains and a defensive
tool for understanding and detecting them.

The transformer-based DGA uses attention mechanisms to model the
statistical structure of real domain names, producing outputs that
are harder to detect than traditional DGA approaches. Three generation
algorithms provide different tradeoffs between domain frequency,
entropy, and stealth. The adversarial DGA module goes further — it
generates domains specifically designed to evade statistical detection
methods, making it a useful benchmark for testing defensive ML models.

The classifier performs DGA family attribution — given a suspicious
domain, it identifies the likely malware family behind it. This
makes KRAKEN useful as both a research tool (generating new attack
data) and a defensive tool (testing and training detection models).