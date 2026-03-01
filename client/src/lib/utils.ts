// src/lib/utils.ts

import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Safely merges Tailwind classes without style conflicts.
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Maps the user's competitive tier to a specific Tailwind text color.
 */
export function getTierColor(tier: string): string {
  switch (tier.toLowerCase()) {
    case "diamond":
      return "text-cyan-500";
    case "platinum":
      return "text-emerald-500";
    case "gold":
      return "text-yellow-500";
    case "silver":
      return "text-slate-400";
    case "bronze":
      return "text-amber-600";
    default:
      return "text-gray-500";
  }
}