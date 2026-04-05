# Bradobrei Party Backend

Бэкенд информационной системы сети барбершопов `Bradobrei Party` на `Go + Gin + GORM + PostgreSQL/PostGIS`.

## Что уже есть

- REST API для пользователей, сотрудников, салонов, услуг, материалов, бронирований, отзывов, платежей и отчётов.
- Swagger UI для локальной разработки.
- JWT-аутентификация с поддержкой `Bearer <token>` и raw JWT в dev-сценариях.
- Автомиграция таблиц через GORM.
- Подготовка под PostGIS: `salons.location` хранится как `geometry(Point,4326)`.
- E2E и unit-тесты backend-слоя.
- Базовая инфраструктура для HTML/PDF-отчётов через `internal/reports` и Gotenberg.

## Структура

- `cmd/api/main.go` — точка входа API.
- `cmd/report_example/main.go` — demo-команда для рендера HTML/PDF отчёта без HTTP endpoint.
- `internal/models` — ORM-модели и report view-models для печатных документов.
- `internal/dto` — DTO запросов и ответов.
- `internal/handlers` — HTTP-обработчики.
- `internal/services` — прикладная логика.
- `internal/repository` — доступ к данным.
- `internal/reports` — HTML-шаблоны отчётов, CSS и клиент Gotenberg.
- `tests` — e2e/integration тесты через HTTP и PostgreSQL.
- `test_artifacts` — сохранённые JSON-артефакты ответов.
- `docs` — сгенерированная Swagger-документация.

## Отчёты 2.2.x из ТЗ

- `2.2.1 Реестр персонала`
  - `GET /api/v1/reports/employees`
  - HTML/PDF модель: `models.EmployeeRegistryReportDocument`
  - шаблон: `internal/reports/templates/employees.html`

- `2.2.2 Аналитический отчёт об операционной активности филиалов`
  - `GET /api/v1/reports/salon-activity`
  - HTML/PDF модель: `models.SalonActivityReportDocument`
  - шаблон: `internal/reports/templates/salon_activity.html`

- `2.2.3 Статистика востребованности услуг`
  - `GET /api/v1/reports/service-popularity`
  - сейчас учитываются бронирования со статусами `CONFIRMED`, `IN_PROGRESS`, `COMPLETED`

- `2.2.4 Отчёт о производительности и ресурсо-затратности мастеров`
  - `GET /api/v1/reports/master-activity`
  - HTML/PDF модель: `models.MasterActivityReportDocument`
  - шаблон: `internal/reports/templates/master_activity.html`

- `2.2.5 Журнал мониторинга качества обслуживания и обратной связи`
  - `GET /api/v1/reports/reviews`
  - `POST /api/v1/reviews`
  - HTML/PDF модель: `models.ReviewsReportDocument`
  - шаблон: `internal/reports/templates/reviews.html`

- `2.2.6 Ведомость движения товарно-материальных ценностей`
  - `GET /api/v1/reports/inventory-movement`
  - параметры: `from`, `to`, `salon_id`

- `2.2.7 Анализ клиентской лояльности и удержания`
  - `GET /api/v1/reports/client-loyalty`
  - параметры: `from`, `to`
  - `bonus_status` пока считается эвристически, пока отдельная бонусная модель не введена

- `2.2.8 Реестр отменённых и нереализованных бронирований`
  - `GET /api/v1/reports/cancelled-bookings`
  - параметры: `from`, `to`

- `2.2.9 Сводный финансовый отчёт по статьям издержек`
  - `GET /api/v1/reports/financial-summary`
  - параметры: `from`, `to`, `salon_id`
  - строки отчёта пока формируются из статусов `payments`, так как отдельной модели статей расходов ещё нет

## Зависимости

Основные Go-зависимости уже зафиксированы в `go.mod`.

Скачать зависимости:

```bash
go mod download
```

Если нужно пересобрать Swagger:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g ./cmd/api/main.go -o ./docs
```

## Настройка окружения

Минимальные переменные лежат в `backend/.env`.

Важно:

- `DB_*` — подключение к PostgreSQL/PostGIS
- `PORT` — порт API
- `JWT_SECRET` — ключ подписи JWT
- `GOTENBERG_URL` — адрес сервиса Gotenberg для генерации PDF

Пример:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=bradobrei
DB_SSLMODE=disable

GIN_MODE=debug
PORT=9000
JWT_SECRET=your-super-secret-jwt-key-change-in-production
GOTENBERG_URL=http://localhost:3000
```

## База данных

Используется PostgreSQL с PostGIS.

Минимально нужно:

```sql
CREATE DATABASE bradobrei;
```

При старте приложение само пытается выполнить:

```sql
CREATE SCHEMA IF NOT EXISTS public;
CREATE EXTENSION IF NOT EXISTS postgis;
```

Если база создана вручную и у пользователя не хватает прав, эти команды нужно выполнить отдельно под пользователем с нужными привилегиями.

## Локальный запуск API

Из папки `backend`:

```bash
go run ./cmd/api
```

Сборка бинарника:

```bash
go build -o bradobrei-api ./cmd/api
```

## Swagger UI

Документация доступна по адресам:

```text
http://localhost:9000/swagger/index.html
http://localhost:9000/docs
```

Для защищённых endpoint:

1. Выполните `POST /api/v1/auth/login`
2. Скопируйте токен
3. Нажмите `Authorize`
4. Вставьте:

```text
Bearer <ваш_jwt_токен>
```

## Тесты

### Быстрые unit-тесты

```bash
go test ./internal/...
```

### E2E и интеграционные тесты

```bash
go test ./tests -v -timeout 60s
```

### Весь backend

```bash
go test ./... -v -timeout 120s
```

Что уже покрыто:

- auth
- employees
- salons
- bookings
- payments
- reviews
- все 9 report endpoint’ов
- helper-функции middleware, report parsing, coordinates normalization, salon IDs normalization

Артефакты успешных сценариев сохраняются в:

```text
backend/test_artifacts/api_outputs.json
```

## Smoke-скрипт для отчётов

Для ручной проверки всех 9 JSON-отчётов и новых складских операций используйте bash-сценарий:

```bash
./scripts/run_reports_smoke.sh
```

Что делает сценарий:

- логинится под уже существующим администратором `admin/password`
- создаёт салон, услугу, материал и норму расхода
- создаёт закупку материала через `POST /api/v1/material-expenses`, чтобы пополнить склад
- нанимает мастера
- регистрирует клиента
- создаёт подтверждённое, отменённое и нереализованное бронирования
- создаёт платежи и отзыв
- вызывает все 9 report endpoint’ов
- отдельно вызывает `POST /api/v1/services/{id}/use` как ручное списание материалов по услуге
- сохраняет JSON-ответы в `backend/test_artifacts/report_smoke`

Зависимости локального smoke-сценария:

- `curl`
- `jq`

Если API поднят не на `http://localhost:9000`, можно передать адрес первым аргументом:

```bash
./scripts/run_reports_smoke.sh http://localhost:8080/api/v1
```

Или через переменные окружения:

```bash
API_BASE_URL=http://localhost:8080/api/v1 ADMIN_USERNAME=admin ADMIN_PASSWORD=password ./scripts/run_reports_smoke.sh
```

Старый PowerShell-сценарий оставлен в репозитории как вспомогательный, но для массивов JSON и повторяемых прогонов предпочтительнее bash-вариант.

## Gotenberg

Для генерации PDF используется отдельный контейнер Gotenberg. Основной `docker-compose.yml` при этом не меняется.

Отдельный compose-файл для локальной разработки:

```text
docker-compose.gotenberg.yml
```

Запуск из корня репозитория:

```bash
docker compose -f docker-compose.gotenberg.yml up -d
```

Остановка:

```bash
docker compose -f docker-compose.gotenberg.yml down
```

В dev-режиме сервис будет доступен на:

```text
http://localhost:3000
```

Важно:

- backend отправляет в Gotenberg `index.html` и дополнительные assets, например `report.css`
- Gotenberg складывает все загруженные файлы в одну плоскую директорию
- в HTML нужно ссылаться на asset по имени файла, например `report.css`, без подпапок

Официальная документация Gotenberg:

- Installation: https://gotenberg.dev/docs/getting-started/installation
- HTML to PDF: https://gotenberg.dev/docs/convert-with-chromium/convert-html-to-pdf

## `internal/reports`

Папка `internal/reports` отвечает за:

- HTML-шаблоны печатных форм
- CSS для печатного документа
- клиент к Gotenberg
- рендер HTML в `[]byte`
- рендер HTML -> PDF в `[]byte`

Сейчас там уже есть:

- `client.go` — Go-клиент для Gotenberg
- `renderer.go` — рендерер HTML/PDF
- `templates/base.html` — базовый layout
- `templates/report.css` — общие стили печатного документа
- `templates/employees.html`
- `templates/salon_activity.html`
- `templates/master_activity.html`
- `templates/reviews.html`

## Пример HTML/PDF без endpoint

Для проверки шаблона без подключения отдельного API endpoint есть demo-команда:

```bash
go run ./cmd/report_example
```

Она:

- собирает пример отчёта `2.2.1`
- сохраняет HTML в `test_artifacts/employees_report_example.html`
- пытается получить PDF через Gotenberg и сохранить его в `test_artifacts/employees_report_example.pdf`

Если нужно только HTML:

```bash
go run ./cmd/report_example --skip-pdf
```

Если Gotenberg уже поднят:

```bash
docker compose -f ../docker-compose.gotenberg.yml up -d
go run ./cmd/report_example
```

## Полезные команды

Форматирование:

```bash
gofmt -w ./cmd ./internal ./tests
```

Проверка сборки:

```bash
go build ./...
```

## Частые проблемы

### `schema for creating objects is not selected` / `SQLSTATE 3F000`

Обычно это значит:

- в БД нет схемы `public`
- у пользователя БД пустой `search_path`
- `.env` указывает не на тот экземпляр PostgreSQL

### `type "geometry" does not exist` / `SQLSTATE 42704`

Обычно это значит, что в текущей базе не включён `postgis`.

Решение:

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
```

### `go build` на Windows падает из-за кеша

Иногда мешает `%LOCALAPPDATA%\\go-build`.

Попробуйте:

```bash
go clean -cache
```

или локальный кеш:

```powershell
New-Item -ItemType Directory -Force '.gocache' | Out-Null
$env:GOCACHE=(Resolve-Path '.gocache')
go test ./... -v -timeout 120s
```
