#!/usr/bin/env python3
"""Render src/content/diagrams/<slug>.mmd to PNG via headless Chromium.

  .venv/bin/python scripts/render_diagrams.py [slug ...]

Uses the vendored mermaid.min.js (scripts/vendor/) — no network needed.
Output: frontend/assets/img/diagrams/<slug>.png (2x scale, dark theme,
brand colors). Exits non-zero if any diagram fails to parse.
"""
import os
import sys

from playwright.sync_api import sync_playwright

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC_DIR = os.path.join(ROOT, "src", "content", "diagrams")
OUT_DIR = os.path.join(ROOT, "frontend", "assets", "img", "diagrams")
MERMAID = os.path.join(ROOT, "scripts", "vendor", "mermaid.min.js")

HTML = """<!DOCTYPE html><html><head><style>
  html, body {{ margin: 0; padding: 24px; background: #040404; }}
</style></head><body>
<div class="mermaid">{src}</div>
<script>{mermaid}</script>
<script>
  mermaid.initialize({{
    startOnLoad: false, securityLevel: 'loose', theme: 'dark',
    themeVariables: {{
      background: '#040404', primaryColor: '#16351f',
      primaryBorderColor: '#18d26e', primaryTextColor: '#d7e0ea',
      lineColor: '#2ea36b', secondaryColor: '#0d1f14',
      tertiaryColor: '#0d1117', fontSize: '15px',
      fontFamily: 'MesloLGS Nerd Font Mono, Liberation Mono, monospace',
    }},
    flowchart: {{ curve: 'basis', padding: 12 }},
  }});
</script></body></html>"""


def main():
    only = sys.argv[1:]
    mmds = sorted(f for f in os.listdir(SRC_DIR) if f.endswith(".mmd"))
    if not mmds:
        print("no .mmd files in", SRC_DIR)
        sys.exit(1)
    os.makedirs(OUT_DIR, exist_ok=True)
    with open(MERMAID) as fh:
        mermaid_js = fh.read()
    failed = []
    with sync_playwright() as pw:
        browser = pw.chromium.launch()
        for f in mmds:
            slug = f[:-4]
            if only and slug not in only:
                continue
            with open(os.path.join(SRC_DIR, f)) as fh:
                src = fh.read()
            page = browser.new_page(
                viewport={"width": 1400, "height": 900}, device_scale_factor=2)
            errors = []
            page.on("pageerror", lambda e: errors.append(str(e)))
            page.set_content(HTML.format(src=src.replace("</", "<\\/"),
                                         mermaid=mermaid_js))
            try:
                page.evaluate("mermaid.run()")
                page.wait_for_selector(".mermaid svg", timeout=15000)
                page.wait_for_timeout(600)
                svg = page.locator(".mermaid svg").first
                if "error" in (svg.get_attribute("aria-roledescription") or "").lower() \
                        or page.locator(".mermaid .error-text").count():
                    raise RuntimeError("mermaid syntax error")
                out = os.path.join(OUT_DIR, f"{slug}.png")
                svg.screenshot(path=out)
                print(f"{slug}: -> {out} ({os.path.getsize(out) // 1024} KB)")
            except Exception as e:
                failed.append(slug)
                print(f"{slug}: FAILED ({e}; js errors: {errors[:2]})")
            page.close()
        browser.close()
    if failed:
        sys.exit(1)


if __name__ == "__main__":
    main()
