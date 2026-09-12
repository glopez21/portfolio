#!/usr/bin/env python3
"""Render real terminal output as branded screenshot PNGs.

  .venv/bin/python scripts/capture_terminal.py [slug ...]

Commands are listed in scripts/terminal_captures.json. Each command is
actually executed (FORCE_COLOR set so rich/typer keep their colors even
when piped); its ANSI-colored output is parsed and re-rendered as a
terminal window PNG at frontend/assets/img/screenshots/<slug>.png.
Only add commands that are safe/read-only — the output must be real.
"""
import glob
import json
import os
import re
import subprocess
import sys

from PIL import Image, ImageDraw, ImageFont

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OUT_DIR = os.path.join(ROOT, "frontend", "assets", "img", "screenshots")
CONFIG = os.path.join(ROOT, "scripts", "terminal_captures.json")

BG = (13, 17, 23)
BAR = (22, 27, 34)
FG_DEFAULT = (200, 210, 220)
GREEN = (24, 210, 110)
MUTED = (110, 125, 140)
MAX_LINES = 80
COLS = 100
FONT_SIZE = 19
LINE_H = 27

FONT_CANDIDATES = [
    "/usr/share/fonts/TTF/MesloLGSNerdFontMono-Regular.ttf",
    "/usr/share/fonts/TTF/MesloLGSDZNerdFontMono-Regular.ttf",
    "/usr/share/fonts/liberation/LiberationMono-Regular.ttf",
]

# xterm 256-color palette -----------------------------------------------------

STD = [
    (0, 0, 0), (205, 0, 0), (0, 205, 0), (205, 205, 0),
    (0, 0, 238), (205, 0, 205), (0, 205, 205), (229, 229, 229),
    (127, 127, 127), (255, 0, 0), (0, 255, 0), (255, 255, 0),
    (92, 92, 255), (255, 0, 255), (0, 255, 255), (255, 255, 255),
]


def color_256(n):
    if n < 16:
        return STD[n]
    if n < 232:
        n -= 16
        r, g, b = (n // 36), ((n // 6) % 6), (n % 6)
        f = lambda c: 0 if c == 0 else 55 + c * 40
        return (f(r), f(g), f(b))
    g = 8 + (n - 232) * 10
    return (g, g, g)


def default_font():
    for pat in FONT_CANDIDATES:
        for path in glob.glob(pat):
            return ImageFont.truetype(path, FONT_SIZE)
    return ImageFont.load_default(FONT_SIZE)


ANSI_RE = re.compile(r"\x1b\[([0-9;]*)m")
STRIP_RE = re.compile(r"\x1b\][^\x07]*\x07|\x1b\[[0-9;?]*[A-LN-Za-ln-z]|\x1b[@-_]")


def parse_ansi_line(line):
    """-> list of (char, fg, bg, bold) cells for one line."""
    cells = []
    fg, bg, bold = FG_DEFAULT, None, False
    pos = 0
    for m in ANSI_RE.finditer(line):
        for ch in line[pos:m.start()]:
            cells.append((ch, fg, bg, bold))
        params = [p for p in m.group(1).split(";") if p != ""]
        i = 0
        while i < len(params):
            p = int(params[i])
            if p == 0:
                fg, bg, bold = FG_DEFAULT, None, False
            elif p == 1:
                bold = True
            elif p == 22:
                bold = False
            elif 30 <= p <= 37:
                fg = color_256(p - 30)
            elif p == 90 <= p or (p == 90 or p in range(90, 98)):
                fg = color_256(p - 90 + 8)
            elif 40 <= p <= 47:
                bg = color_256(p - 40)
            elif p == 39:
                fg = FG_DEFAULT
            elif p == 49:
                bg = None
            elif p == 38 and i + 1 < len(params):
                if params[i + 1] == "5" and i + 2 < len(params):
                    fg = color_256(int(params[i + 2])); i += 2
                elif params[i + 1] == "2" and i + 4 < len(params):
                    fg = tuple(int(x) for x in params[i + 2:i + 5]); i += 4
            elif p == 48 and i + 1 < len(params):
                if params[i + 1] == "5" and i + 2 < len(params):
                    bg = color_256(int(params[i + 2])); i += 2
                elif params[i + 1] == "2" and i + 4 < len(params):
                    bg = tuple(int(x) for x in params[i + 2:i + 5]); i += 4
            i += 1
        pos = m.end()
    for ch in line[pos:]:
        cells.append((ch, fg, bg, bold))
    return cells


def run_capture(entry):
    env = dict(os.environ, FORCE_COLOR="1", COLUMNS=str(COLS),
               TERM="xterm-256color", LINES=str(MAX_LINES))
    p = subprocess.run(entry["cmd"], shell=True, cwd=entry.get("cwd", ROOT),
                       capture_output=True, text=True, env=env, timeout=60)
    text = (p.stdout + ("\n" + p.stderr if p.stderr.strip() else "")).rstrip()
    lines = STRIP_RE.sub("", text).split("\n")
    if len(lines) > MAX_LINES:
        lines = lines[: MAX_LINES - 1] + ["…  (output truncated)"]
    return lines


def render(slug, lines, font):
    prompt = f"$ {lines[0][:0]}"  # keep pattern: we draw the command above
    bar_h = 46
    pad = 22
    cw = font.getbbox("M")[2] - font.getbbox("M")[0] or font.getbbox("M")[2]
    char_w = font.getlength("M")
    img_w = int(pad * 2 + char_w * COLS) + 2
    img_h = bar_h + pad * 2 + LINE_H * (len(lines) + 1)
    img = Image.new("RGB", (img_w, img_h), BG)
    d = ImageDraw.Draw(img)

    # Title bar
    d.rectangle([0, 0, img_w, bar_h], fill=BAR)
    for i, c in enumerate([(255, 95, 86), (255, 189, 46), (39, 201, 63)]):
        d.ellipse([16 + i * 26, bar_h // 2 - 7, 30 + i * 26, bar_h // 2 + 7], fill=c)
    d.text((16 + 3 * 26 + 12, bar_h // 2 - 12), f"4rch3.io — {slug}",
           font=font, fill=MUTED)

    y = bar_h + pad
    # Command echo line (green prompt)
    d.text((pad, y), "$", font=font, fill=GREEN)
    cmd = TERMINAL_CMDS.get(slug, "")
    d.text((pad + int(char_w * 2), y), cmd, font=font, fill=FG_DEFAULT)
    y += LINE_H

    for line in lines:
        x = pad
        for ch, fg, bg, bold in parse_ansi_line(line):
            if ch == "\t":
                ch = "    "
            if bg:
                d.rectangle([x, y - 3, x + char_w * len(ch), y + LINE_H - 6], fill=bg)
            if ch.strip():
                col = tuple(min(255, c + 30) for c in fg) if bold else fg
                d.text((x, y), ch, font=font, fill=col)
            x += char_w * len(ch)
        y += LINE_H

    # Border accent
    d.rectangle([0, 0, img_w - 1, img_h - 1], outline=(24, 210, 110, 60), width=1)
    out = os.path.join(OUT_DIR, f"{slug}.png")
    img.save(out, "PNG", optimize=True)
    print(f"{slug}: {len(lines)} lines -> {out} ({os.path.getsize(out) // 1024} KB)")


TERMINAL_CMDS = {}


def main():
    global TERMINAL_CMDS
    with open(CONFIG) as fh:
        entries = json.load(fh)
    TERMINAL_CMDS = {e["slug"]: e["cmd"] for e in entries}
    only = sys.argv[1:] or list(TERMINAL_CMDS)
    os.makedirs(OUT_DIR, exist_ok=True)
    font = default_font()
    for e in entries:
        if e["slug"] in only:
            lines = run_capture(e)
            render(e["slug"], lines, font)


if __name__ == "__main__":
    main()
