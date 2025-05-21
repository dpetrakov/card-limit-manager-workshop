# Card-Limit Manager – project-wide Agent rules

# PROJECT_CONTEXT

Основные артефакты репозитория:

```
  ─ docs/
    ├─ project_overview.md          # цели проекта, границы, связи
    ├─ architecture.md              # КЛЮЧЕВОЙ документ: диаграммы + описание компонентов
    ├─ deployment_strategy.md       # подход к развёртыванию
    ├─ database.dbml                # КЛЮЧЕВОЙ документ: ER-схема БД (DBML)
    ├─ requirements.md              # функциональные и нефункциональные требования проекта (R-1, R-2...)
    └─specs/                        # подробные спецификации (по каждому R-…)
      ├─ CLM.md                     # детальные спецификации для Limit-Request Service (код = CLM)
      └─ NOTIF.md                   # детальные спецификации для Notification Service (код = NOTIF)
  ─ tasks/tasks.md                  # номерной список задач и статусы
  ─ services/
    ├─ clm/                         # Go-сервис Limit-Request Service (код = CLM)
    └─ notif/                       # Go-сервис Notification Service (код = NOTIF)
  ─ frontend/                       # React 18 (Vite + TS) SPA
  ─ database/migrations/            # SQL-миграции (sql-migrate style)
  ─ integration/examples/           # примеры вызовов внешних систем (READ-ONLY!)
```

Весь стек разворачивается локально одной командой:
`docker compose up` — поднимаются Postgres 16, backend-сервисы, frontend, Adminer.

# SPEC_STYLE

- Архитектурные диаграммы — Mermaid.
  - Блоки: однострочные названия без кавычек/скобок.
  - Стрелки помечать протоколом (REST, gRPC, Kafka).
  - В схеме показываем только системы Card-Limit Manager и внешние интеграции;
    сервисные службы (Postgres, Keycloak и т.п.) НЕ изображаем.
- Диаграммы последовательностей — PlantUML.
  - Участники = компоненты из `architecture.md`.
  - При наличии алгоритма шаг помечается `ALG_<№ алгоритма>`.
- После каждой диаграммы обязательна таблица шагов со столбцами:
  | # | Actor → Actor | Data | Contract/Schema | Algorithm |
- Все OpenAPI-описания — версия 3.1, формат JSON, файл в `docs/specs/`.
- Любые упоминания схем БД ссылаются на `docs/database.dbml`
  (пример: `frame_id uuid [ref: > frames.id, note: "Frame link"]`).

# ANALYST_TASKS

1. Спецификация должна быть достаточной для реализации;
   лишняя «бумажная» детализация запрещена.
2. Каждый R-… обязан порождать отдельный файл спецификации в `docs/specs/`.

# DEVELOPMENT_TASKS

1. Определи микросервис по префиксу задачи (`D-CLM-*`, `D-NOTIF-*`).
2. Изучи артефакты: `architecture.md`, `requirements.md`, `specs/R-<code>-.md`, `database.dbml`, `database/migrations/`.
3. Реализуй функциональность в `services/<code>/`, следуя Clean Architecture:
   `cmd/`, `internal/`, `pkg/`.
   - Go 1.22 + Gin
   - Postgres driver pgx
   - Migrations через sql-migrate
4. Напиши unit-tests ≥ 80 % покрытия (go test -cover, jest for React).
5. Убедись, что `docker compose up` и все тесты проходят без ошибок.
6. Логирование: только английский, уровни INFO/WARN/ERROR.

# VALIDATION_RULES

1. Каждое требование R-… имеет спецификацию и отражено в OpenAPI/BDBL.
2. Диаграммы и шаги строго соответствуют architecture.md.
3. Все поля API сопоставлены с БД; несоответствия = TODO.
4. Примеры в integration/examples/ валидны и неизменны.

# ERROR_HANDLING

- Несоответствие архитектуре или схеме БД → пересмотреть и исправить.
- Отсутствие алгоритма при сложной логике → добавить TODO.
- Не хватает ссылок на схемы/контракты → TODO с вопросом к команде.
- Файл integration/examples/\* изменять запрещено; несоответствие сообщаем владельцу.
- PlantUML/Mermaid не рендерится → исправить синтаксис.

# CODE_QUALITY_STANDARDS

1. Unit-tests ≥ 80 % (за исключением boilerplate) и проходят в CI.
2. Линтеры: golangci-lint (Go), eslint + prettier (React).
3. Комментарии только для нетривиальной логики.
4. Логи без PII, на английском.

# IDE_RULES

- Всегда проверяй наличие описания (requirement/spec) перед генерацией кода.
- Если документа нет — спроси, генерировать ли skeleton.
- При scaffold’е соблюдай структуру: backend/, frontend/, docs/, database/.
- Генерируй docker-compose файл при первом появлении нового сервиса или БД.
- После изменения схемы БД предлагай сгенерировать файл миграции.
- **AGENTS.md:** перед работой с файлами
  - docs/specs/**, services/**, frontend/\*\*
    проверяй наличие `AGENTS.md`.  
    — Инструкции из самого глубокого `AGENTS.md` имеют приоритет  
    — Они переопределяют любые правила вышележащего `AGENTS.md`,  
     кроме прямых указаний из пользовательского запроса.  
    — Если в `AGENTS.md` есть программные проверки (скрипты, `make test` и т.д.),  
     запусти их и убедись, что они проходят до коммита изменений.

# AGENTS_RULES

- С scope-ом всю подпапку, где лежит файл, и глубже.
- Инструкции относятся только к файлам внутри этого scope.
- В случае конфликта: глубже-лежащий `AGENTS.md` > вышележащий > `AGENTS.md`.
- Примеры расположения:  
   • docs/specs/AGENTS.md — уточнения для спецификаций (формат, naming).  
   • services/AGENTS.md — требования к Go-кодингу, Git-workflow.  
   • frontend/AGENTS.md — гайдлайны по React/TS, ESLint, Storybook и т.д.
- Все указания по стилю, структуре, тестированию и CI-командам из `AGENTS.md`
  обязательно выполняй для затрагиваемых файлов.
