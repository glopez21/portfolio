#!/usr/bin/env python3
"""Generate the social/OG image (1200x630 PNG) for 4rch3.io.

  python3 scripts/gen_og_image.py [out_png]

Matches the site's dark-cyber brand: #040404 background, CRT-green
(#18d26e) glowing wordmark, pillar tags, scanlines + faint grid.
Requires Pillow.
"""
import glob
import os
import sys

from PIL import Image, ImageDraw, ImageFilter, ImageFont

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OUT = sys.argv[1] if len(sys.argv) > 1 else os.path.join(ROOT, "frontend", "assets", "img", "og-image.png")

W, H = 1200, 630
BG = (4, 4, 4)
GREEN = (24, 210, 110)
MINT = (143, 242, 196)
MUTED = (120, 145, 130)
PILLARS = [
    ("CYBERSECURITY", (24, 210, 110)),
    ("AI/ML", (34, 211, 238)),
    ("PYTHON", (232, 193, 90)),
    ("RUST", (255, 138, 101)),
]

FONT_CANDIDATES = [
    "/usr/share/fonts/TTF/MesloLGSNerdFontMono-Regular.ttf",
    "/usr/share/fonts/TTF/MesloLGSDZNerdFontMono-Regular.ttf",
    "/usr/share/fonts/liberation/LiberationMono-Regular.ttf",
    "/usr/share/fonts/noto/NotoSansMono-Regular.ttf",
]


def mono_font(size):
    for pat in FONT_CANDIDATES:
        for path in glob.glob(pat):
            return ImageFont.truetype(path, size)
    for path in glob.glob("/usr/share/fonts/**/*Mono*", recursive=True):
        if path.endswith((".ttf", ".otf")):
            return ImageFont.truetype(path, size)
    return ImageFont.load_default(size)


def glow_text(base, xy, text, font, fill, glow_fill, radii=(6, 14, 28)):
    """Draw text with layered blurred glow passes beneath the sharp pass."""
    for radius in radii:
        layer = Image.new("RGBA", base.size, (0, 0, 0, 0))
        d = ImageDraw.Draw(layer)
        d.text(xy, text, font=font, fill=glow_fill)
        layer = layer.filter(ImageFilter.GaussianBlur(radius))
        base.alpha_composite(layer)
    d = ImageDraw.Draw(base)
    d.text(xy, text, font=font, fill=fill)


def main():
    img = Image.new("RGBA", (W, H), BG + (255,))

    # Faint terminal grid.
    grid = Image.new("RGBA", (W, H), (0, 0, 0, 0))
    gd = ImageDraw.Draw(grid)
    for x in range(0, W, 40):
        gd.line([(x, 0), (x, H)], fill=GREEN + (14,), width=1)
    for y in range(0, H, 40):
        gd.line([(0, y), (W, y)], fill=GREEN + (14,), width=1)
    img.alpha_composite(grid)
    draw = ImageDraw.Draw(img)

    # Brand wordmark.
    f_brand = mono_font(150)
    glow_text(img, (84, 150), "4rch3", f_brand, GREEN + (242,),
              GREEN + (110,))
    bbox = draw.textbbox((84, 150), "4rch3", font=f_brand)
    io_x = bbox[2] + 18
    f_io = mono_font(86)
    glow_text(img, (io_x, 150 + 34), ".io", f_io, MINT + (235,),
              MINT + (90,), radii=(4, 10, 20))

    # Greek sub-tagline.
    f_tag = mono_font(34)
    draw.text((90, 350), "αρχή · η πρώτη αρχή", font=f_tag, fill=MUTED + (255,))

    # Pillar tags with color swatch dots.
    f_pill = mono_font(28)
    x = 90
    y = 470
    for label, color in PILLARS:
        dot_r = 7
        draw.ellipse([x, y + 13, x + 2 * dot_r, y + 13 + 2 * dot_r],
                     fill=color + (255,))
        draw.text((x + 2 * dot_r + 14, y), label, font=f_pill,
                  fill=color + (235,))
        w = draw.textbbox((0, 0), label, font=f_pill)[2]
        x += 2 * dot_r + 14 + w + 46

    # Scanlines (CRT feel).
    scan = Image.new("RGBA", (W, H), (0, 0, 0, 0))
    sd = ImageDraw.Draw(scan)
    for y in range(0, H, 4):
        sd.line([(0, y), (W, y)], fill=(0, 0, 0, 26), width=1)
    img.alpha_composite(scan)

    # Vignette edge glow.
    vin = Image.new("RGBA", (W, H), (0, 0, 0, 0))
    vd = ImageDraw.Draw(vin)
    vd.rectangle([0, 0, W, H], outline=GREEN + (90,), width=3)
    vin = vin.filter(ImageFilter.GaussianBlur(12))
    img.alpha_composite(vin)
    img.convert("RGB").save(OUT, "PNG", optimize=True)
    print(f"wrote {OUT} ({os.path.getsize(OUT)} bytes)")


if __name__ == "__main__":
    main()
