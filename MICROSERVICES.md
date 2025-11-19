# WaqfWise Microservices Architecture

## 🏗️ Architecture Overview

WaqfWise menggunakan **microservices architecture** yang terdiri dari 6 layanan independen yang berkomunikasi melalui REST API, dengan API Gateway sebagai single entry point.

```
┌─────────────────────────────────────────────────────────────┐
│                       API Gateway (8000)                     │
│                    Single Entry Point                        │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
┌───────▼────────┐  ┌──────▼─────────┐  ┌──────▼─────────┐
│ Auth Service   │  │ Payment Service │  │Campaign Service│
│   Port 8001    │  │   Port 8002     │  │   Port 8003    │
└────────────────┘  └─────────────────┘  └────────────────┘
        │                   │                   │
┌───────▼────────┐  ┌──────▼─────────┐  ┌──────▼─────────┐
│ Asset Service  │  │Analytics Service│  │Integration Svc │
│   Port 8004    │  │   Port 8005     │  │   Port 8006    │
│                │  │  [Enterprise]   │  │  [Enterprise]  │
└────────────────┘  └─────────────────┘  └────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
┌───────▼────────┐  ┌──────▼─────────┐  ┌──────▼─────────┐
│   PostgreSQL   │  │     Redis       │  │   Prometheus   │
│   + PostGIS    │  │   + Sessions    │  │   + Grafana    │
└────────────────┘  └─────────────────┘  └────────────────┘
```

## 📦 Services Breakdown

### 1. **Auth Service** (Port 8001)
**Responsibilities:**
- User authentication & authorization
- JWT token generation & validation
- Role-based access control (RBAC)
- Multi-factor authentication (MFA) support
- Password management & security

**Key Features:**
- ✅ JWT with refresh tokens
- ✅ MFA with TOTP (Google Authenticator compatible)
- ✅ Password hashing with bcrypt
- ✅ Role-based permissions (Admin, Nazir, Donor, Auditor, Operator)
- ✅ Session management with Redis

**Endpoints:**
```
POST   /api/v1/auth/register          - User registration
POST   /api/v1/auth/login             - User login
POST   /api/v1/auth/refresh           - Refresh access token
GET    /api/v1/auth/profile           - Get user profile
POST   /api/v1/auth/change-password   - Change password
POST   /api/v1/auth/mfa/setup         - Setup MFA
POST   /api/v1/auth/mfa/enable        - Enable MFA
POST   /api/v1/auth/mfa/disable       - Disable MFA
```

---

### 2. **Payment Service** (Port 8002)
**Responsibilities:**
- Payment processing & gateway integration
- Transaction ledger (double-entry bookkeeping)
- Fraud detection & prevention
- Donation management
- Receipt generation

**Key Features:**
- ✅ Multi-gateway support (Midtrans, Xendit)
- ✅ Double-entry bookkeeping ledger
- ✅ Fraud detection with risk scoring
- ✅ Payment method abstraction (Credit Card, Bank Transfer, E-Wallet, QRIS, VA)
- ✅ Recurring donation support
- ✅ Payment logs & audit trail

**Endpoints:**
```
POST   /api/v1/donations              - Create donation
GET    /api/v1/donations/:id          - Get donation details
GET    /api/v1/donations/user/:userId - Get user donations
POST   /api/v1/payments/callback      - Payment gateway callback
GET    /api/v1/ledger/campaign/:id    - Get campaign ledger
```

---

### 3. **Campaign Service** (Port 8003)
**Responsibilities:**
- Wakaf campaign CRUD operations
- Campaign search & filtering
- Goal tracking & milestones
- Campaign analytics
- Featured & urgent campaigns

**Key Features:**
- ✅ Campaign lifecycle management
- ✅ Multiple campaign types (Land, Building, Cash, Education, Healthcare, General)
- ✅ Goal tracking with progress calculation
- ✅ Milestone management
- ✅ Campaign status workflow (Draft → Active → Completed/Cancelled)
- ✅ Search with filters (type, status, location)

**Endpoints:**
```
POST   /api/v1/campaigns              - Create campaign
GET    /api/v1/campaigns              - List campaigns (with filters)
GET    /api/v1/campaigns/:id          - Get campaign details
PUT    /api/v1/campaigns/:id          - Update campaign
DELETE /api/v1/campaigns/:id          - Delete campaign
GET    /api/v1/campaigns/featured     - Get featured campaigns
POST   /api/v1/campaigns/:id/milestones - Add milestone
```

---

### 4. **Asset Service** (Port 8004)
**Responsibilities:**
- Wakaf asset management
- Geospatial data (PostGIS integration)
- Document management (certificates, legal docs)
- Asset valuation tracking
- Location-based queries

**Key Features:**
- ✅ Asset CRUD with geolocation
- ✅ PostGIS for spatial queries
- ✅ Document upload & management
- ✅ Valuation history timeline
- ✅ Multiple asset types (Land, Building, Vehicle, Equipment)
- ✅ GeoJSON boundary support

**Endpoints:**
```
POST   /api/v1/assets                 - Create asset
GET    /api/v1/assets                 - List assets
GET    /api/v1/assets/:id             - Get asset details
PUT    /api/v1/assets/:id             - Update asset
GET    /api/v1/assets/nearby          - Find assets nearby (PostGIS)
POST   /api/v1/assets/:id/documents   - Upload document
POST   /api/v1/assets/:id/valuation   - Add valuation
```

---

### 5. **Analytics Service** (Port 8005) 🔒 *Enterprise Only*
**Responsibilities:**
- Business intelligence & reporting
- Custom report builder
- Data export (PDF, Excel, CSV)
- Dashboard metrics
- Trend analysis

**Key Features:**
- ✅ Real-time analytics dashboard
- ✅ Custom report builder
- ✅ Export to multiple formats
- ✅ Donor insights & behavior
- ✅ Campaign performance metrics
- ✅ Financial reports

**Endpoints:**
```
GET    /api/v1/analytics/dashboard    - Dashboard metrics
GET    /api/v1/analytics/campaigns    - Campaign analytics
GET    /api/v1/analytics/donors       - Donor insights
POST   /api/v1/analytics/reports      - Generate custom report
GET    /api/v1/analytics/export       - Export data
```

---

### 6. **Integration Service** (Port 8006) 🔒 *Enterprise Only*
**Responsibilities:**
- Webhook management
- Third-party API integrations (Accounting, CRM)
- Data synchronization
- API connectors
- Event streaming

**Key Features:**
- ✅ Webhook CRUD & delivery
- ✅ Integration templates (Xero, QuickBooks, Salesforce)
- ✅ OAuth2 flow for integrations
- ✅ Retry mechanism with exponential backoff
- ✅ Event logging & audit

**Endpoints:**
```
POST   /api/v1/webhooks               - Create webhook
GET    /api/v1/webhooks               - List webhooks
PUT    /api/v1/webhooks/:id           - Update webhook
DELETE /api/v1/webhooks/:id           - Delete webhook
POST   /api/v1/integrations           - Connect integration
GET    /api/v1/integrations           - List integrations
POST   /api/v1/sync                   - Trigger data sync
```

---

### 7. **API Gateway** (Port 8000) 🌐 *Single Entry Point*
**Responsibilities:**
- Request routing to appropriate services
- Load balancing
- Rate limiting
- Request/response logging
- CORS handling

**Features:**
- ✅ Reverse proxy to all services
- ✅ Centralized logging
- ✅ Health check aggregation
- ✅ CORS middleware
- ✅ Service discovery

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+ with PostGIS
- Redis 7+

### 1. Clone Repository
```bash
git clone https://github.com/akordium-id/waqfwise.git
cd waqfwise
```

### 2. Install Dependencies
```bash
make -f Makefile.microservices deps
```

### 3. Start Infrastructure (Docker)
```bash
# Community Edition (Core Services)
make -f Makefile.microservices docker-up

# Enterprise Edition (All Services)
make -f Makefile.microservices docker-up-enterprise
```

### 4. Build Services Locally
```bash
# Build all services
make -f Makefile.microservices build-all

# Build specific service
make -f Makefile.microservices build SERVICE=auth-service
```

### 5. Run Service Locally
```bash
./bin/auth-service
```

---

## 🔧 Development

### Build & Run
```bash
# Build all services
make -f Makefile.microservices build-all

# Run specific service
make -f Makefile.microservices run SERVICE=auth-service
```

### Testing
```bash
# Run tests
make -f Makefile.microservices test

# Generate coverage report
make -f Makefile.microservices test-coverage
```

### Docker Development
```bash
# Build Docker images
make -f Makefile.microservices docker-build

# Start all services
make -f Makefile.microservices docker-up

# View logs
make -f Makefile.microservices docker-logs

# Health check
make -f Makefile.microservices health-check

# Stop services
make -f Makefile.microservices docker-down
```

---

## 📊 Monitoring

### Prometheus (Metrics)
- URL: `http://localhost:9090`
- All services expose `/metrics` endpoint

### Grafana (Dashboards)
- URL: `http://localhost:3000`
- Default credentials: `admin` / `admin`

### Health Checks
Each service has `/health` endpoint:
```bash
curl http://localhost:8001/health
```

---

## 🏛️ Architecture Principles

### 1. **Clean Architecture**
```
handler → service → repository → database
```

### 2. **Dependency Injection**
All dependencies injected via constructors

### 3. **Interface-Based Design**
Repository and service layers use interfaces

### 4. **Error Handling**
Custom `AppError` type with HTTP status codes

### 5. **Validation**
Request validation with custom validator

### 6. **Logging**
Structured logging with context

---

## 📁 Project Structure

```
waqfwise/
├── cmd/
│   ├── api-gateway/           # API Gateway entry point
│   ├── auth-service/          # Auth service entry point
│   ├── payment-service/       # Payment service entry point
│   ├── campaign-service/      # Campaign service entry point
│   ├── asset-service/         # Asset service entry point
│   ├── analytics-service/     # Analytics service (Enterprise)
│   └── integration-service/   # Integration service (Enterprise)
├── internal/
│   ├── services/
│   │   ├── auth/              # Auth service implementation
│   │   │   ├── dto/           # Data transfer objects
│   │   │   ├── handler/       # HTTP handlers
│   │   │   ├── service/       # Business logic
│   │   │   └── repository/    # Data access layer
│   │   ├── payment/
│   │   ├── campaign/
│   │   └── asset/
│   └── shared/
│       ├── domain/            # Domain models
│       ├── errors/            # Error types
│       ├── response/          # Response helpers
│       └── validator/         # Validation logic
├── Dockerfile.microservices   # Multi-stage Dockerfile
├── docker-compose.microservices.yml
├── Makefile.microservices
└── MICROSERVICES.md
```

---

## 🔐 Security Best Practices

1. ✅ **JWT with short expiry** (1 hour access, 7 days refresh)
2. ✅ **Password hashing** with bcrypt (cost 12)
3. ✅ **MFA support** with TOTP
4. ✅ **Rate limiting** on sensitive endpoints
5. ✅ **HTTPS only** in production
6. ✅ **SQL injection prevention** with parameterized queries
7. ✅ **CORS configuration** for allowed origins
8. ✅ **Secrets management** via environment variables

---

## 🌍 Environment Variables

### Common (All Services)
```env
PORT=8001
DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
REDIS_URL=redis:6379
```

### Auth Service
```env
JWT_SECRET=your-secret-key-change-in-production
```

### Payment Service
```env
MIDTRANS_SERVER_KEY=your-midtrans-key
XENDIT_SECRET_KEY=your-xendit-key
```

---

## 📈 Performance Considerations

1. **Database Connection Pooling** - Configured per service
2. **Redis Caching** - Sessions & frequently accessed data
3. **Horizontal Scaling** - Each service can scale independently
4. **Load Balancing** - API Gateway distributes requests
5. **Asynchronous Processing** - Kafka for background jobs

---

## 🤝 Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

---

## 📜 License

This project is dual-licensed:
- **AGPL-3.0** for Community Edition
- **Commercial License** for Enterprise Edition

See [LICENSE](LICENSE) for details.

---

## 📞 Support

- **Documentation**: [docs.waqfwise.id](https://docs.waqfwise.id)
- **Issues**: [GitHub Issues](https://github.com/akordium-id/waqfwise/issues)
- **Email**: support@waqfwise.id

---

**Built with ❤️ using idiomatic Go and microservices best practices**
