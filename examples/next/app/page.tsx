import { DepartureBoard } from "@/components/DepartureBoard";
import { NextBoat } from "@/components/NextBoat";
import { Notice } from "@/components/Notice";

export default function Page() {
  return (
    <main className="tg-col tg-gap(lg) tg-max-width(page) tg-center tg-pad(md) md:tg-pad(lg)">
      <header className="tg-surface(masthead) tg-col md:tg-row tg-gap(md) tg-pad-b(lg) md:tg-align(end) md:tg-justify(between)">
        <div className="tg-col tg-gap(sm)">
          <h1 className="tg-text(hero)">Skerry Line</h1>
          <p>Ferries from Hallsund harbour to the three outer islands. Timetable for Tuesday.</p>
        </div>
        <NextBoat />
      </header>
      <Notice title="Fog on the sound">
        Lantern Isle crossings run about 25 minutes late until noon. The 11:15 to Sälgrund is cancelled.
      </Notice>
      <DepartureBoard />
      <footer>
        <p className="tg-text(label)">Skerry Line runs three boats. Cars board 15 minutes before departure.</p>
      </footer>
    </main>
  );
}
