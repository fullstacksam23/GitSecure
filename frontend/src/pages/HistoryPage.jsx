import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import EmptyState from "../components/shared/EmptyState";
import PageSkeleton from "../components/shared/PageSkeleton";
import { formatDate } from "../lib/utils";

export default function HistoryPage() {
  const scansQuery = useQuery({
    queryKey: ["history"],
    queryFn: () => api.getScans({ page: 1, pageSize: 20 }),
  });

  if (scansQuery.isLoading) return <PageSkeleton cards={2} rows={8} />;
  if (scansQuery.isError) return <EmptyState title="Unable to load history" description={scansQuery.error.message} />;

  const grouped = (scansQuery.data?.items || []).reduce((acc, item) => {
    const day = new Date(item.created_at).toDateString();
    acc[day] = acc[day] || [];
    acc[day].push(item);
    return acc;
  }, {});

  const entries = Object.entries(grouped);

  return (
    <div className="space-y-6 pb-8">
      <div>
        <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">History</p>
        <h1 className="mt-3 text-3xl font-semibold text-white">Timeline of scan execution</h1>
      </div>

      {entries.length === 0 ? (
        <EmptyState title="No scan history yet" description="Launch a scan to start building the timeline." />
      ) : (
        <div className="space-y-6">
          {entries.map(([day, items]) => (
            <section key={day} className="panel p-6">
              <div className="mb-4 flex items-center justify-between">
                <h2 className="text-lg font-semibold text-white">{day}</h2>
                <span className="text-sm text-slate-500">{items.length} runs</span>
              </div>
              <div className="space-y-3">
                {items.map((item) => (
                  <div key={item.job_id} className="surface flex items-center justify-between gap-4 p-4">
                    <div>
                      <p className="text-sm font-medium text-white">{item.repo}</p>
                      <p className="mt-1 text-xs text-slate-500">{formatDate(item.created_at)}</p>
                    </div>
                    <div className="text-sm text-slate-400">{item.vulnerability_count} findings</div>
                  </div>
                ))}
              </div>
            </section>
          ))}
        </div>
      )}
    </div>
  );
}
