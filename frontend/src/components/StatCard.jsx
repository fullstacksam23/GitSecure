import { ArrowRight, ShieldAlert } from "lucide-react";
import { compactNumber } from "../lib/utils";
import { Card, CardContent } from "./ui/card";

export default function StatCard({ title, value, hint, accent = "sky" }) {
  const tones = {
    sky: "from-sky-100 to-cyan-50 text-sky-700",
    red: "from-red-100 to-rose-50 text-red-700",
    orange: "from-orange-100 to-amber-50 text-orange-700",
    yellow: "from-yellow-100 to-orange-50 text-amber-700",
    green: "from-emerald-100 to-lime-50 text-emerald-700",
  };

  return (
    <Card className="overflow-hidden">
      <CardContent className="p-0">
        <div className={`bg-gradient-to-br ${tones[accent]} p-6`}>
          <div className="flex items-start justify-between">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.18em]">{title}</p>
              <p className="mt-3 font-display text-4xl font-semibold text-slate-950">
                {compactNumber(value)}
              </p>
            </div>
            <div className="rounded-2xl bg-white/80 p-3 text-slate-900 shadow-sm">
              <ShieldAlert className="h-5 w-5" />
            </div>
          </div>
          <div className="mt-6 flex items-center gap-2 text-sm text-slate-700">
            <ArrowRight className="h-4 w-4" />
            <span>{hint}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
