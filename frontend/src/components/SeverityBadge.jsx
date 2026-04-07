import { Badge } from "./ui/badge";
import { normalizeSeverity } from "../lib/utils";

export default function SeverityBadge({ severity }) {
  const normalized = normalizeSeverity(severity);
  const label = normalized === "unknown" ? severity || "Unknown" : normalized;

  const variantMap = {
    critical: "critical",
    high: "high",
    medium: "medium",
    low: "low",
    unknown: "neutral",
  };

  return <Badge variant={variantMap[normalized]}>{label}</Badge>;
}
