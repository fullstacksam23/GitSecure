import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import { SeverityDonut, TrendChart } from "../components/dashboard/OverviewCharts";
import RecentScansPanel from "../components/dashboard/RecentScansPanel";
import StatsCards from "../components/dashboard/StatsCards";
import TopRiskPackages from "../components/dashboard/TopRiskPackages";
import EmptyState from "../components/shared/EmptyState";
import PageSkeleton from "../components/shared/PageSkeleton";

export default function DashboardPage() {
  const summaryQuery = useQuery({
    queryKey: ["dashboard-summary"],
    queryFn: api.getDashboardSummary,
  });

  if (summaryQuery.isLoading) return <PageSkeleton cards={4} rows={3} />;
  if (summaryQuery.isError) {
    return <EmptyState title="Dashboard unavailable" description={summaryQuery.error.message} />;
  }

  const summary = summaryQuery.data;

  return (
    <div className="space-y-6 pb-8">
      <div>
        <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">Overview</p>
        <h1 className="mt-3 text-3xl font-semibold text-white">Repository security, severity-first.</h1>
        <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-400">
          A high-signal dashboard for scan volume, active risk, recent runs, and the packages driving exposure.
        </p>
      </div>

      <StatsCards summary={summary} />

      <div className="grid gap-6 xl:grid-cols-[1.35fr_0.95fr]">
        <RecentScansPanel scans={summary.recent_scans} />
        <SeverityDonut summary={summary} />
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr] items-stretch">
        <TrendChart summary={summary} />
        <TopRiskPackages items={summary.top_risk_packages} />
      </div>
    </div>
  );
}
