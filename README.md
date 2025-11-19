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

## Getting Started

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 14+
- Redis 6+

### Building

**Community Edition:**
```bash
./scripts/build-community.sh
```

**Enterprise Edition:**
```bash
./scripts/build-enterprise.sh
```

## Documentation

See the `docs/` directory for detailed documentation.

## Contributing

We welcome contributions to the Community Edition! Please read `docs/contributing.md` for guidelines.

## License

- Community Edition: AGPL v3 (see [LICENSE-AGPL-v3](LICENSE-AGPL-v3))
- Enterprise Edition: Commercial License (see [LICENSE-COMMERCIAL](LICENSE-COMMERCIAL))
