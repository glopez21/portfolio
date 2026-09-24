(function() {
  "use strict";

  // Rotating brand tagline under the 4rch3.io name.
  // Every variant renders the same message — "arche · the first principle" —
  // in a different script, and reloads cycle to a new one. The base HTML keeps
  // the Greek rendering as the quiet no-JS fallback.

  var VARIANTS = [
    '<span class="brand-greek">αρχή</span> · η πρώτη αρχή', // Greek
    '<span class="brand-greek">архе</span> · первый принцип', // Cyrillic
    '<span class="brand-greek">根源</span> · 第一原理', // Japanese (kanji)
    '<span class="brand-greek">はじまり</span> · だいいちげんり', // Japanese (kana)
    '<span class="brand-greek">起源</span> · 第一原理', // Chinese
    '<span class="brand-greek">𓆣</span> · 𓋹𓊽𓂀', // Ancient Egyptian (symbolic)
    '<span class="brand-greek">𒀭</span> · 𒈨', // Sumerian cuneiform
    '<span class="brand-greek">आदि</span> · प्रथम सिद्धान्त' // Sanskrit
  ];

  var KEY = "4rch3.brandVariant";
  var previous = -1;
  try {
    var raw = sessionStorage.getItem(KEY);
    previous = raw === null ? -1 : parseInt(raw, 10);
  } catch (e) {
    previous = -1;
  }

  var idx = Math.floor(Math.random() * VARIANTS.length);
  if (VARIANTS.length > 1) {
    while (idx === previous) {
      idx = Math.floor(Math.random() * VARIANTS.length);
    }
  }
  try {
    sessionStorage.setItem(KEY, String(idx));
  } catch (e) {}

  var tag = document.querySelector("#header .brand-tag");
  if (tag) {
    tag.innerHTML = VARIANTS[idx];
  }
})();