# 🔗 URL Shortener Service

A modern, production-ready URL shortening service built with Go, featuring clean architecture principles, comprehensive analytics, and a beautiful web interface.

## ✨ Features

- **URL Shortening**: Create short links with optional custom codes
- **Analytics Dashboard**: Real-time click tracking and device statistics
- **Geo-Location Tracking**: Identify click locations using GeoIP
- **Device Detection**: Automatically detect device types (mobile, tablet, desktop)
- **Caching Layer**: Redis-based caching for optimized performance
- **Expiration Support**: Set expiration times for short links
- **Responsive UI**: Beautiful web interface for easy link management
- **Comprehensive Logging**: Structured logging with ZeroLog
- **Error Handling**: Robust error handling with graceful degradation

## 🐳 Quick Start with Docker

### Prerequisites
- Docker & Docker Compose installed
- Port 8080 available (or modify docker-compose.yml)

### Running the Application

```bash
# Clone the repository
git clone https://github.com/yokitheyo/URLShortener.git
cd URLShortener

# Start all services (PostgreSQL, Redis, API)
docker-compose up --build

# Application will be available at http://localhost:8080
```

### Services in Docker Compose

| Service | Port | Details |
|---------|------|---------|
| **API** | 8080 | Main application server |
| **PostgreSQL** | 5432 | Database (shortener_db) |
| **Redis** | 6379 | Cache layer |

### Environment Variables

The Docker Compose setup uses the following environment variables:

```yaml
DB_MASTER: postgres://shortener:shortener@db:5432/shortener?sslmode=disable
REDIS_ADDR: redis:6379
REDIS_PASSWORD: ""
REDIS_DB: 0
SERVER_ADDR: :8080
```

### Stopping the Application

```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## 📡 API Endpoints

### 1. **Web Interface**

```
GET /
```
Returns the main web interface (HTML).

### 2. **Shorten URL**

```
POST /shorten
Content-Type: application/json

{
  "url": "https://very-long-url.com/path/to/resource",
  "custom": "my-short",  // Optional: custom short code
  "expires": 3600        // Optional: expiration in seconds
}
```

**Response (200 OK):**
```json
{
  "short": "abc123",
  "expires": 1747632455
}
```

**Access shortened URL:**
```
GET /s/{short}
```
Redirects to original URL and records a click.

---

### 3. **Analytics**

#### Get Basic Analytics
```
GET /analytics/{short}
```

**Response (200 OK):**
```json
{
  "short": "abc123",
  "original": "https://very-long-url.com/path/to/resource",
  "created_at": 1747632000,
  "expires_at": 1747635600,
  "visit_count": 42
}
```

#### Get Detailed Analytics
```
GET /analytics/{short}/detailed?from=2026-01-01&to=2026-01-31
```

**Response (200 OK):**
```json
{
  "short": "abc123",
  "daily_clicks": {
    "2026-01-17": 15,
    "2026-01-18": 27
  },
  "device_stats": {
    "desktop": 32,
    "mobile": 10,
    "tablet": 0
  },
  "mobile_percentage": 23,
  "total_clicks": 42
}
```

#### Get Recent Clicks
```
GET /analytics/{short}/recent-clicks?limit=10
```

**Response (200 OK):**
```json
{
  "clicks": [
    {
      "occurred_at": 1747632455000,
      "ip": "192.168.1.1",
      "referrer": "google.com",
      "device": "desktop"
    },
    {
      "occurred_at": 1747632400000,
      "ip": "203.0.113.42",
      "referrer": "direct",
      "device": "mobile"
    }
  ],
  "total": 2
}
```

---

### Error Responses

All errors follow this format:

```json
{
  "error": "Error description"
}
```

**Common HTTP Status Codes:**
- `200 OK`: Success
- `400 Bad Request`: Invalid input
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

---

## 🖼️ Visual Demo

### Main Dashboard
![Dashboard](docs/screenshots/1.png)

### Analytics View
![Analytics](docs/screenshots/2.png)

---

### Local Development

```bash
# Install dependencies
go mod download

# Run migrations
goose -dir migrations postgres "your-db-connection-string" up

# Build the application
go build -o shortener ./cmd/api

# Run the server
./shortener
```
