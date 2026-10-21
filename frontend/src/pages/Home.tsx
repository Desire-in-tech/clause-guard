import { useState, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { APIRequestError, analyzeContract } from "../api/client";
import type { AnalysisReport } from "../api/types";
import Disclaimer from "../components/Disclaimer";
import ServerWakeBanner from "../components/ServerWakeBanner";
import UploadPanel from "../components/UploadPanel";
import useServerWake from "../hooks/useServerWake";

export default function Home() {
  const [contractText, setContractText] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isExtracting, setIsExtracting] = useState(false);
  const { isWaking, runWithWake } = useServerWake();
  const navigate = useNavigate();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const text = contractText.trim();
    if (!text || isSubmitting) return;

    setError("");
    setIsSubmitting(true);
    try {
      const report: AnalysisReport = await runWithWake(analyzeContract(text));
      if (report.id) {
        navigate(`/report/${encodeURIComponent(report.id)}`);
      } else {
        navigate("/results", { state: { report } });
      }
    } catch (requestError) {
      if (requestError instanceof APIRequestError && requestError.status === 502) {
        setError("AI analysis is temporarily unavailable. Please try again later.");
      } else if (
        requestError instanceof APIRequestError &&
        requestError.status === 429
      ) {
        setError(
          requestError.message ||
            "The demo AI limit has been reached. Please try again later.",
        );
      } else {
        setError(
          requestError instanceof Error
            ? requestError.message
            : "The contract could not be analyzed. Please try again.",
        );
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="mx-auto max-w-4xl px-5 pb-12 pt-10 sm:px-8 sm:pt-16">
      <div className="mb-8 max-w-2xl">
        <p className="mb-4 text-xs font-bold uppercase tracking-[0.16em] text-primary">
          A clearer read before you sign
        </p>
        <h1 className="max-w-2xl text-4xl font-semibold leading-[1.08] tracking-tight text-black sm:text-6xl">
          Understand the terms in front of you.
        </h1>
        <p className="mt-5 max-w-xl text-base leading-7 text-black/65">
          Paste a contract to see a plain-language summary, clause-by-clause
          risk notes, and suggested rewrites.
        </p>
      </div>

      <form
        onSubmit={handleSubmit}
        className="border border-black/20 bg-white"
        aria-label="Analyze a contract"
      >
        <UploadPanel
          value={contractText}
          onChange={setContractText}
          disabled={isSubmitting}
          onLoadingChange={setIsExtracting}
        />
        <div className="px-5 pb-5 sm:px-6 sm:pb-6">
          <div className="mt-5 space-y-4">
            <ServerWakeBanner visible={isWaking} />
            {error && (
              <p
                className="border-l-2 border-[#bb2634] bg-[#fff0f1] px-4 py-3 text-sm text-[#8d1a25]"
                role="alert"
                data-testid="error-analyze"
              >
                {error}
              </p>
            )}
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <Disclaimer />
                <p className="mt-2 max-w-md text-xs leading-5 text-black/55">
                  Demo uses free third-party AI models. Do not paste confidential
                  or personal information.
                </p>
              </div>
              <button
                type="submit"
                disabled={!contractText.trim() || isSubmitting || isExtracting}
                className="bg-primary px-7 py-3 text-sm font-semibold text-white hover:bg-[#7040f0] disabled:cursor-not-allowed disabled:bg-black/20"
                data-testid="button-analyze"
              >
                {isSubmitting
                  ? "Analyzing... this can take up to a minute"
                  : "Analyze"}
              </button>
            </div>
          </div>
        </div>
      </form>

    </main>
  );
}
