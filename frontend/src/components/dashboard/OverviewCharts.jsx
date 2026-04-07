import {
  Area,
  AreaChart,
  Cell,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { formatDate, severityTheme } from "../../lib/utils";

const severities = ["critical", "high", "medium", "low"];

export function SeverityDonut({ summary }) {
  const data = severities.map((severity) => ({
    name: severityTheme[severity].label,
    value: summary?.severity_distribution?.[severity] || 0,
    color: severityTheme[severity].bar,
  }));

  return (
    <div className="panel min-w-0 w-full p-6">
      <div className="mb-5">
        <h3 className="text-lg font-semibold text-white">Severity Distribution</h3>
        <p className="text-sm text-slate-400">Critical issues stay visually dominant.</p>
      </div>
      <div className="w-full min-w-0 h-[220px]">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie data={data} innerRadius={78} outerRadius={108} dataKey="value" paddingAngle={5}>
              {data.map((entry) => (
                <Cell key={entry.name} fill={entry.color} />
              ))}
            </Pie>
            <Tooltip contentStyle={{ borderRadius: 16, border: "1px solid rgba(255,255,255,0.08)", background: "#0d141d" }} />
          </PieChart>
        </ResponsiveContainer>
      </div>
      <div className="grid grid-cols-2 gap-3">
        {data.map((item) => (
          <div key={item.name} className="surface flex items-center justify-between px-3 py-2 text-sm">
            <span className="text-slate-400">{item.name}</span>
            <span className="text-white">{item.value}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

export function TrendChart({ summary }) {
  return (
    <div className="panel min-w-0 p-6">
      <div className="mb-5">
        <h3 className="text-lg font-semibold text-white">Risk Trend</h3>
        <p className="text-sm text-slate-400">Scans and findings over the last seven days.</p>
      </div>
      <div className="min-w-0 h-[240px]">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={summary?.risk_trend || []}>
            <defs>
              <linearGradient id="trendScans" x1="0" x2="0" y1="0" y2="1">
                <stop offset="5%" stopColor="#22d3ee" stopOpacity={0.4} />
                <stop offset="95%" stopColor="#22d3ee" stopOpacity={0} />
              </linearGradient>
              <linearGradient id="trendVulns" x1="0" x2="0" y1="0" y2="1">
                <stop offset="5%" stopColor="#fb923c" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#fb923c" stopOpacity={0} />
              </linearGradient>
            </defs>
            <XAxis dataKey="date" tickFormatter={(value) => formatDate(value, { month: "short", day: "numeric" })} tick={{ fill: "#94a3b8", fontSize: 12 }} axisLine={false} tickLine={false} />
            <YAxis tick={{ fill: "#94a3b8", fontSize: 12 }} axisLine={false} tickLine={false} />
            <Tooltip contentStyle={{ borderRadius: 16, border: "1px solid rgba(255,255,255,0.08)", background: "#0d141d" }} />
            <Area type="monotone" dataKey="vulnerabilities" stroke="#fb923c" fill="url(#trendVulns)" strokeWidth={2} />
            <Area type="monotone" dataKey="scans" stroke="#22d3ee" fill="url(#trendScans)" strokeWidth={2} />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
