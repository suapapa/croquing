# Croquis King

Real-time web app for weekly croquis (figure-drawing) meetups. Participants join the same lobby URL and see synchronized photos and timers.

## Stack

| Layer | Tech |
|-------|------|
| Backend | Go (`cmd/server`) |
| Frontend | React, Vite, TypeScript (`frontend/`) |

See [docs/PROJECT_PLAN.md](docs/PROJECT_PLAN.md) for architecture and WorkItems. Implementation progress: [docs/progress/PROGRESS.md](docs/progress/PROGRESS.md).

## Backend

Required environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PIXABAY_API_KEY` | yes | — | PixaBay API key |
| `PORT` | | `8080` | HTTP port |
| `CORS_ORIGINS` | | `*` | Allowed origins (e.g. `http://localhost:5173`) |
| `DRAW_DURATION` | | `5m` | Draw round duration |

```sh
export PIXABAY_API_KEY=your-key
go run ./cmd/server
```

## Frontend

```sh
cd frontend
npm install
cp .env.example .env   # optional
npm run dev
```

Set `VITE_API_BASE` in `frontend/.env` if the backend is not on `http://localhost:8080`.

## Development commands

```sh
make progress    # regenerate docs/progress/PROGRESS.md
go test -race ./...
cd frontend && npm run lint && npm run build
```
