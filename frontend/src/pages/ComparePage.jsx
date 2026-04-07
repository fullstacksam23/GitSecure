import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "react-router-dom";
import { api } from "../api/client";
import SeverityBadge from "../components/shared/SeverityBadge";
import EmptyState from "../components/shared/EmptyState";
import PageSkeleton from "../components/shared/PageSkeleton";
import { compactNumber } from "../lib/utils";

export default function ComparePage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const base = searchParams.get("base") || "";
  const target = searchParams.get("target") || "";

  const scansQuery = useQuery({
    queryKey: ["compare-scans-list"],
    queryFn: () => api.getScans({ page: 1, pageSize: 50 }),
  });

  const compareQuery = useQuery({
    queryKey: ["compare", base, target],
    queryFn: () => api.compareScans({ base, target }),
    enabled: Boolean(base && target),
  });

  const options = scansQuery.data?.items || [];

  function update(key, value) {
    const next = new URLSearchParams(searchParams);
    if (value) next.set(key, value);
    else next.delete(key);
    setSearchParams(next);
  }

  return (
    <div className="space-y-6 pb-8">
      <div>
        <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">Compare</p>
        <h1 className="mt-3 text-3xl font-semibold text-white">Scan drift and remediation delta</h1>
      </div>

      <div className="panel grid gap-4 p-6 md:grid-cols-2">
        <select className="h-11 rounded-2xl border border-white/10 bg-white/[0.04] px-4 text-sm text-slate-200" value={base} onChange={(event) => update("base", event.target.value)}>
          <option value="">Select base scan</option>
          {options.map((item) => (
            <option key={item.job_id} value={item.job_id}>
              {item.repo} · {item.job_id.slice(0, 8)}
            </option>
          ))}
        </select>
        <select className="h-11 rounded-2xl border border-white/10 bg-white/[0.04] px-4 text-sm text-slate-200" value={target} onChange={(event) => update("target", event.target.value)}>
          <option value="">Select target scan</option>
          {options.map((item) => (
            <option key={item.job_id} value={item.job_id}>
              {item.repo} · {item.job_id.slice(0, 8)}
            </option>
          ))}
        </select>
      </div>

      {!base || !target ? (
        <EmptyState title="Choose two scans to compare" description="The compare view shows newly introduced, fixed, and persisting vulnerabilities." />
      ) : compareQuery.isLoading ? (
        <PageSkeleton cards={3} rows={3} />
      ) : compareQuery.isError ? (
        <EmptyState title="Unable to compare scans" description={compareQuery.error.message} />
      ) : (
        <>
          <div className="grid gap-4 md:grid-cols-3">
            {[
              ["New", compareQuery.data.new.count],
              ["Fixed", compareQuery.data.fixed.count],
              ["Persisting", compareQuery.data.persisting.count],
            ].map(([label, value]) => (
              <div key={label} className="panel p-5">
                <p className="text-sm text-slate-500">{label}</p>
                <p className="mt-4 text-3xl font-semibold text-white">{compactNumber(value)}</p>
              </div>
            ))}
          </div>

          <div className="grid gap-6 xl:grid-cols-3">
            <CompareBucket title="New vulnerabilities" items={compareQuery.data.new.items} />
            <CompareBucket title="Fixed vulnerabilities" items={compareQuery.data.fixed.items} />
            <CompareBucket title="Persisting vulnerabilities" items={compareQuery.data.persisting.items} />
          </div>
        </>
      )}
    </div>
  );
}

function CompareBucket({ title, items }) {
  return (
    <div className="panel p-6">
      <h3 className="mb-4 text-lg font-semibold text-white">{title}</h3>
      <div className="space-y-3">
        {items.map((item) => (
          <div key={`${title}-${item.id}-${item.package}`} className="surface p-4">
            <div className="flex items-center justify-between gap-3">
              <p className="font-mono text-sm text-white">{item.package}</p>
              <SeverityBadge severity={item.normalized_severity || item.severity} />
            </div>
            <p className="mt-2 text-xs text-slate-500">{item.id}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
