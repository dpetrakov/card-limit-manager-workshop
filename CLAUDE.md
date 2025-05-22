# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Card Limit Manager is a microservices-based system for automating corporate card limit increase requests and approvals. The project uses:

- **Backend**: Go 1.22 + Gin framework (CLM service)
- **Database**: PostgreSQL 16 with sql-migrate for migrations
- **Frontend**: React 18 + TypeScript + Vite (planned)
- **Deployment**: Docker Compose for local development

## Development Commands

### Local Environment Setup
```bash
# Start entire stack (DB + migrations + services)
docker compose up

# Tear down services and volumes
docker compose down -v
```

### Database Management
```bash
# Manual migration run (if needed)
./run_migrations.sh

# Database connection details:
# Host: localhost:5432
# User: postgres
# Password: postgres
# Database: postgres
```

### Go Service Development
```bash
# In services/clm/ directory:
go mod tidy
go build ./cmd/clm
go test ./...
go test -cover ./...

# Run service directly (requires DB)
go run ./cmd/clm/main.go
```

### Testing
- Unit test coverage target: ≥ 80%
- Run tests with: `go test -cover ./...`
- Test files should follow `*_test.go` pattern

## Architecture

### Core Services
- **CLM Service** (`services/clm/`): Main limit request management service
  - REST API with Gin framework
  - Handles CRUD operations for limit requests
  - Manages approval workflow state machine
  - Clean Architecture: `cmd/`, `internal/`, `pkg/`

### Database Schema
- See `database/migrations/001_initial_schema.sql` for current schema
- Key tables: `users`, `limit_requests`, `approval_steps`, `audit_log`
- Uses UUIDs as primary keys
- Supports multi-stage approval workflow

### Approval Workflow
1. Employee creates request → PENDING_TEAM_LEAD
2. Team Lead approval → PENDING_RISK_OFFICER (or PENDING_CFO if amount > 1M RUB)
3. Risk Officer/CFO approval → COMPLETED
4. Any rejection → REJECTED

## Key Files and Conventions

### Project Structure
```
docs/                   # Architecture, requirements, specs
├── architecture.md     # Core system design
├── requirements.md     # Functional requirements (R-1, R-2, etc.)
└── specs/             # Detailed specifications per service

services/clm/          # Go microservice
├── cmd/clm/           # Main application entry point
├── internal/          # Private application code
└── pkg/               # Public library code

database/migrations/   # SQL migration files
```

### Documentation Standards
- Architecture diagrams use Mermaid format
- Sequence diagrams use PlantUML
- API specifications in OpenAPI 3.1 JSON format
- Database schema documented in DBML format (`docs/database.dbml`)

### Code Standards
- Go: Follow Clean Architecture patterns
- Error handling: Use structured logging in JSON format
- Database: Use pgx driver for PostgreSQL
- No database connection pooling configuration needed for local development

## Local Development Workflow

1. Ensure Docker is running
2. Run `docker compose up` to start all services
3. Services available at:
   - CLM API: http://localhost:8080
   - PostgreSQL: localhost:5432
   - Adminer: http://localhost:8081 (if configured)

## Important Notes

- Always check AGENTS.md files in subdirectories for specific coding standards
- Database migrations run automatically via Docker Compose
- All logging should be in English, INFO level or higher
- JWT authentication integration planned but not yet implemented
- Frontend React application is planned but not yet implemented