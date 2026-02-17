# AI Agent Guidelines for Training Organiser

This document provides context, commands, and standards for AI agents working on this codebase.

## Project Overview
This is a monolithic full-stack application with:
- **Backend:** Go (Golang) using Chi router, PostgreSQL (sqlc + goose), and Swaggo.
- **Frontend:** Server-side rendered HTML using [Templ](https://templ.guide) and [HTMX](https://htmx.org) for interactivity.
- **Styling:** Tailwind CSS (via CDN or built asset).

## 1. Build, Test, & Lint Commands

### Root Directory
- **Run Server:** `go run cmd/api-server/main.go`
- **Generate Templates:** `templ generate` (Run this after modifying `.templ` files)
  - *Watch Mode:* `templ generate --watch`
- **Database Generation:** `sqlc generate` (Run after modifying `sql/queries` or `schema`)
- **Database Migrations:** `goose postgres <DB_URL> up` (Migration files in `sql/schema`)
- **API Documentation:** `swag init -g cmd/api-server/main.go` (Run after modifying API annotations)
- **Test:** `go test ./...`
  - *Agent Note:* Use standard Go testing package `testing`.
- **Lint:** Use `go vet ./...` or `staticcheck ./...` (if available).

## 2. Code Style & Conventions

### Backend (Go)
- **Formatting:** Always run `gofmt`.
- **Structure:**
  - `cmd/api-server/`: Entry point.
  - `internal/`: Private application code.
  - `internal/database/`: SQLC generated code. **DO NOT EDIT MANUALLY.**
  - `internal/api/`: HTTP Handlers and routing.
  - `internal/views/`: Templ components (`.templ` files).
- **Naming:**
  - Handlers: `Handle<Action><Resource>` (e.g., `HandleRegisterUser`).
  - Methods hang off `*ApiConfig` struct to access DB context.
- **Error Handling:**
  - Use `respondWithError(w, code, message, err)` helper.
  - Return early after errors.

### Frontend (Templ + HTMX)
- **Components:** Place in `internal/views/`. Use PascalCase for component names.
- **Interactivity:** Use HTMX attributes (e.g., `hx-get`, `hx-post`, `hx-target`) for dynamic behavior.
- **Styling:** Use Tailwind CSS utility classes directly in elements.

### General
- **Environment Variables:**
  - Copy `.env.example` to `.env`.
- **Git:** Commit messages should be concise and descriptive.

## 3. Rules & Context
- **No Tests:** The project currently lacks tests. Be cautious when refactoring.
- **Generated Code:** Never modify files in `internal/database/`, `internal/views/*_templ.go`, or `docs/`. Modify the source (`sql/`, `.templ`, annotations) and regenerate.
