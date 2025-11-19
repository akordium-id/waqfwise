# WaqfWise

Modern Wakaf (Islamic Endowment) Management Platform

## Overview

WaqfWise is a comprehensive platform for managing Islamic endowments (wakaf), designed to serve nazirs (waqf managers), donors, and institutions.

## Editions

### Community Edition (AGPL v3)
Open-source core features including:
- Authentication & Authorization
- Campaign Management
- Basic Payment Processing
- Donor Management
- Nazir Management
- Basic Reporting

### Enterprise Edition (Commercial License)
Advanced features for large institutions:
- Multi-tenancy Architecture
- Advanced Analytics Engine
- White-labeling
- Third-party Integrations (SAP, etc.)
- Payment Aggregation & Fraud Detection
- Geospatial Asset Tracking

## Technology Stack

### Backend Infrastructure
- **Framework**: Gin (high-performance HTTP router)
- **Database**: PostgreSQL 16 with GORM ORM
- **Migrations**: golang-migrate
- **Cache**: Redis 7 (session management, rate limiting, caching)
- **Message Queue**: Apache Kafka
- **Authentication**: JWT + OAuth2 (Google, Facebook)
- **Payment Gateways**: Midtrans + Xendit
- **Monitoring**: Prometheus + Grafana
- **Logging**: Zap (structured logging)

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Docker & Docker Compose (recommended)
- PostgreSQL 16+ (if not using Docker)
- Redis 7+ (if not using Docker)
- Apache Kafka (if not using Docker)
- Make (for build automation)

### Quick Start with Docker

1. **Clone the repository**
```bash
git clone https://github.com/akordium-id/waqfwise.git
cd waqfwise
```

2. **Copy environment file**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Start infrastructure services**
```bash
make docker-up
```

4. **Run database migrations**
```bash
make migrate-up
```

5. **Build and run**

**Community Edition:**
```bash
make run-community
```

**Enterprise Edition:**
```bash
make run-enterprise
```

The server will start on http://localhost:8080

### Building

**Community Edition:**
```bash
make build-community
# Binary will be in bin/waqfwise-community
```

**Enterprise Edition:**
```bash
make build-enterprise
# Binary will be in bin/waqfwise-enterprise
```

**Build both:**
```bash
make build
```

## Development

### Available Make Commands

```bash
make help              # Show all available commands
make deps              # Download Go dependencies
make build             # Build both editions
make test              # Run tests
make test-coverage     # Run tests with coverage report
make lint              # Run linters
make fmt               # Format code
make docker-up         # Start Docker services
make docker-down       # Stop Docker services
make migrate-up        # Run database migrations
make migrate-down      # Rollback migrations
make migrate-create    # Create new migration
make dev-setup         # Complete development setup
make clean             # Clean build artifacts
```

### Project Structure

```
waqfwise/
├── cmd/                          # Application entry points
│   ├── waqfwise-community/       # Community edition
│   └── waqfwise-enterprise/      # Enterprise edition
├── config/                       # Configuration files
├── deployments/                  # Deployment configs (Docker, K8s)
│   ├── prometheus/
│   └── grafana/
├── internal/                     # Private application code
│   ├── core/                     # Core domain logic (Community)
│   └── enterprise/               # Enterprise features
├── migrations/                   # Database migrations
│   ├── community/
│   └── enterprise/
├── pkg/                          # Public libraries
│   ├── auth/                     # JWT & OAuth2
│   ├── cache/                    # Redis client & utilities
│   ├── config/                   # Configuration management
│   ├── database/                 # Database connection & GORM
│   ├── logger/                   # Zap logger
│   ├── metrics/                  # Prometheus metrics
│   ├── middleware/               # HTTP middleware
│   ├── payment/                  # Payment gateway integrations
│   ├── queue/                    # Kafka producer/consumer
│   └── server/                   # HTTP server
├── scripts/                      # Build and utility scripts
├── docker-compose.yml            # Docker services
├── Makefile                      # Build automation
└── .env.example                  # Environment variables template
```

### API Endpoints

#### Health & Metrics
- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /metrics` - Prometheus metrics

#### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/google` - Google OAuth2
- `GET /api/v1/auth/google/callback` - Google callback
- `GET /api/v1/auth/facebook` - Facebook OAuth2
- `GET /api/v1/auth/facebook/callback` - Facebook callback

#### Protected Routes (requires JWT)
- `GET /api/v1/profile` - Get user profile
- `POST /api/v1/logout` - Logout

### Configuration

The application can be configured using:
1. YAML configuration files (`config/config.community.yaml` or `config/config.enterprise.yaml`)
2. Environment variables (see `.env.example`)
3. Command-line flags

Environment variables take precedence over YAML configuration.

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./pkg/auth/...
```

### Monitoring

Access monitoring dashboards:
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)

### Docker Services

The `docker-compose.yml` provides:
- PostgreSQL 16 (port 5432)
- Redis 7 (port 6379)
- Kafka + Zookeeper (ports 9092, 29092)
- Prometheus (port 9090)
- Grafana (port 3000)

## Documentation

See the `docs/` directory for detailed documentation.

## Contributing

We welcome contributions to the Community Edition! Please read `docs/contributing.md` for guidelines.

## License

- Community Edition: AGPL v3 (see LICENSE-AGPL-v3.txt)
- Enterprise Edition: Commercial License (see LICENSE-COMMERCIAL.txt)
