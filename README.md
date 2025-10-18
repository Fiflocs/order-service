# Order Service

Микросервис для обработки заказов с использованием NATS Streaming и PostgreSQL.

## Features
- ✅ NATS Streaming subscription
- ✅ In-memory cache
- ✅ HTTP REST API
- ✅ Order processing pipeline
- ⏳ PostgreSQL integration (WIP)

## Quick Start
1. `docker-compose up -d`
2. `go run cmd/server/main.go`
3. `go run cmd/publisher/main.go`

## API
- `GET /orders/{id}` - Get order by ID
- `GET /orders` - Get all orders