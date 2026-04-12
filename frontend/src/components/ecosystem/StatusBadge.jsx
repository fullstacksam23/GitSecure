import { Badge } from "../ui/badge";

const variants = {
  queued: "neutral",
  running: "default",
  completed: "low",
  complete: "low",
  failed: "critical",
};

export default function StatusBadge({ status }) {
  const normalized = String(status || "").toLowerCase();
  const label = normalized ? normalized.charAt(0).toUpperCase() + normalized.slice(1) : "Unknown";
  return <Badge variant={variants[normalized] || "neutral"}>{label}</Badge>;
}
