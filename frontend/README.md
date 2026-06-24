# Croquing — Frontend

React SPA (Vite + TypeScript) for the Croquing lobby UI.

## Setup

```sh
cd frontend
npm install
cp .env.example .env   # optional; defaults to http://localhost:8080
```

## Development

Run the Go backend from the repo root, then start the dev server:

```sh
# repo root
go run ./cmd/server

# another terminal
cd frontend && npm run dev
```

Vite serves at `http://localhost:5173`. Set backend `CORS_ORIGINS=http://localhost:5173` if not using `*`.

## Environment

| Variable        | Default                 | Description         |
| --------------- | ----------------------- | ------------------- |
| `VITE_API_BASE` | `http://localhost:8080` | Backend HTTP origin |

## Scripts

| Command                | Description                 |
| ---------------------- | --------------------------- |
| `npm run dev`          | Vite dev server             |
| `npm run build`        | Production build to `dist/` |
| `npm run lint`         | ESLint                      |
| `npm run format`       | Prettier write              |
| `npm run format:check` | Prettier check              |
