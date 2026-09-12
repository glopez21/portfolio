---
title: neural-sim
pillar: ai-ml
tagline: Brain-inspired spiking neural network simulator with live browser visualization.
status: research
order: 2
featured: false
path: neural-sim
repo: https://git.4rch3.io/61-72-6b-68-e9/neural-sim
stack:
- Python
- NumPy
- Flask
- Three.js
- SNN
highlights:
- Leaky Integrate-and-Fire (LIF) neuron model with STDP learning
- Real-time 3D visualization in browser (Three.js)
- 'Adjustable parameters: leak rate, threshold, learning rate, weights'
- Ties into the SNN research lab as a practical demonstration
---

neural-sim is a brain-inspired spiking neural network simulator
designed to make the concepts behind SNNs tangible and interactive.

The backend implements a Leaky Integrate-and-Fire neuron model with
Spike-Timing-Dependent Plasticity (STDP) learning rules — the biological
mechanism by which synapses strengthen or weaken based on the relative
timing of pre- and post-synaptic spikes. The simulator runs the network
and exposes its state over a Flask API.

The frontend renders the network in real-time Three.js, showing spikes
propagating through the network as visual pulses. Adjustable sliders
control the leak rate, firing threshold, and learning rate — change
a parameter and watch the network's behavior change instantly.