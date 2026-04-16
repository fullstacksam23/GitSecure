import { ArrowUpDown, ShieldAlert, Star } from "lucide-react";
import { useNavigate } from "react-router-dom";
import {
  getRepoJobId,
  getRepoName,
  getRepoRank,
  getRepoRiskScore,
  getRepoStars,
  getRepoStatus,
  getRepoTopSeverity,
  getRepoVulnerabilityCount,
} from "../../lib/ecosystem";
import { cn, formatRisk } from "../../lib/utils";
import SeverityBadge from "../shared/SeverityBadge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";
import ProgressBar from "./ProgressBar";
import StatusBadge from "./StatusBadge";

const columns = [
  { id: "repo", label: "Repository" },
  { id: "stars", label: "Stars" },
  { id: "rank", label: "Rank" },
  { id: "status", label: "Vulnerability status" },
  { id: "vulnerability_count", label: "Vulnerabilities" },
  { id: "top_severity", label: "Top severity" },
];

export default function RepoTable({ batchId, items = [], sortBy, sortOrder, onSort, highlightedRepo }) {
  const navigate = useNavigate();

  function handleOpen(repo) {
    const jobId = getRepoJobId(repo);
    if (!jobId) return;
    navigate(`/scans/${jobId}?batch=${encodeURIComponent(batchId)}`);
  }

  return (
    <div className="panel overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow className="border-white/10 hover:bg-transparent">
            {columns.map((column) => (
              <TableHead key={column.id}>
                <button
                  type="button"
                  className="inline-flex items-center gap-2 text-left text-xs uppercase tracking-[0.18em] text-slate-500 transition hover:text-white"
                  onClick={() => onSort(column.id)}
                >
                  {column.label}
                  <ArrowUpDown
                    className={cn(
                      "h-3.5 w-3.5",
                      sortBy === column.id ? "text-cyan-300" : "text-slate-600"
                    )}
                  />
                  {sortBy === column.id ? (
                    <span className="text-[10px] text-cyan-300">{sortOrder === "asc" ? "ASC" : "DESC"}</span>
                  ) : null}
                </button>
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((item) => {
            const repoName = getRepoName(item);
            const topSeverity = getRepoTopSeverity(item);
            const risk = getRepoRiskScore(item);
            const hasJob = Boolean(getRepoJobId(item));
            const isHighlighted = repoName === highlightedRepo;
            const vulnerabilityCount = getRepoVulnerabilityCount(item);
            const repoStatus = getRepoStatus(item);

            return (
              <TableRow
                key={item.id || repoName}
                className={cn(
                  "border-white/10",
                  hasJob ? "cursor-pointer hover:bg-white/[0.04]" : "opacity-80",
                  isHighlighted && "bg-red-500/[0.06]"
                )}
                onClick={() => handleOpen(item)}
              >
                <TableCell className="min-w-[240px]">
                  <div className="flex items-start gap-3">
                    <div className="mt-1 rounded-2xl border border-white/10 bg-white/[0.04] p-2">
                      <ShieldAlert className={cn("h-4 w-4", isHighlighted ? "text-red-300" : "text-cyan-300")} />
                    </div>
                    <div className="min-w-0">
                      <p className="truncate font-medium text-white">{repoName}</p>
                      <p className="mt-1 text-xs text-slate-500">
                        {hasJob ? "Open scan detail" : "Waiting for scan job"}
                      </p>
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <div className="inline-flex items-center gap-2 text-slate-200">
                    <Star className="h-4 w-4 text-yellow-300" />
                    {getRepoStars(item).toLocaleString()}
                  </div>
                </TableCell>
                <TableCell className="min-w-[180px]">
                  <div className="space-y-2">
                    <div className="flex items-center justify-between text-sm text-slate-300">
                      <span>#{getRepoRank(item) || "-"}</span>
                      <span>{formatRisk(risk)} risk</span>
                    </div>
                    <ProgressBar value={risk} className="max-w-[150px]" />
                  </div>
                </TableCell>
                <TableCell>
                  <StatusBadge status={repoStatus} issueCount={vulnerabilityCount} />
                </TableCell>
                <TableCell className="text-slate-100">{vulnerabilityCount}</TableCell>
                <TableCell>
                  <SeverityBadge severity={topSeverity} issueCount={vulnerabilityCount} status={repoStatus} />
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}
