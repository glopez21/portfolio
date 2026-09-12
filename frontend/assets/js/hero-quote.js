!(function($) {
  "use strict";

  // Hero tagline: fetch a random quote from the ai-quotes API and show it under
  // the 4rch3 name. Underline (1-3 random words) keeps the template aesthetic.
  // If the API is unreachable or empty, the static tagline in the HTML stays.

  function chooseUnderlineIndexes(words, maxUnderline) {
    var eligible = [];
    $.each(words, function(i, w) {
      if (w.length >= 3) eligible.push(i);
    });
    if (eligible.length === 0) {
      eligible = words.map(function(w, i) { return i; });
    }
    var count = Math.min(eligible.length, maxUnderline);
    var chosen = [];
    while (chosen.length < count && eligible.length > 0) {
      var pick = Math.floor(Math.random() * eligible.length);
      chosen.push(eligible.splice(pick, 1)[0]);
    }
    return chosen;
  }

  function renderQuote(quote) {
    var $h = $('#header h2');
    if (!$h.length) return;
    var text = $.trim((quote && quote.text) || "");
    if (!text) return;

    var words = text.split(/\s+/).filter(function(w) { return w !== ""; });
    if (words.length === 0) return;

    var maxUnderline = Math.min(words.length, 3);
    var underlines = chooseUnderlineIndexes(words, maxUnderline);

    $h.empty();
    $.each(words, function(i, w) {
      if (i > 0) $h.append(document.createTextNode(" "));
      if ($.inArray(i, underlines) !== -1) {
        $h.append($('<span>').text(w));
      } else {
        $h.append(document.createTextNode(w));
      }
    });

    var author = $.trim((quote && quote.author) || "");
    if (author) {
      $h.append(document.createTextNode(" "));
      $h.append($('<span>', { "class": "hero-author" }).text("\u2014 " + author));
    }
  }

  $(document).ready(function() {
    $.getJSON("/api/quotes/random")
      .done(renderQuote)
      .fail(function() {
        // keep the static tagline
      });
  });

})(jQuery);