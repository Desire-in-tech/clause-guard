export type Risk = "green" | "yellow" | "red";

export type ClauseType =
  | "payment"
  | "intellectual_property"
  | "termination"
  | "liability"
  | "confidentiality"
  | "non_compete"
  | "jurisdiction"
  | "scope_of_work"
  | "revisions"
  | "other";

export interface ClauseAnalysis {
  id: string;
  index: number;
  text: string;
  type: ClauseType;
  risk: Risk;
  reason: string;
  suggested_rewrite: string;
}

export interface AnalysisReport {
  id: string;
  summary: string;
  overall_risk: Risk;
  clauses: ClauseAnalysis[];
}

export interface ChatClause {
  id: string;
  text: string;
}

export interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

export interface ChatRequest {
  question: string;
  clauses: ChatClause[];
  history: ChatMessage[];
}

export interface ChatResponse {
  answer: string;
  citations: string[];
}

export interface APIError {
  error: string;
}

export interface ExtractResponse {
  text: string;
  filename: string;
  characters: number;
}
