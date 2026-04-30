# Mini Edge API Platform

A lightweight API built in Go that simulates an edge/gateway layer with in-memory caching, external API integration, and performance optimization.

The project focuses on reducing latency, improving request efficiency, and simulating real-world API gateway behavior.

---

## Overview

This project simulates an intermediate layer (API Gateway / Edge Server) responsible for:

- Reducing external API calls
- Improving response time via caching
- Handling internal and external data sources
- Demonstrating concurrency-safe in-memory cache design

---

## Technologies

- Go (Golang)
- Gin Web Framework
- net/http
- sync.RWMutex (concurrency control)
- Go testing & benchmarking tools

---

## ⚙️ Architecture

Request flow:

Client
↓

API Handler
↓

Cache Layer
↓ (MISS)

External API / Mock Data
↓

Cache Store
↓

Client Response

---

## Endpoints

### GET /health
Health check endpoint to verify service availability.

### GET /products
Returns a list of products (mock data).

Features:
- In-memory cache
- TTL-based expiration
- Reduces processing time on repeated requests

### GET /external-data
Fetches data from an external API (JSONPlaceholder) acting as a proxy.

Features:
- HTTP client injection (Dependency Injection)
- Caching layer to reduce external calls
- Returns response directly to the client

---
## Cache System

The project uses a custom in-memory cache with TTL support.

### Features:

- Thread-safe (`sync.RWMutex`)
- Key-value storage
- TTL expiration per entry
- Lazy eviction on access
- Background cleanup goroutine (optional)

### Cache behavior:

- `HIT` → returns cached response
- `MISS` → fetches data and stores it in cache

---

## HTTP Headers

### X-Cache

Indicates response origin:

- `X-Cache: HIT` → served from cache
- `X-Cache: MISS` → fetched from source

### Cache-Control

```http
Cache-Control: public, max-age=30
```
Allows client-side caching for 30 seconds.

---

## Benchmarks

Cache performance was measured using Go benchmarks:

| Benchmark                  | ns/op (avg) | B/op | allocs/op |
|--------------------------|------------|------|------------|
| Cache Get (Hit)          | ~23 ns/op  | 0 B  | 0 allocs   |
| Cache Get (Miss)         | ~17 ns/op  | 0 B  | 0 allocs   |
| Cache Set                | ~493 ns/op | 336 B| 3 allocs   |

### Performance Analysis

- **GET operations (Hit/Miss)** are extremely fast due to direct in-memory map access (O(1)), with minimal overhead and zero allocations.
- **SET operations** are more expensive because they involve memory allocations and map updates, which introduces additional latency and garbage collector pressure.
- The cache is optimized for **read-heavy workloads**, which is typical in edge/gateway architectures.

Overall, the results demonstrate a system designed for low-latency reads, where write cost is accepted as a trade-off for fast retrieval and simplicity of design.

---

### Project Structure
cmd/server
internal/api
internal/cache
internal/handlers
internal/middleware
internal/routes

---

### How to run
```
go mod tidy
go run ./cmd/server
```

Server runs at:
```
http://localhost:8080
```

---

### Testing

Run unit tests:
```
go test ./internal/cache -v
```

Run benchmarks:
```
go test ./internal/cache -bench=. -benchmem
```
---

### Future improvements
- Cache based on query parameters
- Request timeout and retry strategy
- Structured logging
- Redis distributed cache
- Cache metrics (hit ratio, latency tracking)
- LRU eviction policy

---

## Purpose

This project was built for learning and demonstrating:

Backend architecture design
Caching strategies (edge-like behavior)
Concurrency in Go
Dependency injection patterns
Performance analysis with benchmarks

---

### Author

Fernanda Ishida
