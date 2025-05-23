# Инструкция по использованию Card Limit Manager

## 🚀 Запуск системы

1. **Клонируйте репозиторий:**

   ```bash
   git clone https://github.com/dpetrakov/card-limit-manager-workshop.git
   cd card-limit-manager-workshop
   ```

2. **Запустите все сервисы:**

   ```bash
   docker compose up
   ```

3. **Дождитесь запуска всех контейнеров:**
   - PostgreSQL (база данных)
   - Migrator (применение миграций)
   - CLM Service (backend API)
   - Frontend (React приложение)
   - Nginx Proxy (роутинг)

## 📝 Создание заявки на повышение лимита

1. **Откройте браузер** и перейдите на [http://localhost](http://localhost)

2. **Заполните форму:**

   - **Amount** - сумма лимита (например: 1000)
   - **Currency** - валюта (например: USD)
   - **Justification** - обоснование (например: "Business trip to conference")
   - **Desired Date** - желаемая дата (например: 2024-12-31)

3. **Нажмите "Submit Request"**

4. **Результат:**
   - При успехе: отобразится информация о созданной заявке
   - При ошибке: отобразится сообщение об ошибке

## 🗄️ Просмотр данных в базе

### Просмотр всех заявок:

```bash
docker exec clm_postgres psql -U postgres -d postgres -c "
SELECT lr.id, lr.amount, lr.currency, lr.justification, lr.status, u.name as user_name, lr.created_at
FROM limit_requests lr
JOIN users u ON lr.user_id = u.id
ORDER BY lr.created_at DESC;"
```

### Просмотр пользователей:

```bash
docker exec clm_postgres psql -U postgres -d postgres -c "
SELECT id, external_id, name, role FROM users ORDER BY role;"
```

### Подключение к базе данных:

```bash
docker exec -it clm_postgres psql -U postgres -d postgres
```

## 👥 Тестовые пользователи

Система автоматически создает следующих тестовых пользователей:

| Роль             | Имя               | Email                | Назначение             |
| ---------------- | ----------------- | -------------------- | ---------------------- |
| **EMPLOYEE**     | Test User         | test@example.com     | Создает заявки         |
| **TEAM_LEAD**    | Team Lead User    | teamlead@example.com | Первый уровень аппрува |
| **RISK_OFFICER** | Risk Officer User | risk@example.com     | Второй уровень аппрува |
| **CFO**          | CFO User          | cfo@example.com      | Финальный аппрув       |

## 🔗 Доступные эндпоинты

- **Главная страница**: [http://localhost](http://localhost)
- **Health check**: [http://localhost/healthz](http://localhost/healthz)
- **API**: [http://localhost/api/v1/requests](http://localhost/api/v1/requests)

## 🛠️ Разработка

### Просмотр логов:

```bash
# Все сервисы
docker compose logs -f

# Конкретный сервис
docker logs clm_service -f
docker logs clm_frontend -f
docker logs clm_nginx_proxy -f
```

### Перезапуск сервиса:

```bash
# Пересборка и перезапуск CLM сервиса
docker compose up -d --build clm-service

# Перезапуск без пересборки
docker compose restart clm-service
```

### Остановка системы:

```bash
docker compose down
```

### Полная очистка (включая данные):

```bash
docker compose down -v
```

## 🐛 Устранение неполадок

### Проблема: Ошибка 500 при отправке формы

**Решение:** Проверьте логи CLM сервиса:

```bash
docker logs clm_service --tail=20
```

### Проблема: Форма не загружается

**Решение:** Проверьте статус всех контейнеров:

```bash
docker compose ps
```

### Проблема: База данных недоступна

**Решение:** Проверьте статус PostgreSQL:

```bash
docker logs clm_postgres --tail=20
```

## 📊 Мониторинг

### Проверка состояния системы:

```bash
# Статус всех контейнеров
docker compose ps

# Использование ресурсов
docker stats

# Health check nginx
curl http://localhost/healthz
```
