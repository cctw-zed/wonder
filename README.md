# Wonder - Full-Stack Application

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Frontend](https://img.shields.io/badge/Frontend-Modern%20Stack-brightgreen.svg)](#frontend)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](#docker-deployment)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](#)

A modern, scalable full-stack application with clean architecture separation between frontend and backend components.

## 🏗️ Project Architecture

Wonder follows a **microservice-oriented architecture** with clear separation between frontend and backend:

```
wonder/
├── backend/          # Go backend service with DDD architecture
├── frontend/         # Modern frontend application
├── README.md         # This file - project overview
└── docker-compose.yml # Full-stack deployment (planned)
```

## 🚀 Quick Start

### Prerequisites

- **Docker & Docker Compose** (recommended for full-stack setup)
- **Go 1.24+** (for backend development)
- **Node.js 18+** (for frontend development)

### Full-Stack Setup with Docker

```bash
# Clone the repository
git clone https://github.com/your-username/wonder.git
cd wonder

# Start the full stack (backend + frontend + monitoring)
docker-compose up -d

# View services
docker-compose ps
```

### Development Setup

#### Backend Development
```bash
cd backend/
source .envrc
go mod download
go run cmd/server/main.go
```
**📖 [Backend Documentation](backend/README.md)**

#### Frontend Development
```bash
cd frontend/
npm install
npm run dev
```
**📖 [Frontend Documentation](frontend/README.md)** *(to be created)*

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