# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Backend (Go)
```bash
# Run server (default: postgres://localhost/dischord?sslmode=disable)
cd backend && go run main.go

# Run all tests (requires a separate test database)
cd backend && TEST_DATABASE_URL="postgres://localhost/dischord_test?sslmode=disable" go test ./...

# Run a specific test package
cd backend && TEST_DATABASE_URL="..." go test -v ./handlers
cd backend && TEST_DATABASE_URL="..." go test -v ./integration_tests

# Run a single test function
cd backend && TEST_DATABASE_URL="..." go test -v ./handlers -run TestCreateUser
```

### Frontend (Node/npm)
```bash
cd frontend && npm run dev      # Dev server on :5173
cd frontend && npm run build    # TypeScript check + Vite production build
cd frontend && npm run preview  # Preview production build
```

## Architecture

### Backend

`main.go` → reads `DATABASE_URL` env var, calls `store.Open()`, passes store to `router.New()`, listens on `:8080`.

`store/store.go` — all SQL lives here. `Open()` calls `ApplySchema()` on startup (idempotent DDL). `TruncateAll()` is used in tests. The `GetPost()` query joins the `votes` table and returns an aggregate `votes` field.

`router/router.go` — wires HTTP method+path patterns to handler structs. Uses Go 1.22+ `http.ServeMux` path patterns (`GET /users/{id}`).

`handlers/` — one file per resource. Each handler struct holds a `*store.Store`. `generateID()` produces 32-char random hex. `writeJSON()` writes `Content-Type: application/json`.

`models/models.go` — shared structs used by both handlers and store. `Server.MemberIDs` is populated in memory (not stored); only `post_ids` come from DB.

**No auth layer** — `author_id` / `owner_id` are trusted values from request bodies.

### Frontend

`App.tsx` owns top-level state: current user, list of server IDs, selected server, and active tab. User ID and server IDs are persisted to `localStorage` (`dischord_user_id`, `dischord_server_ids`). On mount, it loads the stored user ID and verifies it with `GET /users/{id}`.

`api/client.ts` — thin `fetch` wrapper. All calls use the base path `/api`, which Vite proxies to `http://localhost:8080` (stripping the `/api` prefix).

`components/Sidebar.tsx` — friends list, servers list, and user panel. Add friend / join server by pasting an ID.

`components/ServerView.tsx` — renders the selected server. Switches between Posts and Messages tabs. Maintains a local user cache to avoid refetching author details. Uses an `isCancelled` flag to prevent stale state updates after unmount.

`components/PostCard.tsx` — displays a post with vote controls (up/neutral/down). Shows edit/delete controls only when `post.author_id === currentUser.id`.

## Database Schema Summary

| Table | Key columns |
|---|---|
| users | id, username, email |
| servers | id, name, owner_id |
| posts | id, server_id, author_id, title, body |
| votes | (post_id, author_id) PK, vote INT |
| friends | (user_id, friend_id) PK — bidirectional rows |
| messages | id, server_id, author_id, content |

The schema is applied automatically via `store.ApplySchema()` on backend startup. The raw SQL is at `store/schema.sql` for reference.

## Testing Conventions

Handler tests use a `testStore(t)` helper (defined in `handlers/testdb_test.go`) that opens a connection to `TEST_DATABASE_URL` and calls `t.Cleanup(store.TruncateAll)` to reset state between tests. Integration tests in `integration_tests/` follow the same pattern with their own `testdb_test.go`.
