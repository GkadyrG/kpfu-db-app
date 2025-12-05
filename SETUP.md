# Инструкция по запуску проекта

## Быстрый старт с Docker Compose

Это самый простой способ запустить проект. Требуется только Docker и Docker Compose.

### Шаг 1: Проверка установки Docker

```bash
docker --version
docker-compose --version
```

Если Docker не установлен, скачайте его с https://www.docker.com/products/docker-desktop

### Шаг 2: Запуск проекта

```bash
cd d:\kfu\my-kpfu-db-app
docker-compose up --build
```

### Шаг 3: Открыть приложение

Откройте браузер и перейдите по адресу: **http://localhost:8080**

### Остановка проекта

Нажмите `Ctrl+C` в терминале или выполните:

```bash
docker-compose down
```

Для полной очистки (включая данные БД):

```bash
docker-compose down -v
```

---

## Локальная разработка (без Docker)

Требуется Go 1.23+ и PostgreSQL.

### Шаг 1: Установка PostgreSQL

Установите PostgreSQL 17 или запустите в Docker:

```bash
docker run -d \
  --name postgres-shipments \
  -e POSTGRES_USER=shipment_user \
  -e POSTGRES_PASSWORD=shipment_pass \
  -e POSTGRES_DB=shipments_db \
  -p 5432:5432 \
  postgres:17
```

### Шаг 2: Инициализация БД

```bash
psql -h localhost -U shipment_user -d shipments_db -f 0_init.sql
```

Или через Docker:

```bash
docker exec -i postgres-shipments psql -U shipment_user -d shipments_db < 0_init.sql
```

### Шаг 3: Настройка переменных окружения

Создайте файл `.env` или экспортируйте переменную:

```bash
export DB_URL="postgres://shipment_user:shipment_pass@localhost:5432/shipments_db?sslmode=disable"
```

Windows PowerShell:

```powershell
$env:DB_URL="postgres://shipment_user:shipment_pass@localhost:5432/shipments_db?sslmode=disable"
```

### Шаг 4: Установка зависимостей

```bash
go mod download
```

### Шаг 5: Запуск приложения

```bash
go run cmd/main.go
```

Или через Makefile:

```bash
make run
```

### Шаг 6: Открыть приложение

Откройте браузер: **http://localhost:8080**

---

## Структура страниц

### Главная страница (/)
- Отображение всех таблиц
- CRUD операции
- Результат хранимой процедуры

### VIEW (/view)
- Объединенные данные из трех таблиц

### Динамическое отображение (/dynamic)
- Выбор таблицы из выпадающего списка
- Динамическое отображение данных

### Задача 1 (/task-1)
- Параметр: город (по умолчанию "Казань")
- Два метода: SQL и ORM

### Задача 2 (/task-2)
- Отгрузки текущего года
- Оконные функции для расчета доли

### Задача 3 (/task-3)
- Два метода: кванторный SQL и record-based

---

## API Endpoints

### CRUD Детали
- `GET /` - главная страница (включает все данные)
- `POST /api/parts` - создать
- `PUT /api/parts/:code` - обновить
- `DELETE /api/parts/:code` - удалить

### CRUD Покупатели
- `POST /api/customers` - создать
- `PUT /api/customers/:id` - обновить
- `DELETE /api/customers/:id` - удалить

### CRUD Отгрузки
- `POST /api/shipments` - создать
- `PUT /api/shipments/:warehouse/:doc` - обновить
- `DELETE /api/shipments/:warehouse/:doc` - удалить

### Задачи
- `GET /api/task-1/sql?city=Казань`
- `GET /api/task-1/orm?city=Казань`
- `GET /api/task-2`
- `GET /api/task-3/sql`
- `GET /api/task-3/record`

### Другое
- `GET /api/table/:name` - динамические данные (parts/customers/shipments)
- `GET /api/procedure/:customer_id` - хранимая процедура

---

## Тестирование API

### Примеры с curl

```bash
# Получить данные задачи 1
curl "http://localhost:8080/api/task-1/sql?city=Казань"

# Создать новую деталь
curl -X POST http://localhost:8080/api/parts \
  -H "Content-Type: application/json" \
  -d '{"part_code":"D999","part_type":"покупная","name":"Тест","unit":"шт","plan_price":50.00}'

# Удалить деталь
curl -X DELETE http://localhost:8080/api/parts/D999

# Вызвать хранимую процедуру
curl http://localhost:8080/api/procedure/1
```

---

## Проверка базы данных

### Подключение к PostgreSQL

```bash
# Через Docker
docker exec -it my-kpfu-db-app-db-1 psql -U shipment_user -d shipments_db

# Локально
psql -h localhost -U shipment_user -d shipments_db
```

### Полезные SQL запросы

```sql
-- Просмотр всех таблиц
\dt

-- Просмотр данных
SELECT * FROM parts;
SELECT * FROM customers;
SELECT * FROM shipments;
SELECT * FROM v_full_shipment_info;

-- Проверка триггеров
SELECT * FROM shipments_audit;

-- Вызов процедуры
CALL p_customer_shipment_summary(1, NULL, NULL);

-- Вызов функций
SELECT fn_customer_count_by_city('Казань');
SELECT * FROM fn_shipments_in_range('2025-01-01', '2025-12-31');
```

---

## Решение проблем

### Порт 5432 занят

Если PostgreSQL уже запущен локально, измените порт в `compose.yaml`:

```yaml
ports:
  - "5433:5432"  # Используем 5433 вместо 5432
```

И обновите `DB_URL` в environment:

```yaml
DB_URL: "postgres://shipment_user:shipment_pass@db:5432/shipments_db?sslmode=disable"
```

### Порт 8080 занят

Измените порт приложения в `compose.yaml`:

```yaml
ports:
  - "8081:8080"  # Используем 8081 вместо 8080
```

### База данных не инициализируется

Убедитесь, что файл `0_init.sql` находится в корне проекта.

Для переинициализации:

```bash
docker-compose down -v  # Удалить volumes
docker-compose up --build
```

### Ошибка подключения к БД

Подождите несколько секунд - PostgreSQL может инициализироваться.

Проверьте логи:

```bash
docker-compose logs db
```

---

## Разработка

### Изменение кода

После изменения Go кода:

```bash
docker-compose up --build
```

Или в режиме разработки:

```bash
go run cmd/main.go
```

### Изменение HTML шаблонов

В режиме разработки шаблоны перезагружаются автоматически при перезапуске.

В Docker нужно пересобрать:

```bash
docker-compose up --build
```

### Изменение SQL схемы

1. Обновите `0_init.sql`
2. Пересоздайте БД:

```bash
docker-compose down -v
docker-compose up --build
```

---

## Дополнительная информация

- **Документация PostgreSQL**: https://www.postgresql.org/docs/
- **Документация Gin**: https://gin-gonic.com/docs/
- **Документация pgx**: https://pkg.go.dev/github.com/jackc/pgx/v5
- **Bootstrap**: https://getbootstrap.com/docs/4.5/

---

## Контакты

Для вопросов и предложений создавайте issues в репозитории проекта.

