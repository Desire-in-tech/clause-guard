import type { ClauseAnalysis } from "../api/types";
import RiskBadge from "./RiskBadge";

interface ClauseCardProps {
  clause: ClauseAnalysis;
}

function formatType(type: string) {
  return type.replaceAll("_", " ");
}

export default function ClauseCard({ clause }: ClauseCardProps) {
  return (
    <article
      id={`clause-${clause.id}`}
      className="scroll-mt-6 border border-black/15 bg-white"
      data-testid={`card-clause-${clause.id}`}
    >
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-black/10 px-5 py-4">
        <div className="flex items-center gap-3">
          <span className="grid size-8 place-items-center border border-black/15 text-sm font-semibold tabular-nums">
            {clause.index}
          </span>
          <h3 className="text-sm font-semibold capitalize text-black">
            {formatType(clause.type)}
          </h3>
        </div>
        <RiskBadge risk={clause.risk} />
      </div>
      <div className="space-y-5 px-5 py-5">
        <blockquote className="border-l-2 border-primary/50 pl-4 text-sm leading-6 text-black/75">
          {clause.text}
        </blockquote>
        <div>
          <h4 className="mb-1 text-xs font-bold uppercase tracking-[0.12em] text-black/55">
            Why it matters
          </h4>
          <p className="text-sm leading-6 text-black">{clause.reason}</p>
        </div>
        <div className="border-l-2 border-primary bg-[#f8f6ff] px-4 py-3">
          <h4 className="mb-1 text-xs font-bold uppercase tracking-[0.12em] text-black/55">
            Possible rewrite
          </h4>
          <p className="text-sm leading-6 text-black">{clause.suggested_rewrite}</p>
        </div>
      </div>
    </article>
  );
}
