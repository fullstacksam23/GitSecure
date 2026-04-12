import { useQuery, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft, Download, GitCompareArrows } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
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
import { formatCommitHash, formatDate, formatRisk } from "../lib/utils";
import { groupVulnerabilitiesByPackage } from "../lib/vulnerability-groups";
import { useDebouncedValue } from "../hooks/useDebouncedValue";

const allowedSortColumns = new Set(["package", "vulnCount", "severity", "risk", "ecosystem"]);

export default function ScanDetailPage() {
  const { jobId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [searchParams, setSearchParams] = useSearchParams();
  const [searchInput, setSearchInput] = useState(searchParams.get("search") || "");
  const debouncedSearch = useDebouncedValue(searchInput, 300);
  const lastStatusRef = useRef("");

  const severity = (searchParams.get("severity") || "").split(",").filter(Boolean);
  const severityKey = severity.join(",");
  const ecosystem = searchParams.get("ecosystem") || "";
  const fixState = searchParams.get("fix_state") || "";
  const sortBy = allowedSortColumns.has(searchParams.get("sort_by")) ? searchParams.get("sort_by") : "severity";
  const sortOrder = searchParams.get("sort_order") || "desc";
  const selectedId = searchParams.get("selected") || "";
  const selectedPackage = searchParams.get("selected_package") || "";
  const batchId = searchParams.get("batch") || "";
  const vulnerabilitiesQueryKey = useMemo(
    () => ["vulns", "all", jobId, severityKey, ecosystem, fixState, debouncedSearch],
    [jobId, severityKey, ecosystem, fixState, debouncedSearch]
  );
  const isActiveStatus = (status) => ["queued", "running"].includes(String(status || "").toLowerCase());

  const scanQuery = useQuery({
    queryKey: ["scan", jobId],
    queryFn: () => api.getScan(jobId),
    enabled: Boolean(jobId),
    staleTime: 0,
    refetchInterval: (query) => (isActiveStatus(query.state.data?.status) ? 2000 : false),
    refetchIntervalInBackground: true,
  });

  const vulnsQuery = useQuery({
    queryKey: vulnerabilitiesQueryKey,
    queryFn: () =>
      api.getAllVulnerabilities({
        pageSize: 100,
        severity: severityKey,
        ecosystem,
        fixState,
        search: debouncedSearch,
        jobId,
      }),
    enabled: Boolean(jobId),
    staleTime: 0,
    placeholderData: (previous) => previous,
    refetchInterval: () => (isActiveStatus(scanQuery.data?.status) ? 3000 : false),
    refetchIntervalInBackground: true,
  });

  const scan = scanQuery.data;
  const rows = vulnsQuery.data?.items || [];
  const facets = vulnsQuery.data?.facets || { ecosystems: [], fix_states: [] };
  const groupedPackages = useMemo(() => groupVulnerabilitiesByPackage(rows), [rows]);

  useEffect(() => {
    const next = new URLSearchParams(searchParams);
    if (debouncedSearch) next.set("search", debouncedSearch);
    else next.delete("search");
    if (debouncedSearch !== searchParams.get("search")) setSearchParams(next, { replace: true });
  }, [debouncedSearch, searchParams, setSearchParams]);

  useEffect(() => {
    const nextStatus = String(scan?.status || "").toLowerCase();
    if (!nextStatus || nextStatus === lastStatusRef.current) return;

    lastStatusRef.current = nextStatus;

    queryClient.invalidateQueries({ queryKey: ["scans"] });
    queryClient.invalidateQueries({ queryKey: ["sidebar-scans"] });
    queryClient.invalidateQueries({ queryKey: ["history"] });
    queryClient.invalidateQueries({ queryKey: ["compare-scans-list"] });
    queryClient.invalidateQueries({ queryKey: ["dashboard-summary"] });

    if (nextStatus === "completed") {
      queryClient.invalidateQueries({ queryKey: vulnerabilitiesQueryKey });
      queryClient.refetchQueries({ queryKey: vulnerabilitiesQueryKey, exact: true });
    }
  }, [queryClient, scan?.status, vulnerabilitiesQueryKey]);

  function updateParam(key, value) {
    const next = new URLSearchParams(searchParams);
    if (value) next.set(key, value);
    else next.delete(key);
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
    const headers = ["Package", "Unique Vulnerabilities", "Highest Severity", "Highest Risk (%)", "Fix State", "Ecosystem"];
    const csvRows = groupedPackages.map((item) => [
      item.package,
      item.vulnCount,
      item.highestSeverity,
      formatRisk(item.highestRisk).toFixed(2),
      item.fixState,
      item.ecosystem,
    ]);
    const serializeCell = (value) => `"${String(value ?? "").replaceAll('"', '""')}"`;
    const blob = new Blob([[headers, ...csvRows].map((row) => row.map(serializeCell).join(",")).join("\n")], { type: "text/csv;charset=utf-8;" });
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${jobId}-packages.csv`;
    link.click();
  }

  if (scanQuery.isLoading) return <PageSkeleton cards={3} rows={8} />;
  if (scanQuery.isError) return <EmptyState title="Unable to load scan" description={scanQuery.error.message} />;

  return (
    <div className="space-y-6 pb-8">
      <div className="panel p-6">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div>
            {batchId ? (
              <Button variant="ghost" size="sm" className="mb-4" onClick={() => navigate(`/ecosystem/batches/${batchId}`)}>
                <ArrowLeft className="h-4 w-4" />
                Back to batch
              </Button>
            ) : null}
            <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">Scan Details</p>
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
        {isActiveStatus(scan.status) ? (
          <div className="mt-4 flex items-center justify-between gap-3 rounded-[20px] border border-cyan-400/20 bg-cyan-400/10 px-4 py-3 text-sm text-cyan-100">
            <div>
              <p className="font-medium text-cyan-50">Scan in progress</p>
              <p className="text-cyan-100/80">Status refreshes automatically every few seconds. Vulnerabilities will appear here as soon as the scan finishes.</p>
            </div>
            <div className="h-5 w-5 animate-spin rounded-full border-2 border-cyan-200/40 border-t-cyan-100" aria-hidden="true" />
          </div>
        ) : null}
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
          {isActiveStatus(scan.status) && !rows.length && vulnsQuery.isFetching ? (
            <div className="rounded-[24px] border border-white/8 bg-white/[0.03] px-5 py-4 text-sm text-slate-400">
              Waiting for vulnerability results. We're checking the backend automatically.
            </div>
          ) : null}
          <VulnerabilityTable
            items={rows}
            sortBy={sortBy}
            sortOrder={sortOrder}
            onSort={handleSort}
            onSelect={(item) => {
              const next = new URLSearchParams(searchParams);
              next.set("selected", item.id);
              next.set("selected_package", item.package || "");
              setSearchParams(next);
            }}
          />
          <div className="flex items-center justify-between gap-4 rounded-[24px] border border-white/8 bg-white/[0.03] px-5 py-4">
            <p className="text-sm text-slate-400">
              Showing {groupedPackages.length} package{groupedPackages.length === 1 ? "" : "s"} across {rows.length} matched record{rows.length === 1 ? "" : "s"}.
            </p>
            <p className="text-xs uppercase tracking-[0.22em] text-slate-500">
              Deduplicated by CVE ID inside each package
            </p>
          </div>
        </>
      )}

      <VulnerabilityDrawer
        jobId={jobId}
        vulnerabilityId={selectedId}
        packageName={selectedPackage}
        open={Boolean(selectedId)}
        onOpenChange={(open) => {
          if (open) return;
          const next = new URLSearchParams(searchParams);
          next.delete("selected");
          next.delete("selected_package");
          setSearchParams(next);
        }}
      />
    </div>
  );
}
