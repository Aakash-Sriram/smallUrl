# 🚀 smallUrl — Minimal, Secure & Containerized URL Shortener

A compact, production-ready URL shortening service built with Go (Gin), Redis, and an optional nginx frontend. Designed for clarity, minimalism, and containerized deployment workflows.

---

## 📑 Table of Contents

1. Overview
2. Features
3. Architecture
4. How It Works
   - POST /shorten — flow
   - GET /:id — flow
5. API Reference
6. Configuration
7. Running the Project
   - Using Docker Compose (recommended)
   - Local Development (without Docker)
8. Building the Docker Image
9. Troubleshooting & Logs
10. Future Improvements
11. License

---

## 🧭 Overview

`smallUrl` is a lightweight URL-shortener service designed for simplicity, speed, and containerizable deployment. It exposes a minimal HTTP API:

- `POST /shorten` — create a short URL
- `GET /:id` — resolve a slug and redirect

Internally, a Go server handles routing and slug generation, while Redis maintains slug → long URL mappings. A static frontend can be served via nginx (optional), making the project suitable for local demos or cloud deployment.

## ✨ Features

- ⚡ Fast Go backend built on the Gin framework
- 🔐 Cryptographically secure slug generation
- 🗄️ Redis-backed storage for O(1) lookups
- ⏳ Configurable TTL (24h by default)
- 🐳 Full Docker Compose environment
- 🧱 Multi-stage Dockerfile for optimized builds
- 🧼 Clean API with simple request/response contracts

## 🏗️ Architecture

```text
               ┌───────────────┐
               │    Frontend   │  (optional: served by nginx)
               └───────┬───────┘
                       │  HTTP
               ┌───────▼────────┐
               │     nginx       │  ← terminates HTTP, proxies /api → app
               └───────┬────────┘
                       │ internal Docker network
               ┌───────▼────────┐
               │   Go API App    │  ← slug creation, redirect logic
               └───────┬────────┘
                       │ redis:6379
               ┌───────▼────────┐
               │     Redis       │  ← slug → long URL store
               └─────────────────┘
```

## 🔄 How It Works

### POST /shorten — flow

Client sends JSON (example):

```json
{ "url": "https://example.com" }
```

Steps:

1. Server validates the payload.
2. Generates an 8-character slug via `crypto/rand`.
3. Stores the mapping in Redis with a 24h TTL.
4. Responds with JSON containing the shortened URL.

##  License

Apache 2.0. Feel free to fork, modify, deploy, and build on top of smallUrl.
```

### GET /:id — flow

1. Client requests a slug, e.g., `/aB93xYpQ`.
2. Server queries Redis:
   - If key found → returns a `302` redirect.
   - If key missing → returns `404` JSON.
3. Browser follows redirect to the original URL.

## 📘 API Reference

### POST /shorten

**Request**

```http
POST /shorten
Content-Type: application/json

Body:
{
  "url": "https://example.com"
}
```

**Response (200)**

```json
{ "shortUrl": "http://localhost/<slug>" }
```

**Errors**

| Status | Meaning |
|--------|---------|
| 400    | Invalid body |
| 500    | Failed to generate slug or Redis error |

### GET /:id

Redirects to the original URL.

**Example**

```http
GET /Af92pLmQ  → 302 → https://example.com
```

**Errors**

| Status | Meaning |
|--------|---------|
| 404    | Slug not found |
| 500    | Redis failure |

## ⚙️ Configuration

The server reads the following environment variables:

| Variable   | Description                    | Default            |
|------------|--------------------------------|--------------------|
| BASE_URL   | Prefix for returned shortUrl   | `http://localhost` |
| PORT       | HTTP port for Go app           | `9808`             |
| REDIS_ADDR | Redis address                  | `localhost:6379`   |

**Example `.env`:**

```bash
BASE_URL=https://short.example.com
PORT=9808
REDIS_ADDR=redis:6379
```

## 🏃‍♂️ Running the Project

### 🚀 Using Docker Compose (recommended)

From project root:

```bash
docker compose up --build
```

This launches:

- Redis
- The Go app
- Nginx (serving frontend + proxying /api)

Visit:

- Frontend: `http://localhost`
- API via proxy: `POST http://localhost/api/shorten`
- Redirect: `http://localhost/<slug>`

Stop environment:

```bash
docker compose down
```

### 🧪 Local Development (without Docker)

Start Redis:

```bash
redis-server
```

Run the Go backend:

```bash
go run .
```

Call (example):

```http
POST http://localhost:9808/shorten
```

## 🐳 Building the Docker Image

The project includes an optimized multi-stage Dockerfile:

```bash
docker build -t smallurl:latest .
```

Producing:

- Large build stage (Go compiler)
- Minimal runtime stage (Alpine + your binary)

## 🔍 Troubleshooting & Logs

View logs:

```bash
docker compose logs -f app
docker compose logs -f redis
docker compose logs -f nginx
```

### Common Issues

| Issue                          | Fix |
|--------------------------------|-----|
| Backend cannot reach Redis     | Ensure `REDIS_ADDR=redis:6379` inside Compose |
| nginx not serving frontend     | Run `npm run build` to populate `frontend/dist` |
| 404 on slug                    | TTL expired or slug never created |
| CORS issues                    | Use nginx reverse proxy (`/api/ → app:9808`) |

## 🧭 Future Improvements

- Persistent storage instead of TTL-only Redis
- Rate limiting & abuse prevention
- Custom slug support
- Analytics (click counts, IPs, etc.)
- Full test suite
- CI/CD pipelines

## 📄 License

Apache 2.0. Feel free to fork, modify, deploy, and build on top of smallUrl.
