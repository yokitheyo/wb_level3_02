# --- Stage 1: Build ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник
RUN go build -o shortener ./cmd/api

# --- Stage 2: Run ---
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и папки
COPY --from=builder /app/shortener .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations

# Порты
EXPOSE 8080

# Запуск
CMD ["./shortener"]
