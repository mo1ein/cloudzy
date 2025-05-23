# 🌤️ Cloudzy – Scalable Microservices Weather & Pricing Platform

Cloudzy is a modern, modular, and production-grade microservices platform built in Go. Designed with real-world scalability and clean architecture in mind, it provides weather and pricing APIs via a gRPC-based gateway. It serves as both a robust starting point for distributed systems and a showcase of Go’s power in microservice orchestration.

---

## 🚀 Features

✅ **Clean Architecture** – Separates core domain, service logic, and delivery for maintainability and testability
✅ **gRPC APIs** – High-performance communication between services using gRPC and Protobuf
✅ **Modular Design** – Each service (weather, pricing, gateway) is independently runnable and deployable
✅ **Retry Interceptors** – gRPC retry logic with exponential backoff
✅ **Configurable & Extensible** – Easily adapt new services or extend current APIs
✅ **Makefile Driven** – Streamlined build, lint, proto-generation, and dev workflow
✅ **Dockerized** – Local development and deployment ready using `docker-compose`
✅ **Centralized Proto Definitions** – Shared and versioned API definitions under `/proto`

---

## 📁 Project Structure

```
cloudzy/
├── gateway/        # Entrypoint gateway service (API aggregator)
├── pricing/        # Pricing microservice
├── weather/        # Weather microservice
├── proto/          # Shared Protobuf definitions + generated code
├── Makefile        # CLI build and proto tools
└── docker-compose.yml
```

## 🛠️ Getting Started

### 📦 Prerequisites

- Docker & Docker Compose
- `protoc` and `protoc-gen-go` / `protoc-gen-go-grpc`

### 🔧 Build & Run Locally

```bash
# Generate gRPC code from proto
make proto

# Start all services with Docker Compose
docker-compose up --build
```

## 🧰 Development Tips

- Use `make generate` to regenerate code after modifying `.proto` files.
- Use `make run-gateway`, `make run-weather`, `make run-pricing` to run services individually.

---

## 📌 Why This Project?

This project was designed as a **technical demonstration of scalable Go microservice design**. It highlights the use of:

- Clear service boundaries
- Shared contracts via protobuf
- Production-grade reliability features (retries, logging)
- Developer-first workflows using Makefile and Docker

---
