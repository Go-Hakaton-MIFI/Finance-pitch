 # Документация API для фронтенд-разработчиков

## Базовый URL
```
http://localhost:8089/api/v1
```

## Аутентификация
Все защищенные эндпоинты требуют JWT токен в заголовке Authorization:
```
Authorization: Bearer <token>
```

## Публичные эндпоинты

### Регистрация
```
POST /registration
Content-Type: application/json

{
    "login": string,
    "password": string,
    "user_type": string
}
```

### Вход
```
POST /login
Content-Type: application/json

{
    "login": string,
    "password": string
}
```

### Получение типов пользователей
```
GET /subject_types
```

## Защищенные эндпоинты

### Транзакции

#### Получение списка транзакций
```
GET /transactions
```

#### Создание транзакции
```
POST /transactions
Content-Type: application/json

{
    // Данные транзакции
}
```

#### Подготовка транзакции
```
POST /transactions/prepared
Content-Type: application/json

{
    // Данные для подготовки
}
```

#### Получение транзакции по ID
```
GET /transactions/{id}
```

#### Удаление транзакции
```
DELETE /transactions/{id}
```

### Категории

#### Получение всех категорий
```
GET /categories
```

### Статусы

#### Получение всех статусов
```
GET /trans_statuses
```

### Аналитика

#### Динамика по периоду
```
POST /analytics/dynamics/by-period?period=<period>
Content-Type: application/json

{
    // Параметры фильтрации
}
```

#### Динамика по типу
```
POST /analytics/dynamics/by-type?trans_type=<type>
Content-Type: application/json

{
    // Параметры фильтрации
}
```

#### Сравнение доходов и расходов
```
POST /analytics/compare-income-expense
Content-Type: application/json

{
    // Параметры фильтрации
}
```

#### Сводка по статусам
```
POST /analytics/status-summary
Content-Type: application/json

{
    // Параметры фильтрации
}
```

#### Сводка по банкам
```
POST /analytics/banks-summary
Content-Type: application/json

{
    // Параметры фильтрации
}
```

#### Сводка по категориям
```
POST /analytics/categories-summary
Content-Type: application/json

{
    // Параметры фильтрации
}
```

#### Генерация PDF отчета по банкам
```
POST /analytics/banks-summary/report
Content-Type: application/json

{
    // Параметры фильтрации
}
```

#### Скачивание отчета
```
GET /analytics/banks-summary/report/{rep_id}
```

## Запуск бэкенда

1. Убедитесь, что у вас установлен Docker и Docker Compose
2. Создайте файл `.env` в корневой директории проекта со следующими параметрами:
```
DB_ENGINE=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_clean_template_db
DB_PORT=5432
DB_HOST=postgres-container-go

S3_ROOT_USER=development_minio_key
S3_ROOT_PASSWORD=development_minio_secret
S3_URL=http://minio:9000

IMAGE_BUCKET_NAME=images

GOTENBERG_API_URL=http://test_gotenberg:3000
GOTENBERG_PDF_CONVERTER_URL=/forms/chromium/convert/html

APP_ADDRESS=0.0.0.0
APP_PORT=8089

DEBUG=true

JWT_PRIVATE_KEY=<ваш_приватный_ключ>
JWT_PUBLIC=<ваш_публичный_ключ>
```

3. Запустите проект с помощью Docker Compose:
```bash
docker-compose up -d
```

## Миграции базы данных

Миграции выполняются автоматически при запуске контейнера с приложением. Если вам нужно выполнить миграции вручную:

1. Подключитесь к контейнеру с базой данных:
```bash
docker exec -it postgres-container-go psql -U postgres -d go_clean_template_db
```

2. Выполните SQL-скрипты из директории `migrations/`

## CORS

Бэкенд настроен на работу с CORS и принимает запросы с любого источника (`Access-Control-Allow-Origin: *`). Поддерживаются следующие методы:
- GET
- POST
- PUT
- DELETE
- OPTIONS

Поддерживаемые заголовки:
- Content-Type
- Authorization

Основные пути API (префикс /api/v1):
Публичные маршруты:
POST /api/v1/registration — регистрация пользователя
POST /api/v1/login — вход пользователя
GET /api/v1/subject_types — типы пользователей
Защищённые маршруты (требуется JWT):
GET /api/v1/transactions — получить список транзакций
POST /api/v1/transactions — создать транзакцию
POST /api/v1/transactions/prepared — подготовить транзакцию
GET /api/v1/transactions/{id} — получить транзакцию по id
DELETE /api/v1/transactions/{id} — удалить транзакцию
GET /api/v1/categories — получить все категории
GET /api/v1/trans_statuses — получить все статусы транзакций
Аналитика:
POST /api/v1/analytics/dynamics/by-period — динамика по периоду
POST /api/v1/analytics/dynamics/by-type — динамика по типу
POST /api/v1/analytics/compare-income-expense — сравнение доходов/расходов
POST /api/v1/analytics/status-summary — сводка по статусам
POST /api/v1/analytics/banks-summary — сводка по банкам
POST /api/v1/analytics/categories-summary — сводка по категориям
POST /api/v1/analytics/banks-summary/report — сгенерировать PDF-отчёт по банкам
GET /api/v1/analytics/banks-summary/report/{rep_id} — скачать PDF-отчёт

