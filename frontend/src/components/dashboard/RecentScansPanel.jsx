import { useNavigate } from "react-router-dom";
import { formatCommitHash, formatDate, toPercent } from "../../lib/utils";
import SeverityBadge from "../shared/SeverityBadge";
import StatusIndicator from "../shared/StatusIndicator";

const bars = ["critical", "high", "medium", "low"];

export default function RecentScansPanel({ scans = [] }) {
  const navigate = useNavigate();

  return (
    <div className="panel p-6">
      <div className="mb-5">
        <h3 className="text-lg font-semibold text-white">Recent Scan Jobs</h3>
        <p className="text-sm text-slate-400">Latest runs with severity distribution at a glance.</p>
      </div>
      <div className="space-y-4">
        {scans.map((scan) => {
          const total = scan.vulnerability_count || 0;
          return (
            <button
              key={scan.job_id}
              className="surface w-full p-4 text-left transition hover:-translate-y-0.5 hover:border-white/20 hover:bg-white/[0.06]"
              onClick={() => navigate(`/scans/${scan.job_id}`)}
            >
              <div className="flex flex-wrap items-center justify-between gap-3">
                <div>
                  <p className="font-medium text-white">{scan.repo}</p>
                  <div className="mt-1 flex items-center gap-2 text-xs text-slate-500">
                    <span className="font-mono">{formatCommitHash(scan.commit_hash)}</span>
                    <span>{formatDate(scan.created_at, { dateStyle: "medium", timeStyle: "short" })}</span>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <SeverityBadge severity={scan.top_severity} issueCount={total} status={scan.status} />
                  <StatusIndicator status={scan.status} />
                </div>
              </div>
              <div className="mt-4 overflow-hidden rounded-full bg-white/[0.06]">
                <div className="flex h-2">
                  {bars.map((severity) => (
                    <div
                      key={severity}
                      className={`h-full ${
                        severity === "critical" ? "bg-red-500" : severity === "high" ? "bg-orange-400" : severity === "medium" ? "bg-yellow-400" : "bg-emerald-400"
                      }`}
                      style={{ width: `${toPercent(scan.severity_counts?.[severity] || 0, total)}%` }}
                    />
                  ))}
                </div>
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
}
