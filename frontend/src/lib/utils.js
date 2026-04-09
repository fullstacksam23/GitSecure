import { clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs) {
  return twMerge(clsx(inputs));
}

export function formatDate(value, options = {}) {
  if (!value) return "N/A";
  const hasGranularOptions = ["weekday", "year", "month", "day", "hour", "minute", "second"].some(
    (key) => key in options
  );
  const defaults = hasGranularOptions
    ? {}
    : {
        dateStyle: "medium",
        timeStyle: "short",
      };
  return new Intl.DateTimeFormat("en-US", {
    ...defaults,
    ...options,
  }).format(new Date(value));
}

export function compactNumber(value) {
  return new Intl.NumberFormat("en-US", {
    notation: "compact",
    maximumFractionDigits: 1,
  }).format(value ?? 0);
}

export function normalizeRisk(value) {
  const numeric = Number(value);
  if (!Number.isFinite(numeric)) return 0;
  const scaled = numeric <= 1 ? numeric * 100 : numeric;
  return Math.max(0, Math.min(100, scaled));
}

export function formatRisk(value) {
  return Math.round(normalizeRisk(value) * 100) / 100;
}

export function getSeverityRiskFallback(severity) {
  switch (normalizeSeverity(severity)) {
    case "critical":
      return 95;
    case "high":
      return 80;
    case "medium":
      return 55;
    case "low":
      return 25;
    default:
      return 0;
  }
}

export function normalizeSeverity(value) {
  const text = String(value || "").toLowerCase();
  if (text.includes("critical")) return "critical";
  if (text.includes("high")) return "high";
  if (text.includes("medium") || text.includes("moderate")) return "medium";
  if (text.includes("low") || text.includes("negligible")) return "low";
  return "unknown";
}

export function severityOrder(value) {
  switch (normalizeSeverity(value)) {
    case "critical":
      return 1;
    case "high":
      return 2;
    case "medium":
      return 3;
    case "low":
      return 4;
    default:
      return 5;
  }
}

export const severityTheme = {
  critical: {
    label: "Critical",
    dot: "bg-red-500",
    text: "text-red-300",
    ring: "ring-red-500/30",
    badge: "border-red-500/20 bg-red-500/10 text-red-300",
    accent: "bg-red-500",
    soft: "bg-red-500/14",
    bar: "#ef4444",
  },
  high: {
    label: "High",
    dot: "bg-orange-400",
    text: "text-orange-300",
    ring: "ring-orange-400/30",
    badge: "border-orange-400/20 bg-orange-400/10 text-orange-200",
    accent: "bg-orange-400",
    soft: "bg-orange-400/14",
    bar: "#fb923c",
  },
  medium: {
    label: "Medium",
    dot: "bg-yellow-400",
    text: "text-yellow-200",
    ring: "ring-yellow-400/30",
    badge: "border-yellow-400/20 bg-yellow-400/10 text-yellow-100",
    accent: "bg-yellow-400",
    soft: "bg-yellow-400/14",
    bar: "#facc15",
  },
  low: {
    label: "Low",
    dot: "bg-emerald-400",
    text: "text-emerald-300",
    ring: "ring-emerald-400/30",
    badge: "border-emerald-400/20 bg-emerald-400/10 text-emerald-200",
    accent: "bg-emerald-400",
    soft: "bg-emerald-400/14",
    bar: "#34d399",
  },
  unknown: {
    label: "Unknown",
    dot: "bg-slate-500",
    text: "text-slate-300",
    ring: "ring-slate-500/30",
    badge: "border-white/10 bg-white/5 text-slate-300",
    accent: "bg-slate-500",
    soft: "bg-white/8",
    bar: "#94a3b8",
  },
};

export function formatCommitHash(value) {
  if (!value) return "No commit";
  return value.slice(0, 8);
}

export function toPercent(value, total) {
  if (!total) return 0;
  return Math.round((value / total) * 100);
}
