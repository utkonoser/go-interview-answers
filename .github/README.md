# GitHub Actions

Автоматическая проверка Go-кода при push и pull request в `main`, `master`, `develop`.

## Workflow: Go Tests (`workflows/go-tests.yml`)

| Шаг | Что делает |
|-----|------------|
| Setup Go 1.24 | Установка Go, кэш модулей |
| `go test -v` | Unit-тесты `go-interview-tasks/tests/...` |
| Coverage | `go test -coverprofile` + `go tool cover` |
| Benchmarks | `go test -bench=.` |
| `go fmt` | Проверка форматирования |
| `go vet` | Статический анализ |

### Build Docs (`workflows/build-docs.yml`)

При push в `main` (кроме изменений только в `docs/`):

- сборка PDF через WeasyPrint
- коммит `docs/*.pdf` ботом `github-actions[bot]`

## Локальная проверка

Перед push — из корня проекта:

```bash
./scripts/check.sh
```

Ручной запуск тех же команд:

```bash
cd go-interview-tasks
go mod download
go test -v ./tests/...
go test -v -coverprofile=coverage.out ./tests/... ./strings/...
go tool cover -func=coverage.out
go test -bench=. -benchmem ./tests/...
go fmt ./...
go vet ./...
```

## Структура

```
.github/
├── workflows/
│   ├── go-tests.yml      # CI pipeline
│   └── build-docs.yml    # сборка PDF в docs/
└── README.md             # этот файл

scripts/
└── check.sh            # локальный аналог CI
```

## Требования

- Go 1.24+ (как в CI)
- Для PDF-сборки отдельно: Python 3, `./scripts/build_docs.sh`

## Troubleshooting

**Тесты падают в CI**
1. Посмотреть логи в Actions
2. Воспроизвести локально: `./scripts/check.sh`

**`check.sh` не запускается**
```bash
chmod +x scripts/check.sh
# запускать из корня репозитория, где лежит go.mod
```

**Ошибка форматирования**
```bash
cd go-interview-tasks && go fmt ./...
```
