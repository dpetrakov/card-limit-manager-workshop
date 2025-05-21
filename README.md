# Card Limit Manager – AI Workshop

Этот репозиторий содержит материалы и исходный код демонстрационного проекта, созданного в рамках воркшопа по **архитектурному проектированию и разработке с использованием AI в IDE Cursor/Windsurf**.

## 📌 Цель

Создать микросервисную систему **Card Limit Manager** — решение для автоматизации подачи и согласования заявок на повышение лимитов корпоративных карт.  
В процессе демонстрируются лучшие практики взаимодействия с AI-ассистентом на всех этапах: от требований и архитектуры до кода, тестов и развёртывания.

---

## 🧱 Архитектура

- **Frontend**: React 18 + TypeScript + Vite
- **Backend**: Go 1.22 + Gin (2 микросервиса)
- **Database**: PostgreSQL 16
- **DevTools**: Docker Compose, sql-migrate, Adminer, JWT, Gopls, Air
- **AI**: Cursor/Windsurf IDE + `AGENTS.md` + голосовой ввод требований

---

## 🚀 Быстрый старт

```bash
git clone https://github.com/dpetrakov/card-limit-manager-workshop.git
cd card-limit-manager-workshop
docker compose up
```

После старта:

* UI доступен на [http://localhost:5173](http://localhost:5173)
* Backend API — [http://localhost:8080](http://localhost:8080)
* Postgres — localhost:5432 (user: `postgres`, pass: `postgres`)
* Adminer — [http://localhost:8081](http://localhost:8081)

---

## 📁 Структура проекта

```bash
.
├── docs/                   # документация (требования, архитектура, спецификации, схема БД)
│   ├── architecture.md
│   ├── database.dbml
│   └── requirements/
├── services/               # микросервисы на Go (clm/, notif/)
├── frontend/               # SPA на React + Vite
├── database/migrations/    # SQL-миграции
├── integration/examples/   # контракты внешних систем (read-only)
├── .devcontainer/          # devcontainer config для VS Code / Cursor / Windsurf
├── docker-compose.yml
└── AGENTS.md               # правила генерации и кодирования для Cursor / Windsurf
```


---

## 🧪 Покрытие и тестирование

* Unit-тесты (Go, React) ≥ 80% покрытия
* Интеграционные тесты (Go + Postgres)
* End-to-End сценарии (Playwright)
* CI-ready: тесты запускаются локально и в pipeline

---

## 🤖 AGENTS.md

Внутри ключевых каталогов (`docs/specs/`, `services/`, `frontend/`) находятся файлы `AGENTS.md`, которые содержат локальные инструкции по оформлению спецификаций и коду. Эти инструкции имеют приоритет над `.cursorrules`.

---

## 📚 Лицензия

MIT License. Этот проект используется в обучающих целях.


