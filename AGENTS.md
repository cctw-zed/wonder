# AGENTS.md - Wonder Project Multi-Agent Guide

This document defines how AI agents collaborate when working in the Wonder project repository. It complements `CLAUDE.md` by emphasizing coordinated workflows, context detection, and shared responsibilities across different AI assistants.

## 🧬 Instruction Inheritance

`AGENTS.md` files follow a hierarchical inheritance model similar to `CLAUDE.md`:
- The root `AGENTS.md` applies to the entire repository and is always in effect.
- Nested directories may define their own `AGENTS.md` to add context-specific expectations.
- Directory-specific guides **extend** the root rules; they never replace them. If instructions conflict, follow the most specific directory guide while still honoring root-level requirements.
- When multiple nested guides exist, apply the rules from the nearest ancestor up to the root.

Codex and other agents must explicitly acknowledge which `AGENTS.md` files apply to their current working directory and confirm that inheritance has been respected.

## 🤝 Shared Mission

All agents contribute to the same engineering goals:
- Maintain architectural integrity across backend and frontend services
- Keep documentation in sync with implementation changes
- Preserve high code quality through testing and review discipline
- Communicate assumptions, risks, and gaps explicitly in responses

## 🧭 Context-Aware Collaboration

Follow the same context-detection rules described in `CLAUDE.md`:
- **`backend/**` changes** → Apply `backend/CLAUDE.md` guidance and Go/DDD standards
- **`frontend/**` changes** → Apply `frontend/CLAUDE.md` guidance and component-driven patterns
- **Root-level changes** → Apply project-wide practices from `CLAUDE.md` and consider full-stack impacts
- **Cross-cutting work** → Combine relevant context rules and document coordination points

Agents must announce the context they are operating in, reference the appropriate guidelines, and confirm that required tests or validations have been considered before finishing a task.

## 🗂️ Role Alignment

Although individual agents may specialize differently, they share common expectations:
- Review relevant documentation (`docs/status.md`, `docs/tasks/tasks.md`, `docs/technical.md`, `docs/architecture.mermaid`) before significant work
- Reference file paths with line numbers when discussing code
- Escalate ambiguous requirements or conflicting instructions to the user
- Default to English for documentation, comments, and commit messages

When multiple agents touch the same area, they must highlight dependencies, note pending follow-ups, and avoid duplicating effort.

## 🔄 Handoff Protocols

To ensure continuity between agents:
1. **Summarize Progress** – Capture what was changed, why it was changed, and remaining risks
2. **List Next Actions** – Provide clear guidance for the next agent (tests to run, docs to update, reviews required)
3. **Attach Evidence** – Reference relevant diffs, tests, or logs that justify readiness
4. **Flag Blockers Early** – Identify missing information or external dependencies as soon as they appear

## 🧪 Quality and Testing Expectations

Adhere to the testing mandates from `CLAUDE.md` and specialized context files:
- Run or recommend the required test suites after code changes (`go test ./...`, `npm test`, etc.)
- Call out tests that were skipped or require follow-up, including reasons
- Ensure configuration changes are validated across environments if applicable

## 📝 Documentation Discipline

- Update `docs/status.md` and `docs/tasks/tasks.md` to reflect progress and remaining work
- Keep architectural notes current when structural changes occur
- Cross-reference related documents to support future sessions and restore context quickly

## 📣 Communication Standards

- Be concise but thorough; focus on risks, decisions, and verifications
- Use bullet lists for readability when summarizing findings or steps
- Provide actionable next steps instead of generic advice
- Confirm adherence to repository conventions before handing work back to the user

By following this guide alongside the existing `CLAUDE.md` and context-specific instructions, AI agents can collaborate effectively and maintain consistent outcomes across the Wonder project.
