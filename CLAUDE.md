# CLAUDE.md - Wonder Project Root

This file provides guidance to Claude Code when working from the project root directory.

## 🤖 Context-Aware Development

When working from the root directory, Claude Code should **automatically determine the appropriate development context** based on the files being modified:

### 📁 File-Based Context Detection

#### Backend Development Context
**Trigger**: When modifying files in `backend/` directory
**Apply**: Backend development guidelines from `backend/CLAUDE.md`

Examples:
- `backend/internal/**/*.go` → Use Go/DDD development patterns
- `backend/cmd/**/*.go` → Use Go application patterns
- `backend/pkg/**/*.go` → Use Go package development patterns
- `backend/configs/**/*.yaml` → Use backend configuration patterns
- `backend/docs/**/*.md` → Use backend documentation standards

#### Frontend Development Context
**Trigger**: When modifying files in `frontend/` directory
**Apply**: Frontend development guidelines from `frontend/CLAUDE.md`

Examples:
- `frontend/src/**/*.{js,ts,jsx,tsx}` → Use frontend component patterns
- `frontend/public/**/*` → Use frontend asset management
- `frontend/package.json` → Use frontend dependency management
- `frontend/docs/**/*.md` → Use frontend documentation standards

#### Project-Level Context
**Trigger**: When modifying root-level files
**Apply**: Project-wide guidelines from this file

Examples:
- `README.md` → Use project overview documentation standards
- `docker-compose.yml` → Use full-stack deployment patterns
- `.gitignore` → Use project-wide ignore patterns
- Root-level configuration files

## 🔄 Context Switching Rules

### Automatic Context Detection
```
IF modifying files in backend/**
  THEN use backend/CLAUDE.md guidelines
  AND apply Go development patterns
  AND follow DDD architecture principles

ELSE IF modifying files in frontend/**
  THEN use frontend/CLAUDE.md guidelines
  AND apply frontend development patterns
  AND follow component architecture principles

ELSE IF modifying root-level files
  THEN use project-level guidelines
  AND consider full-stack implications
```

### Multi-Context Changes
When making changes that span multiple contexts:

1. **Full-Stack Features**: Apply both backend and frontend guidelines
2. **API Changes**: Update backend API + frontend integration
3. **Configuration**: Consider impact on both services
4. **Documentation**: Update relevant documentation in both contexts

## 🛠️ Root Directory Commands

### Project-Level Operations
```bash
# Full-stack development setup
docker-compose up -d

# Run backend from root
cd backend && source .envrc && go run cmd/server/main.go

# Run frontend from root (when ready)
cd frontend && npm run dev

# Full project testing
cd backend && ./scripts/test.sh all
cd frontend && npm test

# Project-wide linting/formatting
cd backend && go fmt ./...
cd frontend && npm run lint
```

### Development Workflow from Root
```bash
# 1. Start services
docker-compose up -d postgres redis  # Start dependencies

# 2. Backend development
cd backend/
source .envrc
go run cmd/server/main.go
# (In another terminal for backend changes)

# 3. Frontend development
cd frontend/
npm run dev
# (In another terminal for frontend changes)
```

## 📋 Root-Level Development Guidelines

### File Modification Guidelines

#### When modifying `backend/**` files:
- **MUST** follow backend/CLAUDE.md guidelines
- **MUST** run backend tests: `cd backend && go test ./...`
- **MUST** ensure Go code compilation: `cd backend && go build ./...`
- **SHOULD** verify backend service still runs properly

#### When modifying `frontend/**` files:
- **MUST** follow frontend/CLAUDE.md guidelines
- **MUST** run frontend tests: `cd frontend && npm test`
- **MUST** ensure frontend builds: `cd frontend && npm run build`
- **SHOULD** verify frontend app still runs properly

#### When modifying root files:
- **MUST** consider impact on both backend and frontend
- **MUST** update documentation if architecture changes
- **SHOULD** test full-stack integration if needed

### Cross-Context Changes

#### API Development Workflow:
1. **Backend**: Add/modify API endpoints in backend/
2. **Testing**: Test API with backend/api.http
3. **Frontend**: Update frontend API integration
4. **Integration**: Test end-to-end functionality
5. **Documentation**: Update both backend and frontend docs

#### Configuration Changes:
1. **Environment**: Update relevant config files
2. **Backend**: Update backend configs if needed
3. **Frontend**: Update frontend configs if needed
4. **Docker**: Update docker-compose.yml if needed
5. **Documentation**: Update setup instructions

## 🎯 Context-Specific Task Management

### Backend Tasks (from root):
```bash
# Navigate to backend context
cd backend/

# Apply backend development workflow
source .envrc
# ... follow backend/CLAUDE.md guidelines
```

### Frontend Tasks (from root):
```bash
# Navigate to frontend context
cd frontend/

# Apply frontend development workflow
# ... follow frontend/CLAUDE.md guidelines
```

### Full-Stack Tasks:
- Use project-level perspective
- Consider both services simultaneously
- Apply integration testing approaches
- Update project-wide documentation

## 🚨 Important Context Rules

### File Path Detection
Claude Code should automatically detect context based on file paths:

- **`backend/**`** → Auto-apply backend development standards
- **`frontend/**`** → Auto-apply frontend development standards
- **Root files** → Apply project-level standards

### Testing Requirements
When modifying files from root directory:

```bash
# Backend changes - run backend tests
cd backend && source .envrc && go test ./...

# Frontend changes - run frontend tests
cd frontend && npm test

# Full-stack changes - run both test suites
cd backend && source .envrc && go test ./...
cd frontend && npm test
```

### Documentation Updates
- **Backend changes** → Update backend docs
- **Frontend changes** → Update frontend docs
- **Architecture changes** → Update root README.md
- **API changes** → Update both backend/api.http and frontend API docs

---

**Key Principle**: Always apply the most specific development context available based on the files being modified, while maintaining awareness of full-stack implications.