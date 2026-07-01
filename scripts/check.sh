#!/bin/bash

# Скрипт для локальной проверки кода (аналогично CI)
# Запуск: ./scripts/check.sh

set -e

echo "🔍 Запуск локальных проверок кода..."
echo "=================================="

# Проверяем, что мы в корне проекта
if [ ! -f "go.mod" ]; then
    echo "❌ Ошибка: запустите скрипт из корня проекта"
    exit 1
fi

# Переходим в папку с задачами
cd go-interview-tasks

echo "📦 Загрузка зависимостей..."
go mod download

echo "🧪 Запуск тестов..."
go test -v ./tests/...

echo "🔎 Code review (solution)..."
go test -tags=solution -race ./code-review/tests/...

echo "📊 Запуск тестов с покрытием..."
go test -v -coverprofile=coverage.out ./tests/... ./strings/...

echo "📈 Анализ покрытия..."
go tool cover -func=coverage.out

echo "⚡ Запуск бенчмарков..."
go test -bench=. -benchmem ./tests/...

echo "🎨 Проверка форматирования..."
if [ "$(go fmt ./... | wc -l)" -gt 0 ]; then
    echo "❌ Код не отформатирован. Запустите: go fmt ./..."
    exit 1
else
    echo "✅ Код отформатирован корректно"
fi

echo "🔍 Статический анализ (go vet)..."
go vet ./...

echo "✨ Проверка golangci-lint пропущена (не используется в проекте)"

echo "=================================="
echo "🎉 Все проверки завершены успешно!"
echo "💡 Теперь можно делать push в GitHub"
