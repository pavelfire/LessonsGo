# Book Shop API

REST API для книжного магазина с тремя уровнями доступа:
- **Анонимный пользователь**: просмотр и фильтрация книг.
- **Авторизованный пользователь**: все, что доступно анонимному + корзина и checkout.
- **Администратор**: CRUD для категорий и книг.

Проект реализован на Go (стандартный `net/http`) и поддерживает:
- хранение в PostgreSQL (основной режим),
- in-memory хранилище (для локальной разработки без БД).

---

## Возможности

- Регистрация и логин по `email + password`.
- JWT-аутентификация (`Bearer` токен).
- CRUD категорий (только админ).
- CRUD книг (только админ).
- Фильтрация книг по одной или нескольким категориям.
- Книги с нулевым остатком не выдаются в публичный список.
- Корзина для авторизованных пользователей.
- Покупка книг через `checkout` без платежных данных.
- Защита от гонок при последней копии книги.
- Автоосвобождение корзины через 30 минут.

---

## Технологии

- Go 1.25
- PostgreSQL 16
- JWT (`github.com/golang-jwt/jwt/v5`)
- Bcrypt (`golang.org/x/crypto/bcrypt`)
- Docker / Docker Compose
- Testify для тестов

---

## Структура проекта

```text
testbookapi/
├── cmd/server/main.go                 # входная точка приложения
├── internal/
│   ├── api/handler.go                 # HTTP роуты и middleware
│   ├── api/handler_test.go            # интеграционные API тесты
│   ├── auth/service.go                # регистрация, логин, JWT
│   ├── category/service.go            # бизнес-логика категорий
│   ├── book/service.go                # бизнес-логика книг
│   ├── cart/service.go                # бизнес-логика корзины/checkout
│   ├── domain/                        # модели, интерфейсы, ошибки
│   └── storage/
│       ├── memory/store.go            # in-memory хранилище
│       └── postgres/store.go          # PostgreSQL хранилище
├── migrations/001_init.sql            # SQL миграция схемы
├── docker-compose.yml
├── Dockerfile
└── Makefile
```

---

## Быстрый старт

### 1) Запуск в Docker (рекомендуется)

```bash
make dc
```

Поднимется:
- API на `http://localhost:8080`
- PostgreSQL на `localhost:5432`

### 2) Локальный запуск без Docker

```bash
make run
```

По умолчанию это режим in-memory (данные не сохраняются между рестартами), если не передан `DATABASE_URL`.

### 3) Локальный запуск с PostgreSQL

```bash
DATABASE_URL="postgres://bookshop:bookshop@localhost:5432/bookshop?sslmode=disable" \
JWT_SECRET="dev-secret-change-me" \
make run
```

---

## Переменные окружения и флаги

### ENV

- `DATABASE_URL` — строка подключения к PostgreSQL.
- `JWT_SECRET` — секрет для подписи JWT.

### CLI флаги

- `-addr` (по умолчанию `:8080`) — адрес HTTP сервера.
- `-database-url` — переопределяет `DATABASE_URL`.
- `-jwt-secret` — переопределяет `JWT_SECRET`.
- `-memory` — принудительно использовать in-memory хранилище.

---

## Make команды

- `make build` — сборка бинарника.
- `make run` — запуск API.
- `make dc` — запуск через Docker Compose.
- `make test` — запуск тестов.
- `make lint` — `go vet`.
- `make tidy` — `go mod tidy`.

---

## Аутентификация

Для защищенных маршрутов передавайте:

```http
Authorization: Bearer <token>
```

Токен выдается через:
- `POST /auth/register`
- `POST /auth/login`

---

## Роли и права

### Аноним
- `GET /books`
- `GET /books/{id}`
- `GET /categories`
- `GET /categories/{id}`

### Авторизованный пользователь
- все права анонима
- `GET /cart`
- `POST /cart/items`
- `DELETE /cart/items/{book_id}`
- `POST /cart/checkout`

### Администратор
- все права пользователя
- `POST /categories`
- `PUT /categories/{id}`
- `DELETE /categories/{id}`
- `POST /books`
- `PUT /books/{id}`
- `DELETE /books/{id}`

Пользователи становятся админами только через БД:

```sql
UPDATE users SET is_admin = true WHERE email = 'admin@example.com';
```

---

## API (REST)

### Auth

#### `POST /auth/register`
Регистрация пользователя.

Request:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response `201`:
```json
{
  "token": "<jwt>"
}
```

#### `POST /auth/login`
Логин пользователя.

Request:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response `200`:
```json
{
  "token": "<jwt>"
}
```

---

### Categories

#### `GET /categories`
Список категорий.

#### `GET /categories/{id}`
Одна категория по ID.

#### `POST /categories` (admin)
Создать категорию.

Request:
```json
{
  "name": "Fiction"
}
```

#### `PUT /categories/{id}` (admin)
Обновить категорию.

#### `DELETE /categories/{id}` (admin)
Удалить категорию (если нет связанных книг).

---

### Books

#### `GET /books`
Список доступных книг (только `stock > 0`).

Параметры:
- `category=1&category=2` или `categories=1,2`
- `limit` (по умолчанию 50, максимум 100)
- `offset` (по умолчанию 0)

Пример:
`GET /books?category=1&category=2&limit=20&offset=0`

#### `GET /books/{id}`
Получить книгу по ID (если в наличии).

#### `POST /books` (admin)
Создать книгу.

Request:
```json
{
  "title": "Domain-Driven Design",
  "year_published": 2003,
  "author": "Eric Evans",
  "price_usd": 49.99,
  "category_id": 1,
  "stock": 10
}
```

#### `PUT /books/{id}` (admin)
Обновить книгу.

Важно: `stock` изменять нельзя (валидационная ошибка).

#### `DELETE /books/{id}` (admin)
Удалить книгу.

---

### Cart

#### `GET /cart` (auth)
Получить текущую корзину пользователя.

#### `POST /cart/items` (auth)
Добавить книгу в корзину.

Request:
```json
{
  "book_id": 42
}
```

#### `DELETE /cart/items/{book_id}` (auth)
Удалить книгу из корзины.

#### `POST /cart/checkout` (auth)
Завершить покупку (без данных карты).

Response `200`:
```json
{
  "status": "completed"
}
```

---

## Логика стока и корзины

- При добавлении книги в корзину:
  - проверяется наличие,
  - `stock` уменьшается сразу (резерв),
  - создается запись в корзине с `expires_at = now + 30m`.

- Если пользователь удаляет позицию из корзины:
  - `stock` возвращается.

- Если пользователь делает checkout:
  - корзина очищается,
  - резерв считается проданным.

- Если позиция в корзине истекла (30 минут):
  - она автоматически освобождается,
  - `stock` возвращается.

Фоновая очистка запускается в `main` и выполняется раз в минуту.

---

## Ошибки и коды ответов

Базовые коды:
- `200` — успех.
- `201` — создано.
- `204` — удалено.
- `400` — ошибка валидации / бизнес-ограничения.
- `401` — неавторизован.
- `403` — недостаточно прав.
- `404` — сущность не найдена.
- `500` — внутренняя ошибка.

Формат ошибки:

```json
{
  "error": "human readable message"
}
```

---

## Примеры curl

### Регистрация
```bash
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Логин
```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Создать категорию (admin)
```bash
curl -s -X POST http://localhost:8080/categories \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Programming"}'
```

### Создать книгу (admin)
```bash
curl -s -X POST http://localhost:8080/books \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title":"The Go Programming Language",
    "year_published":2015,
    "author":"Alan A. A. Donovan",
    "price_usd":39.99,
    "category_id":1,
    "stock":5
  }'
```

### Фильтрация книг по нескольким категориям
```bash
curl -s "http://localhost:8080/books?category=1&category=2&limit=20&offset=0"
```

### Добавить книгу в корзину
```bash
curl -s -X POST http://localhost:8080/cart/items \
  -H "Authorization: Bearer <USER_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"book_id":1}'
```

### Checkout
```bash
curl -s -X POST http://localhost:8080/cart/checkout \
  -H "Authorization: Bearer <USER_TOKEN>"
```

---

## Тестирование

Запуск:

```bash
make test
```

В тестах покрыты:
- регистрация/логин,
- доступы по ролям,
- CRUD операций,
- ограничения на stock,
- конкурентный сценарий на последнюю копию,
- сценарий с истечением корзины.

---

## Ограничения и заметки

- Количество книг в корзине: не более 1 копии каждой книги (по `book_id`).
- Вложенные категории не поддерживаются (иерархия плоская).
- Изменение `stock` через update книги запрещено по требованиям.
- В dev-режиме in-memory данные не персистентны.

