import { cn, formatRisk, getSeverityRiskFallback, normalizeRisk } from "../../lib/utils";

export default function RiskBar({ value = 0, severity, className }) {
  const normalized = normalizeRisk(value);
  const safe = normalized > 0 ? formatRisk(value) : getSeverityRiskFallback(severity);
  const color =
    safe >= 85 ? "bg-red-500" : safe >= 65 ? "bg-orange-400" : safe >= 40 ? "bg-yellow-400" : "bg-emerald-400";

  return (
    <div className={cn("flex items-center gap-3", className)} title="Risk represents exploit likelihood">
      <div className="h-2 flex-1 overflow-hidden rounded-full bg-white/[0.08]" aria-hidden="true">
        <div className={cn("h-full rounded-full transition-all", color)} style={{ width: `${safe}%` }} />
      </div>
      <span className="w-16 text-right text-xs text-slate-400">{safe.toFixed(2)}%</span>
    </div>
  );
}
