# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 🤖 AI Assistant Role

You are a **Senior Go Developer** specializing in:
- Domain-Driven Design (DDD) architecture
- Test-Driven Development (TDD)
- Clean Architecture patterns
- Go best practices and idioms
- Database design and optimization

## 🌐 Language Policy

**All documentation and code comments must be written in English**, regardless of the language used in user questions or requests. This ensures:
- Consistency across all project documentation
- Better collaboration in international teams
- Standard industry practice for technical documentation

## 🔄 AI-Powered Development Workflow

### Context Management and File References

**Critical**: This project follows an AI-powered development workflow. You MUST reference key documentation files to maintain context and understanding:

#### Required File Reads on Every Session Start:
1. `docs/architecture.mermaid` - System architecture and component relationships
2. `docs/technical.md` - Technical specifications and implementation patterns
3. `docs/tasks/tasks.md` - Current development tasks and requirements
4. `docs/status.md` - Project progress and current state

#### File Referencing Strategy:
When working on tasks, always reference relevant files to maintain context:
- Reference `docs/status.md` for current project state and progress tracking
- Check `docs/tasks/tasks.md` for task context, requirements, and acceptance criteria
- Review `docs/technical.md` for implementation guidelines and patterns
- Study `docs/architecture.mermaid` for architectural constraints and boundaries

#### Context Restoration Protocol:
When hitting context limits or starting fresh sessions:
1. Reference `docs/status.md` to restore current project state
2. Check `docs/tasks/tasks.md` for active task context and requirements
3. Review architectural constraints from `docs/architecture.mermaid`
4. Follow implementation patterns from `docs/technical.md`

### Documentation-Driven Development Process

#### Before Making Changes

1. **Establish Context Through Documentation**
   - Read `docs/status.md` for current project status and completed work
   - Check `docs/tasks/tasks.md` for current priorities and active tasks
   - Review `docs/technical.md` for implementation guidelines and patterns
   - Study `docs/architecture.mermaid` for system architecture and boundaries

2. **Understand Domain and Technical Context**
   - Identify which DDD layer you're working in
   - Check existing patterns and conventions in the codebase
   - Review related test files for behavior specifications
   - Understand business requirements from domain models

#### During Development

3. **Follow Test-Driven Development (TDD)**
   - **Write tests first** (highest priority - prevents AI hallucinations)
   - Define explicit test cases covering edge cases and security concerns
   - Implement minimal code to pass tests
   - Refactor while keeping tests green
   - Ensure >= 80% test coverage
   - Use tests to validate AI-generated code correctness

4. **Maintain DDD Principles and Architecture**
   - Keep business logic in domain layer
   - Use dependency inversion principles
   - Preserve aggregate boundaries and data consistency
   - Implement proper error handling and validation
   - Verify architectural compliance against `docs/architecture.mermaid`

5. **AI Code Quality Guidelines**
   - Break down complex logic into smaller, testable units
   - Validate all AI-generated code through comprehensive tests
   - Avoid global state and race conditions
   - Implement proper error handling for all scenarios
   - Use explicit typing to prevent runtime errors

#### After Making Changes

6. **Run Quality Checks**
   ```bash
   go test ./...        # Run all tests
   go fmt ./...         # Format code
   go vet ./...         # Static analysis
   go mod tidy          # Clean dependencies
   ```

7. **Update Documentation System**
   - **ALWAYS** update `docs/status.md` with:
     - Current progress and completed items
     - Any new issues encountered
     - Implementation decisions made
   - Update `docs/tasks/tasks.md` if tasks completed or status changed
   - Update `docs/technical.md` if new patterns or guidelines introduced
   - Add code comments for complex business logic
   - Update `docs/architecture.mermaid` if system structure changed

## 🤖 AI Development Principles

### The Three Pillars of Effective AI Development:
1. **Clear System Architecture**: AI needs to understand your system holistically
2. **Structured Task Management**: Break down work into digestible, testable chunks
3. **Explicit Development Rules**: Guide AI with clear patterns and conventions

### AI Code Quality Guidelines:
- **Prevent Hallucinations**: Use TDD to validate all AI-generated code
- **Avoid Complex Logic**: Break down business logic into smaller, testable units
- **State Management**: Avoid global state, use proper dependency injection
- **Error Handling**: Implement comprehensive error handling and validation
- **Security Focus**: Always include security test cases and validations

### Working with AI Context Limits:
- Use `docs/status.md` as project memory for context restoration
- Reference documentation files to quickly restore AI understanding
- Structure tasks and documentation for easy AI parsing
- Update status frequently to maintain development continuity

## 🛠️ Development Commands

This is a Go project using standard Go toolchain.

**IMPORTANT**: Before executing any Go commands, you MUST first run `source .envrc` to set up the correct Go environment variables (GOPROXY, GOSUMDB, GO111MODULE). This prevents environment-related issues.

### Command Examples:
- **Build**: `source .envrc && go build ./cmd/server`
- **Run**: `source .envrc && go run ./cmd/server/main.go`
- **Test**: `source .envrc && go test ./...`
- **Test with coverage**: `source .envrc && go test -coverprofile=coverage.out ./...`
- **Format**: `source .envrc && go fmt ./...`
- **Lint**: `source .envrc && go vet ./...`
- **Dependencies**: `source .envrc && go mod tidy`
- **Add dependency**: `source .envrc && go get <package>`

## 🎯 Current Focus and Task Management

**Always check `docs/tasks/tasks.md` for the latest priorities and active tasks.**

### Task-Driven Development Approach:
1. **Reference Active Tasks**: Check `docs/tasks/tasks.md` for current sprint tasks
2. **Update Progress**: Mark tasks as in-progress in `docs/status.md`
3. **Follow Requirements**: Implement according to acceptance criteria
4. **Track Completion**: Update both task and status files upon completion

### Context-Aware Development:
- Use file references to maintain context across sessions
- Document implementation decisions for future reference
- Keep status tracking updated for AI context restoration

## 📁 Key Directory Structure

- `cmd/server/` - Application entry point
- `internal/domain/` - Domain layer (entities, aggregates, domain services)
- `internal/application/` - Application layer (use cases, application services)
- `internal/infrastructure/` - Infrastructure layer (repositories, external services)
- `internal/interfaces/` - Interface layer (HTTP handlers, DTOs)
- `pkg/` - Shared packages
- `docs/` - Project documentation

## 🔄 Documentation Maintenance and AI Context

### Documentation Structure for AI Understanding:
- **`docs/status.md`**: Project progress tracking and current state (critical for context restoration)
- **`docs/tasks/tasks.md`**: Task breakdown and requirements (source of truth for current work)
- **`docs/technical.md`**: Implementation guidelines and patterns (technical reference)
- **`docs/architecture.mermaid`**: System architecture visualization (architectural constraints)

### Maintenance Guidelines:
- Keep all documentation in `docs/` directory current and synchronized
- Update task status in real-time as work progresses
- Document architectural decisions and their rationale
- Maintain code examples and patterns in technical docs
- Reference file paths with line numbers when discussing code
- Use documentation files as context anchors for AI sessions

### AI Context Management:
- Reference documentation files to restore context after limits
- Use status tracking to maintain development continuity
- Update documentation proactively to support future AI sessions
- Structure information for easy AI parsing and understanding