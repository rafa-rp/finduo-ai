# GEMINI.md

This file serves as a guide for instructions and configurations for AI agents (such as Antigravity/Gemini) collaborating on the development of the **finduo-ai** project.

## 1. Project Context
* **Name**: finduo-ai
* **Description**: AI-powered split finance app — hands-on project for mastering modern AI tooling (agents, MCP, RAG, evals).

## 2. Tech Stack
* **Backend**: Go (version 1.26.2)
* **Frontend**: React + Next.js

## 3. Git Workflow
* **Base Branch**: All development work must start from the `develop` branch.
* **Branch Naming**: Branch names must always start with `feat/` (for features) or `fix/` (for bug fixes).
* **Pre-commit Verification**: Run tests and linters (e.g., `go test ./...` or frontend checks) before committing to ensure everything is correct.
* **Commit History**: Squash commits to maintain a clean history before merging.
* **Commit Messages**: Use Conventional Commits format (e.g., `feat: add something`, `fix: fix something`, `docs: update documentation`).

## 4. Default Guidelines
* **Code Style**: Idiomatic Go, formatted using `gofmt` / `goimports`.
* **Unit Testing**: Write unit tests for new features and critical business logic (e.g., using `go test` for Go backend, and appropriate test runner for frontend).
* **Environment Variables**: Never commit secrets or credentials. Always document new environment variables in `.env.example` or equivalent templates.
* **Security & Execution**: Always ask for permission before running potentially destructive commands or infrastructure changes.
* **Language**: All communication in Brazilian Portuguese (PT-BR), but all code, documentation, commits, and comments must be in English.

## 5. Behavioral Rules
1. **Scope Limit**: Never refactor code that is outside the scope of the requested task.
2. **File Modification Preference**: Prefer editing existing files rather than creating new files.
3. **Dependencies**: Do not install any new package or library without asking for permission first.
4. **Out-of-scope Bugs**: If you find any bug or issue outside the scope of your current task, point it out to the user but do not fix it without explicit approval.
5. **Code Review & Approval**: Always present the code diff or new code contents in the output so the user can review and approve changes prior to committing.

