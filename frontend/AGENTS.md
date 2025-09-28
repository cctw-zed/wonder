# AGENTS.md - Frontend Collaboration Guide

This guide extends the root `AGENTS.md` rules for work performed within `frontend/`. All root-level collaboration, testing, documentation, and communication expectations remain in effect.

## 🔁 Inheritance Reminder

- Apply the root `AGENTS.md` instructions first.
- Layer in the frontend-specific rules below.
- If future subdirectories add their own `AGENTS.md`, follow the most specific file applicable to your working directory.

## 🎯 Frontend Agent Focus

Agents operating in `frontend/` must:
- Confirm alignment with `frontend/CLAUDE.md` before making changes.
- Reference component architecture, state management patterns, and testing strategy defined in the frontend documentation and codebase.
- Call out any required coordination with backend agents when API contracts or shared models change.

## 🧪 Quality Expectations

- Recommend or run `npm test`, `npm run lint`, `npm run type-check`, and `npm run build` after code modifications unless explicitly deferred by the user.
- Highlight unmet test obligations or skipped checks with rationale and next steps.

## 📚 Documentation & Context

- Update frontend-specific docs (e.g., `frontend/README.md`, component notes) when behavior, APIs, or workflows change.
- Reference relevant files with line numbers when handing off work or summarizing findings.
- Sync progress with project-wide docs in `docs/` when frontend work impacts broader milestones.

## 🤝 Cross-Agent Coordination

- Note dependencies on backend updates, design assets, or environment configuration.
- Surface outstanding questions or risks early for downstream agents.
- Provide clear handoff notes, including pending tests or verifications, when concluding a task.

By following these rules alongside the inherited guidance, frontend agents ensure consistent, high-quality collaboration throughout the Wonder project.
