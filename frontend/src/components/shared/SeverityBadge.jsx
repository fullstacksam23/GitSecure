import { Badge } from "../ui/badge";
import { normalizeSeverity, severityTheme } from "../../lib/utils";
import { NO_KNOWN_VULNERABILITIES_LABEL } from "../../lib/ecosystem";

export default function SeverityBadge({ severity, issueCount }) {
  if (issueCount === 0) {
    return (
      <Badge variant="low" className="gap-2 normal-case tracking-normal">
        <span>{NO_KNOWN_VULNERABILITIES_LABEL}</span>
        <span className="rounded-full border border-emerald-400/20 bg-emerald-400/10 px-2 py-0.5 text-[10px] text-emerald-100">
          0 issues
        </span>
      </Badge>
    );
  }

  const normalized = normalizeSeverity(severity);
  return <Badge variant={normalized}>{severityTheme[normalized].label}</Badge>;
}
