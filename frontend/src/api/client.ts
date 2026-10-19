import type {
  AnalysisReport,
  APIError,
  ChatRequest,
  ChatResponse,
} from "./types";

const API_BASE_URL = (import.meta.env.VITE_API_URL || "/api").replace(/\/$/, "");

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
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
    throw new Error(message);
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

export function extractFile(): Promise<never> {
  return request<never>("/extract", { method: "POST" });
}
