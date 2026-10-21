import { useEffect, useMemo, useState } from "react";
import { Link, useLocation, useParams } from "react-router-dom";
import { getReport } from "../api/client";
import type { AnalysisReport } from "../api/types";
import ChatPanel from "../components/ChatPanel";
import ClauseCard from "../components/ClauseCard";
import RiskBadge from "../components/RiskBadge";
import ServerWakeBanner from "../components/ServerWakeBanner";
import useServerWake from "../hooks/useServerWake";

interface RouteState {
  report?: AnalysisReport;
}

function isMissingReportError(error: unknown) {
  if (!(error instanceof Error)) return false;
  const message = error.message.toLowerCase();
  return (
    message.includes("404") ||
    message.includes("not found") ||
    message.includes("missing") ||
    message.includes("expired")
  );
}

export default function Results() {
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const routeState = location.state as RouteState | null;
  const unsavedReport = routeState?.report;
  const [report, setReport] = useState<AnalysisReport | null>(null);
  const [isLoading, setIsLoading] = useState(Boolean(id));
  const [error, setError] = useState("");
  const [isMissing, setIsMissing] = useState(false);
  const [copyMessage, setCopyMessage] = useState("");
  const [retryCount, setRetryCount] = useState(0);
  const { isWaking, runWithWake } = useServerWake();

  useEffect(() => {
    if (!id) {
      setReport(unsavedReport ?? null);
      setIsMissing(!unsavedReport);
      setError("");
      setIsLoading(false);
      return;
    }

    let isCurrent = true;
    setReport(null);
    setError("");
    setIsMissing(false);
    setIsLoading(true);

    runWithWake(getReport(id))
      .then((savedReport) => {
        if (isCurrent) setReport(savedReport);
      })
      .catch((requestError: unknown) => {
        if (!isCurrent) return;
        if (isMissingReportError(requestError)) {
          setIsMissing(true);
        } else {
          setError(
            requestError instanceof Error
              ? requestError.message
              : "The report could not be loaded. Please try again.",
          );
        }
      })
      .finally(() => {
        if (isCurrent) setIsLoading(false);
      });

    return () => {
      isCurrent = false;
    };
  }, [id, retryCount, runWithWake, unsavedReport]);

  const shareId = report?.id || id || "";
  const shareUrl = useMemo(
    () =>
      shareId
        ? `${window.location.origin}/report/${encodeURIComponent(shareId)}`
        : "",
    [shareId],
  );

  async function copyShareLink() {
    if (!shareUrl) return;
    setCopyMessage("");
    try {
      await navigator.clipboard.writeText(shareUrl);
      setCopyMessage("Link copied.");
    } catch {
      setCopyMessage("Could not copy the link. Please try again.");
    }
  }

  if (isLoading) {
    return (
      <main className="mx-auto min-h-[55dvh] max-w-6xl px-5 py-10 sm:px-8 sm:py-14">
        <ServerWakeBanner visible={isWaking} />
        <div className="mt-6 max-w-3xl space-y-4" aria-label="Loading report">
          <div className="h-4 w-36 bg-[#eeeaf8]" />
          <div className="h-10 w-3/4 bg-[#eeeaf8]" />
          <div className="h-24 w-full bg-[#eeeaf8]" />
          <div className="h-40 w-full bg-[#eeeaf8]" />
        </div>
      </main>
    );
  }

  if (isMissing || (!report && !error)) {
    return (
      <main className="mx-auto min-h-[55dvh] max-w-3xl px-5 py-14 sm:px-8 sm:py-20">
        <p className="mb-3 text-xs font-bold uppercase tracking-[0.15em] text-primary">
          Report unavailable
        </p>
        <h1 className="text-3xl font-semibold tracking-tight sm:text-4xl">
          This report may have expired.
        </h1>
        <p className="mt-4 max-w-xl text-base leading-7 text-black/65">
          Reports are kept for 7 days. If you still have the contract, you can
          analyze it again to create a new report.
        </p>
        <Link
          to="/"
          className="mt-7 inline-flex bg-primary px-5 py-3 text-sm font-semibold text-white"
          data-testid="link-analyze-again"
        >
          Analyze a contract
        </Link>
      </main>
    );
  }

  if (error) {
    return (
      <main className="mx-auto min-h-[55dvh] max-w-3xl px-5 py-14 sm:px-8 sm:py-20">
        <p className="mb-3 text-xs font-bold uppercase tracking-[0.15em] text-primary">
          Report unavailable
        </p>
        <h1 className="text-3xl font-semibold tracking-tight sm:text-4xl">
          We couldn’t load this report.
        </h1>
        <p className="mt-4 text-sm leading-6 text-black/65" role="alert">
          {error}
        </p>
        <div className="mt-7 flex flex-wrap items-center gap-5">
          <button
            type="button"
            onClick={() => setRetryCount((current) => current + 1)}
            className="bg-primary px-5 py-3 text-sm font-semibold text-white"
            data-testid="button-retry-report"
          >
            Try again
          </button>
          <Link to="/" data-testid="link-back-home">
            Back to ClauseGuard
          </Link>
        </div>
      </main>
    );
  }

  if (!report) return null;

  return (
    <main className="mx-auto max-w-6xl px-5 pb-12 pt-8 sm:px-8 sm:pt-12">
      <div className="mb-6 flex flex-wrap items-center justify-between gap-4">
        <Link to="/" className="text-sm font-medium" data-testid="link-new-analysis">
          ← New analysis
        </Link>
        <div className="flex flex-wrap items-center gap-3">
          <button
            type="button"
            onClick={copyShareLink}
            disabled={!shareUrl}
            className="border border-black/25 px-4 py-2 text-sm font-semibold text-black hover:border-primary disabled:cursor-not-allowed disabled:opacity-45"
            data-testid="button-copy-share-link"
          >
            Copy share link
          </button>
          {copyMessage && (
            <span
              className="text-sm text-black/65"
              role="status"
              data-testid="status-copy-link"
            >
              {copyMessage}
            </span>
          )}
        </div>
      </div>

      <ServerWakeBanner visible={isWaking} />

      <section className="border-y border-black/15 py-7 sm:py-9" aria-labelledby="report-title">
        <p className="mb-3 text-xs font-bold uppercase tracking-[0.15em] text-primary">
          Contract review
        </p>
        <div className="flex flex-col gap-5 sm:flex-row sm:items-end sm:justify-between">
          <div className="max-w-3xl">
            <h1 id="report-title" className="text-3xl font-semibold tracking-tight sm:text-5xl">
              Your contract, in plain language.
            </h1>
            <p
              className="mt-4 max-w-3xl text-base leading-7 text-black/70"
              data-testid="text-report-summary"
            >
              {report.summary}
            </p>
            <p className="mt-3 text-xs leading-5 text-black/55">
              AI-generated analysis can be wrong. Information only, not legal advice.
            </p>
          </div>
          <div className="flex shrink-0 flex-col gap-2">
            <span className="text-xs font-bold uppercase tracking-[0.12em] text-black/55">
              Overall risk
            </span>
            <RiskBadge risk={report.overall_risk} />
          </div>
        </div>
      </section>

      <div className="mt-8 grid items-start gap-8 lg:grid-cols-[minmax(0,1fr)_340px]">
        <section aria-labelledby="clauses-title">
          <div className="mb-4 flex flex-wrap items-baseline justify-between gap-3">
            <h2 id="clauses-title" className="text-xl font-semibold tracking-tight">
              Clause-by-clause
            </h2>
            <p className="text-xs text-black/55">
              {report.clauses.length}{" "}
              {report.clauses.length === 1 ? "clause" : "clauses"} reviewed
            </p>
          </div>
          {report.clauses.length > 0 ? (
            <div className="space-y-4">
              {report.clauses.map((clause) => (
                <ClauseCard key={clause.id} clause={clause} />
              ))}
            </div>
          ) : (
            <p className="border border-dashed border-black/25 px-5 py-8 text-sm text-black/65">
              No individual clauses were returned for this report.
            </p>
          )}
        </section>
        <aside className="space-y-5 lg:sticky lg:top-6">
          <ChatPanel clauses={report.clauses} />
          <p className="border-l-2 border-black/20 px-4 py-3 text-xs leading-5 text-black/60">
            Reports are deleted automatically after 7 days
          </p>
        </aside>
      </div>
    </main>
  );
}
