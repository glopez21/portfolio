!(function($) {
  "use strict";

  function slugFromQuery() {
    var params = new URLSearchParams(window.location.search);
    return params.get('slug') || null;
  }

  function escapeHtml(s) {
    return $('<div>').text(s).html();
  }

  function renderParagraphs(body) {
    var blocks = (body || '').split(/\n\s*\n+/);
    var $frag = $('<div>');
    $.each(blocks, function(i, block) {
      block = block.trim();
      if (!block) return;
      var lines = block.split('\n');
      var isList = lines.every(function(l) {
        var t = l.trim();
        return t.length > 0 && (t.indexOf('- ') === 0 || t.indexOf('* ') === 0);
      });
      if (isList) {
        var $ul = $('<ul>');
        $.each(lines, function(j, l) {
          $ul.append($('<li>', { text: l.trim().replace(/^[-*]\s+/, '') }));
        });
        $frag.append($ul);
      } else {
        $frag.append($('<p>', { text: lines.join(' ').trim() }));
      }
    });
    return $frag;
  }

  function renderProject(p) {
    $('#portfolio-title').text(p.title);
    document.title = p.title + ' - Portfolio';
    $('#portfolio-hero').attr('src', 'assets/img/project-logos/' + p.slug + '.svg').attr('alt', p.title + ' logo');

    if (p.tagline) {
      $('#portfolio-tagline').text(p.tagline);
    }

    var $info = $('#portfolio-info');
    $info.append($('<h3>', { text: 'Project information' }));

    var $ul = $('<ul>');
    $ul.append($('<li>').append($('<strong>', { text: 'Category' })).append(': ' + escapeHtml(p.pillar)));
    $ul.append($('<li>').append($('<strong>', { text: 'Status' })).append(': ' + escapeHtml(p.status)));
    if (p.stack && p.stack.length) {
      $ul.append($('<li>').append($('<strong>', { text: 'Stack' })).append(': ' + escapeHtml(p.stack.join(', '))));
    }
    if (p.repo) {
      $ul.append($('<li>')
        .append($('<strong>', { 'text': 'Repository' }))
        .append(': ')
        .append($('<a>', { 'href': p.repo, 'target': '_blank', 'rel': 'noopener', 'text': 'source' })));
    }
    if (p.live) {
      $ul.append($('<li>')
        .append($('<strong>', { 'text': 'Live' }))
        .append(': ')
        .append($('<a>', { 'href': p.live, 'target': '_blank', 'rel': 'noopener', 'text': 'site' })));
    }
    $info.append($ul);

    var $highlights = $('#portfolio-highlights');
    $highlights.append($('<h3>', { text: 'Highlights' }));
    var $hl = $('<ul>');
    $.each(p.highlights || [], function(i, h) {
      $hl.append($('<li>', { text: h }));
    });
    $highlights.append($hl);

    var $about = $('#portfolio-about');
    $about.append($('<h3>', { text: 'About' })).append(renderParagraphs(p.body));
  }

  function showError() {
    $('#portfolio-title').text('Project not found');
    document.title = 'Project not found - Portfolio';
    $('#portfolio-error').removeClass('d-none');
  }

  $(document).ready(function() {
    $.getJSON('/api/projects/')
      .done(function(data) {
        var slug = slugFromQuery();
        var projects = (data && data.projects) || [];
        var p = null;
        if (slug) {
          p = projects.find(function(x) { return x.slug === slug; });
          if (!p) {
            showError();
            return;
          }
        } else if (projects.length) {
          p = projects[0];
        }
        if (!p) {
          showError();
        } else {
          renderProject(p);
        }
      })
      .fail(function() {
        showError();
      });
  });

})(jQuery);