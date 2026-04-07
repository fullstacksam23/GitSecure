import { useQuery } from "@tanstack/react-query";
import { Download, GitCompareArrows } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import { api } from "../api/client";
import FilterBar from "../components/jobs/FilterBar";
import VulnerabilityDrawer from "../components/jobs/VulnerabilityDrawer";
import VulnerabilityTable from "../components/jobs/VulnerabilityTable";
import CopyButton from "../components/shared/CopyButton";
import EmptyState from "../components/shared/EmptyState";
import PageSkeleton from "../components/shared/PageSkeleton";
import SeverityBadge from "../components/shared/SeverityBadge";
import StatusIndicator from "../components/shared/StatusIndicator";
import { Button } from "../components/ui/button";
import { formatCommitHash, formatDate } from "../lib/utils";
import { useDebouncedValue } from "../hooks/useDebouncedValue";

export default function ScanDetailPage() {
  const { jobId } = useParams();
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const [searchInput, setSearchInput] = useState(searchParams.get("search") || "");
  const debouncedSearch = useDebouncedValue(searchInput, 300);

  const page = Number(searchParams.get("page") || 1);
  const severity = (searchParams.get("severity") || "").split(",").filter(Boolean);
  const ecosystem = searchParams.get("ecosystem") || "";
  const fixState = searchParams.get("fix_state") || "";
  const sortBy = searchParams.get("sort_by") || "created_at";
  const sortOrder = searchParams.get("sort_order") || "desc";
  const selectedId = searchParams.get("selected") || "";

  const scanQuery = useQuery({
    queryKey: ["scan", jobId],
    queryFn: () => api.getScan(jobId),
    enabled: Boolean(jobId),
    refetchInterval: (query) => ["queued", "running"].includes(query.state.data?.status) ? 5000 : false,
  });

  const vulnsQuery = useQuery({
    queryKey: ["vulns", jobId, page, severity.join(","), ecosystem, fixState, debouncedSearch, sortBy, sortOrder],
    queryFn: () =>
      api.getVulnerabilities({
        page,
        pageSize: 50,
        severity: severity.join(","),
        ecosystem,
        fixState,
        search: debouncedSearch,
        jobId,
        sortBy,
        sortOrder,
      }),
    enabled: Boolean(jobId),
    placeholderData: (previous) => previous,
  });

  const scan = scanQuery.data;
  const rows = vulnsQuery.data?.items || [];
  const pagination = vulnsQuery.data?.pagination;
  const facets = vulnsQuery.data?.facets || { ecosystems: [], fix_states: [] };

  useEffect(() => {
    const next = new URLSearchParams(searchParams);
    if (debouncedSearch) next.set("search", debouncedSearch);
    else next.delete("search");
    if (debouncedSearch !== searchParams.get("search")) setSearchParams(next, { replace: true });
  }, [debouncedSearch, searchParams, setSearchParams]);

  function updateParam(key, value) {
    const next = new URLSearchParams(searchParams);
    if (value) next.set(key, value);
    else next.delete(key);
    if (key !== "selected") next.set("page", "1");
    setSearchParams(next);
  }

  function toggleSeverity(value) {
    const next = severity.includes(value) ? severity.filter((item) => item !== value) : [...severity, value];
    updateParam("severity", next.join(","));
  }

  function handleSort(column) {
    const nextOrder = sortBy === column && sortOrder === "desc" ? "asc" : "desc";
    updateParam("sort_by", column);
    updateParam("sort_order", nextOrder);
  }

  function exportCsv() {
    const headers = ["Package", "Version", "Severity", "Risk", "Fix State", "Ecosystem"];
    const csvRows = rows.map((item) => [item.package, item.version, item.normalized_severity || item.severity, item.risk, item.fix_state, item.ecosystem]);
    const blob = new Blob([[headers, ...csvRows].map((row) => row.join(",")).join("\n")], { type: "text/csv;charset=utf-8;" });
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${jobId}-vulnerabilities.csv`;
    link.click();
  }

  if (scanQuery.isLoading) return <PageSkeleton cards={3} rows={8} />;
  if (scanQuery.isError) return <EmptyState title="Unable to load scan" description={scanQuery.error.message} />;

  return (
    <div className="space-y-6 pb-8">
      <div className="panel p-6">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div>
            <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">Scan Detail</p>
            <h1 className="mt-3 text-3xl font-semibold text-white">{scan.repo}</h1>
            <div className="mt-3 flex flex-wrap items-center gap-3 text-sm text-slate-400">
              <span className="font-mono text-slate-200">{formatCommitHash(scan.commit_hash)}</span>
              <CopyButton value={scan.commit_hash} label="Commit hash copied" />
              <span>{formatDate(scan.created_at)}</span>
            </div>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <SeverityBadge severity={scan.severity_counts?.critical ? "critical" : scan.severity_counts?.high ? "high" : scan.severity_counts?.medium ? "medium" : scan.severity_counts?.low ? "low" : "unknown"} />
            <StatusIndicator status={scan.status} />
            <Button variant="outline" onClick={exportCsv}>
              <Download className="h-4 w-4" />
              Export CSV
            </Button>
            <Button variant="secondary" onClick={() => navigate(`/compare?base=${jobId}`)}>
              <GitCompareArrows className="h-4 w-4" />
              Compare
            </Button>
          </div>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-4">
        {[
          ["Findings", scan.vulnerability_count],
          ["Critical", scan.severity_counts?.critical || 0],
          ["High", scan.severity_counts?.high || 0],
          ["Ecosystems", scan.ecosystems?.length || 0],
        ].map(([label, value]) => (
          <div key={label} className="panel p-5">
            <p className="text-sm text-slate-500">{label}</p>
            <p className="mt-4 text-3xl font-semibold text-white">{value}</p>
          </div>
        ))}
      </div>

      <FilterBar
        severity={severity}
        onSeverityToggle={toggleSeverity}
        ecosystem={ecosystem}
        onEcosystemChange={(value) => updateParam("ecosystem", value)}
        ecosystems={facets.ecosystems}
        fixState={fixState}
        onFixStateChange={(value) => updateParam("fix_state", value)}
        fixStates={facets.fix_states}
        search={searchInput}
        onSearchChange={setSearchInput}
      />

      {vulnsQuery.isError ? (
        <EmptyState title="Unable to load vulnerabilities" description={vulnsQuery.error.message} />
      ) : (
        <>
          <VulnerabilityTable items={rows} sortBy={sortBy} sortOrder={sortOrder} onSort={handleSort} onSelect={(item) => updateParam("selected", item.id)} />
          <div className="flex items-center justify-between">
            <p className="text-sm text-slate-500">
              Page {pagination?.page || 1} of {pagination?.total_pages || 1}
            </p>
            <div className="flex gap-2">
              <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => updateParam("page", String(page - 1))}>
                Previous
              </Button>
              <Button variant="outline" size="sm" disabled={page >= (pagination?.total_pages || 1)} onClick={() => updateParam("page", String(page + 1))}>
                Next
              </Button>
            </div>
          </div>
        </>
      )}

      <VulnerabilityDrawer jobId={jobId} vulnerabilityId={selectedId} open={Boolean(selectedId)} onOpenChange={(open) => !open && updateParam("selected", "")} />
    </div>
  );
}
