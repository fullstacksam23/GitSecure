import * as React from "react";
import { cva } from "class-variance-authority";
import { cn } from "../../lib/utils";

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-2xl text-sm font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-cyan-400 text-slate-950 hover:bg-cyan-300",
        secondary: "border border-white/10 bg-white/[0.05] text-white hover:bg-white/[0.08]",
        outline: "border border-white/10 bg-transparent text-slate-200 hover:bg-white/[0.05]",
        ghost: "text-slate-300 hover:bg-white/[0.05] hover:text-white",
        danger: "border border-red-500/20 bg-red-500/10 text-red-200 hover:bg-red-500/16",
      },
      size: {
        default: "h-11 px-4",
        sm: "h-9 rounded-xl px-3",
        lg: "h-12 px-5",
        icon: "h-10 w-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
);

const Button = React.forwardRef(({ className, variant, size, ...props }, ref) => (
  <button
    ref={ref}
    className={cn(buttonVariants({ variant, size, className }))}
    {...props}
  />
));

Button.displayName = "Button";

export { Button, buttonVariants };
