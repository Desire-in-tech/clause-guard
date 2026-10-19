# ClauseGuard

ClauseGuard is a small portfolio MVP that gives freelancers and remote workers a mock, plain-language review of pasted contract clauses. It is informational only, not legal advice.

## Stack

- Go, chi, and the official MongoDB Go driver
- React, Vite, TypeScript, Tailwind CSS, and React Router
- Mock AI by default; in-memory report storage unless MongoDB is configured

## Run locally

Install Go and Node.js, then install the frontend packages:

```bash
cd frontend
npm install
cd ..
cp .env.example .env
./dev.sh
```

Open <http://localhost:5173>. The API health check is at <http://localhost:8080/api/health>. The app runs without `.env` or MongoDB; reports are then kept in memory until the API stops.

## Environment

See `.env.example` for all optional variables. `VITE_API_URL` can point to an API base such as `http://localhost:8080/api`; when empty, the frontend uses Vite's `/api` proxy. `LLM_PROVIDER` defaults to `mock`; `openrouter` is an unimplemented stub. Set `MONGODB_URI` to use MongoDB, and `REPORT_TTL_DAYS` to change the default seven-day report lifetime.

## Checks

```bash
cd backend && go vet ./... && go test ./...
cd ../frontend && npm run build
```

## Structure

- `backend/internal/` — config, HTTP API, clause splitting, analysis, chat, mock LLM, and report stores
- `frontend/src/` — typed API client, pages, and UI components
- `docs/api.md` — endpoint request and response formats
