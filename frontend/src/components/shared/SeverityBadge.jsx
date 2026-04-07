import { Badge } from "../ui/badge";
import { normalizeSeverity, severityTheme } from "../../lib/utils";

export default function SeverityBadge({ severity }) {
  const normalized = normalizeSeverity(severity);
  return <Badge variant={normalized}>{severityTheme[normalized].label}</Badge>;
}
