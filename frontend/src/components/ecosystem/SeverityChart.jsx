import { Bar, BarChart, Cell, Pie, PieChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { severityTheme } from "../../lib/utils";

const keys = ["critical", "high", "medium", "low"];

export default function SeverityChart({ counts = {}, variant = "bar" }) {
  const data = keys.map((key) => ({
    key,
    label: severityTheme[key].label,
    value: Number(counts[key] || 0),
    color: severityTheme[key].bar,
  }));

  return (
    <div className="panel p-6">
      <div className="mb-5">
        <h3 className="text-lg font-semibold text-white">Severity Breakdown</h3>
        <p className="text-sm text-slate-400">Calculated only from repositories in this batch.</p>
      </div>

      <div className="h-[260px] min-w-0">
        <ResponsiveContainer width="100%" height="100%">
          {variant === "pie" ? (
            <PieChart>
              <Pie data={data} dataKey="value" nameKey="label" innerRadius={70} outerRadius={104} paddingAngle={4}>
                {data.map((entry) => (
                  <Cell key={entry.key} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip contentStyle={{ borderRadius: 16, border: "1px solid rgba(255,255,255,0.08)", background: "#0d141d" }} />
            </PieChart>
          ) : (
            <BarChart data={data} barSize={46}>
              <XAxis dataKey="label" tick={{ fill: "#94a3b8", fontSize: 12 }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fill: "#94a3b8", fontSize: 12 }} axisLine={false} tickLine={false} allowDecimals={false} />
              <Tooltip contentStyle={{ borderRadius: 16, border: "1px solid rgba(255,255,255,0.08)", background: "#0d141d" }} />
              <Bar dataKey="value" radius={[16, 16, 6, 6]}>
                {data.map((entry) => (
                  <Cell key={entry.key} fill={entry.color} />
                ))}
              </Bar>
            </BarChart>
          )}
        </ResponsiveContainer>
      </div>

      <div className="mt-5 grid gap-3 sm:grid-cols-2">
        {data.map((item) => (
          <div key={item.key} className="surface flex items-center justify-between px-4 py-3">
            <div className="flex items-center gap-3">
              <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: item.color }} />
              <span className="text-sm text-slate-300">{item.label}</span>
            </div>
            <span className="text-sm font-medium text-white">{item.value}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
