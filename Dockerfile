# --- Stage 1: Build ---
FROM golang:1.23-alpine AS builder

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

# Устанавливаем необходимые пакеты
RUN apk add --no-cache wget unzip ca-certificates

# Копируем бинарник, конфиг и папки
COPY --from=builder /app/shortener .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY ./migrations ./migrations

COPY ./internal/geoip/IP2LOCATION-LITE-DB3.CSV /app/geoip/IP2LOCATION-LITE-DB3.CSV


# Порты
EXPOSE 8080

# Запуск
CMD ["./shortener"]
