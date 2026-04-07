import { ShieldAlert } from "lucide-react";

export default function EmptyState({ title, description, icon: Icon = ShieldAlert }) {
  return (
    <div className="panel grid-overlay flex min-h-[240px] flex-col items-center justify-center gap-4 p-10 text-center">
      <div className="rounded-2xl border border-white/10 bg-white/[0.05] p-4">
        <Icon className="h-6 w-6 text-slate-300" />
      </div>
      <div className="space-y-2">
        <h3 className="text-lg font-semibold text-white">{title}</h3>
        <p className="max-w-lg text-sm leading-6 text-slate-400">{description}</p>
      </div>
    </div>
  );
}
