import { Search } from "lucide-react";
import { Input } from "../ui/input";
import { cn } from "../../lib/utils";

const severities = ["critical", "high", "medium", "low"];

export default function FilterBar({
  severity = [],
  onSeverityToggle,
  ecosystem = "",
  onEcosystemChange,
  ecosystems = [],
  fixState = "",
  onFixStateChange,
  fixStates = [],
  search = "",
  onSearchChange,
}) {
  return (
    <div className="panel sticky top-[92px] z-20 p-4">
      <div className="grid gap-3 xl:grid-cols-[1.3fr_220px_220px]">
        <div className="flex flex-wrap items-center gap-2">
          {severities.map((item) => (
            <button
              key={item}
              className={cn(
                "rounded-full border px-3 py-2 text-xs uppercase tracking-[0.18em] transition",
                severity.includes(item)
                  ? "border-white/20 bg-white/[0.08] text-white"
                  : "border-white/10 bg-transparent text-slate-400 hover:border-white/20 hover:text-white"
              )}
              onClick={() => onSeverityToggle(item)}
            >
              {item}
            </button>
          ))}
        </div>

        <select
          value={ecosystem}
          onChange={(event) => onEcosystemChange(event.target.value)}
          className="h-11 rounded-2xl border border-white/10 bg-white/[0.04] px-4 text-sm text-slate-200 outline-none"
        >
          <option value="">All ecosystems</option>
          {ecosystems.map((item) => (
            <option key={item.value} value={item.value}>
              {item.value} ({item.count})
            </option>
          ))}
        </select>

        <select
          value={fixState}
          onChange={(event) => onFixStateChange(event.target.value)}
          className="h-11 rounded-2xl border border-white/10 bg-white/[0.04] px-4 text-sm text-slate-200 outline-none"
        >
          <option value="">All fix states</option>
          {fixStates.map((item) => (
            <option key={item.value} value={item.value}>
              {item.value} ({item.count})
            </option>
          ))}
        </select>
      </div>

      <div className="relative mt-3">
        <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-500" />
        <Input value={search} onChange={(event) => onSearchChange(event.target.value)} className="pl-10" placeholder="Search package, CVE, or summary" />
      </div>
    </div>
  );
}
