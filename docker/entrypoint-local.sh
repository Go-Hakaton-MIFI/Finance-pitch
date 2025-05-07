#!/bin/sh
set -e  # Остановка при ошибке

echo "Waiting 2 secodns"
sleep 2

echo "Applying database migrations..."

goose -dir /migrations postgres "host=$DB_HOST port=$DB_PORT user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" up

echo "Starting application..."
exec air -c /app/.air.toml
