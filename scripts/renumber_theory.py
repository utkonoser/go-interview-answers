#!/usr/bin/env python3
"""Renumber ## N. headers in theory/*.md sequentially and rebuild questions.md."""

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
THEORY = ROOT / "theory"
HEADER_RE = re.compile(r"^(## )(\d+)(\. .+)$")


def main() -> None:
    files = sorted(THEORY.glob("[0-9]*.md"))
    index: list[tuple[str, int, str]] = []
    n = 1

    for path in files:
        text = path.read_text(encoding="utf-8")
        lines = text.splitlines(keepends=True)
        out: list[str] = []
        changed = False

        for line in lines:
            m = HEADER_RE.match(line.rstrip("\n"))
            if m:
                title = m.group(3)[2:]  # after ". "
                index.append((path.name, n, title))
                new_line = f"## {n}. {title}\n"
                if new_line != line:
                    changed = True
                out.append(new_line)
                n += 1
            else:
                out.append(line)

        if changed:
            path.write_text("".join(out), encoding="utf-8")
            print(f"updated {path.name}")

    # questions.md
    sections: dict[str, list[tuple[int, str]]] = {}
    for fname, num, title in index:
        sections.setdefault(fname, []).append((num, title))

    lines = [
        "# Полный список вопросов\n",
        "\n",
        "Типичные вопросы по Go и смежным темам на собеседованиях. Ответы — в соответствующих файлах `theory/`.\n",
        "\n",
    ]
    for path in files:
        fname = path.name
        if fname not in sections:
            continue
        title = path.read_text(encoding="utf-8").split("\n", 1)[0].lstrip("# ").strip()
        lines.append(f"## {title} (`{fname}`)\n")
        lines.append("\n")
        for num, q in sections[fname]:
            lines.append(f"{num}. {q}\n")
        lines.append("\n")

    for subdir, prefix, title in (
        ("magnit-tech", "M", "Magnit Tech"),
        ("avito", "A", "Avito"),
    ):
        company_dir = THEORY / subdir
        company_files = sorted(company_dir.glob("[0-9]*.md")) if company_dir.is_dir() else []
        if not company_files:
            continue
        lines.append(f"## {title} (`{subdir}/`)\n")
        lines.append("\n")
        lines.append(f"Отдельный блок вопросов с собесов {title}. Ответы — в `theory/{subdir}/`.\n")
        lines.append("\n")
        for path in company_files:
            text = path.read_text(encoding="utf-8")
            for line in text.splitlines():
                m = re.match(rf"^## ({prefix}\d+[a-z]?)\. (.+)$", line)
                if m:
                    lines.append(f"- **{m.group(1)}** {m.group(2)} → `{path.name}`\n")
            lines.append("\n")

    (THEORY / "questions.md").write_text("".join(lines), encoding="utf-8")
    print(f"questions.md: {n - 1} questions")


if __name__ == "__main__":
    main()
