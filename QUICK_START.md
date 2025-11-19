# 🚀 WaqfWise Quick Start Guide

Choose your architecture and get started in minutes!

---

## 🎯 Choose Your Architecture

WaqfWise supports **two architectures** - choose based on your needs:

### **Option A: Monolith** ⭐ *Recommended for getting started*

**Best for:**
- Small to medium teams (< 10 developers)
- MVP / Early stage startups
- Traffic < 100K daily active users
- Budget conscious ($10-20/month)
- Need fast iteration

**Pros:**
✅ Simple deployment (single binary)
✅ Fast development (no service orchestration)
✅ Low cost infrastructure
✅ Easy debugging
✅ ACID transactions

**Get Started:**
```bash
# Start with Docker Compose
make -f Makefile.monolith docker-up

# Or build and run locally
make -f Makefile.monolith build
make -f Makefile.monolith run
```

📖 **Full Guide:** [MONOLITH.md](MONOLITH.md)

---

### **Option B: Microservices** 🏗️

**Best for:**
- Large teams (> 15 developers)
- Scale-up / Enterprise
- Traffic > 500K daily active users
- Need independent scaling per service
- Multiple product teams

**Pros:**
✅ Independent scaling
✅ Fault isolation
✅ Independent deployment
✅ Technology flexibility
✅ Team autonomy

**Get Started:**
```bash
# Start all services with Docker Compose
make -f Makefile.microservices docker-up

# Or build all services
make -f Makefile.microservices build-all
```

📖 **Full Guide:** [MICROSERVICES.md](MICROSERVICES.md)

---

## ⚡ Quick Start (5 Minutes)

### **Monolith - Docker Compose** (Easiest)

```bash
# 1. Clone repository
git clone https://github.com/akordium-id/waqfwise.git
cd waqfwise

# 2. Copy environment file
cp .env.monolith.example .env

# 3. Start everything with one command
make -f Makefile.monolith docker-up

# 4. Done! Services running at:
#    App:        http://localhost:8080
#    Health:     http://localhost:8080/health
#    Prometheus: http://localhost:9090
#    Grafana:    http://localhost:3000
```

### **Microservices - Docker Compose**

```bash
# 1. Clone repository
git clone https://github.com/akordium-id/waqfwise.git
cd waqfwise

# 2. Start all microservices
make -f Makefile.microservices docker-up

# 3. Done! Services running at:
#    API Gateway: http://localhost:8000
#    Auth:        http://localhost:8001
#    Payment:     http://localhost:8002
#    Campaign:    http://localhost:8003
#    Asset:       http://localhost:8004
```

---

## 📋 Prerequisites

- **Docker** & Docker Compose (for containerized deployment)
- **Go 1.21+** (for local development)
- **PostgreSQL 15+** (included in Docker Compose)
- **Redis 7+** (included in Docker Compose)

---

## 🎓 Next Steps

### **After starting the application:**

1. **Test the API:**
   ```bash
   # Health check
   curl http://localhost:8080/health  # Monolith
   curl http://localhost:8000/health  # Microservices (API Gateway)

   # Register a user
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "email": "test@example.com",
       "password": "SecurePass123!",
       "name": "Test User"
     }'
   ```

2. **Explore the documentation:**
   - [Monolith Guide](MONOLITH.md)
   - [Microservices Guide](MICROSERVICES.md)
   - [API Documentation](docs/api/)

3. **Setup monitoring:**
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3000 (admin/admin)

4. **Read the code:**
   - Start with `cmd/` for entry points
   - Explore `internal/services/` for business logic
   - Check `internal/shared/` for shared utilities

---

## 🆚 Comparison Table

| Feature | Monolith | Microservices |
|---------|----------|---------------|
| **Deployment** | Single binary | 7 services |
| **Complexity** | Low | High |
| **Cost** | $10-20/month | $50-100/month |
| **Latency** | < 10ms | ~15-150ms |
| **Scaling** | Vertical | Horizontal per service |
| **Debugging** | Easy | Complex |
| **Team Size** | < 10 | > 15 |
| **Tech Stack** | Unified (Go) | Flexible per service |
| **Transactions** | ACID | Eventual consistency |
| **Deployment Time** | ~1 min | ~5-10 min |

---

## 🔄 Migration Path

**Don't worry about choosing wrong!** WaqfWise is designed for easy migration:

```
Start: Monolith (simple)
  ↓
Phase 1: Modular Monolith (clean boundaries)
  ↓
Phase 2: Extract high-traffic services
  ↓
Phase 3: Selective microservices
```

The same clean architecture works for both!

---

## 🛠️ Common Commands

### **Monolith**
```bash
make -f Makefile.monolith help        # Show all commands
make -f Makefile.monolith build       # Build binary
make -f Makefile.monolith test        # Run tests
make -f Makefile.monolith docker-up   # Start with Docker
make -f Makefile.monolith docker-logs # View logs
```

### **Microservices**
```bash
make -f Makefile.microservices help        # Show all commands
make -f Makefile.microservices build-all   # Build all services
make -f Makefile.microservices docker-up   # Start all services
make -f Makefile.microservices docker-logs # View logs
make -f Makefile.microservices health-check # Check all services
```

---

## 🐛 Troubleshooting

**Port already in use?**
```bash
# For monolith (port 8080)
lsof -i :8080
kill -9 <PID>

# Or change port in .env
PORT=8081
```

**Database connection failed?**
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Test connection
psql "postgres://waqfwise:waqfwise@localhost:5432/waqfwise_community"
```

**More help:**
- [Monolith Troubleshooting](MONOLITH.md#troubleshooting)
- [Microservices Troubleshooting](MICROSERVICES.md#troubleshooting)
- [GitHub Issues](https://github.com/akordium-id/waqfwise/issues)

---

## 📚 Documentation

- **[MONOLITH.md](MONOLITH.md)** - Complete monolith guide
- **[MICROSERVICES.md](MICROSERVICES.md)** - Complete microservices guide
- **[README.md](README.md)** - Project overview
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - How to contribute

---

## 💬 Support

- **Email**: support@waqfwise.id
- **Issues**: https://github.com/akordium-id/waqfwise/issues
- **Discussions**: https://github.com/akordium-id/waqfwise/discussions

---

**Happy coding! 🚀**

**Remember:** Start simple (monolith) → Scale when needed → Extract to microservices
