# Nikita Selin

**Email:** utkonoser@gmail.com | **Telegram:** @ctrlshiftesc

[GitHub](https://github.com/utkonoser) | [LinkedIn](https://www.linkedin.com/in/nikita-selin28)

---

## Work Experience

**Senior Go Developer** | [Premier.One](https://premier.one/) | May 2024 — Present

Online cinema team. Feed service (Go) replacing a legacy Python service:

- Led client API migration to feeds (Strangler Pattern, Kafka, Debezium, response comparison) — phased rollout of core client traffic
- Rewrote key high-traffic endpoints (showcases, favorites, card groups, next best item, etc.): legacy reverse-engineering, response parity — latency 300–400 ms → 10–80 ms; improved main showcase UX
- Optimized watch history endpoint: SQL rewrite and indexes — DB LA ~5 → ~0.5, response ~180 ms → ~25 ms; reduced peak DB load
- Designed and implemented feeds admin REST API — content management endpoints (shelves, showcases, etc.)
- **Tech Stack:** Go, PostgreSQL, KeyDB, Kafka, Debezium, Docker, Grafana, Prometheus, Loki, GitLab

**Middle Go Developer** | [AgentAPP](https://agentapp.ru/) | November 2023 — May 2024

Car-data integrations team. Services for collecting automotive information:

- Built Go services for data collection and normalization for client integrations
- Designed an in-house image analysis module (pixel-level segmentation and pattern matching) to automate visual verification in data ingestion pipelines — improved integration reliability
- Configured resilient ingestion via proxies, retries, and error handling under rate limits
- Maintained and evolved production integrations
- **Tech Stack:** Go, PostgreSQL, MongoDB, RabbitMQ, Docker, REST API

**Junior Go Developer** | [Wildberries](https://www.wildberries.ru/) | February 2023 — November 2023

Partner pickup points team (Tickets, Writeoff services):

- Supported production Tickets/Writeoff services: incidents, stability under load
- Wrote SQL reports and complex data exports with business calculations
- Enhanced Telegram bot for pickup point owners
- **Tech Stack:** Go, gRPC, PostgreSQL, Kafka, Docker, Telegram Bot API, YouTracker

---

## Personal Projects

**[ConvertArt](https://convertart.ru/)** | 2025 — Present

Free stateless service in production: image → crystal mosaic, cross-stitch, paint-by-numbers, machine embroidery. No signup — files are processed in-memory and not persisted.

- Designed microservices architecture (gateway, mosaic-api, stitch workers behind nginx); built backend, frontend, and image processing pipelines; CI/CD (GitHub Actions → VPS) and monitoring
- **Tech Stack:** Go, Docker, nginx, GitHub Actions, JavaScript, Grafana, Prometheus, Loki

---

## Education

**Saint Petersburg Electrotechnical University "LETI"** named after V.I. Ulyanov (Lenin)
Faculty of Information-Measuring and Biotechnical Systems, Instrumentation Engineering | 2017

---

## Achievements

**[VKCup 2022/23](https://cups.online/ru/contests/VKCup2022/)** — Top 16 in Go category

---

## Skills

**Stack:** Go | PostgreSQL, MongoDB, KeyDB | Docker, Linux, Git, GitLab CI | gRPC, REST API | Kafka, RabbitMQ, Debezium | Grafana, Prometheus, Loki

**Languages:** Russian (Native), English (B1, fluent technical reading)

---

## Interests

Play and watch basketball; enjoy watching American football and soccer
