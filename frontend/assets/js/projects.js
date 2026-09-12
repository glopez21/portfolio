!(function($) {
  "use strict";

  var PILLAR_LABELS = {
    "cybersecurity": "Cybersecurity",
    "ai-ml": "AI/ML",
    "python": "Python",
    "rust": "Rust",
    "homelab": "Homelab"
  };

  function pillarLabel(pillar) {
    return PILLAR_LABELS[pillar] || pillar.charAt(0).toUpperCase() + pillar.slice(1);
  }

  var GH_SVG = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/></svg>';

  function renderGrid(projects) {
    var $container = $('.portfolio-container');
    var $flters = $('#portfolio-flters');

    var pillars = [];
    $.each(projects, function(i, p) {
      if ($.inArray(p.pillar, pillars) === -1) {
        pillars.push(p.pillar);
      }
    });

    $.each(pillars, function(i, pillar) {
      $flters.append($('<li>', {
        'data-filter': '.filter-' + pillar,
        text: pillarLabel(pillar)
      }));
    });

    $.each(projects, function(i, p) {
      var isLab = p.pillar === 'homelab';
      var logo = 'assets/img/project-logos/' + p.slug + '.svg';
      var $wrap = $('<div>', { 'class': 'portfolio-wrap' + (isLab ? ' card-panel' : '') });
      if (isLab) {
        $wrap
          .append($('<i>', { 'class': 'bx bx-server card-panel-icon' }))
          .append($('<h3>', { 'class': 'card-panel-title', text: p.title }))
          .append($('<div>', { 'class': 'card-panel-text', text: p.tagline }));
      } else {
        $wrap.append($('<img>', { 'src': logo, 'class': 'img-fluid project-logo', 'alt': p.title }));
      }
      var itemCls = isLab
        ? 'portfolio-item filter-homelab lab-item'
        : 'col-lg-4 col-md-6 portfolio-item filter-' + p.pillar;
      var $item = $('<div>', { 'class': itemCls })
        .append($wrap);
      $wrap
        .append($('<div>', { 'class': 'portfolio-info' })
          .append($('<h3>', { text: p.title }))
          .append($('<p>', { text: pillarLabel(p.pillar) }))
          .append($('<div>', { 'class': 'portfolio-links' })
            .append($('<a>', { 'href': '#', 'class': 'portfolio-repo', 'title': 'Repository' })
              .html(GH_SVG))
            .append($('<a>', { 'href': 'portfolio-details.html?slug=' + p.slug, 'data-gall': 'portfolioDetailsGallery', 'data-vbtype': 'iframe', 'class': 'venobox', 'title': 'Portfolio Details' })
              .append($('<i>', { 'class': 'bx bx-link' })))));
      $container.append($item);
    });

    var firstFilter = $flters.children('li').first().data('filter') || '*';

    var portfolioIsotope = $container.isotope({
      itemSelector: '.portfolio-item',
      layoutMode: 'fitRows',
      filter: firstFilter
    });
    $flters.children('li').first().addClass('filter-active');

    $('#portfolio-flters li').on('click', function() {
      $("#portfolio-flters li").removeClass('filter-active');
      $(this).addClass('filter-active');
      portfolioIsotope.isotope({
        filter: $(this).data('filter')
      });
    });

    $('.venobox').venobox({ 'share': false });
  }

  function showError() {
    $('#portfolio-error').removeClass('d-none');
  }

  $(document).ready(function() {
    $.getJSON('/api/projects/')
      .done(function(data) {
        if (data && data.projects && data.projects.length) {
          renderGrid(data.projects);
        } else {
          showError();
        }
      })
      .fail(function() {
        showError();
      });
  });

})(jQuery);