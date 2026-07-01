#!/usr/bin/env python3
"""Break inline numbered lists (1) 2) 3) ...) onto separate lines in theory/*.md."""

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent / "theory"
ITEM = re.compile(r"\d+\)")


def needs_fix(line: str) -> bool:
    if line.startswith("#") or line.startswith("```"):
        return False
    return len(ITEM.findall(line)) >= 2


def fix_line(line: str) -> str:
    if not needs_fix(line):
        return line

    s = re.sub(r": (1\)) ", r":\n\n\1 ", s := line, count=1)
    s = re.sub(r"; (\d+\)) ", r";\n\n\1 ", s)
    s = re.sub(r"\) (\d+\)) ", r")\n\n\1 ", s)
    s = re.sub(r"\. (\d+\)) ", r".\n\n\1 ", s)
    return s


def main() -> None:
    for path in sorted(ROOT.glob("[0-9]*.md")):
        text = path.read_text(encoding="utf-8")
        lines = text.splitlines()
        new_lines = [fix_line(line) for line in lines]
        if new_lines != lines:
            path.write_text("\n".join(new_lines) + ("\n" if text.endswith("\n") else ""), encoding="utf-8")
            print(path.name)


if __name__ == "__main__":
    main()
