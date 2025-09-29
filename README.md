# Wonder Project

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Next.js](https://img.shields.io/badge/Next.js-14+-black.svg)](https://nextjs.org)
[![Docker](https://img.shields.io/badge/Docker-Layered%20Architecture-blue.svg)](#layered-deployment)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](#)

A modern full-stack application with **layered deployment architecture** for optimal development experience and data persistence.

## 🏗️ Layered Architecture

Wonder uses a **layered deployment architecture** that separates concerns for better development experience:

- **Infrastructure Layer**: PostgreSQL, Elasticsearch, Grafana (persistent data services)
- **Monitoring Layer**: Prometheus, Logstash, Kibana, cAdvisor (stateless monitoring services)
- **Application Layer**: Backend API (Go) and Frontend Web (Next.js) (frequently updated services)

```
wonder/
├── backend/                    # Go API service
├── frontend/                   # Next.js web application
├── scripts/                    # Development automation scripts
├── docker-compose.*.yaml       # Layered Docker configurations
├── Makefile                    # Development commands
└── README.md                   # This file
```

This separation ensures:
- ✅ Data is preserved when rebuilding applications
- ✅ Monitoring dashboards and configurations persist
- ✅ Quick rebuilds don't affect the entire stack
- ✅ Independent scaling and maintenance

## 🚀 Quick Start

### Prerequisites

- **Docker & Docker Compose** (required)
- **Go 1.24+** (for backend development)
- **Node.js 18+** (for frontend development)

### First Time Setup

```bash
# Clone and navigate to project
git clone <repository-url>
cd wonder

# Set up complete development environment
make setup
```

This will start all services in the correct order and show service URLs when ready.

### Daily Development Workflow

```bash
# Check what's running
make status

# After making backend changes
make rebuild-backend

# After making frontend changes
make rebuild-frontend

# After making changes to both
make rebuild

# View service URLs
make urls

# Check logs
make logs
```

## 📋 Available Commands

### 🚀 Environment Management
- `make setup` - Set up complete development environment
- `make status` - Show status of all services
- `make urls` - Show all service URLs
- `make logs` - Show logs from all services

### 🔄 Quick Development
- `make rebuild` - Rebuild applications (preserves data)
- `make rebuild-backend` - Rebuild only backend service
- `make rebuild-frontend` - Rebuild only frontend service
- `make rebuild-no-cache` - Rebuild without Docker cache

### 🧹 Cleanup
- `make stop-apps` - Stop application services only
- `make stop-all` - Stop all services (preserves data)
- `make clean-all` - Complete reset (⚠️ removes all data)

### 🧪 Testing
- `make test` - Run all tests
- `make test-backend` - Run backend tests only
- `make test-frontend` - Run frontend tests only

## 🎯 Service Overview

### Backend Service
- **Technology**: Go with Gin framework
- **Architecture**: Domain-Driven Design (DDD)
- **Features**: JWT Authentication, RESTful API, Monitoring
- **Port**: `8080`
- **Health Check**: `http://localhost:8080/health`

### Frontend Application
- **Technology**: Next.js with TypeScript
- **Architecture**: Component-based architecture
- **Features**: Responsive design, API integration
- **Port**: `3001`
- **Status**: ✅ *Ready*

### Monitoring Stack
- **Prometheus**: Metrics collection (`http://localhost:9090`)
- **Grafana**: Dashboard visualization (`http://localhost:3000`)
- **ELK Stack**: Centralized logging
- **Status**: ✅ *Ready*

## 📊 Current Status

| Component | Status | Coverage | Documentation |
|-----------|--------|----------|---------------|
| Backend API | ✅ Ready | 93.4% | [Complete](backend/README.md) |
| Frontend | ✅ Ready | - | [Complete](frontend/README.md) |
| Monitoring | ✅ Ready | - | [Available](backend/docs/monitoring.md) |
| E2E Integration | 🚧 Planned | - | Pending |

## 🛠️ Development Workflow

### 1. Backend-First Development
The backend service is fully implemented and provides:
- Complete REST API with authentication
- Comprehensive monitoring and logging
- Database migrations and seeding
- Full test coverage (93.4%)

### 2. Frontend Integration
Frontend development should:
- Connect to backend API at `http://localhost:8080/api/v1/`
- Use JWT tokens for authentication
- Follow component-based architecture
- Implement responsive design

### 3. Full-Stack Testing
- Backend has comprehensive unit, integration, and E2E tests
- Frontend should implement similar testing strategy
- End-to-end integration tests across both services

## 🔗 API Integration

### Backend API Endpoints
```bash
# Health check
GET http://localhost:8080/health

# User authentication
POST http://localhost:8080/api/v1/users/register
POST http://localhost:8080/api/v1/users/login
GET http://localhost:8080/api/v1/users/profile

# Metrics
GET http://localhost:8080/metrics
```

**📋 [Complete API Documentation](backend/api.http)**

### Authentication Flow
1. Frontend sends credentials to `/api/v1/users/login`
2. Backend returns JWT token
3. Frontend includes token in Authorization header
4. Backend validates token for protected routes

## 🐳 Docker Deployment

### Development Environment
```bash
# Start backend services only
cd backend/
docker-compose up -d

# Start full stack (when frontend is ready)
docker-compose -f docker-compose.yml -f docker-compose.frontend.yml up -d
```

### Production Deployment
```bash
# Build and deploy
docker-compose -f docker-compose.prod.yml up -d
```

## 📚 Documentation Structure

```
├── README.md                     # This overview document
├── backend/
│   ├── README.md                # Backend service documentation
│   ├── CLAUDE.md               # Backend AI development guide
│   ├── docs/                   # Technical documentation
│   └── api.http               # API testing examples
└── frontend/
    ├── README.md              # Frontend application documentation
    ├── CLAUDE.md             # Frontend AI development guide
    └── docs/                 # Frontend-specific documentation
```

## 🎯 Development Priorities

### Immediate (Sprint 1)
- [ ] Frontend framework selection and setup
- [ ] Basic component structure implementation
- [ ] API integration with backend
- [ ] Authentication flow implementation

### Near-term (Sprint 2-3)
- [ ] Frontend testing setup
- [ ] Full-stack E2E testing
- [ ] Production deployment configuration
- [ ] Performance optimization

### Long-term
- [ ] Advanced monitoring and analytics
- [ ] CI/CD pipeline implementation
- [ ] Progressive Web App features
- [ ] Mobile application considerations

## 🧪 Testing Strategy

### Backend Testing (✅ Complete)
- **Unit Tests**: Domain and application layer logic
- **Integration Tests**: Database and external service interactions
- **E2E Tests**: Complete API workflow testing
- **Coverage**: 93.4%

### Frontend Testing (🚧 Planned)
- **Unit Tests**: Component logic and utilities
- **Integration Tests**: Component interactions and API calls
- **E2E Tests**: Complete user workflows
- **Target Coverage**: 80%+

### Full-Stack Testing (🚧 Planned)
- **Integration**: Frontend ↔ Backend API integration
- **User Flows**: Complete authentication and data workflows
- **Performance**: Load testing and optimization

## 🤝 Contributing

### Development Setup
1. **Clone and setup**: Follow quick start guide above
2. **Backend changes**: See [backend development guide](backend/CLAUDE.md)
3. **Frontend changes**: See [frontend development guide](frontend/CLAUDE.md)
4. **Testing**: Ensure all tests pass before committing
5. **Documentation**: Update relevant docs with changes

### Code Quality Standards
- **Backend**: Go best practices, DDD principles, comprehensive testing
- **Frontend**: Modern JavaScript/TypeScript, component architecture, accessibility
- **Integration**: Consistent error handling, proper API contracts
- **Documentation**: Keep all documentation current and comprehensive

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Backend API**: [http://localhost:8080](http://localhost:8080)
- **Frontend App**: [http://localhost:3001](http://localhost:3001)
- **Monitoring**: [http://localhost:3000](http://localhost:3000) *(Grafana)*
- **Metrics**: [http://localhost:9090](http://localhost:9090) *(Prometheus)*

---

**Wonder** - Modern full-stack application with clean architecture and comprehensive tooling.