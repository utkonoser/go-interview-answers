# Go Interview Preparation

Материалы для подготовки к собеседованиям по Go: теория, практические задачи, генерация PDF.

[![Go Tests](https://github.com/utkonoser/interviews/workflows/Go%20Tests/badge.svg)](https://github.com/utkonoser/interviews/actions/workflows/go-tests.yml)

## Структура проекта

```
interviews/
├── theory/                 # Теория: вопросы и ответы по темам
│   ├── 01-basics.md
│   ├── ...
│   ├── 20-kafka.md
│   └── questions.md        # Полный список вопросов
├── go-interview-tasks/     # Практические задачи (LeetCode-style)
│   ├── strings/
│   └── tests/
├── cv/                     # Резюме (markdown-исходники)
│   ├── Nikita_Selin_Resume.md
│   └── Nikita_Selin_Resume_EN.md
├── scripts/                # Скрипты проверки и сборки документов
│   ├── check.sh            # Локальные Go-проверки (как в CI)
│   ├── build_docs.sh       # Сборка всех PDF
│   ├── build_interview_pdf.py
│   └── convert_resume_to_pdf.py
├── docs/                   # PDF (генерируются CI при push в main)
├── requirements.txt        # Python-зависимости для сборки PDF
├── .venv/                  # Python venv (в .gitignore)
└── go.mod
```

## Теория

20 тематических файлов в `theory/`:

| Тема | Файл |
|------|------|
| Основы Go | `01-basics.md` |
| Структуры данных | `02-data-structures.md` |
| Конкурентность, каналы, интерфейсы | `03`–`05` |
| Память, ошибки, тесты, модули | `06`–`09` |
| Практические задачи (код) | `10-practical-tasks.md` |
| System design, БД, HTTP, perf | `11`–`14` |
| Security, networking, OS | `15`–`17` |
| БД (общая теория), SOLID, Kafka | `18`–`20` |

Полный оглавление — в `theory/questions.md`.

## Практические задачи

6 задач по строкам в `go-interview-tasks/strings/`:

- Palindrome (LeetCode #125)
- Reverse String (#344)
- Backspace Compare (#844)
- Remove Stars (#2390)
- Remove Adjacent Duplicates (#1047)
- Valid Parentheses (#20)

## Быстрый старт

```bash
git clone <repository-url>
cd interviews

# Теория — читать markdown в theory/
# Практика — запустить тесты
./scripts/check.sh

# Собрать PDF (теория + резюме)
./scripts/build_docs.sh
# → docs/interview-guide.pdf
# → docs/Nikita_Selin_Resume.pdf
```

### Python-окружение (для PDF)

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
python scripts/build_interview_pdf.py
python scripts/convert_resume_to_pdf.py cv/Nikita_Selin_Resume.md
```

## CI/CD

GitHub Actions (`go-tests.yml`) при каждом push/PR:

- unit-тесты и покрытие
- бенчмарки
- `go fmt` и `go vet`

Локальный аналог: `./scripts/check.sh`

## PDF в `docs/`

При push в `main` workflow [build-docs.yml](.github/workflows/build-docs.yml) пересобирает PDF из `theory/` и `cv/` и коммитит в `docs/`:

- `interview-guide.pdf`
- `Nikita_Selin_Resume.pdf`
- `Nikita_Selin_Resume_EN.pdf`

Локально: `./scripts/build_docs.sh`

## Как готовиться

1. Пройти теорию по темам в `theory/` (или сгенерировать PDF и слушать/читать)
2. Решить задачи в `go-interview-tasks/`, прогнать тесты
3. Повторить слабые места по `theory/questions.md`
