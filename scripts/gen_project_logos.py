#!/usr/bin/env python3
"""Generate a Matrix-style SVG logo per project from src/content/projects/*.md.

  python3 scripts/gen_project_logos.py [content_dir] [out_dir]

Each logo is a dark tile with a pillar-colored concept mark (see LOGOS
below), the project title, and a pillar tag. Swap any logo by editing its
entry in LOGOS — every mark is drawn in a 100x100 coordinate system using
currentColor, so it always inherits the pillar accent.
"""
import glob
import html
import os
import re
import sys
import textwrap

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CONTENT_DIR = sys.argv[1] if len(sys.argv) > 1 else os.path.join(ROOT, "src", "content", "projects")
OUT_DIR = sys.argv[2] if len(sys.argv) > 2 else os.path.join(ROOT, "frontend", "assets", "img", "project-logos")

PILLAR_COLOR = {
    "cybersecurity": "#18d26e",
    "ai-ml": "#22d3ee",
    "python": "#e8c15a",
    "rust": "#ff8a65",
    "homelab": "#2dd4bf",
}

PILLAR_TAG = {
    "cybersecurity": "CYBERSECURITY",
    "ai-ml": "AI/ML",
    "python": "PYTHON",
    "rust": "RUST",
    "homelab": "HOMELAB",
}

# Concept marks, drawn in a 100x100 space with currentColor.
LOGOS = {
    # --- cybersecurity ---
    "threatpulse": (
        '<circle cx="50" cy="50" r="40" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="50" cy="50" r="20" fill="none" stroke="currentColor" stroke-width="3" opacity="0.6"/>'
        '<line x1="50" y1="50" x2="84" y2="50" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
    ),
    "augur": (
        '<circle cx="50" cy="42" r="30" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M50 72 V82" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M28 82 H72" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M18 20 L28 30 M82 20 L72 30" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
    ),
    "sentinel": (
        '<path d="M16 50 Q50 20 84 50 Q50 80 16 50 Z" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="50" cy="50" r="9" fill="currentColor"/>'
    ),
    "sentryd": (
        '<path d="M50 14 L80 27 V47 C80 69 66 83 50 87 C34 83 20 69 20 47 V27 Z" fill="none" stroke="currentColor" stroke-width="3" stroke-linejoin="round"/>'
        '<path d="M37 49 L47 59 L63 41" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>'
    ),
    "sleeper": (
        '<path d="M64 12 A38 38 0 1 0 88 68 A30 30 0 1 1 64 12 Z" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="24" cy="24" r="3" fill="currentColor"/><circle cx="80" cy="20" r="2.5" fill="currentColor"/>'
    ),
    "shadowsim": (
        '<path d="M16 74 A42 42 0 0 1 84 74" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M26 74 A30 30 0 0 1 74 74" fill="none" stroke="currentColor" stroke-width="3" opacity="0.7"/>'
        '<path d="M36 74 A18 18 0 0 1 64 74" fill="none" stroke="currentColor" stroke-width="3" opacity="0.5"/>'
        '<line x1="50" y1="74" x2="86" y2="38" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
    ),
    "kraken": (
        '<path d="M14 26 Q32 18 50 26 T86 26 M86 26 Q92 40 82 48" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M16 84 Q26 56 15 32" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M50 84 Q56 52 45 26" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M84 84 Q74 56 85 32" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<circle cx="15" cy="32" r="3.5" fill="currentColor"/><circle cx="45" cy="26" r="3.5" fill="currentColor"/><circle cx="85" cy="32" r="3.5" fill="currentColor"/>'
    ),
    "logsentry": (
        '<path d="M22 26 H78 M22 40 H60 M22 54 H70 M22 68 H50" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<line x1="70" y1="54" x2="78" y2="54" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<circle cx="84" cy="84" r="4" fill="currentColor"/>'
    ),
    "ph4nt0m": (
        '<path d="M30 48 A20 20 0 0 1 70 48 V84 L60 74 L50 84 L40 74 L30 84 Z" fill="none" stroke="currentColor" stroke-width="3" stroke-linejoin="round"/>'
        '<circle cx="42" cy="52" r="3.5" fill="currentColor"/><circle cx="58" cy="52" r="3.5" fill="currentColor"/>'
    ),
    "vestige": (
        '<path d="M50 14 L78 38 L62 84 H38 L22 38 Z" fill="none" stroke="currentColor" stroke-width="3" stroke-linejoin="round"/>'
        '<path d="M38 38 L50 54 L43 68" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
    ),
    "lab": (
        '<rect x="26" y="32" width="48" height="26" rx="6" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="36" cy="45" r="3" fill="currentColor"/><circle cx="46" cy="45" r="3" fill="currentColor"/><circle cx="56" cy="45" r="3" fill="currentColor"/>'
        '<rect x="26" y="64" width="48" height="26" rx="6" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="36" cy="77" r="3" fill="currentColor"/><circle cx="46" cy="77" r="3" fill="currentColor"/><circle cx="56" cy="77" r="3" fill="currentColor"/>'
    ),
    # --- ai-ml / rust-infer etc ---
    "3v4l": (
        '<path d="M38 20 L22 38 L38 56" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>'
        '<path d="M62 20 L78 38 L62 56" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>'
        '<path d="M48 28 L64 38 L48 48 Z" fill="currentColor"/>'
    ),
    "analyst-agent": (
        '<rect x="22" y="52" width="38" height="26" rx="6" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="34" cy="63" r="3" fill="currentColor"/><circle cx="48" cy="63" r="3" fill="currentColor"/>'
        '<path d="M41 52 V42 M34 46 H48" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<circle cx="72" cy="28" r="11" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<line x1="81" y1="37" x2="90" y2="46" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
    ),
    "neural-sim": (
        '<circle cx="50" cy="18" r="7" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="82" cy="50" r="7" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="50" cy="82" r="7" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="18" cy="50" r="7" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="50" cy="50" r="7" fill="currentColor"/>'
        '<path d="M50 25 V43 M50 57 V75 M25 50 H43 M57 50 H75" fill="none" stroke="currentColor" stroke-width="2" opacity="0.55"/>'
        '<path d="M56 25 L73 43 M56 75 L73 57 M44 25 L27 43 M44 75 L27 57" fill="none" stroke="currentColor" stroke-width="2" opacity="0.55"/>'
    ),
    "r0ut3r": (
        '<rect x="28" y="46" width="44" height="26" rx="6" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="44" cy="59" r="3" fill="currentColor"/><circle cx="56" cy="59" r="3" fill="currentColor"/>'
        '<path d="M20 30 A14 14 0 0 1 34 18" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M66 18 A14 14 0 0 1 80 30" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M26 22 A20 20 0 0 1 46 10" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" opacity="0.55"/>'
    ),
    "r4g": (
        '<rect x="26" y="16" width="48" height="68" rx="5" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M36 34 H56 M36 48 H64" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M72 20 V40 M62 30 H82" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
    ),
    "rust-infer": (
        '<rect x="28" y="34" width="44" height="44" rx="5" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<rect x="38" y="44" width="24" height="24" rx="3" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M34 20 V28 M50 18 V28 M66 20 V28 M34 84 V92 M50 84 V92 M66 84 V92 M18 48 H26 M18 64 H26 M74 48 H82 M74 64 H82" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M64 12 V24 M58 18 H70" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
    ),
    # --- rust ---
    "bl4ck1c3": (
        '<path d="M50 12 L80 26 V54 L50 68 L20 54 V26 Z" fill="none" stroke="currentColor" stroke-width="3" stroke-linejoin="round"/>'
        '<path d="M50 26 V52 M34 34 L66 46" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<circle cx="50" cy="39" r="3.5" fill="currentColor"/>'
    ),
    "hermes": (
        '<path d="M58 10 L24 54 H44 L40 92 L76 44 H54 Z" fill="none" stroke="currentColor" stroke-width="3" stroke-linejoin="round"/>'
        '<circle cx="62" cy="24" r="3" fill="currentColor"/>'
    ),
    "m0rpheus-rs": (
        '<path d="M64 12 A38 38 0 1 0 88 68 A30 30 0 1 1 64 12 Z" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M68 74 L72 82 L80 86 L72 90 L68 98 L64 90 L56 86 L64 82 Z" fill="currentColor"/>'
    ),
    "n3xus-flow": (
        '<circle cx="24" cy="30" r="8" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="76" cy="30" r="8" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="50" cy="76" r="8" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M31 36 L42 68 M69 36 L58 68 M24 22 V14 M76 22 V14 M50 68 V56" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
    ),
    "n3xusdb": (
        '<ellipse cx="50" cy="22" rx="26" ry="9" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M24 22 V68 A26 9 0 0 0 76 68 V22" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M24 44 A26 9 0 0 0 76 44" fill="none" stroke="currentColor" stroke-width="3" opacity="0.55"/>'
    ),
    # --- python ---
    "4sk": (
        '<path d="M20 26 Q20 16 30 16 H70 Q80 16 80 26 V50 Q80 60 70 60 H52 L44 70 V60 H30 Q20 60 20 50 Z" fill="none" stroke="currentColor" stroke-width="3" stroke-linejoin="round"/>'
        '<path d="M42 30 A10 10 0 0 1 58 32 A9 9 0 0 1 44 46" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="42" cy="52" r="3" fill="currentColor"/>'
    ),
    "jobtracker": (
        '<rect x="32" y="22" width="36" height="62" rx="5" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M42 16 H58 V24 H42 Z" fill="none" stroke="currentColor" stroke-width="3" stroke-linejoin="round"/>'
        '<path d="M40 46 L48 54 L61 38" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>'
    ),
    "s3arc": (
        '<path d="M50 14 L86 31 V69 L50 86 L14 69 V31 Z" fill="none" stroke="currentColor" stroke-width="3" stroke-linejoin="round"/>'
        '<path d="M14 31 L50 49 L86 31 M50 49 V86" fill="none" stroke="currentColor" stroke-width="2.5" opacity="0.6"/>'
    ),
    "termvault": (
        '<circle cx="50" cy="52" r="30" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<circle cx="50" cy="52" r="16" fill="none" stroke="currentColor" stroke-width="2.5" opacity="0.55"/>'
        '<path d="M50 22 V32 M50 72 V82 M20 52 H30 M70 52 H80" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M50 62 V70 H42 M50 62 V70" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
    ),
    "transitflow-nyc": (
        '<rect x="16" y="36" width="68" height="26" rx="8" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M28 49 H44 M56 49 H72" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<path d="M50 36 V26" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>'
        '<circle cx="50" cy="22" r="4" fill="none" stroke="currentColor" stroke-width="3"/>'
        '<path d="M24 74 H76 M20 66 H80" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" opacity="0.55"/>'
    ),
}


def parse_frontmatter(raw):
    m = re.match(r"^---\s*\n(.*?)\n---", raw, re.S)
    if not m:
        return {}
    front = {}
    for line in m.group(1).splitlines():
        kv = re.match(r"^(\w+):\s*(.*)$", line)
        if kv:
            front[kv.group(1)] = kv.group(2).strip().strip('"\'')
    return front


def wrap_title(title):
    lines = textwrap.wrap(title, width=24) or [""]
    if len(lines) > 2:
        lines = lines[:2]
        lines[1] = lines[1].rstrip(".,;:") + "\u2026"
    return lines


def logo_svg(project):
    title = project.get("title", "Untitled")
    slug = project.get("slug", "")
    pillar = project.get("pillar", "ai-ml")
    color = PILLAR_COLOR.get(pillar, PILLAR_COLOR["ai-ml"])
    tag = PILLAR_TAG.get(pillar, pillar.upper())
    lines = wrap_title(title)
    esc = html.escape

    mark = LOGOS.get(slug, "")
    assert mark, "no logo mark for slug %r" % slug
    k = 2.17
    mark_node = '<g transform="translate(%r %r) scale(%r)">' % (400 - 50 * k, 270 - 50 * k, k) + mark + "</g>"

    title_nodes = []
    if len(lines) == 1:
        title_nodes.append(
            '<text x="400" y="424" text-anchor="middle" font-family="\'Share Tech Mono\', monospace" '
            'font-size="49" fill="#e6edf3">%s</text>' % esc(lines[0])
        )
    else:
        title_nodes.append(
            '<text x="400" y="416" text-anchor="middle" font-family="\'Share Tech Mono\', monospace" '
            'font-size="41" fill="#e6edf3">%s</text>' % esc(lines[0])
        )
        title_nodes.append(
            '<text x="400" y="454" text-anchor="middle" font-family="\'Share Tech Mono\', monospace" '
            'font-size="41" fill="#e6edf3">%s</text>' % esc(lines[1])
        )

    return (
        '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 600" '
        'style="color:%s" role="img" aria-label="%s logo">\n'
        "  <defs>\n"
        '    <linearGradient id="bg" x1="0" y1="0" x2="1" y2="1">\n'
        '      <stop offset="0" stop-color="#0b0f17"/>\n'
        '      <stop offset="1" stop-color="#101828"/>\n'
        "    </linearGradient>\n"
        '    <radialGradient id="glow" cx="0.5" cy="0.5" r="0.5">\n'
        '      <stop offset="0" stop-color="%s" stop-opacity="0.14"/>\n'
        '      <stop offset="1" stop-color="%s" stop-opacity="0"/>\n'
        "    </radialGradient>\n"
        '    <pattern id="grid" width="24" height="24" patternUnits="userSpaceOnUse">\n'
        '      <circle cx="2" cy="2" r="1" fill="%s" opacity="0.06"/>\n'
        "    </pattern>\n"
        "  </defs>\n"
        '  <rect width="800" height="600" fill="url(#bg)"/>\n'
        '  <rect width="800" height="600" fill="url(#glow)"/>\n'
        '  <rect width="800" height="600" fill="url(#grid)"/>\n'
        '  <path d="M56 84 V56 H84" fill="none" stroke="%s" stroke-width="4" opacity="0.55"/>\n'
        '  <path d="M744 56 H716 V84" fill="none" stroke="%s" stroke-width="4" opacity="0.55"/>\n'
        '  <path d="M56 544 V572 H84" fill="none" stroke="%s" stroke-width="4" opacity="0.55"/>\n'
        '  <path d="M744 572 V544 H716" fill="none" stroke="%s" stroke-width="4" opacity="0.55"/>\n'
        '  <text x="400" y="72" text-anchor="middle" font-family="\'Share Tech Mono\', monospace" '
        'font-size="15" letter-spacing="6" fill="%s" opacity="0.85">%s</text>\n'
        "  %s\n"
        '  <line x1="300" y1="382" x2="500" y2="382" stroke="%s" stroke-width="2" opacity="0.5"/>\n'
        "%s"
        '  <text x="400" y="570" text-anchor="middle" font-family="\'Share Tech Mono\', monospace" '
        'font-size="14" fill="#18d26e" opacity="0.22">4rch3.io</text>\n'
        "</svg>\n"
    ) % (
        color,
        esc(title),
        color,
        color,
        color,
        color, color, color, color,
        color,
        tag,
        mark_node,
        color,
        "\n".join(title_nodes),
    )


def main():
    os.makedirs(OUT_DIR, exist_ok=True)
    written = 0
    for path in sorted(glob.glob(os.path.join(CONTENT_DIR, "*.md"))):
        slug = os.path.splitext(os.path.basename(path))[0]
        with open(path, encoding="utf-8") as fh:
            project = parse_frontmatter(fh.read())
        if not project or "title" not in project:
            continue
        project["slug"] = slug
        try:
            svg = logo_svg(project)
        except AssertionError as exc:
            sys.stderr.write("skipping %s: %s\n" % (slug, exc))
            continue
        out_path = os.path.join(OUT_DIR, slug + ".svg")
        with open(out_path, "w", encoding="utf-8") as fh:
            fh.write(svg)
        written += 1
    print("wrote %d logos -> %s" % (written, OUT_DIR))
    return 0 if written else 1


if __name__ == "__main__":
    sys.exit(main())