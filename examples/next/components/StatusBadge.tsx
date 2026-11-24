import { cn } from "@/lib/cn";
import type { Status } from "@/lib/departures";

// A lookup table is a variable, so each name goes through cn() to be read by the checker.
const badge: Record<Status, string> = {
  ok: cn("tg-badge(ok)"),
  late: cn("tg-badge(late)"),
  cancelled: cn("tg-badge(cancelled)"),
};

const words: Record<Status, string> = { ok: "On time", late: "Late", cancelled: "Cancelled" };

export function StatusBadge({ status, className }: { status: Status; className?: string }) {
  return <span className={cn(badge[status], className)}>{words[status]}</span>;
}
