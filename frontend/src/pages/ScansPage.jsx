import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { api } from "../api/client";
import ScanList from "../components/jobs/ScanList";
import EmptyState from "../components/shared/EmptyState";
import PageSkeleton from "../components/shared/PageSkeleton";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";

export default function ScansPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const page = Number(searchParams.get("page") || 1);
  const repo = searchParams.get("repo") || "";
  const status = searchParams.get("status") || "";

  const scansQuery = useQuery({
    queryKey: ["scans", page, repo, status],
    queryFn: () => api.getScans({ page, pageSize: 12, repo, status }),
  });

  function updateParam(key, value) {
    const next = new URLSearchParams(searchParams);
    if (value) next.set(key, value);
    else next.delete(key);
    if (key !== "page") next.set("page", "1");
    setSearchParams(next);
  }

  if (scansQuery.isLoading) return <PageSkeleton cards={3} rows={6} />;
  if (scansQuery.isError) return <EmptyState title="Unable to load scans" description={scansQuery.error.message} />;

  const { items, pagination } = scansQuery.data;

  return (
    <div className="space-y-6 pb-8">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">All Scans</p>
          <h1 className="mt-3 text-3xl font-semibold text-white">Operational scan inventory</h1>
        </div>
        <div className="flex gap-3">
          <Input defaultValue={repo} placeholder="Filter by repository" onKeyDown={(event) => event.key === "Enter" && updateParam("repo", event.currentTarget.value)} />
          <select className="app-select" value={status} onChange={(event) => updateParam("status", event.target.value)}>
            <option value="">All statuses</option>
            <option value="queued">Queued</option>
            <option value="running">Running</option>
            <option value="completed">Completed</option>
            <option value="failed">Failed</option>
          </select>
        </div>
      </div>

      {items.length ? <ScanList items={items} /> : <EmptyState title="No scans found" description="Try broadening the repository or status filters." />}

      <div className="flex items-center justify-between">
        <p className="text-sm text-slate-500">
          Page {pagination.page} of {pagination.total_pages || 1}
        </p>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => updateParam("page", String(page - 1))}>
            Previous
          </Button>
          <Button variant="outline" size="sm" disabled={page >= pagination.total_pages} onClick={() => updateParam("page", String(page + 1))}>
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
