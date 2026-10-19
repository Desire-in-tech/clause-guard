import type { Risk } from "../api/types";

interface RiskBadgeProps {
  risk: Risk;
}

const riskStyle: Record<Risk, string> = {
  green: "border-[#157347]/30 bg-[#eaf6ef] text-[#125c39]",
  yellow: "border-[#987000]/30 bg-[#fff6d9] text-[#745600]",
  red: "border-[#bb2634]/30 bg-[#fff0f1] text-[#a21e2b]",
};

const riskLabel: Record<Risk, string> = {
  green: "Green · Lower risk",
  yellow: "Yellow · Review",
  red: "Red · Higher risk",
};

export default function RiskBadge({ risk }: RiskBadgeProps) {
  return (
    <span
      className={`inline-flex w-fit items-center border px-3 py-1 text-xs font-semibold ${riskStyle[risk]}`}
      data-testid={`status-risk-${risk}`}
    >
      {riskLabel[risk]}
    </span>
  );
}
