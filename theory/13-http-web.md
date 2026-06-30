# HTTP и веб-разработка

## 91. Как работает `net/http` пакет?

`net/http` — стандартный пакет для HTTP клиентов и серверов. Сервер: `http.ListenAndServe()`, `http.HandleFunc()`, `http.Handler` интерфейс. Клиент: `http.Get()`, `http.Post()`, `http.Client` для настройки. Обрабатывает HTTP протокол, управляет соединениями, парсит заголовки. `http.Request` содержит метод, URL, заголовки, тело. `http.ResponseWriter` для записи ответа. Основа веб-разработки в Go.

## 92. В чем разница между `http.Handle` и `http.HandleFunc`?

`http.Handle(pattern, handler)` принимает `http.Handler` интерфейс (метод `ServeHTTP`). `http.HandleFunc(pattern, handler)` принимает функцию `func(http.ResponseWriter, *http.Request)`. `HandleFunc` — удобная обертка, конвертирующая функцию в `http.Handler`. Функционально эквивалентны, выбор зависит от стиля: функции — `HandleFunc`, структуры с состоянием — `Handle` с кастомным `Handler`.

## 93. Как создать middleware для HTTP handlers?

Middleware — функция, оборачивающая handler для добавления функциональности (логирование, аутентификация, CORS). Паттерн: `func(http.HandlerFunc) http.HandlerFunc` или `func(http.Handler) http.Handler`. Цепочка middleware через композицию функций. Примеры: логирование запросов, проверка авторизации, добавление заголовков, rate limiting. Можно использовать библиотеки типа `gorilla/mux` для удобной композиции.

## 94. Как работает `context` в HTTP запросах?

Каждый HTTP запрос имеет `context`, доступный через `r.Context()`. Автоматически отменяется при закрытии соединения клиентом. Используется для: передачи request-scoped данных (request ID, user info), таймаутов, отмены долгих операций. Передавать `ctx` во все асинхронные операции. `http.Request.WithContext()` создает новый запрос с контекстом. Критично для корректной обработки отмены запросов.

## 95. Как реализовать graceful shutdown HTTP сервера?

Использовать `http.Server.Shutdown(ctx)` вместо прямого закрытия. Процесс: перехватить сигнал (SIGTERM, SIGINT), вызвать `server.Shutdown()` с таймаутом, сервер перестанет принимать новые запросы и дождется завершения текущих (до таймаута). После таймаута принудительно закрывает соединения. Позволяет обработать in-flight запросы, критично для production без потери данных.

