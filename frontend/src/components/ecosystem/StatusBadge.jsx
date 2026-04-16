import { Badge } from "../ui/badge";
import {
  NO_KNOWN_VULNERABILITIES_LABEL,
  NO_KNOWN_VULNERABILITIES_STATUS,
} from "../../lib/ecosystem";

const variants = {
  queued: "neutral",
  running: "default",
  completed: "low",
  complete: "low",
  failed: "critical",
  [NO_KNOWN_VULNERABILITIES_STATUS]: "low",
};

export default function StatusBadge({ status, issueCount }) {
  const normalized = String(status || "").toLowerCase();
  const isNoKnownVulnerabilities = normalized === NO_KNOWN_VULNERABILITIES_STATUS;
  const label = isNoKnownVulnerabilities
    ? NO_KNOWN_VULNERABILITIES_LABEL
    : normalized
      ? normalized.charAt(0).toUpperCase() + normalized.slice(1)
      : "Unknown";

  return (
    <Badge
      variant={variants[normalized] || "neutral"}
      className={isNoKnownVulnerabilities ? "gap-2 normal-case tracking-normal" : undefined}
    >
      <span>{label}</span>
      {isNoKnownVulnerabilities && issueCount === 0 ? (
        <span className="rounded-full border border-emerald-400/20 bg-emerald-400/10 px-2 py-0.5 text-[10px] text-emerald-100">
          0 issues
        </span>
      ) : null}
    </Badge>
  );
}
