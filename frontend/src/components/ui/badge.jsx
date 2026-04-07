import { cva } from "class-variance-authority";
import { cn } from "../../lib/utils";

const badgeVariants = cva(
  "inline-flex items-center rounded-full border px-2.5 py-1 text-[11px] font-medium uppercase tracking-[0.18em]",
  {
    variants: {
      variant: {
        default: "border-white/10 bg-white/[0.06] text-slate-100",
        critical: "border-red-500/20 bg-red-500/10 text-red-300",
        high: "border-orange-400/20 bg-orange-400/10 text-orange-200",
        medium: "border-yellow-400/20 bg-yellow-400/10 text-yellow-100",
        low: "border-emerald-400/20 bg-emerald-400/10 text-emerald-200",
        neutral: "border-white/10 bg-white/[0.05] text-slate-300",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
);

export function Badge({ className, variant, ...props }) {
  return <div className={cn(badgeVariants({ variant, className }))} {...props} />;
}
