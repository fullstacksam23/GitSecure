import { Shield, ShieldAlert, ShieldCheck, Wrench } from "lucide-react";
import { compactNumber } from "../../lib/utils";

const stats = [
  { key: "total_scans", label: "Total Scans", icon: Shield },
  { key: "critical", label: "Critical Vulns", icon: ShieldAlert },
  { key: "high", label: "High Vulns", icon: ShieldAlert },
  { key: "packages_fixed", label: "Packages Fixed", icon: Wrench },
];

export default function StatsCards({ summary }) {
  return (
    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      {stats.map((stat) => {
        const Icon = stat.icon || ShieldCheck;
        return (
          <div key={stat.key} className="panel group p-5 transition hover:-translate-y-0.5 hover:border-white/20">
            <div className="flex items-center justify-between">
              <p className="text-sm text-slate-400">{stat.label}</p>
              <div className="rounded-2xl border border-white/10 bg-white/[0.04] p-2 text-slate-300">
                <Icon className="h-4 w-4" />
              </div>
            </div>
            <p className="mt-6 text-3xl font-semibold text-white">{compactNumber(summary?.[stat.key] || 0)}</p>
          </div>
        );
      })}
    </div>
  );
}
