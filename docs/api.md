# ClauseGuard API

The API listens on port `8080` by default. JSON errors use the form `{"error":"message"}`. JSON requests are limited to 1 MiB; uploads are limited to 5 MB; contract and extracted text are limited to 100,000 characters. No authentication is used in this portfolio MVP.

## `GET /api/health`

Returns:

```json
{
  "status": "ok",
  "llm_provider": "mock",
  "store": "memory"
}
```

`store` is `mongodb`, `memory`, or `unavailable`.

## `POST /api/analyze`

Request:

```json
{ "contract_text": "First clause.\n\nSecond clause." }
```

Clauses are split at blank lines. The default maximum is 40 clauses; configure this with `MAX_CLAUSES`. Longer contracts return `400` with a message naming the configured maximum. With OpenRouter enabled, requests that exceed the daily budget or per-IP hourly limit return `429`. If the provider is temporarily unavailable after its retries, the API returns `502`:

```json
{ "error": "Demo limit reached for today. Please try again tomorrow." }
```

Provider failures return:

```json
{ "error": "AI analysis is temporarily unavailable. Please try again later." }
```

Response:

```json
{
  "id": "URL-safe-report-id",
  "summary": "Plain-language mock summary.",
  "overall_risk": "yellow",
  "clauses": [
    {
      "id": "clause-1",
      "index": 1,
      "text": "First clause.",
      "type": "other",
      "risk": "green",
      "reason": "No common risk keyword was detected by the mock reviewer.",
      "suggested_rewrite": "Clarify the parties’ expectations in plain language."
    }
  ]
}
```

Risk values: `green`, `yellow`, `red`. Clause types: `payment`, `intellectual_property`, `termination`, `liability`, `confidentiality`, `non_compete`, `jurisdiction`, `scope_of_work`, `revisions`, `other`. If report storage fails, analysis is still returned with an empty `id`.

## `GET /api/reports/{id}`

Returns the saved report in the same shape as `/api/analyze`. Missing and expired reports return `404`:

```json
{ "error": "This report is missing or has expired." }
```

## `POST /api/chat`

Request:

```json
{
  "question": "What does the payment clause say?",
  "clauses": [{ "id": "clause-1", "text": "Payment is due within 30 days." }],
  "history": [{ "role": "user", "content": "What does this say?" }]
}
```

History roles are `user` and `assistant`. Response:

```json
{
  "answer": "The mock answer references the matching clause.",
  "citations": ["clause-1"]
}
```

Questions are limited to 1,000 characters, and only the last six history messages are sent to the AI provider. Provider rate limits and the configured daily/IP limits return `429`; temporary provider failures return `502`.

## `POST /api/extract`

Upload one file as multipart form data in the `file` field. Supported extensions are `.txt`, `.docx`, and `.pdf`; PDF and DOCX contents are checked against their file signatures. DOCX extraction reads `word/document.xml` and limits decompressed XML to 5 MB. PDF extraction reads selectable text only and does not use OCR.

Successful response:

```json
{
  "text": "Payment is due within 30 days.",
  "filename": "agreement.txt",
  "characters": 30
}
```

Unsupported extensions return `415`; malformed, empty, unreadable, or overlong extracted text returns a JSON error. Uploads over 5 MB return `413`. Scanned or image-only PDFs return:

```json
{ "error": "No selectable text found. Scanned PDFs are not supported yet." }
```

Uploaded files are processed in memory and are not stored.

## CORS and storage

When `FRONTEND_ORIGIN` is empty, CORS allows all origins for local development. Set it to one exact origin to restrict access. When `MONGODB_URI` is empty, reports use in-memory storage. MongoDB storage creates a TTL index on `created_at`; the default expiry is seven days and can be changed with `REPORT_TTL_DAYS`.

## OpenRouter configuration and privacy

`LLM_PROVIDER` defaults to `mock`. To enable OpenRouter, set `LLM_PROVIDER=openrouter`, `OPENROUTER_API_KEY`, and `MODEL_ID`; the API exits at startup if the key or primary model is missing. `MODEL_FALLBACKS` is an optional comma-separated list. Each provider attempt has a `LLM_TIMEOUT_SECONDS` timeout (default 45), with a 100-second overall request deadline. `MAX_CLAUSES` defaults to 40, `LLM_DAILY_REQUEST_BUDGET` to 40 requests per UTC day, and `AI_RATE_LIMIT_PER_HOUR` to 10 per IP per hour. These in-memory request limits apply only to the real provider.

Contract clauses are treated as untrusted prompt data. The demo uses free third-party AI models; do not send confidential or personal information. AI-generated analysis can be wrong and is not legal advice.
