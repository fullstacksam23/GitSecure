import { ArrowUpRight } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";

export default function DashboardCard({ title, value, description, accent, hint }) {
  return (
    <Card className="overflow-hidden">
      <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-4">
        <div>
          <p className="text-sm text-slate-400">{title}</p>
          <CardTitle className="mt-3 text-3xl">{value}</CardTitle>
        </div>
        <div className={`rounded-2xl border border-white/10 px-3 py-2 text-xs uppercase tracking-[0.2em] ${accent}`}>
          <ArrowUpRight className="h-4 w-4" />
        </div>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-slate-300">{description}</p>
        {hint ? <p className="mt-3 text-xs uppercase tracking-[0.18em] text-slate-500">{hint}</p> : null}
      </CardContent>
    </Card>
  );
}
