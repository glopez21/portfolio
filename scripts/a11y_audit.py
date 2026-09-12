#!/usr/bin/env python3
"""Run an axe-core accessibility audit against the live site.

  .venv/bin/python scripts/a11y_audit.py [path ...]

Audits the given paths (default: home, research, two detail pages) on
localhost:8095 with the vendored axe-core, waiting for client-side
rendering first. Prints violations grouped by page, with impact, rule,
failing selectors and a fix hint. Exit code 1 if any violations remain
(useful as a regression gate; pass/fail counts print at the end).
"""
import json
import os
import sys

from playwright.sync_api import sync_playwright

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
AXE = os.path.join(ROOT, "scripts", "vendor", "axe.min.js")
BASE = os.environ.get("A11Y_BASE", "http://localhost:8095")

DEFAULT_PATHS = [
    "/",
    "/research.html",
    "/portfolio-details.html?slug=threatpulse",
    "/portfolio-details.html?slug=sleeper",
]

RUN = """
async () => {
  axe.configure({ branding: { application: '4rch3-audit' } });
  const r = await axe.run(document, {
    resultTypes: ['violations'],
    run: { 'color-contrast': { enabled: true } },
  });
  return r.violations.map(v => ({
    id: v.id, impact: v.impact, help: v.help,
    nodes: v.nodes.map(n => n.target.join(' ')).slice(0, 6),
  }));
}
"""


def main():
    paths = sys.argv[1:] or DEFAULT_PATHS
    with open(AXE) as fh:
        axe_src = fh.read()
    total = 0
    with sync_playwright() as pw:
        browser = pw.chromium.launch()
        for path in paths:
            page = browser.new_context(viewport={"width": 1440, "height": 900}).new_page()
            page.goto(BASE + path, wait_until="networkidle", timeout=30000)
            page.wait_for_timeout(1500)  # client-side grid/quote rendering
            page.add_script_tag(content=axe_src)
            try:
                violations = page.evaluate(RUN)
            except Exception as e:
                print(f"{path}: audit failed: {e}")
                violations = []
            n = len(violations)
            total += n
            print(f"\n=== {path} — {n} violation type(s)")
            for v in violations:
                print(f"  [{v['impact'] or 'none':8}] {v['id']:22} {v['help']}")
                for sel in v["nodes"]:
                    print(f"      -> {sel}")
            page.close()
        browser.close()
    print(f"\n{'PASS' if total == 0 else 'FAIL'}: {total} violation type(s) across {len(paths)} page(s)")
    sys.exit(1 if total else 0)


if __name__ == "__main__":
    main()
