# AGENTS.md - Backend Collaboration Guide

This guide extends the root `AGENTS.md` instructions for all work performed within `backend/`. Continue to follow the repository-wide collaboration, testing, documentation, and communication standards defined at the root.

## 🔁 Inheritance Reminder

- Apply the root `AGENTS.md` rules before adding the backend-specific expectations below.
- If deeper subdirectories introduce their own `AGENTS.md`, always follow the most specific applicable file while honoring ancestor requirements.

## 🎯 Backend Agent Focus

Agents operating in `backend/` must:
- Align with `backend/CLAUDE.md` and `backend/CODEX.md` prior to making changes.
- Respect the domain-driven design (DDD) boundaries and layering described in the documentation.
- Coordinate with frontend agents when API schemas, contracts, or shared DTOs change.

## 🧪 Quality Expectations

- Recommend or run `source .envrc && go test ./...`, `source .envrc && go build ./...`, and other mandatory checks described in `backend/CLAUDE.md` after code modifications unless explicitly deferred by the user.
- Call out skipped tests or checks with justification and next steps.

## 📚 Documentation & Context

- Update backend-specific docs (e.g., `backend/README.md`, `backend/docs/**`) when behavior, APIs, or workflows change.
- Reference relevant files with line numbers when handing off work or summarizing findings.
- Sync project-wide documentation in `docs/` when backend work affects cross-team milestones or integrations.

## 🤝 Cross-Agent Coordination

- Surface dependencies on infrastructure, database migrations, or environment configuration early.
- Provide clear handoff notes, including pending verifications or open questions, when concluding a task.
- Flag API or contract changes so downstream agents (frontend, integrations, QA) can respond promptly.

By layering these backend-focused expectations on top of the root guidance, agents ensure reliable, high-quality collaboration throughout the Wonder project.
