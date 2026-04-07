import { useEffect, useState } from "react";
import { Search, SlidersHorizontal } from "lucide-react";
import EmptyState from "../components/EmptyState";
import LoadingState from "../components/LoadingState";
import PageHeader from "../components/PageHeader";
import VulnerabilityTable from "../components/VulnerabilityTable";
import { Card, CardContent } from "../components/ui/card";
import { Input } from "../components/ui/input";
import { useAsyncData } from "../hooks/useAsyncData";
import { api } from "../services/api";

const severityOptions = ["all", "critical", "high", "medium", "low"];

export default function VulnerabilitiesPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [appliedSearch, setAppliedSearch] = useState("");
  const [severity, setSeverity] = useState("all");

  const query = useAsyncData(
    () =>
      api.getVulnerabilities({
        page,
        pageSize: 10,
        search: appliedSearch,
        severity,
      }),
    [page, appliedSearch, severity]
  );

  useEffect(() => {
    const timeout = window.setTimeout(() => {
      setPage(1);
      setAppliedSearch(search.trim());
    }, 300);

    return () => window.clearTimeout(timeout);
  }, [search]);

  const items = query.data?.items || [];
  const pagination = query.data?.pagination || { page, total_pages: 1 };

  return (
    <div className="space-y-6">
      <PageHeader
        eyebrow="Findings"
        title="Explore vulnerabilities"
        description="Search by advisory ID or package, filter by severity, and drill into the vulnerability metadata stored by the backend."
      />

      <Card>
        <CardContent className="grid gap-4 p-6 md:grid-cols-[1fr_220px]">
          <div className="relative">
            <Search className="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
            <Input
              placeholder="Search by package or vulnerability ID"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              className="pl-11"
            />
          </div>

          <div className="relative">
            <SlidersHorizontal className="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
            <select
              className="flex h-11 w-full rounded-xl border border-input bg-white/85 pl-11 pr-4 text-sm shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
              value={severity}
              onChange={(event) => {
                setPage(1);
                setSeverity(event.target.value);
              }}
            >
              {severityOptions.map((option) => (
                <option key={option} value={option}>
                  {option === "all" ? "All severities" : `${option} severity`}
                </option>
              ))}
            </select>
          </div>
        </CardContent>
      </Card>

      {query.loading ? <LoadingState rows={8} /> : null}

      {!query.loading && query.error ? (
        <EmptyState
          title="Vulnerability API not available yet"
          description={`${query.error}. Implement the documented GET /vulnerabilities endpoints and this page will render live data immediately.`}
        />
      ) : null}

      {!query.loading && !query.error && items.length === 0 ? (
        <EmptyState
          title="No vulnerabilities matched your filters"
          description="Try a broader search, clear the severity filter, or start a new repository scan."
        />
      ) : null}

      {!query.loading && !query.error && items.length > 0 ? (
        <Card>
          <CardContent className="p-0">
            <div className="p-6">
              <VulnerabilityTable
                items={items}
                page={pagination.page}
                totalPages={pagination.total_pages}
                onPageChange={setPage}
              />
            </div>
          </CardContent>
        </Card>
      ) : null}
    </div>
  );
}
