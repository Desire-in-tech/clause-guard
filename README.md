# ClauseGuard

ClauseGuard is a portfolio MVP that helps freelancers and remote workers review service agreements in plain language. It separates a contract into clauses, assigns each clause a **green**, **yellow**, or **red** risk level, explains potential concerns, and—where appropriate—suggests balanced wording to discuss with the client. Users can also ask questions about the analyzed clauses and review the citations returned by the assistant.

> **Important:** ClauseGuard is an educational/demo tool, not a lawyer or a substitute for legal advice. AI-generated analysis may be incomplete or incorrect. Do not rely on it to make legal decisions, and consult a qualified lawyer about important agreements.

## Live demo

- **Frontend:** [https://clause-guard-frontend-olive.vercel.app](https://clause-guard-frontend-olive.vercel.app)
- **Backend API:** [https://clauseguard-api-z0he.onrender.com](https://clauseguard-api-z0he.onrender.com)
- **API health check:** [https://clauseguard-api-z0he.onrender.com/api/health](https://clauseguard-api-z0he.onrender.com/api/health)

The health endpoint reports whether the API is responding and identifies the configured AI provider and report-storage status. The API root itself is not a homepage; opening the backend root URL may return a not-found response. Use the frontend to interact with the app.

### Try it

1. Open the [live frontend](https://clause-guard-frontend-olive.vercel.app).
2. Paste a short sample agreement into the editor, or upload a supported text-based contract file.
3. Review the extracted text before submitting it for analysis.
4. Read the overall summary and each clause's risk label, explanation, and suggested rewrite (where one is provided).
5. Use the clause Q&A/chat feature to ask a focused question about the analyzed text.

Use a fictional or public sample agreement for testing. **Do not upload confidential client contracts, personal information, credentials, or sensitive business information.** The demo sends contract text to a third-party AI provider when the OpenRouter provider is enabled.

## Features

- **Plain-language contract review** aimed at freelancers and remote workers acting as service providers.
- **Clause-by-clause risk labels:** green (generally fair/standard), yellow (worth clarifying or negotiating), and red (seriously one-sided or potentially high-impact).
- **Reasons and proposed wording:** explanations describe the potential concern; non-green clauses may include a balanced suggested rewrite. Green clauses do not need a rewrite.
- **Questions about analyzed clauses:** chat answers are grounded in the clauses provided to it and include clause-ID citations.
- **File text extraction:** upload `.txt`, `.docx`, or `.pdf` files up to 5 MB. Extracted text is shown for review before analysis.
- **Report retrieval:** reports can be saved in MongoDB and retrieved by report ID. When MongoDB is not configured, the API falls back to in-memory storage.
- **Configurable AI provider:** a mock provider is the default for local development; the deployed demo is configured to use OpenRouter.
- **Basic request safeguards:** configurable maximum clause count, provider-request budget, and per-IP hourly limits for AI analysis/chat.

## Technology stack

| Layer | Technologies |
| --- | --- |
| Frontend | React, TypeScript, Vite, Tailwind CSS, React Router |
| Backend | Go, Chi HTTP router |
| AI integration | OpenRouter-compatible chat completions API; mock provider for development/testing |
| Report storage | MongoDB via the official Go driver; in-memory fallback |
| Frontend hosting | Vercel |
| Backend hosting | Render |
| Version control | Git and GitHub |

## How it works

1. **Input:** The user pastes contract text or uploads a supported file.
2. **Extraction and segmentation:** Uploaded file text is extracted by the API; contract text is split into clauses.
3. **Analysis:** The backend asks the configured provider to classify clauses and explain risks from the service provider's perspective.
4. **Validation:** The backend validates the structured model response, ensures it maps to the submitted clause IDs, and derives the overall risk from the clause risks.
5. **Results and follow-up:** The frontend displays the report and can send clause-grounded follow-up questions to the chat endpoint.
6. **Storage:** When MongoDB is available, reports are saved with a configurable retention period (seven days by default). Without MongoDB, reports are held in memory and disappear when the API process restarts.

## Risk labels and supported clause categories

| Label | Meaning in this demo |
| --- | --- |
| **Green** | Generally fair and standard; no rewrite is suggested by the AI reviewer. |
| **Yellow** | Unfavorable, vague, or worth clarifying or negotiating. |
| **Red** | Potentially seriously one-sided or high impact; review carefully before agreeing. |

Supported clause categories are payment, intellectual property, termination, liability, confidentiality, non-compete, jurisdiction, scope of work, revisions, and other.

These labels are guidance for prioritizing human review, not legal determinations. A label does not establish whether a term is valid or enforceable in a particular jurisdiction.

## API overview

The API is hosted at [https://clauseguard-api-z0he.onrender.com](https://clauseguard-api-z0he.onrender.com). API routes are prefixed with `/api`.

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `GET` | `/api/health` | Health and configured provider/storage status |
| `POST` | `/api/analyze` | Analyze contract text and return a report |
| `GET` | `/api/reports/{id}` | Retrieve a saved report by ID |
| `POST` | `/api/chat` | Ask a question about provided clauses |
| `POST` | `/api/extract` | Extract text from an uploaded `.txt`, `.docx`, or `.pdf` file |

For request and response examples, validation rules, limits, and error formats, see [API documentation](docs/api.md).

### Health check example

Open the [live health endpoint](https://clauseguard-api-z0he.onrender.com/api/health) in a browser. A typical healthy response includes `"status": "ok"`, the configured `llm_provider`, and `store` status. The exact provider/storage values depend on deployment configuration and dependency availability.

### Example analysis request

```bash
curl -X POST "https://clauseguard-api-z0he.onrender.com/api/analyze" \
  -H "Content-Type: application/json" \
  -d '{
    "contract_text": "Payment is due within 30 days of receipt of invoice.\n\nEither party may terminate this agreement with 14 days written notice."
  }'
```

The request returns a report ID, summary, overall risk, and an array of analyzed clauses. Clause text is split using blank lines. The default maximum is 40 clauses and can be changed with `MAX_CLAUSES`.

## Run locally

### Prerequisites

- Git
- Go version compatible with `backend/go.mod` (currently Go 1.26)
- Node.js and npm compatible with the Vite project

### 1. Clone the repository

```bash
git clone https://github.com/Desire-in-tech/clause-guard.git
cd clause-guard
```

### 2. Install frontend dependencies

```bash
cd frontend
npm install
cd ..
```

### 3. Configure environment variables (optional)

Copy the example environment file:

```bash
cp .env.example .env
```

The project can run without a `.env` file, an OpenRouter API key, or MongoDB: it defaults to the mock AI provider and in-memory report storage. To use real AI analysis locally, configure the OpenRouter and other variables described below. Keep `.env` private and never commit API keys or database credentials.

### 4. Start the backend and frontend

From the repository root:

```bash
./dev.sh
```

Then open [http://localhost:5173](http://localhost:5173). The local API health check is [http://localhost:8080/api/health](http://localhost:8080/api/health). The Vite development server proxies `/api` requests to the local Go API when `VITE_API_URL` is empty.

If the script is not executable in your environment, run `bash dev.sh`.

## Environment variables

The root [`.env.example`](.env.example) lists the supported variables. Set values in the deployment platform's environment settings rather than committing secrets.

| Variable | Purpose | Default / notes |
| --- | --- | --- |
| `PORT` | Port used by the Go HTTP server | `8080` locally; Render supplies the production value |
| `FRONTEND_ORIGIN` | Allowed frontend origin for CORS | Empty locally; set to the exact Vercel origin in production, without a trailing slash or path |
| `LLM_PROVIDER` | AI provider selection | `mock`; set to `openrouter` for live AI analysis |
| `OPENROUTER_API_KEY` | OpenRouter credential | Required when `LLM_PROVIDER=openrouter`; keep secret |
| `MODEL_ID` | Primary OpenRouter model ID | Required when `LLM_PROVIDER=openrouter` |
| `MODEL_FALLBACKS` | Optional comma-separated fallback model IDs | Empty by default |
| `LLM_TIMEOUT_SECONDS` | Timeout for an individual provider attempt | `45` seconds by default; overall request deadline is 100 seconds |
| `MAX_CLAUSES` | Maximum clauses accepted for analysis | `40` |
| `LLM_DAILY_REQUEST_BUDGET` | Maximum real-provider requests per UTC day | `40` |
| `AI_RATE_LIMIT_PER_HOUR` | Real-provider requests per IP per hour | `10`; applies to analysis and chat |
| `MONGODB_URI` | MongoDB connection string | Empty means in-memory report storage |
| `MONGODB_DB` | MongoDB database name | `clauseguard` |
| `REPORT_TTL_DAYS` | Saved report lifetime | `7` days |
| `VITE_API_URL` | Frontend API base URL | Empty for local Vite proxy; for production use the Render API base ending in `/api` |

Request budgets and rate limits are in-memory controls, not a complete abuse-prevention system. They reset when the backend process restarts and are not shared across multiple backend instances.

## Deployment configuration

The live demo uses separate frontend and backend deployments:

- **Vercel frontend:** [clause-guard-frontend-olive.vercel.app](https://clause-guard-frontend-olive.vercel.app)
- **Render backend:** [clauseguard-api-z0he.onrender.com](https://clauseguard-api-z0he.onrender.com)

For a deployment of your own:

### Vercel

- Connect the GitHub repository and configure the project root as `frontend`.
- Build command: `npm run build`.
- Output directory: `dist`.
- Set `VITE_API_URL` to your backend API base, for example `https://YOUR-RENDER-SERVICE.onrender.com/api`.
- The included [`frontend/vercel.json`](frontend/vercel.json) rewrites client-side routes to `index.html`.

Vite environment variables are embedded at build time. Redeploy the frontend after changing `VITE_API_URL`.

### Render

- Connect the same GitHub repository and set the root directory to `backend`.
- Build command: `go build -o app ./cmd/api`.
- Start command: `./app`.
- Health-check path: `/api/health`.
- Set `FRONTEND_ORIGIN` to the exact deployed frontend origin, for example `https://your-project.vercel.app`.
- Configure `LLM_PROVIDER`, `OPENROUTER_API_KEY`, `MODEL_ID`, and `MONGODB_URI` through Render's environment settings for the intended production behavior.

Do not put secrets in `VITE_*` variables: frontend values are public to visitors. Keep OpenRouter and MongoDB credentials exclusively in backend environment settings.

Free hosting tiers can have cold starts, resource limits, or intermittent availability. A successful health check confirms that the API responds; it does not guarantee that the AI provider or database will always be available.

## File uploads and data handling

- Supported file extensions: `.txt`, `.docx`, and `.pdf`.
- Maximum upload size: 5 MB.
- Scanned/image-only PDFs are not supported because PDF extraction reads selectable text and does not perform OCR.
- Uploaded files are processed in memory and are not saved as uploaded files by the API.
- Contract text is sent to the configured AI provider for analysis/chat when OpenRouter is enabled.
- Saved reports may contain contract text and analysis. When MongoDB is configured, report retention defaults to seven days; without MongoDB, reports are stored in process memory only.
- The API does not implement user accounts or authentication. Do not treat report IDs or the public demo as a secure private document-management system.

Use fictional, public, or appropriately redacted sample contracts when trying the public demo.

## Development checks

Run the backend checks:

```bash
cd backend
go vet ./...
go test ./...
```

Build the frontend:

```bash
cd frontend
npm run build
```

## Repository structure

```text
clause-guard/
├── backend/
│   ├── cmd/api/             # Go API entry point
│   └── internal/
│       ├── analysis/        # Analysis service and risk/result models
│       ├── chat/            # Clause-grounded Q&A
│       ├── config/          # Environment-based configuration
│       ├── contracts/       # Contract clause structures
│       ├── httpapi/         # HTTP routes, handlers, CORS, limits
│       ├── llm/              # Mock and OpenRouter providers
│       └── store/            # MongoDB and in-memory report stores
├── frontend/
│   └── src/
│       ├── components/      # Reusable UI and clause cards
│       ├── lib/             # API integration and frontend helpers
│       └── pages/           # Main and report/results pages
├── docs/
│   └── api.md               # Endpoint details and payload examples
├── .env.example             # Documented environment-variable template
├── dev.sh                  # Starts the local API and Vite dev server
└── README.md
```

## Scope and limitations

ClauseGuard is intentionally a small portfolio demonstration, not a production legal platform. It has no user accounts, authentication, payment flow, or guarantee of legal accuracy. It does not determine whether a clause is enforceable under a particular country's law. Uploaded files are not OCR-processed, and the application does not support scanned/image-only PDFs. The public deployment may be subject to free-tier cold starts and provider request limits.

## License

No license is specified in this repository at the time of writing. Unless a license is added, do not assume that the source code is available for unrestricted reuse.
