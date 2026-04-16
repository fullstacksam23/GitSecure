import { useNavigate } from "react-router-dom";
import { formatCommitHash, formatDate } from "../../lib/utils";
import SeverityBadge from "../shared/SeverityBadge";
import StatusIndicator from "../shared/StatusIndicator";
import { Button } from "../ui/button";

export default function ScanList({ items = [], title, compact = false }) {
  const navigate = useNavigate();

  return (
    <div className="panel p-6">
      {title ? (
        <div className="mb-5">
          <h3 className="text-lg font-semibold text-white">{title}</h3>
        </div>
      ) : null}
      <div className="space-y-3">
        {items.map((item) => (
          <div key={item.job_id} className="surface flex flex-wrap items-center justify-between gap-4 p-4">
            <div className="min-w-0">
              <p className="truncate text-sm font-medium text-white">{item.repo}</p>
              <div className="mt-1 flex flex-wrap items-center gap-2 text-xs text-slate-500">
                <span className="font-mono">{formatCommitHash(item.commit_hash)}</span>
                <span>{formatDate(item.created_at, { dateStyle: compact ? "short" : "medium", timeStyle: "short" })}</span>
              </div>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              <SeverityBadge
                severity={item.top_severity}
                issueCount={Number(item.vulnerability_count || 0)}
                status={item.status}
              />
              <StatusIndicator status={item.status} />
              <Button variant="outline" size="sm" onClick={() => navigate(`/scans/${item.job_id}`)}>
                View Scan
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
