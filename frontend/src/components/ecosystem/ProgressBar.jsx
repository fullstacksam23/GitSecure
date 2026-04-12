import { cn } from "../../lib/utils";

export default function ProgressBar({ value = 0, label, tone = "default", className }) {
  const trackClassName =
    tone === "success"
      ? "bg-emerald-400"
      : tone === "warning"
        ? "bg-orange-400"
        : "bg-cyan-400";

  return (
    <div className={cn("space-y-2", className)}>
      {label ? <div className="flex items-center justify-between gap-3 text-xs text-slate-400">{label}</div> : null}
      <div className="h-2.5 overflow-hidden rounded-full bg-white/[0.06]">
        <div
          className={cn("h-full rounded-full transition-all duration-500", trackClassName)}
          style={{ width: `${Math.max(0, Math.min(100, value))}%` }}
        />
      </div>
    </div>
  );
}
