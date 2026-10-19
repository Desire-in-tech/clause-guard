import { useState, type FormEvent } from "react";
import type {
  AnalysisReport,
  ChatMessage as APIChatMessage,
} from "../api/types";
import { chatAboutContract } from "../api/client";
import ServerWakeBanner from "./ServerWakeBanner";
import useServerWake from "../hooks/useServerWake";

interface ChatPanelProps {
  clauses: AnalysisReport["clauses"];
}

interface DisplayMessage extends APIChatMessage {
  citations?: string[];
}

export default function ChatPanel({ clauses }: ChatPanelProps) {
  const [question, setQuestion] = useState("");
  const [messages, setMessages] = useState<DisplayMessage[]>([]);
  const [error, setError] = useState("");
  const [isSending, setIsSending] = useState(false);
  const { isWaking, runWithWake } = useServerWake();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const submittedQuestion = question.trim();
    if (!submittedQuestion || isSending) return;

    setError("");
    setIsSending(true);
    try {
      const response = await runWithWake(
        chatAboutContract({
          question: submittedQuestion,
          clauses: clauses.map(({ id, text }) => ({ id, text })),
          history: messages.map(({ role, content }) => ({ role, content })),
        }),
      );
      setMessages((current) => [
        ...current,
        { role: "user", content: submittedQuestion },
        {
          role: "assistant",
          content: response.answer,
          citations: response.citations,
        },
      ]);
      setQuestion("");
    } catch (requestError) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Your question could not be sent. Please try again.",
      );
    } finally {
      setIsSending(false);
    }
  }

  return (
    <section className="border border-black/15 bg-white" aria-labelledby="chat-title">
      <div className="border-b border-black/15 px-5 py-4">
        <p className="mb-1 text-xs font-bold uppercase tracking-[0.14em] text-primary">
          Ask about this contract
        </p>
        <h2 id="chat-title" className="text-lg font-semibold tracking-tight">
          Questions, with context
        </h2>
      </div>
      <div className="space-y-4 px-5 py-4">
        <ServerWakeBanner visible={isWaking} />
        {messages.length === 0 ? (
          <p className="text-sm leading-6 text-black/65">
            Ask a question about a clause or the terms as a whole.
          </p>
        ) : (
          <ol className="space-y-4" aria-live="polite">
            {messages.map((message, index) => (
              <li
                key={`${message.role}-${index}`}
                className={
                  message.role === "user"
                    ? "ml-5 border-l-2 border-black/20 pl-3"
                    : "mr-2 border-l-2 border-primary pl-3"
                }
                data-testid={`message-chat-${index}`}
              >
                <p className="mb-1 text-[11px] font-bold uppercase tracking-[0.12em] text-black/50">
                  {message.role === "user" ? "You" : "ClauseGuard"}
                </p>
                <p className="whitespace-pre-wrap text-sm leading-6 text-black">
                  {message.content}
                </p>
                {message.citations && message.citations.length > 0 && (
                  <p className="mt-2 flex flex-wrap gap-x-3 gap-y-1 text-xs text-black/65">
                    <span>Cited clauses:</span>
                    {message.citations.map((citation) => {
                      const citedClause = clauses.find(
                        (clause) => clause.id === citation,
                      );
                      const label = citedClause
                        ? `Clause ${citedClause.index}`
                        : `Clause ${citation}`;

                      return (
                        <a
                          key={citation}
                          href={`#clause-${citation}`}
                          data-testid={`link-citation-${citation}`}
                        >
                          {label}
                        </a>
                      );
                    })}
                  </p>
                )}
              </li>
            ))}
          </ol>
        )}
        {error && (
          <p
            className="border-l-2 border-[#bb2634] bg-[#fff0f1] px-3 py-2 text-sm text-[#8d1a25]"
            role="alert"
            data-testid="error-chat"
          >
            {error}
          </p>
        )}
        <form onSubmit={handleSubmit} className="space-y-3">
          <label
            htmlFor="contract-question"
            className="block text-xs font-semibold uppercase tracking-[0.1em] text-black/60"
          >
            Your question
          </label>
          <textarea
            id="contract-question"
            value={question}
            onChange={(event) => setQuestion(event.target.value)}
            placeholder="What does the payment clause mean?"
            rows={3}
            className="w-full resize-y border border-black/25 bg-white px-3 py-3 text-sm leading-6 text-black placeholder:text-black/40 focus:border-primary focus:outline-none"
            disabled={isSending}
            data-testid="input-chat-question"
          />
          <button
            type="submit"
            disabled={!question.trim() || isSending}
            className="w-full bg-primary px-4 py-3 text-sm font-semibold text-white hover:bg-[#7040f0] disabled:cursor-not-allowed disabled:bg-black/20"
            data-testid="button-send-question"
          >
            {isSending ? "Sending…" : "Ask question"}
          </button>
        </form>
      </div>
    </section>
  );
}
