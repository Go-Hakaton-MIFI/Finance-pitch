# Финансовый бэкенд

## Структура проекта

```
.
├── cmd/                    # Точка входа приложения
├── internal/              # Внутренний код приложения
│   ├── delivery/         # Слои доставки (HTTP, gRPC)
│   ├── domain/           # Бизнес-логика и модели
│   ├── infrastructure/   # Внешние сервисы и БД
│   └── migrations/       # Миграции базы данных
├── frontend/             # React фронтенд
├── config/              # Конфигурационные файлы
├── pkg/                 # Переиспользуемые пакеты
└── docker/              # Docker конфигурация
```

## Основные компоненты

### Бэкенд (Go)

#### Архитектура
- Чистая архитектура с разделением на слои:
  - Delivery (HTTP handlers)
  - Domain (бизнес-логика)
  - Infrastructure (репозитории, внешние сервисы)

#### База данных
- PostgreSQL
- Миграции с использованием goose
- Основные таблицы:
  - transactions
  - categories
  - transaction_statuses
  - prepared_transactions

#### API Endpoints
- `/api/v1/transactions` - управление транзакциями
- `/api/v1/categories` - управление категориями
- `/api/v1/trans_statuses` - управление статусами
- `/api/v1/analytics` - аналитика и статистика

### Фронтенд (React)

#### Основные компоненты
- Dashboard - главная страница с общей статистикой
- Transactions - управление транзакциями
- Analytics - аналитика и графики

#### Состояние приложения
- Управление состоянием через React Hooks
- API взаимодействие через axios
- Форматирование дат и чисел для RU локали

## Важные моменты для разработки

### Бэкенд

1. **Типы транзакций**
   - credit (доход)
   - debit (расход)

2. **Форматы данных**
   - Даты: ISO 8601 с часовым поясом
   - Суммы: decimal(15,5)
   - ИНН: 10 или 12 цифр
   - Телефон: 11 цифр

3. **Валидация**
   - Проверка типов пользователей (individual/legal)
   - Валидация ИНН и телефона
   - Проверка соответствия категорий типу транзакции

### Фронтенд

1. **Форматирование**
   - Даты: DD.MM.YYYY, HH:mm
   - Суммы: с разделителями тысяч и символом ₽
   - Пустые значения: отображаются как "-"

2. **Обработка ошибок**
   - Валидация форм на клиенте
   - Обработка API ошибок
   - Отображение уведомлений пользователю

3. **Фильтрация**
   - По датам
   - По типу транзакции
   - По категории
   - По статусу

## Разработка

### Запуск проекта
```bash
# Запуск в Docker
docker-compose -f docker/docker-compose-local.yml up -d

# Запуск фронтенда
cd frontend
npm install
npm start

# Запуск бэкенда
go run cmd/api/main.go
```

### Миграции
```bash
# Создание миграции
goose create migration_name sql

# Применение миграций
goose up
```

### Тестирование
```bash
# Запуск тестов
go test ./...

# Запуск линтера
golangci-lint run
```

## Безопасность

1. **Валидация данных**
   - Проверка всех входящих данных
   - Санитизация пользовательского ввода
   - Защита от SQL-инъекций

2. **Аутентификация**
   - JWT токены
   - Защищенные эндпоинты
   - Проверка прав доступа

## Мониторинг

1. **Логирование**
   - Уровни: info, error, debug
   - Структурированные логи
   - Трейсинг запросов

2. **Метрики**
   - Время ответа API
   - Количество запросов
   - Ошибки и исключения

## Дальнейшее развитие

1. **Планируемые улучшения**
   - Кэширование часто запрашиваемых данных
   - Оптимизация запросов к БД
   - Добавление новых типов аналитики

2. **Масштабирование**
   - Горизонтальное масштабирование
   - Балансировка нагрузки
   - Репликация БД

## Конфигурация

### Переменные окружения
```bash
# База данных
DB_ENGINE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_clean_template_db

# Сервер
APP_ADDRESS=0.0.0.0
APP_PORT=8089

# JWT
JWT_PUBLIC=your_public_key
JWT_PRIVATE_KEY=your_private_key

# S3/MinIO
S3_URL=http://minio:9000
S3_ROOT_USER=development_minio_key
S3_ROOT_PASSWORD=development_minio_secret
IMAGE_BUCKET_NAME=images

# Приложение
DEBUG=false
APP_NAME=finance-backend
APP_VERSION=1.0.0
```

### Docker
- Используется многоэтапная сборка
- Основные сервисы:
  - app: основное приложение
  - db: PostgreSQL
  - minio: S3-совместимое хранилище
  - test_gotenberg: сервис для генерации PDF

### Hot Reload
- Используется Air для hot reload в режиме разработки
- Конфигурация в `.air.toml`

# TODO:
- logger entry
- oapicodegen
- sql queries consts
- validate is admin depends on URL (contains admin)

