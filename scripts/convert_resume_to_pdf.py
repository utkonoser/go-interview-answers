#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Конвертирует Markdown-резюме в PDF (результат в docs/).

Usage:
    python3 scripts/convert_resume_to_pdf.py
    python3 scripts/convert_resume_to_pdf.py cv/Nikita_Selin_Resume.md
"""
import os
import sys

import markdown
from weasyprint import HTML
from weasyprint.text.fonts import FontConfiguration

REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DOCS_DIR = os.path.join(REPO_ROOT, "docs")
DEFAULT_INPUT = os.path.join(REPO_ROOT, "cv", "Nikita_Selin_Resume.md")


def main() -> int:
    if len(sys.argv) > 1:
        input_file = os.path.abspath(sys.argv[1])
    else:
        input_file = DEFAULT_INPUT

    if not os.path.exists(input_file):
        print(f"Error: File '{input_file}' not found!")
        return 1

    os.makedirs(DOCS_DIR, exist_ok=True)
    output_file = os.path.join(DOCS_DIR, os.path.basename(input_file).replace(".md", ".pdf"))

    with open(input_file, encoding="utf-8") as f:
        md_content = f.read()

    html_content = markdown.markdown(md_content, extensions=["extra"])

    html_doc = f"""
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        @page {{
            size: A4;
            margin: 0.8cm;
        }}
        body {{
            font-family: 'Arial', 'Helvetica', sans-serif;
            font-size: 10pt;
            line-height: 1.22;
            margin: 0;
            padding: 0;
        }}
        h1 {{
            font-size: 19pt;
            margin: 0 0 3.5pt 0;
            padding: 0;
        }}
        h2 {{
            font-size: 13pt;
            margin: 4.5pt 0 2.5pt 0;
            padding: 0;
            border-bottom: 1px solid #ccc;
            padding-bottom: 1.2pt;
        }}
        p {{
            margin: 1.2pt 0;
            padding: 0;
            font-size: 10pt;
        }}
        ul {{
            margin: 0.8pt 0;
            padding-left: 14pt;
        }}
        li {{
            margin: 0.4pt 0;
            padding: 0;
            font-size: 10pt;
        }}
        strong {{
            font-weight: bold;
            font-size: 10pt;
        }}
        a {{
            color: #0066cc;
            text-decoration: none;
            font-size: 10pt;
        }}
        hr {{
            margin: 2.5pt 0;
            border: none;
            border-top: 1px solid #ccc;
        }}
        h1 + p {{
            margin: 0 0 2.5pt 0;
        }}
        h1 + p + p {{
            margin: 0 0 2.5pt 0;
        }}
        h2 + p {{
            margin: 1.2pt 0;
        }}
        ul {{
            margin-bottom: 1.8pt;
        }}
        ul + p {{
            margin-top: 4.5pt;
        }}
        p + ul + p {{
            margin-top: 1.8pt;
        }}
    </style>
</head>
<body>
{html_content}
</body>
</html>
"""

    font_config = FontConfiguration()
    HTML(string=html_doc, base_url=os.path.dirname(input_file)).write_pdf(
        output_file,
        font_config=font_config,
    )

    print(f"PDF created: {output_file}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
