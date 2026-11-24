"use client";

import { useState } from "react";
import { departures, routes } from "@/lib/departures";
import { DepartureRow } from "./DepartureRow";

export function DepartureBoard() {
  const [route, setRoute] = useState<string>("All");
  const shown = departures.filter((d) => route === "All" || d.route === route);

  return (
    <section className="tg-col tg-gap(md)">
      <div className="tg-row tg-wrap tg-align(center) tg-justify(between) tg-gap(md)">
        <h2 className="tg-text(heading)">Departures from Hallsund</h2>
        <div className="tg-row tg-wrap tg-gap(sm)" role="group" aria-label="Route">
          {routes.map((r) => (
            <button
              key={r}
              type="button"
              aria-pressed={r === route}
              className="tg-action(filter)"
              onClick={() => setRoute(r)}
            >
              {r}
            </button>
          ))}
        </div>
      </div>
      <ul className="tg-surface(board)">
        {shown.map((d) => (
          <DepartureRow key={d.time + d.route} d={d} />
        ))}
      </ul>
    </section>
  );
}
