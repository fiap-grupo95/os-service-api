# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**OS Service API** is a Go microservice that orchestrates the full lifecycle of vehicle repair shop service orders (Ordens de Serviço). It acts as an integration hub connecting Entity Service (customers/vehicles/inventory), Billing Service (estimates/payments), Execution Service (repair tracking), and Notification Service.

## Commands

All commands run inside Docker containers via `docker-compose`.

```bash
# First-time setup (copies .env-example to .env, builds containers, waits for readiness)
make init

# Day-to-day development
make up           # Start containers (build + background)
make down         # Stop and remove containers
make logs         # Follow app container logs

# Testing (runs inside dev container)
make test         # Run all tests with verbose output
make coverage     # Run tests + coverage summary in terminal
make coverage-html    # Run tests + open HTML coverage report in browser
make coverage-core    # Coverage focused on usecase and handler layers only

# Swagger docs (regenerate after changing API comments)
make swag-generate

# Kubernetes (local Minikube)
make deploy-local # Deploy to Minikube
make run-jobs     # Run migration jobs
make local-api    # Port-forward service to localhost:8080
make get-all      # Check deployment status
make delete       # Remove all k8s resources
```

To run a single test package:
```bash
docker-compose exec dev go test ./internal/usecase/... -v -run TestFunctionName
```

Swagger UI is available at `http://localhost:8080/swagger/index.html` while the app is running.

## Architecture

The project follows **Hexagonal Architecture** (Ports & Adapters) with **DDD** principles, organized in `internal/`:

```
internal/
├── domain/          # Business entities, value objects, state machine enums
├── usecase/         # Application orchestration; interfaces/ defines all dependency contracts
├── adapter/
│   ├── http/        # Gin handlers, routes, middleware, DTOs, gateway HTTP clients
│   └── persistence/ # MongoDB repositories
└── infrastructure/  # MongoDB connection, New Relic init, k8s manifests
```

**Dependency flow:** `adapter/http → usecase → domain`, with `adapter/persistence` and `adapter/http/gateway` implementing interfaces defined in `usecase/interfaces/`.

### Key Patterns

- **State Machine**: ServiceOrder and AdditionalRepair entities follow strict status transitions. Current statuses for ServiceOrder: `RECEBIDA → EM_DIAGNOSTICO → AGUARDANDO_APROVACAO → APROVADO → EM_EXECUCAO → FINALIZADA → ENTREGUE`. Status logic lives in the domain layer.

- **Gateway Pattern**: External service calls (Entity API, Billing Service, Execution Service) are abstracted through interfaces in `usecase/interfaces/` and implemented in `adapter/http/gateway/`. This enables mock injection in tests.

- **Repository Pattern**: MongoDB repositories in `adapter/persistence/mongodb/` implement interfaces defined in the usecase layer.

- **DTOs**: Request/response objects are defined in `adapter/http/dto/` and kept separate from domain entities.

### MongoDB Configuration

- Read preference `SecondaryPreferred` for list queries; `Primary` for critical operations (approvals/rejections)
- Connection pool: min 10, max 50 connections
- Collections: `service_orders`, `additional_repairs`, `users`

## Environment

Copy `.env-example` to `.env` (done automatically by `make init`). Key variables:

```
MONGODB_URI=mongodb://mongodb:27017/os-service-db
MONGODB_DATABASE=os-service-db
GIN_MODE=debug
JWT_SECRET=chave_muito_segura
JWT_TTL=24h
NEW_RELIC_LICENSE_KEY=<key>
```

The `mockoon` container (port 8083) serves mock responses for external service integrations during local development. Mock definitions live in `mockoon/`.

## Testing

Mocks are generated with `go.uber.org/mock`. Mock files follow the naming pattern `mock_*.go` alongside the interfaces they mock. Tests use `testify` for assertions. The CI pipeline runs the full test suite on every push via `.github/workflows/go.yml`.

## CI/CD

- **`go.yml`**: Runs on every push/PR — builds, generates Swagger, runs tests
- **`ci-cd-pipeline.yml`**: Triggered on merge to `main` — builds Docker image, pushes to `mandaapag03/os-service-api` on Docker Hub, deploys to AWS EKS (us-east-1)

Docker image tag convention for test builds: `os-service-api:<version>-[<description>-]test-<number>`
