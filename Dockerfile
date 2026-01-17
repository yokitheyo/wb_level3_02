# --- Stage 1: Build ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o shortener ./cmd/api

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# --- Stage 2: Run ---
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache bash ca-certificates

COPY --from=builder /app/shortener ./
COPY --from=builder /app/config.yaml ./
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/geoip/IP2LOCATION-LITE-DB3.CSV ./geoip/IP2LOCATION-LITE-DB3.CSV

COPY --from=builder /go/bin/goose /usr/local/bin/goose

EXPOSE 8080

CMD ["./shortener"]
