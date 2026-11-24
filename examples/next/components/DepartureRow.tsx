import type { Departure } from "@/lib/departures";
import { StatusBadge } from "./StatusBadge";

export function DepartureRow({ d }: { d: Departure }) {
  return (
    <li className="tg-surface(row) tg-row tg-wrap tg-align(center) tg-gap(sm) tg-pad(md) md:tg-gap(md)">
      <time className="tg-text(time) tg-grow md:tg-fixed md:tg-width(slot) md:tg-order(first)">{d.time}</time>
      <StatusBadge status={d.status} className="md:tg-order(last)" />
      <div className="tg-col tg-width(full) md:tg-grow">
        <span className="tg-text(heading)">{d.to}</span>
        <span>
          {d.route}, {d.vessel}. {d.note}.
        </span>
      </div>
    </li>
  );
}
