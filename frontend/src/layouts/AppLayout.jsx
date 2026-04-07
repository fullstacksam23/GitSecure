import { Menu, Radar, Shield, Sparkles } from "lucide-react";
import { NavLink, Outlet } from "react-router-dom";
import { useState } from "react";
import { Button } from "../components/ui/button";
import { Card } from "../components/ui/card";
import { cn } from "../lib/utils";

const navItems = [
  { to: "/dashboard", label: "Dashboard", icon: Radar },
  { to: "/scan", label: "Scan Repo", icon: Sparkles },
  { to: "/vulnerabilities", label: "Vulnerabilities", icon: Shield },
];

export default function AppLayout() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <div className="min-h-screen">
      <div className="page-shell">
        <div className="grid gap-6 lg:grid-cols-[280px_minmax(0,1fr)]">
          <aside className="hidden lg:block">
            <Sidebar />
          </aside>

          <main className="min-w-0">
            <div className="mb-4 flex items-center justify-between lg:hidden">
              <div>
                <p className="text-xs font-semibold uppercase tracking-[0.24em] text-sky-700">
                  GitSecure
                </p>
                <h2 className="font-display text-2xl font-semibold text-slate-950">
                  Security dashboard
                </h2>
              </div>
              <Button variant="outline" size="icon" onClick={() => setMobileOpen((open) => !open)}>
                <Menu className="h-5 w-5" />
              </Button>
            </div>

            {mobileOpen ? (
              <div className="mb-4 lg:hidden">
                <Sidebar onNavigate={() => setMobileOpen(false)} />
              </div>
            ) : null}

            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
}

function Sidebar({ onNavigate }) {
  return (
    <Card className="sticky top-6 overflow-hidden">
      <div className="border-b border-border/80 bg-gradient-to-br from-slate-950 via-slate-900 to-sky-950 p-6 text-white">
        <p className="text-xs font-semibold uppercase tracking-[0.24em] text-sky-200">GitSecure</p>
        <h1 className="mt-3 font-display text-2xl font-semibold">Vulnerability Command Center</h1>
        <p className="mt-3 text-sm leading-6 text-slate-300">
          Monitor scans, inspect findings, and keep repository risk visible.
        </p>
      </div>

      <nav className="space-y-2 p-4">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            onClick={onNavigate}
            className={({ isActive }) =>
              cn(
                "flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-semibold transition",
                isActive
                  ? "bg-slate-950 text-white shadow-sm"
                  : "text-slate-600 hover:bg-secondary hover:text-slate-950"
              )
            }
          >
            <item.icon className="h-4 w-4" />
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>

      <div className="px-4 pb-4">
        <div className="rounded-2xl bg-gradient-to-br from-amber-100 to-orange-50 p-4">
          <p className="text-xs font-semibold uppercase tracking-[0.18em] text-amber-700">
            Backend status
          </p>
          <p className="mt-2 text-sm leading-6 text-slate-700">
            Scan creation is implemented today. Dashboard read endpoints are documented in
            `frontend/docs/frontend-api-contract.md`.
          </p>
        </div>
      </div>
    </Card>
  );
}
