import type {
  AnalysisReport,
  APIError,
  ChatRequest,
  ChatResponse,
  ExtractResponse,
} from "./types";

const API_BASE_URL = (import.meta.env.VITE_API_URL || "/api").replace(/\/$/, "");

export class APIRequestError extends Error {
  readonly status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "APIRequestError";
    this.status = status;
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers);
  if (!(init?.body instanceof FormData) && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers,
  });

  const data = (await response.json().catch(() => null)) as
    | T
    | APIError
    | null;

  if (!response.ok) {
    const apiError =
      data && typeof data === "object" && "error" in data
        ? (data as APIError)
        : null;
    const message =
      apiError
        ? apiError.error
        : `Request failed with status ${response.status}.`;
    throw new APIRequestError(message, response.status);
  }

  if (data === null) {
    throw new Error("The server returned an empty response.");
  }

  return data as T;
}

export function analyzeContract(contractText: string): Promise<AnalysisReport> {
  return request<AnalysisReport>("/analyze", {
    method: "POST",
    body: JSON.stringify({ contract_text: contractText }),
  });
}

export function getReport(id: string): Promise<AnalysisReport> {
  return request<AnalysisReport>(`/reports/${encodeURIComponent(id)}`);
}

export function chatAboutContract(
  body: ChatRequest,
): Promise<ChatResponse> {
  return request<ChatResponse>("/chat", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export function extractFile(file: File): Promise<ExtractResponse> {
  const body = new FormData();
  body.append("file", file);
  return request<ExtractResponse>("/extract", {
    method: "POST",
    body,
  });
}
