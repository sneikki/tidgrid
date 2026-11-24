import { departures, expectedTime, now } from "@/lib/departures";

export function NextBoat() {
  const d = departures.find((x) => x.status !== "cancelled" && x.time >= now);
  if (!d) return null;
  return (
    <div className="tg-surface(next) tg-col tg-gap(xs) tg-pad(md) md:tg-min-width(panel)">
      <span className="tg-text(label)">Next boat</span>
      <span className="tg-text(clock)">{expectedTime(d)}</span>
      <span>
        {d.to}, {d.vessel}
      </span>
    </div>
  );
}
