import { ArrowLeft, FolderGit2, Radar, ShieldAlert, Sparkles } from "lucide-react";
import { startTransition, useEffect, useState } from "react";
import { useNavigate, useParams, useSearchParams } from "react-router-dom";
import DashboardCard from "../components/ecosystem/DashboardCard";
import ProgressBar from "../components/ecosystem/ProgressBar";
import RepoTable from "../components/ecosystem/RepoTable";
import SeverityChart from "../components/ecosystem/SeverityChart";
import StatusBadge from "../components/ecosystem/StatusBadge";
import VulnerabilityList from "../components/ecosystem/VulnerabilityList";
import EmptyState from "../components/shared/EmptyState";
import PageSkeleton from "../components/shared/PageSkeleton";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import {
  useEcosystemBatchQuery,
  useEcosystemBatchReposQuery,
  useEcosystemBatchSummaryQuery,
} from "../hooks/useEcosystemQueries";
import { useDebouncedValue } from "../hooks/useDebouncedValue";
import {
  getBatchId,
  getBatchLanguage,
  getBatchProgress,
  getMostVulnerableRepo,
  getRepoName,
  getRepoVulnerabilityPreview,
  getSummaryAverageRisk,
  getSummarySeverityBreakdown,
  getSummaryTotalRepos,
  getSummaryTotalVulnerabilities,
  isActiveBatchStatus,
} from "../lib/ecosystem";
import { compactNumber, formatDate, formatRisk } from "../lib/utils";

const allowedSortColumns = new Set(["repo", "stars", "rank", "status", "vulnerability_count", "top_severity"]);

export default function BatchDetailPage() {
  const { batchId } = useParams();
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const [searchInput, setSearchInput] = useState(searchParams.get("search") || "");
  const debouncedSearch = useDebouncedValue(searchInput, 300);

  const page = Number(searchParams.get("page") || 1);
  const status = searchParams.get("status") || "";
  const severity = searchParams.get("severity") || "";
  const sortBy = allowedSortColumns.has(searchParams.get("sort_by")) ? searchParams.get("sort_by") : "rank";
  const sortOrder = searchParams.get("sort_order") || "asc";

  const batchQuery = useEcosystemBatchQuery(batchId);
  const shouldPoll = isActiveBatchStatus(batchQuery.data?.status);
  const summaryQuery = useEcosystemBatchSummaryQuery(batchId, shouldPoll);
  const reposQuery = useEcosystemBatchReposQuery(
    batchId,
    {
      page,
      pageSize: 20,
      search: debouncedSearch,
      status,
      severity,
      sortBy,
      sortOrder,
    },
    shouldPoll
  );

  useEffect(() => {
    const next = new URLSearchParams(searchParams);
    if (debouncedSearch) next.set("search", debouncedSearch);
    else next.delete("search");
    if (debouncedSearch !== searchParams.get("search")) {
      startTransition(() => setSearchParams(next, { replace: true }));
    }
  }, [debouncedSearch, searchParams, setSearchParams]);

  function updateParam(key, value) {
    const next = new URLSearchParams(searchParams);
    if (value) next.set(key, value);
    else next.delete(key);
    if (key !== "page") next.set("page", "1");
    startTransition(() => setSearchParams(next));
  }

  function handleSort(column) {
    const nextOrder = sortBy === column && sortOrder === "desc" ? "asc" : "desc";
    updateParam("sort_by", column);
    updateParam("sort_order", nextOrder);
  }

  const isLoading = batchQuery.isLoading || summaryQuery.isLoading || reposQuery.isLoading;
  if (isLoading) return <PageSkeleton cards={4} rows={8} />;
  if (batchQuery.isError) return <EmptyState title="Unable to load batch" description={batchQuery.error.message} />;
  if (summaryQuery.isError) return <EmptyState title="Unable to load batch summary" description={summaryQuery.error.message} />;
  if (reposQuery.isError) return <EmptyState title="Unable to load repositories" description={reposQuery.error.message} />;

  const batch = batchQuery.data;
  const summary = summaryQuery.data;
  const repoResult = reposQuery.data;
  const progress = getBatchProgress(batch);
  const severityCounts = getSummarySeverityBreakdown(summary);
  const totalRepos = getSummaryTotalRepos(summary, batch);
  const totalVulnerabilities = getSummaryTotalVulnerabilities(summary);
  const averageRisk = getSummaryAverageRisk(summary);
  const mostVulnerableRepo = getMostVulnerableRepo(summary);
  const highlightedRepoName = getRepoName(mostVulnerableRepo);
  const previewVulnerabilities = getRepoVulnerabilityPreview(mostVulnerableRepo);
  const insight = mostVulnerableRepo
    ? `${highlightedRepoName} currently leads this batch with ${mostVulnerableRepo.vulnerability_count || 0} findings.`
    : "Most vulnerable repository will appear once batch aggregation is available.";

  return (
    <div className="space-y-6 pb-8">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <Button variant="ghost" size="sm" className="mb-4" onClick={() => navigate("/ecosystem/batches")}>
            <ArrowLeft className="h-4 w-4" />
            Back to batches
          </Button>
          <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">Batch Detail</p>
          <h1 className="mt-3 text-3xl font-semibold text-white">{getBatchId(batch)}</h1>
          <div className="mt-3 flex flex-wrap items-center gap-3 text-sm text-slate-400">
            <span>{getBatchLanguage(batch)}</span>
            <span>Created {formatDate(batch.created_at)}</span>
            <span>Completed {formatDate(batch.completed_at)}</span>
            <StatusBadge status={batch.status} />
          </div>
        </div>
        <div className="panel min-w-[280px] p-5">
          <div className="flex items-center justify-between gap-3">
            <div>
              <p className="text-sm text-slate-400">Batch progress</p>
              <p className="mt-1 text-2xl font-semibold text-white">{progress.percent}%</p>
            </div>
            <Radar className="h-8 w-8 text-cyan-300" />
          </div>
          <ProgressBar
            value={progress.percent}
            className="mt-4"
            label={
              <>
                <span>{progress.completed} completed</span>
                <span>{progress.total} repositories</span>
              </>
            }
            tone={batch.status === "completed" ? "success" : "default"}
          />
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-5">
        <DashboardCard title="Total Repositories" value={compactNumber(totalRepos)} description="Repositories included only in this batch." accent="bg-cyan-400/10 text-cyan-200" hint="Batch scoped" />
        <DashboardCard title="Total Vulnerabilities" value={compactNumber(totalVulnerabilities)} description="All findings aggregated from repositories in this batch." accent="bg-orange-400/10 text-orange-200" hint="No cross-batch rollup" />
        <DashboardCard title="Critical + High" value={severityCounts.critical + severityCounts.high} description="Priority issues requiring the fastest response." accent="bg-red-500/10 text-red-200" hint="Urgent exposure" />
        <DashboardCard title="Average Risk Score" value={`${formatRisk(averageRisk)}%`} description="Average normalized risk signal across this batch." accent="bg-emerald-400/10 text-emerald-200" hint="Score average" />
        <DashboardCard title="Most Vulnerable Repo" value={highlightedRepoName} description={insight} accent="bg-white/5 text-slate-200" hint="Highest exposure" />
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
        <SeverityChart counts={severityCounts} variant="bar" />

        <div className="panel p-6">
          <div className="mb-5 flex items-center justify-between gap-4">
            <div>
              <h3 className="text-lg font-semibold text-white">Most Vulnerable Repository</h3>
              <p className="text-sm text-slate-400">A spotlight panel for the riskiest repository inside this batch.</p>
            </div>
            <Sparkles className="h-5 w-5 text-cyan-300" />
          </div>

          {mostVulnerableRepo ? (
            <>
              <div className="surface mb-4 p-4">
                <div className="flex items-center justify-between gap-4">
                  <div>
                    <p className="text-base font-semibold text-white">{highlightedRepoName}</p>
                    <p className="mt-1 text-sm text-slate-400">
                      {mostVulnerableRepo.vulnerability_count || 0} findings across this batch snapshot
                    </p>
                  </div>
                  {mostVulnerableRepo.job_id ? (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => navigate(`/scans/${mostVulnerableRepo.job_id}?batch=${encodeURIComponent(batchId)}`)}
                    >
                      Open scan
                    </Button>
                  ) : null}
                </div>
              </div>
              <VulnerabilityList
                items={previewVulnerabilities}
                emptyMessage="Top vulnerability previews will appear here if the summary endpoint includes them."
              />
            </>
          ) : (
            <EmptyState title="No repository insight yet" description="The summary endpoint has not reported a most vulnerable repository for this batch." />
          )}
        </div>
      </div>

      <div className="panel p-5">
        <div className="flex flex-wrap items-end justify-between gap-4">
          <div>
            <p className="text-sm uppercase tracking-[0.22em] text-cyan-300">Repositories</p>
            <h2 className="mt-2 text-2xl font-semibold text-white">Per-batch repository inventory</h2>
            <p className="mt-2 text-sm text-slate-400">
              Sort, search, and filter repositories within this batch only.
            </p>
          </div>
          <div className="flex flex-wrap gap-3">
            <Input value={searchInput} onChange={(event) => setSearchInput(event.target.value)} placeholder="Search repositories" className="w-[240px]" />
            <select className="h-11 rounded-2xl border border-white/10 bg-white/[0.04] px-4 text-sm text-slate-200" value={status} onChange={(event) => updateParam("status", event.target.value)}>
              <option value="">All scan statuses</option>
              <option value="queued">Queued</option>
              <option value="running">Running</option>
              <option value="completed">Completed</option>
              <option value="failed">Failed</option>
            </select>
            <select className="h-11 rounded-2xl border border-white/10 bg-white/[0.04] px-4 text-sm text-slate-200" value={severity} onChange={(event) => updateParam("severity", event.target.value)}>
              <option value="">All severities</option>
              <option value="critical">Critical</option>
              <option value="high">High</option>
              <option value="medium">Medium</option>
              <option value="low">Low</option>
            </select>
          </div>
        </div>
      </div>

      {repoResult.items.length ? (
        <RepoTable
          batchId={batchId}
          items={repoResult.items}
          sortBy={sortBy}
          sortOrder={sortOrder}
          onSort={handleSort}
          highlightedRepo={highlightedRepoName}
        />
      ) : (
        <EmptyState title="No repositories found" description="Try broadening the repo search or clearing the status and severity filters." />
      )}

      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="inline-flex items-center gap-2 rounded-2xl border border-white/10 bg-white/[0.03] px-4 py-3 text-sm text-slate-400">
          <FolderGit2 className="h-4 w-4 text-cyan-300" />
          Showing {repoResult.items.length} repositories on this page
        </div>
        <div className="inline-flex items-center gap-3">
          <div className="inline-flex items-center gap-2 rounded-2xl border border-white/10 bg-white/[0.03] px-4 py-3 text-sm text-slate-400">
            <ShieldAlert className="h-4 w-4 text-cyan-300" />
            Page {repoResult.pagination.page} of {repoResult.pagination.total_pages || 1}
          </div>
          <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => updateParam("page", String(page - 1))}>
            Previous
          </Button>
          <Button variant="outline" size="sm" disabled={page >= (repoResult.pagination.total_pages || 1)} onClick={() => updateParam("page", String(page + 1))}>
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
