# ClauseGuard

ClauseGuard is a small portfolio MVP that gives freelancers and remote workers a plain-language review of pasted contract clauses or text extracted from uploaded files. It uses a mock reviewer by default and can optionally use OpenRouter. It is informational only, not legal advice.

## Stack

- Go, chi, and the official MongoDB Go driver
- React, Vite, TypeScript, Tailwind CSS, and React Router
- Mock AI by default; optional OpenRouter provider; in-memory report storage unless MongoDB is configured

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

You can upload `.txt`, `.docx`, or `.pdf` files up to 5 MB. Extracted text appears in the editor for review before analysis. Scanned/image-only PDFs are not supported. Files are processed in memory and are not stored.

## Environment

See `.env.example` for the full list of variables. `VITE_API_URL` can point to an API base such as `http://localhost:8080/api`; when empty, the frontend uses Vite's `/api` proxy.

`LLM_PROVIDER` defaults to `mock`, so the app starts without an OpenRouter key. To use the real provider, set `LLM_PROVIDER=openrouter`, `OPENROUTER_API_KEY`, and `MODEL_ID`. `MODEL_FALLBACKS` is an optional comma-separated list of fallback model IDs. `LLM_TIMEOUT_SECONDS` defaults to 45 seconds per attempt; a complete request has a 100-second deadline. `MAX_CLAUSES` defaults to 40.

The OpenRouter demo allows up to `LLM_DAILY_REQUEST_BUDGET` real-provider requests per UTC day (default 40) and `AI_RATE_LIMIT_PER_HOUR` requests per IP per hour (default 10). The per-IP limit applies to `/api/analyze` and `/api/chat`; mock-provider requests do not use these limits. Set `MONGODB_URI` to use MongoDB, and `REPORT_TTL_DAYS` to change the default seven-day report lifetime.

OpenRouter uses free third-party AI models in this demo. Do not paste confidential or personal information. AI-generated analysis can be wrong and is informational only, not legal advice.

## Checks

```bash
cd backend && go vet ./... && go test ./...
cd ../frontend && npm run build
```

## Structure

- `backend/internal/` — config, HTTP API, clause splitting, analysis, chat, mock LLM, and report stores
- `frontend/src/` — typed API client, pages, and UI components
- `docs/api.md` — endpoint request and response formats
