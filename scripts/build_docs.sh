#!/bin/bash
# Сборка PDF-документов в docs/
# Запуск из корня: ./scripts/build_docs.sh

set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [ ! -d .venv ]; then
    echo "Creating .venv..."
    python3 -m venv .venv
fi

echo "Installing Python dependencies..."
.venv/bin/pip install -q -r requirements.txt

echo "Building interview guide..."
.venv/bin/python scripts/build_interview_pdf.py

if [ -f cv/Nikita_Selin_Resume.md ]; then
    echo "Building CV PDFs..."
    .venv/bin/python scripts/convert_resume_to_pdf.py cv/Nikita_Selin_Resume.md
    .venv/bin/python scripts/convert_resume_to_pdf.py cv/Nikita_Selin_Resume_EN.md
fi

echo "Done. Output: docs/"
