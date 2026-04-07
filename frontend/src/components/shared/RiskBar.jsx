import { cn } from "../../lib/utils";

export default function RiskBar({ value = 0, className }) {
  const safe = Math.max(0, Math.min(100, Math.round(value)));
  const color =
    safe >= 85 ? "bg-red-500" : safe >= 65 ? "bg-orange-400" : safe >= 40 ? "bg-yellow-400" : "bg-emerald-400";

  return (
    <div className={cn("flex items-center gap-3", className)}>
      <div className="h-2 flex-1 overflow-hidden rounded-full bg-white/[0.08]">
        <div className={cn("h-full rounded-full transition-all", color)} style={{ width: `${safe}%` }} />
      </div>
      <span className="w-10 text-right text-xs text-slate-400">{safe}</span>
    </div>
  );
}
