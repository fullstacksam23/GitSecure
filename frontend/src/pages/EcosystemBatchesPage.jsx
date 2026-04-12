import { Database, Layers3 } from "lucide-react";
import { startTransition, useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import BatchTable from "../components/ecosystem/BatchTable";
import DashboardCard from "../components/ecosystem/DashboardCard";
import EmptyState from "../components/shared/EmptyState";
import PageSkeleton from "../components/shared/PageSkeleton";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";
import { useEcosystemBatchesQuery } from "../hooks/useEcosystemQueries";
import { useDebouncedValue } from "../hooks/useDebouncedValue";
import { getBatchProgress } from "../lib/ecosystem";

export default function EcosystemBatchesPage() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const [searchInput, setSearchInput] = useState(searchParams.get("search") || "");
  const debouncedSearch = useDebouncedValue(searchInput, 300);

  const page = Number(searchParams.get("page") || 1);
  const status = searchParams.get("status") || "";
  const language = searchParams.get("language") || "";

  const batchesQuery = useEcosystemBatchesQuery({
    page,
    pageSize: 12,
    status,
    language,
    search: debouncedSearch,
  });

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

  if (batchesQuery.isLoading) return <PageSkeleton cards={4} rows={6} />;
  if (batchesQuery.isError) {
    return <EmptyState title="Unable to load ecosystem batches" description={batchesQuery.error.message} />;
  }

  const { items, pagination } = batchesQuery.data;
  const totalRepos = items.reduce((sum, item) => sum + getBatchProgress(item).total, 0);
  const runningCount = items.filter((item) => String(item.status).toLowerCase() === "running").length;

  return (
    <div className="space-y-6 pb-8">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">Ecosystem Scanning</p>
          <h1 className="mt-3 text-3xl font-semibold text-white">Batch-based security coverage</h1>
          <p className="mt-2 max-w-3xl text-sm leading-6 text-slate-400">
            Each batch is an independent aggregation boundary for repository coverage, vulnerability metrics, and scan progress.
          </p>
        </div>
        <div className="flex flex-wrap gap-3">
          <Button onClick={() => navigate("/ecosystem/new-scan")}>New Ecosystem Scan</Button>
          <Input value={searchInput} placeholder="Search by batch id" onChange={(event) => setSearchInput(event.target.value)} className="w-[240px]" />
          <select className="h-11 rounded-2xl border border-white/10 bg-white/[0.04] px-4 text-sm text-slate-200" value={language} onChange={(event) => updateParam("language", event.target.value)}>
            <option value="">All languages</option>
            <option value="go">Go</option>
            <option value="javascript">JavaScript</option>
            <option value="typescript">TypeScript</option>
            <option value="python">Python</option>
            <option value="rust">Rust</option>
          </select>
          <select className="h-11 rounded-2xl border border-white/10 bg-white/[0.04] px-4 text-sm text-slate-200" value={status} onChange={(event) => updateParam("status", event.target.value)}>
            <option value="">All statuses</option>
            <option value="queued">Queued</option>
            <option value="running">Running</option>
            <option value="completed">Completed</option>
          </select>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <DashboardCard title="Visible Batches" value={items.length} description="Current page of ecosystem batches ready for review." accent="bg-cyan-400/10 text-cyan-200" hint="List view" />
        <DashboardCard title="Repositories in View" value={totalRepos} description="Repository count across the currently loaded batches." accent="bg-orange-400/10 text-orange-200" hint="Page scoped" />
        <DashboardCard title="Running Batches" value={runningCount} description="Live progress polling stays active while batches are still scanning." accent="bg-emerald-400/10 text-emerald-200" hint="Auto-refresh" />
        <DashboardCard title="Entry Point" value="Per batch" description="Aggregation never rolls up across unrelated ecosystem batches." accent="bg-white/5 text-slate-200" hint="Isolation enforced" />
      </div>

      {items.length ? (
        <BatchTable items={items} />
      ) : (
        <EmptyState title="No batches found" description="Try broadening the search, status, or language filters." />
      )}

      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="inline-flex items-center gap-2 rounded-2xl border border-white/10 bg-white/[0.03] px-4 py-3 text-sm text-slate-400">
          <Layers3 className="h-4 w-4 text-cyan-300" />
          Page {pagination.page} of {pagination.total_pages || 1}
        </div>
        <div className="inline-flex items-center gap-3">
          <div className="inline-flex items-center gap-2 rounded-2xl border border-white/10 bg-white/[0.03] px-4 py-3 text-sm text-slate-400">
            <Database className="h-4 w-4 text-cyan-300" />
            {pagination.total_items || items.length} total batches
          </div>
          <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => updateParam("page", String(page - 1))}>
            Previous
          </Button>
          <Button variant="outline" size="sm" disabled={page >= (pagination.total_pages || 1)} onClick={() => updateParam("page", String(page + 1))}>
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
