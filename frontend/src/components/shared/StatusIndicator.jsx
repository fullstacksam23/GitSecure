import { Activity, CheckCircle2, LoaderCircle, ShieldAlert, XCircle } from "lucide-react";
import { cn } from "../../lib/utils";

const statusConfig = {
  queued: { label: "Queued", icon: Activity, className: "text-slate-300 bg-white/[0.05]" },
  running: { label: "Running", icon: LoaderCircle, className: "text-cyan-300 bg-cyan-400/10" },
  complete: { label: "Completed", icon: CheckCircle2, className: "text-emerald-300 bg-emerald-400/10" },
  completed: { label: "Completed", icon: CheckCircle2, className: "text-emerald-300 bg-emerald-400/10" },
  failed: { label: "Failed", icon: XCircle, className: "text-red-300 bg-red-500/10" },
};

export default function StatusIndicator({ status }) {
  const config = statusConfig[String(status || "").toLowerCase()] || {
    label: status || "Unknown",
    icon: ShieldAlert,
    className: "text-slate-300 bg-white/[0.05]",
  };
  const Icon = config.icon;

  return (
    <span className={cn("inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-medium", config.className)}>
      <Icon className={cn("h-3.5 w-3.5", status === "running" && "animate-spin")} />
      {config.label}
    </span>
  );
}
