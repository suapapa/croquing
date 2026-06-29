# Croquing

Real-time web app for weekly croquis (figure-drawing) meetups. Participants join the same lobby URL and see synchronized photos and timers.

Project Intro & Demo: [https://croquing.homin.dev/](https://croquing.homin.dev/)

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
| `PIXABAY_API_KEY` | yes | — | Pixabay API key |
| `PORT` | | `8080` | HTTP port |
| `CORS_ORIGINS` | | `*` | Allowed origins (e.g. `http://localhost:5173`) |
| `DRAW_DURATION` | | `5m` | Draw round duration |
| `APP_NAME` | | — | App name shown on the home screen (omitted if empty) |
| `APP_LOGO` | | — | Custom logo image URL (falls back to `/example_logo.png`) |
| `APP_LOGO_LINK` | | `https://homin.dev` | Link to open when clicking the logo |

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

### Production (single server)

Build the SPA and run the Go server from the repo root. When `frontend/dist/` exists, the backend serves the React app and falls back to `index.html` for client-side routes such as `/lobby/:id`.

```sh
make web
export PIXABAY_API_KEY=your-key
./bin/croquing
```

Open `http://localhost:8080`. API and WebSocket stay on the same origin, so you do not need a separate `VITE_API_BASE` in production.

### Development (split origins)

Run the backend and Vite dev server separately:

```sh
# terminal 1
export PIXABAY_API_KEY=your-key
export CORS_ORIGINS=http://localhost:5173
go run ./cmd/server

# terminal 2
cd frontend && npm run dev
```

### Reverse proxy (optional)

For TLS or multiple services, terminate HTTP at nginx/Caddy and proxy:

- `/api/*` and `/ws/*` → Go backend
- `/` and `/lobby/*` → `frontend/dist` static files or the Go server when using `make web`

Example nginx location blocks:

```nginx
location /api/ { proxy_pass http://127.0.0.1:8080; }
location /ws/  { proxy_pass http://127.0.0.1:8080; proxy_http_version 1.1; proxy_set_header Upgrade $http_upgrade; proxy_set_header Connection "upgrade"; }
location /     { try_files $uri $uri/ /index.html; root /path/to/croquing/frontend/dist; }
```

## Development commands

```sh
make progress    # regenerate docs/progress/PROGRESS.md
go test -race ./...
cd frontend && npm run lint && npm run build
```
