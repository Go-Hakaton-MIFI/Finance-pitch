# Dockerfile

# Stage 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Устанавливаем зависимости для сборки
RUN apk update && \
    apk add --no-cache git ca-certificates

# Копируем файлы проекта
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY internal/ internal/

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app

# Устанавливаем необходимые утилиты
RUN apk update && \
    apk add --no-cache postgresql-client dos2unix ca-certificates

# Копируем бинарники и файлы
COPY --from=builder /app/main .
COPY --from=builder /app/internal/migrations ./migrations
COPY --from=builder /app/internal/config ./config

# Копируем и настраиваем entrypoint
COPY docker/entrypoint.sh /app/entrypoint.sh
COPY docker/goose_linux_x86_64 /usr/local/bin/goose

# Устанавливаем права на выполнение и конвертируем CRLF в LF
RUN chmod +x /usr/local/bin/goose && \
    chmod +x /app/entrypoint.sh && \
    dos2unix /app/entrypoint.sh && \
    sed -i 's/\r$//' /app/entrypoint.sh

ENTRYPOINT ["/bin/sh", "/app/entrypoint.sh"]
