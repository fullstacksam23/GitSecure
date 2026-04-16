import { Bell, ChevronRight, Menu, Plus, Search } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useMemo, useState } from "react";
import { api } from "../../api/client";
import { cn, formatDate } from "../../lib/utils";
import SeverityBadge from "../shared/SeverityBadge";
import StatusIndicator from "../shared/StatusIndicator";
import { Button } from "../ui/button";
import { Input } from "../ui/input";

const navItems = [
  { to: "/dashboard", label: "Dashboard" },
  { to: "/ecosystem/batches", label: "Ecosystem Batches" },
  { to: "/scans", label: "All Scans" },
  { to: "/history", label: "History" },
  { to: "/compare", label: "Compare" },
];

export default function AppShell({ children }) {
  const [mobileOpen, setMobileOpen] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const scansQuery = useQuery({
    queryKey: ["sidebar-scans"],
    queryFn: () => api.getScans({ page: 1, pageSize: 8 }),
  });

  const breadcrumbs = useMemo(() => {
    const parts = location.pathname.split("/").filter(Boolean);
    if (parts.length === 0) return ["Dashboard"];
    return parts.map((part) =>
      part.length === 36 ? `Scan ${part.slice(0, 8)}` : part.replace(/-/g, " ").replace(/\b\w/g, (char) => char.toUpperCase())
    );
  }, [location.pathname]);

  const repoItems = scansQuery.data?.items || [];

  return (
    <div className="app-shell">
      <div className="mx-auto flex min-h-screen w-full max-w-[1600px] gap-6 px-4 py-4 lg:px-6">
        <aside className="hidden w-[310px] shrink-0 lg:block">
          <Sidebar items={repoItems} />
        </aside>

        <div className="min-w-0 flex-1">
          <header className="panel mb-6 flex flex-wrap items-center gap-4 px-5 py-4">
            <Button variant="outline" size="icon" className="lg:hidden" onClick={() => setMobileOpen((value) => !value)}>
              <Menu className="h-4 w-4" />
            </Button>

            <div className="flex min-w-0 items-center gap-2 text-sm text-slate-400">
              {breadcrumbs.map((crumb, index) => (
                <div key={`${crumb}-${index}`} className="flex items-center gap-2">
                  {index > 0 ? <ChevronRight className="h-4 w-4 text-slate-600" /> : null}
                  <span className={cn(index === breadcrumbs.length - 1 && "text-white")}>{crumb}</span>
                </div>
              ))}
            </div>

            <div className="relative ml-auto min-w-[240px] flex-1 md:max-w-sm">
              <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500" />
              <Input
                placeholder="Search scans, packages, CVEs"
                className="pl-10"
                onKeyDown={(event) => {
                  if (event.key === "Enter" && event.currentTarget.value.trim()) {
                    navigate(`/scans?repo=${encodeURIComponent(event.currentTarget.value.trim())}`);
                  }
                }}
              />
            </div>

            <Button variant="secondary" size="icon">
              <Bell className="h-4 w-4" />
            </Button>
            <Button onClick={() => navigate("/new-scan")}>
              <Plus className="h-4 w-4" />
              New Scan
            </Button>
          </header>

          {mobileOpen ? (
            <div className="mb-6 lg:hidden">
              <Sidebar items={repoItems} />
            </div>
          ) : null}

          {children || <Outlet />}
        </div>
      </div>
    </div>
  );
}

function Sidebar({ items }) {
  return (
    <div className="panel sticky top-4 h-[calc(100vh-2rem)] overflow-hidden">
      <div className="grid-overlay border-b border-white/10 px-6 py-6">
        <div className="flex items-center gap-3">
          <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-cyan-400/90 text-sm font-semibold text-slate-950">
            G
          </div>
          <div>
            <p className="text-xs uppercase tracking-[0.28em] text-cyan-300">GitSecure</p>
            <h1 className="text-xl font-semibold text-white">Security Platform</h1>
          </div>
        </div>
        <p className="mt-4 text-sm leading-6 text-slate-400">
          Developer-grade vulnerability visibility for repositories, packages, and scan drift.
        </p>
      </div>

      <nav className="space-y-1 p-4">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              cn(
                "flex items-center justify-between rounded-2xl px-4 py-3 text-sm transition",
                isActive ? "bg-white/[0.08] text-white" : "text-slate-400 hover:bg-white/[0.04] hover:text-white"
              )
            }
          >
            <span>{item.label}</span>
            <ChevronRight className="h-4 w-4" />
          </NavLink>
        ))}
      </nav>

      <div className="border-t border-white/10 px-4 py-4">
        <div className="mb-3 flex items-center justify-between">
          <h2 className="text-xs uppercase tracking-[0.24em] text-slate-500">Repositories</h2>
          <span className="text-xs text-slate-500">{items.length}</span>
        </div>

        <div className="space-y-3">
          {items.map((item) => (
            <NavLink
              key={item.job_id}
              to={`/scans/${item.job_id}`}
              className="surface block p-4 transition hover:border-white/20 hover:bg-white/[0.06]"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <p className="truncate text-sm font-medium text-white">{item.repo}</p>
                  <p className="mt-1 text-xs text-slate-500">{formatDate(item.created_at, { dateStyle: "short", timeStyle: "short" })}</p>
                </div>
                <SeverityBadge
                  severity={item.top_severity}
                  issueCount={Number(item.vulnerability_count || 0)}
                  status={item.status}
                />
                </div>
                <div className="mt-3 flex items-center justify-between">
                  <StatusIndicator status={item.status} />
                <span className="text-xs text-slate-400">{item.vulnerability_count} findings</span>
              </div>
            </NavLink>
          ))}
        </div>
      </div>
    </div>
  );
}
