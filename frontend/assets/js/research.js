!(function($) {
  "use strict";

  var TOPICS = [
    {
      id: 'quantum',
      icon: 'bx bx-atom',
      title: 'Quantum Computing',
      status: 'reading',
      meta: 'quantum · two-level systems · error correction',
      text: 'Superposition, entanglement, and interference — computing past the transistor. The physics is settled; the engineering is not. Decoherence, fault-tolerant error correction, and qubit scalability are shaping up as one of the hardest systems problems of the century.',
      terms: ['superposition', 'entanglement', 'two-level systems', 'decoherence', 'error correction'],
      build: 'From research → build: gate-design notes feed the NOMAD / RustLearningProjects sandbox while the theory tooling-up continues.',
      body: [
        'The scale-up problem is really an error problem: correcting decoherence faster than it can corrupt the computation. Most of my reading time goes into the surface-code stabilizer machinery and what it takes to run a single logical qubit.',
        'The tools are maturing, so this is a practical track too — I want to judge a quantum proposal on the physics, not the hype.'
      ],
      links: []
    },
    {
      id: 'liquid-neural-networks',
      icon: 'bx bx-droplet',
      title: 'Liquid Neural Networks',
      status: 'experimenting',
      meta: 'liquid time constants · continuous time · few params',
      text: 'Small neural circuits, inspired by the 302-neuron C. elegans, that treat time as a liquid state instead of a stack of layers. A few dozen sparsely-connected neurons can learn to control, navigate, and adapt in real time — cheap enough for edge devices.',
      terms: ['liquid time constants', 'continuous-time dynamics', 'C. elegans', 'sparse connectivity', 'few parameters'],
      build: 'From research → build: prototyped an LTC cell in the neural-sim sandbox; eyeing the ONNX export path for a lightweight edge-inference version.',
      body: [
        'What pulls me in is the parameter efficiency. A liquid cell encodes dynamics in its time constants rather than its width, so you get sequence memory with a handful of neurons and no attention stack.',
        'The honest questions are stability and interpretability — fixed-step vs adaptive ODE solvers, and what changes when the trained dynamics get exported to a smaller target.'
      ],
      links: []
    },
    {
      id: 'agentic-cybersecurity',
      icon: 'bx bx-shield',
      title: 'Agentic Cybersecurity',
      status: 'experimenting',
      meta: 'perceive · reason · act',
      text: 'Autonomous defenders that perceive, reason, and act: LLM-driven agents that triage alerts, hunt threats, orchestrate incident response, and red-team at machine speed — while the adversary builds the same. Making them trustworthy, observable, and governable is the hard half.',
      terms: ['LLM agents', 'alert triage', 'threat hunting', 'IR orchestration', 'guardrails & audit'],
      build: 'From research → build: the alert-triage loop is prototyped against the ThreatPulse SOC pipeline; agent tool-calling borrows from the mesh work in Bl4ck1c3.',
      body: [
        'The interesting failure mode is trust, not capability. An agent that acts at machine speed has to explain itself, halt on ambiguity, and leave a replayable audit trail.',
        'I collect small evals around two axes: task success on synthetic IR scenarios, and safety — how often an agent takes a destructive action without a human gate.'
      ],
      links: []
    },
    {
      id: 'riemann-zeta',
      icon: 'bx bx-math',
      title: 'The Mystery of the Riemann Zeta Function',
      status: 'reading',
      meta: 'ζ(s) · critical line · 2 centuries, 0 proofs',
      text: 'One equation that encodes the primes. Analytic continuation, the functional equation, and the unresolved silence of the nontrivial zeros — all said to lie on the critical line. The Riemann Hypothesis has survived two centuries and a million-dollar prize.',
      terms: ['ζ(s)', 'analytic continuation', 'functional equation', 'nontrivial zeros', 'prime distribution'],
      build: 'From research → build: paired with a small Rust explorer (RustLearningProjects) that plots ζ across the critical strip, so the zeros stop being abstract.',
      body: [
        'The hypothesis is one door: is there an analytic structure that forces every nontrivial zero onto the critical line? The primes hide inside it, which is why the problem refuses to leave number theory alone.',
        'My level is patient self-study — Euler products and the functional equation are tractable, so that is where the notes live for now.'
      ],
      links: []
    }
  ];

  var LOG = [
    { date: '2026-09-08', text: 'Re-reviewed the decoherence notes from the two-level-system papers; queued the surface-code read and how it overlaps with the NOMAD sandbox.' },
    { date: '2026-09-05', text: 'LNN probe: ONNX export of the LTC cell, adaptive vs fixed-step solver comparison on the same trajectory.' },
    { date: '2026-09-02', text: 'Agentic IR eval harness — measuring destructive-action rate with vs without a human-gate stop.' },
    { date: '2026-08-29', text: 'ζ on the critical strip: first Rust plot in the explorer project, trivial zeros sitting exactly on the real axis.' },
    { date: '2026-08-24', text: 'ThreatPulse triage loop wiring — agent drafts an IR summary, a human approves the action.' },
    { date: '2026-08-19', text: 'Back to basics with analytic continuation; the Riemann sheets are less mystical with a worked example.' }
  ];

  function chips(arr) {
    return $.map(arr, function(t) {
      return '<span class="chip">' + t + '</span>';
    }).join('');
  }

  function renderCards() {
    var $row = $('#research-cards');
    $.each(TOPICS, function(i, t) {
      $row.append(
        '<div class="col-lg-6 mb-4 research-col">' +
          '<div class="card research-card" data-topic="' + t.id + '">' +
            '<a class="venobox read-note" data-vbtype="inline" data-bgcolor="#0c0e10" href="#research-note-' + t.id + '"></a>' +
            '<div class="card-body">' +
              '<i class="' + t.icon + '"></i>' +
              '<h3 class="card-title mt-3">' + t.title + '</h3>' +
              '<p class="card-text">' + t.text + '</p>' +
              '<p class="card-meta">' + t.meta + '</p>' +
              '<div class="card-foot">' +
                '<span class="status-chip status-' + t.status + '">' + t.status + '</span>' +
                '<span class="open-note">open note</span>' +
              '</div>' +
            '</div>' +
          '</div>' +
        '</div>'
      );
    });
  }

  function renderNotes() {
    var $holder = $('#research-notes');
    $.each(TOPICS, function(i, t) {
      var links = $.map(t.links || [], function(l) {
        return '<li><a href="' + l.url + '" target="_blank" rel="noopener">' + l.label + '</a></li>';
      }).join('');
      var $note = $('<div>', { 'class': 'research-note', 'id': 'research-note-' + t.id });
      $note.append(
        '<span class="status-chip status-' + t.status + '">' + t.status + '</span>' +
        '<h3>' + t.title + '</h3>' +
        '<p class="note-meta">' + t.meta + '</p>' +
        '<div class="note-chips">' + chips(t.terms) + '</div>' +
        '<div class="note-body">' + $.map(t.body, function(p) { return '<p>' + p + '</p>'; }).join('') + '</div>' +
        (t.build ? '<div class="note-build">' + t.build + '</div>' : '') +
        (links ? '<ul class="note-links">' + links + '</ul>' : '')
      );
      $holder.append($note);
    });
  }

  function renderLog() {
    var $list = $('#research-log');
    $.each(LOG, function(i, e) {
      $list.append(
        '<li class="log-entry">' +
          '<span class="log-date">' + e.date + '</span>' +
          '<span class="log-text">' + e.text + '</span>' +
        '</li>'
      );
    });
  }

  $(document).ready(function() {
    renderCards();
    renderNotes();
    renderLog();

    $('.venobox').venobox({ 'share': false });

    $('.research-card').on('click', function(e) {
      var $link = $(this).find('.read-note');
      if ($link.length && !$(e.target).is($link)) {
        $link[0].click();
      }
    });
  });

})(jQuery);