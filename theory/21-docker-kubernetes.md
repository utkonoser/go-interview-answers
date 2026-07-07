# Docker и Kubernetes

## 155. Чем контейнер отличается от виртуальной машины?

**VM** эмулирует железо: на гипервизоре свой guest OS, ядро, драйверы — тяжело (гигабайты, минуты на старт). **Контейнер** делит ядро хоста, изолирует процессы через namespaces и cgroups — легче (мегабайты, секунды). VM — сильная изоляция, разные ОС. Контейнер — изоляция процессов, одна ОС, идеален для микросервисов. Go-бинарник + минимальный образ — типичный деплой.

## 156. Что такое image и container?

**Image** — неизменяемый шаблон (слои файловой системы + метаданные: CMD, ENV, EXPOSE). **Container** — запущенный экземпляр image: процесс + writable layer поверх read-only слоёв. Один image → много контейнеров. Image собирается из Dockerfile, хранится в registry (Docker Hub, GitLab Registry, ECR).

## 157. Как устроен Dockerfile? Что такое слои и cache?

Dockerfile — инструкции сборки image. Каждая инструкция (`FROM`, `RUN`, `COPY`) создаёт **слой**. Слои кешируются: если инструкция и контекст не изменились — берётся из cache. Порядок важен: редко меняющиеся шаги (зависимости) — выше, часто меняющиеся (код) — ниже. `COPY` перед `RUN go build` — пересборка только при изменении кода.

## 158. `CMD` vs `ENTRYPOINT`, `COPY` vs `ADD`

**ENTRYPOINT** — фиксированная точка входа контейнера, не перезаписывается полностью аргументами `docker run`. **CMD** — аргументы по умолчанию, легко переопределяются. Паттерн: `ENTRYPOINT ["./app"]` + `CMD ["--config", "prod.yaml"]`. **COPY** — копирует файлы из контекста сборки (рекомендуется). **ADD** — то же + auto-распаковка tar и загрузка по URL (используй редко).

## 159. Volume vs bind mount — когда что

**Volume** — управляемое Docker хранилище, живёт отдельно от контейнера, бэкапится, переносится между хостами. **Bind mount** — примонтированная директория хоста (`/host/path:/container/path`), удобно для dev (hot reload). Volume — для prod данных (БД, файлы). Bind mount — для разработки и конфигов с хоста. Оба переживают перезапуск контейнера.

## 160. Как уменьшить размер образа Go-приложения?

Статическая сборка: `CGO_ENABLED=0 go build -ldflags="-s -w"`. Базовый image: `scratch` (пустой) или `distroless` (минимум libc). Не класть в image: исходники, тесты, toolchain. Multi-stage build: builder с golang — compile, runtime — только бинарник. Итог: 10–20 MB вместо 800 MB на `golang:latest`.

## 161. Что такое multi-stage build и зачем он Go-разработчику?

Несколько `FROM` в одном Dockerfile. Стадия 1: `golang:1.24` — компиляция. Стадия 2: `scratch`/`alpine` — копируем только бинарник `COPY --from=builder /app/main .`. В финальном image нет компилятора, исходников, кеша модулей. Меньше attack surface, быстрее pull в k8s.

## 162. Docker Compose — зачем нужен

Оркестрация нескольких контейнеров на одной машине для **локальной разработки**: app + PostgreSQL + Redis + Kafka одной командой `docker compose up`. Описание в `docker-compose.yml`: сервисы, сети, volumes, env. Не замена k8s в prod, но стандарт для dev/stage окружений.

## 163. Как пробросить переменные окружения и секреты в контейнер

**ENV** в Dockerfile — значения по умолчанию (не для секретов). `docker run -e VAR=value` или `env_file:` в compose. **Secrets** в compose/k8s — отдельно от image, монтируются как файлы (`/run/secrets/db_password`). Никогда не bake секреты в image. В k8s — `Secret` + env или volume mount.

## 164. Почему в контейнере не должно быть root-пользователя

Root в контейнере = root на хосте при escape из namespace (уязвимость ядра, misconfig). **USER** в Dockerfile — запуск от non-root. В k8s: `securityContext.runAsNonRoot: true`. Меньше рисков, требование security compliance. Go-приложение слушает порт >1024 или используй capabilities.

## 165. Зачем нужен Kubernetes? Что он даёт поверх Docker

Docker запускает контейнеры на одной машине. **Kubernetes** — оркестратор кластера: деплой на N нод, self-healing (перезапуск упавших pod), масштабирование, rolling updates, service discovery, балансировка, secrets/config. Декларативно: описываешь желаемое состояние, k8s приводит к нему.

## 166. Pod, Deployment, Service — что делает каждый

**Pod** — минимальная единица: один или несколько контейнеров с общей сетью/volume, ephemeral. **Deployment** — декларативное управление ReplicaSet: желаемое число pod, rolling update, rollback. **Service** — стабильный сетевой endpoint (ClusterIP/NodePort/LoadBalancer) для доступа к pod за балансировщиком. Типичный стек: Deployment → Pod → Service.

## 167. ReplicaSet vs Deployment vs StatefulSet

**ReplicaSet** — поддерживает N идентичных pod (заменён Deployment как API для пользователя). **Deployment** — ReplicaSet + стратегии обновления. **StatefulSet** — для stateful: стабильные имена pod (`app-0`, `app-1`), стабильные volume, упорядоченный старт/стоп. БД, Kafka, ZooKeeper — StatefulSet. Stateless Go API — Deployment.

## 168. ConfigMap и Secret — зачем и чем отличаются

**ConfigMap** — несекретная конфигурация (URLs, feature flags, yaml). **Secret** — чувствительные данные (пароли, токены), base64 в etcd (не шифрование по умолчанию!). Монтируются как env или файлы в pod. Отделяют конфиг от image — один image на все окружения. Secret — с RBAC и encryption at rest в prod.

## 169. Ingress — что это и зачем

**Ingress** — L7 маршрутизация HTTP/HTTPS в кластер: домены, пути, TLS termination. Нужен Ingress Controller (nginx, traefik). Внешний трафик → LoadBalancer → Ingress → Service → Pod. Альтернатива: Service type LoadBalancer на каждый сервис (дорого в облаке).

## 170. Liveness vs Readiness probe

**Liveness** — жив ли контейнер. Fail → kubelet **перезапускает** pod (deadlock, panic loop). **Readiness** — готов ли принимать трафик. Fail → pod **убирается из Service endpoints**, но не перезапускается (прогрев кеша, миграции, зависимость от БД). Go: liveness — `/healthz`, readiness — `/ready` с проверкой БД.

## 171. Rolling update и rollback

**Rolling update** — постепенная замена pod: `maxSurge` (сколько лишних pod), `maxUnavailable` (сколько можно убить). Zero-downtime при readiness probe. **Rollback** — `kubectl rollout undo deployment/app` к предыдущей ReplicaSet. История в `kubectl rollout history`. Стратегия `Recreate` — убить все, поднять новые (даунтайм, редко).

## 172. Что такое requests/limits (CPU, memory)

**requests** — гарантированные ресурсы для scheduling (на какую ноду поместить pod). **limits** — максимум: CPU throttling, OOMKill при превышении memory. Go: типично `requests: 100m CPU, 128Mi` / `limits: 500m, 512Mi` — подбирать по pprof. Без limits — noisy neighbor. Без requests — scheduler кладёт куда попало.

## 173. Что происходит, когда pod падает

Kubelet детектит (probe fail, exit code != 0) → удаляет pod → Deployment controller создаёт новый pod на доступной ноде → readiness OK → Service добавляет в endpoints. Данные в ephemeral container layer теряются — state в volume или внешней БД. Grace period: SIGTERM → graceful shutdown → SIGKILL через `terminationGracePeriodSeconds`.

## 174. Namespace — зачем

Логическое разделение ресурсов в кластере: `dev`, `staging`, `prod`, `team-billing`. Изоляция RBAC, квот, сетевых политик. Не изоляция на уровне железа — один кластер. `kubectl get pods -n prod`. По умолчанию — `default`.

## 175. HPA — горизонтальное автомасштабирование (базово)

**Horizontal Pod Autoscaler** — меняет `replicas` Deployment по метрикам: CPU, memory, custom (RPS из Prometheus). `minReplicas`–`maxReplicas`. Пример: CPU > 70% → добавить pod. Требует metrics-server или custom metrics adapter. Не путать с VPA (вертикальное — размер pod).

## 176. Как Go-сервис деплоится в k8s: healthcheck, graceful shutdown, SIGTERM

1) **Dockerfile** multi-stage → маленький image.

2) **Deployment** + probes на `/healthz` и `/ready`.

3) **SIGTERM** при остановке pod — `signal.Notify` + `server.Shutdown(ctx)` с таймаутом < `terminationGracePeriodSeconds`.

4) **PreStop hook** (опционально) — sleep для drain из Service.

5) **ConfigMap/Secret** для конфига.

6) **resources** по результатам нагрузочного теста.
