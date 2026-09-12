#!/usr/bin/env python3
"""Headless-capture the live homelab dashboards as project screenshots.

  .venv/bin/python scripts/capture_dashboards.py [slug ...]

Pages are listed in DASHBOARDS below. Each capture is saved to
frontend/assets/img/screenshots/<slug>.png. The script also prints a
verdict per page — whether it hit a login wall (visible password input)
and how data-rich the rendered DOM is — so you know which captures need
credentials (set DASH_CREDS='{"slug":["user","pass"]}' to authenticate).
"""
import json
import os
import sys

from playwright.sync_api import sync_playwright

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OUT_DIR = os.path.join(ROOT, "frontend", "assets", "img", "screenshots")

DASHBOARDS = {
    "threatpulse": "https://threatpulse.4rch3.io/",
    "augur": "https://augur.4rch3.io/",
    "eventflow": "https://eventflow.4rch3.io/",
}

CREDS = {}
if os.environ.get("DASH_CREDS"):
    CREDS = json.loads(os.environ["DASH_CREDS"])

DIAG_JS = """
() => {
  const vis = el => el && el.offsetParent !== null;
  const pw = [...document.querySelectorAll('input[type=password]')].some(vis);
  const text = document.body.innerText || '';
  return {
    login: pw,
    charts: document.querySelectorAll('canvas, svg.chart, .chart, .recharts-wrapper, [class*=chart]').length,
    tables: document.querySelectorAll('table tr').length,
    textLen: text.trim().length,
    title: document.title,
  };
}
"""


def main():
    only = sys.argv[1:] or list(DASHBOARDS)
    os.makedirs(OUT_DIR, exist_ok=True)
    with sync_playwright() as pw:
        browser = pw.chromium.launch()
        for slug in only:
            if slug not in DASHBOARDS:
                print(f"unknown slug {slug}")
                continue
            url = DASHBOARDS[slug]
            ctx = browser.new_context(
                viewport={"width": 1440, "height": 900},
                device_scale_factor=2,
                ignore_https_errors=True,
            )
            page = ctx.new_page()
            try:
                page.goto(url, wait_until="networkidle", timeout=30000)
            except Exception as e:
                print(f"{slug}: goto warning: {e}")
            page.wait_for_timeout(2500)

            creds = CREDS.get(slug)
            if creds and slug in DASHBOARDS:
                user, passwd = creds
                try:
                    page.fill("input[type=text], input[name*=user], input[name*=email]",
                              user, timeout=3000)
                    page.fill("input[type=password]", passwd, timeout=3000)
                    page.keyboard.press("Enter")
                    page.wait_for_timeout(4000)
                except Exception:
                    pass

            d = page.evaluate(DIAG_JS)
            out = os.path.join(OUT_DIR, f"{slug}.png")
            page.screenshot(path=out)
            kb = os.path.getsize(out) // 1024
            verdict = "LOGIN WALL" if d["login"] else (
                "data-rich" if (d["charts"] + d["tables"]) > 6 else "sparse")
            print(f"{slug}: {verdict} | title={d['title']!r} charts={d['charts']} "
                  f"table-rows={d['tables']} text={d['textLen']}ch -> {out} ({kb} KB)")
            ctx.close()
        browser.close()


if __name__ == "__main__":
    main()
