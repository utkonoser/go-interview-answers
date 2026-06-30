#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Собирает theory/ в один PDF (docs/interview-guide.pdf).

Usage:
    python3 scripts/build_interview_pdf.py
    python3 scripts/build_interview_pdf.py -o docs/my-guide.pdf
"""
from __future__ import annotations

import argparse
import glob
import os
import sys

import markdown
from weasyprint import HTML
from weasyprint.text.fonts import FontConfiguration

REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DEFAULT_OUTPUT = os.path.join(REPO_ROOT, "docs", "interview-guide.pdf")

SECTIONS: list[tuple[str, list[str], set[str]]] = [
    (
        "Теория",
        ["theory/*.md"],
        {"questions.md"},
    ),
]

HTML_TEMPLATE = """\
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="utf-8">
    <style>
        @page {{
            size: A4;
            margin: 1.2cm 1.4cm;
            @bottom-center {{
                content: counter(page);
                font-size: 7pt;
                color: #888;
            }}
        }}
        body {{
            font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
            font-size: 8pt;
            line-height: 1.28;
            color: #1a1a1a;
        }}
        h1 {{
            font-size: 13pt;
            margin: 0 0 6pt 0;
            padding-bottom: 3pt;
            border-bottom: 1.5px solid #333;
            page-break-after: avoid;
        }}
        h2 {{
            font-size: 9.5pt;
            margin: 8pt 0 4pt 0;
            page-break-after: avoid;
        }}
        h3 {{
            font-size: 8.5pt;
            margin: 6pt 0 3pt 0;
            page-break-after: avoid;
        }}
        h4 {{
            font-size: 8pt;
            margin: 5pt 0 2pt 0;
            page-break-after: avoid;
        }}
        p, li {{
            margin: 1.5pt 0;
        }}
        ul, ol {{
            margin: 3pt 0;
            padding-left: 14pt;
        }}
        li {{
            margin: 0.5pt 0;
        }}
        strong {{
            font-weight: 600;
        }}
        a {{
            color: #0066cc;
            text-decoration: none;
        }}
        hr {{
            margin: 8pt 0;
            border: none;
            border-top: 1px solid #ccc;
        }}
        pre {{
            background: #f5f5f5;
            border: 1px solid #ddd;
            border-radius: 2px;
            padding: 5pt 7pt;
            font-family: "Menlo", "Consolas", monospace;
            font-size: 6.5pt;
            line-height: 1.25;
            overflow-wrap: break-word;
            white-space: pre-wrap;
            page-break-inside: avoid;
        }}
        code {{
            font-family: "Menlo", "Consolas", monospace;
            font-size: 7pt;
            background: #f0f0f0;
            padding: 0 2pt;
            border-radius: 1px;
        }}
        pre code {{
            background: none;
            padding: 0;
        }}
        table {{
            border-collapse: collapse;
            width: 100%;
            margin: 5pt 0;
            font-size: 7pt;
        }}
        th, td {{
            border: 1px solid #ccc;
            padding: 2pt 5pt;
            text-align: left;
        }}
        th {{
            background: #eee;
        }}
        .page-break {{
            page-break-before: always;
        }}
        .section-title {{
            font-size: 16pt;
            text-align: center;
            margin: 0 0 20pt 0;
            padding: 40pt 0;
            page-break-after: always;
        }}
    </style>
</head>
<body>
{body}
</body>
</html>
"""


def collect_section_files(patterns: list[str], exclude_names: set[str]) -> list[str]:
    files: list[str] = []
    seen: set[str] = set()
    for pattern in patterns:
        for path in sorted(glob.glob(os.path.join(REPO_ROOT, pattern))):
            if not path.endswith(".md"):
                continue
            if os.path.basename(path) in exclude_names:
                continue
            real = os.path.realpath(path)
            if real in seen:
                continue
            seen.add(real)
            files.append(path)
    return files


def merge_sections(sections: list[tuple[str, list[str]]]) -> str:
    parts: list[str] = []
    first = True
    for title, files in sections:
        if not files:
            continue
        if not first:
            parts.append('<div class="page-break"></div>')
        first = False
        parts.append(f'<h1 class="section-title">{title}</h1>')
        for i, path in enumerate(files):
            with open(path, encoding="utf-8") as f:
                content = f.read().strip()
            if i > 0:
                parts.append('<div class="page-break"></div>')
            parts.append(content)
    return "\n\n".join(parts)


def markdown_to_html(md_content: str) -> str:
    return markdown.markdown(
        md_content,
        extensions=["extra", "sane_lists", "tables", "fenced_code"],
    )


def build_pdf(sections: list[tuple[str, list[str]]], output_path: str) -> None:
    merged = merge_sections(sections)
    html_body = markdown_to_html(merged)
    html_doc = HTML_TEMPLATE.format(body=html_body)

    font_config = FontConfiguration()
    HTML(string=html_doc, base_url=REPO_ROOT).write_pdf(
        output_path,
        font_config=font_config,
    )


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Собрать theory в один PDF",
    )
    parser.add_argument(
        "-o", "--output",
        default=DEFAULT_OUTPUT,
        help=f"Путь к PDF (по умолчанию: {DEFAULT_OUTPUT})",
    )
    args = parser.parse_args()

    resolved_sections: list[tuple[str, list[str]]] = []
    total = 0
    print("Sections:")
    for title, patterns, exclude in SECTIONS:
        files = collect_section_files(patterns, exclude)
        resolved_sections.append((title, files))
        total += len(files)
        print(f"\n{title} ({len(files)} files):")
        for path in files:
            rel = os.path.relpath(path, REPO_ROOT)
            print(f"  - {rel}")

    if total == 0:
        print("Error: no markdown files found", file=sys.stderr)
        return 1

    output_path = os.path.abspath(args.output)
    os.makedirs(os.path.dirname(output_path) or ".", exist_ok=True)

    print(f"\nMerging {total} files...")
    build_pdf(resolved_sections, output_path)
    size_kb = os.path.getsize(output_path) // 1024
    print(f"PDF created: {output_path} ({size_kb} KB)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
