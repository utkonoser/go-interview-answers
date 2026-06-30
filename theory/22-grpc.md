# gRPC

## 184. Что такое gRPC и чем отличается от REST?

**gRPC** — RPC-фреймворк поверх HTTP/2 с **Protocol Buffers** (бинарная сериализация). REST — обычно JSON over HTTP/1.1, ресурсо-ориентированный. gRPC: строгий контракт (.proto), быстрее и компактнее, bidirectional streaming, встроенные deadlines. REST: проще отладка (curl), шире поддержка браузеров. В микросервисах внутри кластера — часто gRPC; наружу — REST/GraphQL.

## 185. Что такое Protocol Buffers (protobuf)?

Язык-независимый формат описания структур данных в `.proto` файлах. `protoc` генерирует Go-структуры и сериализацию. Бинарный, компактный, быстрый parse. Схема версионируется: не удалять/не менять номера полей, добавлять новые с новыми tag. Альтернативы: JSON Schema, Avro. В Go: `google.golang.org/protobuf`, кодоген через `protoc-gen-go-grpc`.

## 186. Какие типы вызовов есть в gRPC?

1) **Unary** — один запрос, один ответ (как обычный RPC). 2) **Server streaming** — клиент шлёт один запрос, сервер — поток ответов. 3) **Client streaming** — клиент шлёт поток, сервер — один ответ. 4) **Bidirectional streaming** — оба потока одновременно. Streaming для больших данных, real-time, без буферизации всего в память.

## 187. Как работают deadline и cancellation в gRPC?

Клиент задаёт deadline: `ctx, cancel := context.WithTimeout(ctx, 3*time.Second)`. Deadline передаётся в metadata по HTTP/2. Сервер проверяет `ctx.Done()` — отмена при таймауте или cancel клиента. В Go gRPC deadline/cancellation интегрирован с `context.Context`. Важно: прокидывать ctx во все downstream вызовы и БД-запросы.

## 188. Как обрабатывать ошибки в gRPC?

Статус через `status.Error(codes.NotFound, "user not found")`. Коды: `OK`, `InvalidArgument`, `NotFound`, `AlreadyExists`, `DeadlineExceeded`, `Unavailable`, `Internal`. Клиент: `status.FromError(err)` → code + message. **Не** возвращать stack trace клиенту. Retry только на `Unavailable`, `DeadlineExceeded` с backoff. В Go: `google.golang.org/grpc/status`, `codes`.

## 189. Что такое interceptors (middleware) в gRPC?

Аналог HTTP middleware: цепочка до/после handler. **UnaryInterceptor** — логирование, auth, метрики, recovery от panic, tracing. Server: `grpc.UnaryInterceptor(loggingInterceptor)`. Client: для retry, tracing propagation. Несколько interceptors — `grpc.ChainUnaryInterceptor(a, b, c)`.

## 190. gRPC vs REST — когда что выбрать?

**gRPC**: внутренние микросервисы, high throughput, streaming, строгий контракт, polyglot (много языков). **REST**: публичный API, мобильные/веб-клиенты, простота, кеширование HTTP. **Оба**: gRPC-gateway генерирует REST из .proto. Go-бэкенд часто: gRPC внутри, REST наружу.

## 191. Как тестировать gRPC-сервисы в Go?

1) **In-process**: `bufconn` — listener в памяти, без сети. 2) **grpc.NewClient** на тестовый сервер. 3) Моки интерфейсов сервисного слоя (не gRPC напрямую). 4) Интеграционные тесты с реальным сервером на случайном порту. Проверять status codes, metadata, deadline behaviour.
