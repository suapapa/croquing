# AGENTS.md — Croquis King

Guidance for coding agents working in this repository.

## Primary planning document

**Before implementing features, read:**

- **[docs/PROJECT_PLAN.md](docs/PROJECT_PLAN.md)** — product goals, architecture, API sketch, session state machine, and **indexed WorkItems** (001–017 backend, 101–112 frontend).
- **[docs/progress/PROGRESS.md](docs/progress/PROGRESS.md)** — **implementation progress** (done/pending, artifacts, next candidates). Regenerate with `make progress` after editing [`docs/progress/workitems.json`](docs/progress/workitems.json).

When the user asks to implement a specific index (e.g. “003 구현해줘”), open that document, locate the matching WorkItem row, and implement only that scope unless dependencies require otherwise.

## What this is

**Croquis King** is a real-time web app for weekly croquis (5-minute figure-drawing) meetups. Participants join the same lobby URL and see synchronized photos and timers—no video screen sharing.

| Layer | Stack |
|-------|--------|
| Backend | Go |
| Frontend | React (Vite + TypeScript), under `frontend/` — **after backend WorkItems 001–017** |

Module path: `github.com/suapapa/croquis-king`.

### Core behaviors (from plan)

- Admin opens a lobby; joiners see the same page state via **WebSocket**.
- **PixaBay** search/selection through the backend ([API docs](https://pixabay.com/api/docs/)); API key stays server-side.
- Photo order is **shuffled and hidden** until the session starts.
- During each **5-minute** round, show the photo as large as possible with a thin top progress bar and remaining time; hide the photo when time ends.
- **Server-authoritative timer** (`DrawEndsAt`) to prevent drawing before the official start.

## Implementation order

1. **Backend first:** WorkItems **001 → 017** (see [docs/PROJECT_PLAN.md §6–7](docs/PROJECT_PLAN.md#6-workitem-목록)).
2. **Frontend second:** WorkItems **101 → 112** only after backend is done.

Suggested MVP bundles are listed in the plan (§7).

## Target layout

```
croquis-king/
├── cmd/server/           # Go entrypoint (WorkItem 001)
├── internal/
│   ├── config/
│   ├── lobby/            # domain, state machine, store
│   ├── pixabay/
│   ├── timer/
│   ├── ws/
│   └── http/
├── frontend/             # React SPA (WorkItem 101+)
├── docs/
│   ├── PROJECT_PLAN.md   # canonical spec & WorkItem index
│   └── progress/
│       ├── workitems.json  # progress source of truth (edit this)
│       └── PROGRESS.md       # generated dashboard (make progress)
├── go.mod
└── AGENTS.md
```

Current state: **WorkItems 001–007, 009–010 done** — lobby API, PixaBay search, WebSocket hub. See [docs/progress/PROGRESS.md](docs/progress/PROGRESS.md). Next: **008** (real-time snapshot sync).

## Conventions for agents

- **Scope:** Prefer the smallest diff that satisfies the requested WorkItem. Do not implement frontend tasks (101+) while backend items are still open unless explicitly asked.
- **Go:** Follow idiomatic Go. See `.agents/skills/golang-pro/SKILL.md` for project-local expectations (vet, tests with `-race`, layout under `internal/`).
- **State:** Lobby phase and timer are owned by the **server**; clients render snapshots only.
- **Secrets:** `PIXABAY_API_KEY` and admin tokens never belong in frontend code or committed env files.
- **PixaBay:** Show attribution when displaying search results (API guideline).

## Commands (from repo root)

```sh
go run ./cmd/server
make progress                 # regenerate docs/progress/PROGRESS.md
go vet ./...
go test -race ./...
```

Frontend (after WorkItem 101):

```sh
cd frontend && npm install && npm run dev
```

## Documentation hygiene

When a change alters **APIs**, **lobby phases**, **directory layout**, **env vars**, or **WorkItem scope/completion**, update:

- **`docs/PROJECT_PLAN.md`** — if the spec or WorkItem definitions change.
- **`docs/progress/workitems.json`** + **`make progress`** — when a WorkItem status or artifacts change.
- **`AGENTS.md`** — if agent-facing conventions or commands change.
- **`README.md`** — if user-facing run/build steps change (create or update when those steps exist).

Treat doc updates as part of the same change, not optional follow-up.

## WorkItem quick reference

| Range | Phase |
|-------|--------|
| 001–017 | Backend: scaffold → HTTP/WS → PixaBay → session & timer → integration tests |
| 101–112 | Frontend: React → lobby UI → search/select → drawing screen with timer |

Full descriptions and dependencies: **[docs/PROJECT_PLAN.md §6](docs/PROJECT_PLAN.md#6-workitem-목록)**.
