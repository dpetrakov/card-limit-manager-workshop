# Card Limit Manager – AI Workshop

Этот репозиторий содержит материалы и исходный код демонстрационного проекта, созданного в рамках воркшопа по **архитектурному проектированию и разработке с использованием AI в IDE Cursor/Windsurf**.

## 📌 Цель

Создать микросервисную систему **Card Limit Manager** — решение для автоматизации подачи и согласования заявок на повышение лимитов корпоративных карт.  
В процессе демонстрируются лучшие практики взаимодействия с AI-ассистентом на всех этапах: от требований и архитектуры до кода, тестов и развёртывания.

---

## 🧱 Архитектура

- **Frontend**: React 18 + TypeScript + Vite
- **Backend**: Go 1.22 + Gin (CLM сервис)
- **Database**: PostgreSQL 16 с автоматическими миграциями
- **Proxy**: Nginx для роутинга frontend/backend
- **DevTools**: Docker Compose, sql-migrate

---

## 🚀 Быстрый старт

```bash
git clone https://github.com/dpetrakov/card-limit-manager-workshop.git
cd card-limit-manager-workshop
docker compose up
```

После старта:

- **Основное приложение**: [http://localhost](http://localhost) — React форма для создания заявок
- **Health check**: [http://localhost/healthz](http://localhost/healthz) — статус nginx
- **Backend API**: [http://localhost/api/v1/](http://localhost/api/v1/) — CLM сервис
- **Postgres**: localhost:5432 (user: `postgres`, pass: `postgres`)

---

## 👥 Тестовые пользователи

При запуске автоматически создаются тестовые пользователи для демонстрации workflow:

| Роль             | Имя               | Email                | ID                                     |
| ---------------- | ----------------- | -------------------- | -------------------------------------- |
| **EMPLOYEE**     | Test User         | test@example.com     | `a005d32d-6190-477c-b23e-38c44eaaaae0` |
| **TEAM_LEAD**    | Team Lead User    | teamlead@example.com | `b123e567-e89b-12d3-a456-426614174000` |
| **RISK_OFFICER** | Risk Officer User | risk@example.com     | `c789f012-e89b-12d3-a456-426614174001` |
| **CFO**          | CFO User          | cfo@example.com      | `d456c789-e89b-12d3-a456-426614174002` |

Заявки создаются от имени `Test User` (EMPLOYEE).

---

## ✅ Реализованный функционал

### Frontend (React + TypeScript)

- ✅ Форма создания заявки на повышение лимита
- ✅ Валидация полей (сумма, валюта, обоснование, дата)
- ✅ Отправка данных в backend API
- ✅ Отображение результата (успех/ошибка)

### Backend (Go + Gin)

- ✅ REST API эндпоинт `POST /api/v1/requests`
- ✅ Валидация входящих данных
- ✅ Сохранение в PostgreSQL
- ✅ Подключение к базе данных

### Database (PostgreSQL)

- ✅ Автоматические миграции при запуске
- ✅ Таблицы: users, limit_requests, approval_steps, audit_log
- ✅ Тестовые данные (пользователи)

---

## 🔗 API Endpoints

### Создание заявки

```http
POST /api/v1/requests
Content-Type: application/json

{
  "amount": 1000.00,
  "currency": "USD",
  "justification": "Business trip to conference",
  "desired_date": "2024-12-31"
}
```

**Response 201:**

```json
{
  "id": "fa401de5-aaa1-49ec-bcb1-78e5a7fc243f",
  "user_id": "a005d32d-6190-477c-b23e-38c44eaaaae0",
  "amount": 1000,
  "currency": "USD",
  "justification": "Business trip to conference",
  "desired_date": "2024-12-31",
  "status": "PENDING_TEAM_LEAD",
  "created_at": "2025-05-23T07:24:23.337009Z",
  "updated_at": "2025-05-23T07:24:23.337009Z"
}
```

---

## 🗄️ Работа с базой данных

### Просмотр заявок

```bash
docker exec clm_postgres psql -U postgres -d postgres -c "
SELECT lr.id, lr.amount, lr.currency, lr.justification, lr.status, u.name as user_name, lr.created_at
FROM limit_requests lr
JOIN users u ON lr.user_id = u.id
ORDER BY lr.created_at DESC;"
```

### Просмотр пользователей

```bash
docker exec clm_postgres psql -U postgres -d postgres -c "
SELECT id, external_id, name, role FROM users ORDER BY role;"
```

---

## 📁 Структура проекта

```bash
.
├── docs/                   # документация (требования, архитектура, спецификации, схема БД)
│   ├── architecture.md
│   ├── database.dbml
│   └── requirements/
├── services/               # микросервисы на Go
│   └── clm/               # Card Limit Manager сервис
├── frontend/               # SPA на React + Vite
├── database/migrations/    # SQL-миграции
│   ├── 001_initial_schema.sql    # Создание таблиц
│   └── 002_add_test_data.sql     # Тестовые пользователи
├── nginx-proxy-conf/       # конфигурация nginx
├── integration/examples/   # контракты внешних систем (read-only)
├── .devcontainer/          # devcontainer config для VS Code / Cursor / Windsurf
├── docker-compose.yml
└── AGENTS.md               # правила генерации и кодирования для Cursor / Windsurf
```

---

## 🧪 Покрытие и тестирование

- Unit-тесты (Go, React) ≥ 80% покрытия
- Интеграционные тесты (Go + Postgres)
- End-to-End сценарии (Playwright)
- CI-ready: тесты запускаются локально и в pipeline

---

## 🤖 AGENTS.md

Внутри ключевых каталогов (`docs/specs/`, `services/`, `frontend/`) находятся файлы `AGENTS.md`, которые содержат локальные инструкции по оформлению спецификаций и коду. Эти инструкции имеют приоритет над `.cursorrules`.

---

## 📚 Лицензия

MIT License. Этот проект используется в обучающих целях.
