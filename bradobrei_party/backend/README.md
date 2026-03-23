# Bradobrei Party Backend

Бэкенд информационной системы сети барбершопов `Bradobrei Party` на `Go + Gin + GORM + PostgreSQL/PostGIS`.

## Что уже есть

- REST API для пользователей, сотрудников, салонов, услуг, материалов, бронирований, отзывов и отчётов.
- Автомиграция таблиц через GORM.
- Отчёты по части требований `2.2.x` из ТЗ.
- Подготовка под PostGIS: поле `salons.location` хранится как `geometry(Point,4326)`.

## Структура

- `cmd/api/main.go` — точка входа, настройка БД, роутов и middleware.
- `internal/models` — ORM-модели GORM.
- `internal/dto` — DTO запросов и ответов.
- `internal/handlers` — HTTP-эндпоинты.
- `internal/services` — прикладная логика.
- `internal/repository` — доступ к данным.
- `migrations` — SQL-миграции расширений и индексов.

## Отчёты 2.2.x из ТЗ

Ниже перечислены уже заложенные места, где реализуются или подготавливаются требования раздела `2.2`.

- `2.2.1 Реестр персонала`
  - эндпоинт: `GET /api/v1/reports/employees`
  - handler: `internal/handlers/report_handler.go`
  - repository: `internal/repository/report_repo.go`
  - модель: `internal/models/models.go` (`EmployeeProfile`)

- `2.2.2 Аналитический отчёт об операционной активности филиалов`
  - эндпоинт: `GET /api/v1/reports/salon-activity`
  - handler: `internal/handlers/report_handler.go`
  - repository: `internal/repository/report_repo.go`
  - модели: `Salon`, `Booking`, `BookingItem`

- `2.2.3 Статистика востребованности услуг`
  - эндпоинт: `GET /api/v1/reports/service-popularity`
  - handler: `internal/handlers/report_handler.go`
  - repository: `internal/repository/report_repo.go`
  - модель: `Service`

- `2.2.4 Отчёт о производительности и ресурсо-затратности мастеров`
  - эндпоинт: `GET /api/v1/reports/master-activity`
  - handler: `internal/handlers/report_handler.go`
  - repository: `internal/repository/report_repo.go`
  - модели: `Booking`, `EmployeeProfile`, `ServiceMaterial`

- `2.2.5 Журнал мониторинга качества обслуживания и обратной связи`
  - эндпоинт: `GET /api/v1/reports/reviews`
  - дополнительный рабочий эндпоинт: `POST /api/v1/reviews`
  - handler: `internal/handlers/report_handler.go`, `internal/handlers/review_handler.go`
  - модель: `Review`

- `2.2.6 Ведомость движения ТМЦ`
  - пока подготовлена модельная база
  - модели: `Inventory`, `Material`, `ServiceMaterial`
  - следующий шаг: отдельный отчётный запрос/эндпоинт

- `2.2.7 Анализ клиентской лояльности и удержания`
  - пока подготовлена база на уровне модели `Booking`
  - следующий шаг: агрегирующий отчёт по клиентам

- `2.2.8 Реестр отменённых и нереализованных бронирований`
  - пока подготовлена база на уровне `Booking.Status`
  - следующий шаг: отдельный репозиторный запрос и отчётный эндпоинт

- `2.2.9 Сводный финансовый отчёт по статьям издержек`
  - пока не реализован

## Установка зависимостей

Базовые зависимости модуля:

```bash
go mod download
```

Если хотите подключить Swagger UI, понадобятся ещё пакеты:

```bash
go get github.com/swaggo/files github.com/swaggo/gin-swagger github.com/swaggo/swag
go install github.com/swaggo/swag/cmd/swag@latest
```

После этого можно сгенерировать документацию:

```bash
swag init -g ./cmd/api/main.go -o ./docs
```

Сейчас в коде уже есть базовая подготовка в `main.go`, но сами swagger-зависимости в этот репозиторий ещё не подтянуты автоматически.

## Настройка БД

Используется PostgreSQL с PostGIS.

Важно: в `docker-compose.yml` и `.env` сейчас разные учётные данные.

- `docker-compose.yml`: `admin / password`
- `.env`: `postgres / ?aug1337`

Их нужно привести к одному варианту.

Если база была удалена и создана заново вручную, достаточно:

```sql
CREATE DATABASE bradobrei;
```

Схему `public` приложение теперь пытается создать само при старте:

```sql
CREATE SCHEMA IF NOT EXISTS public;
```

Для геоданных PostGIS также нужно расширение:

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
```

Если хотите проверить руками:

```sql
CREATE SCHEMA IF NOT EXISTS public;
CREATE EXTENSION IF NOT EXISTS postgis;
GRANT ALL ON SCHEMA public TO postgres;
```

Если используете пользователя `admin`, замените `postgres` на `admin`.

## Запуск через Docker

Из корня проекта:

```bash
docker compose up --build
```

## Локальный запуск

Из папки `backend`:

```bash
go run ./cmd/api
```

## Сборка

Из папки `backend`:

```bash
go build -o bradobrei-api ./cmd/api
```

Запуск собранного бинарника:

```bash
./bradobrei-api
```

Для Git Bash в Windows:

```bash
./bradobrei-api.exe
```

## Полезные команды

Проверка форматирования:

```bash
gofmt -w ./cmd ./internal
```

Проверка сборки всего модуля:

```bash
go build ./...
```

## Частые проблемы

### Ошибка `schema for creating objects is not selected` / `SQLSTATE 3F000`

Обычно это значит одно из трёх:

- база пересоздана вручную и в ней нет рабочей схемы `public`
- у пользователя БД пустой `search_path`
- `.env` смотрит не в тот экземпляр PostgreSQL

### Ошибка `type "geometry" does not exist` / `SQLSTATE 42704`

Это почти всегда означает, что в текущей базе не включён `postgis`.

Что сделать:

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
```

Если команда не выполняется, значит у вас либо не `PostgreSQL + PostGIS`, либо у пользователя БД нет прав на создание расширений.

Что проверить:

```sql
SELECT current_database();
SHOW search_path;
SELECT schema_name FROM information_schema.schemata;
```

Что сделать:

```sql
CREATE SCHEMA IF NOT EXISTS public;
```

И убедиться, что приложение ходит в ту же БД, что вы смотрите в клиенте.

### Почему `go build` может падать, а `go run` идти дальше

На Windows иногда мешает доступ к кэшу Go в `%LOCALAPPDATA%\\go-build`.
Если такое повторится, попробуйте открыть терминал/IDE с обычными правами пользователя и очистить кэш:

```bash
go clean -cache
```
