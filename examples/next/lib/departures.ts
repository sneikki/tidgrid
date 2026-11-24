export type Status = "ok" | "late" | "cancelled";

export type Departure = {
  time: string;
  route: string;
  to: string;
  vessel: string;
  status: Status;
  note: string;
};

export const routes = ["All", "Outer Loop", "Seal Rock", "Lantern Isle"] as const;

// The board is a fixed Tuesday so the page renders the same every time.
export const now = "09:15";

export const departures: Departure[] = [
  { time: "07:40", route: "Outer Loop", to: "Brännskär", vessel: "Marta K", status: "ok", note: "Departed" },
  { time: "08:20", route: "Seal Rock", to: "Sälgrund", vessel: "Orn", status: "ok", note: "Departed" },
  { time: "09:30", route: "Lantern Isle", to: "Fyrholm", vessel: "Linnea", status: "late", note: "Fog on the sound, +25 min" },
  { time: "10:05", route: "Outer Loop", to: "Brännskär", vessel: "Marta K", status: "ok", note: "Boarding from 09:50" },
  { time: "11:15", route: "Seal Rock", to: "Sälgrund", vessel: "Orn", status: "cancelled", note: "Crane repair at Sälgrund" },
  { time: "12:00", route: "Lantern Isle", to: "Fyrholm", vessel: "Linnea", status: "ok", note: "Cars and bicycles" },
  { time: "13:40", route: "Outer Loop", to: "Brännskär", vessel: "Marta K", status: "ok", note: "Last cars 13:25" },
  { time: "15:10", route: "Seal Rock", to: "Sälgrund", vessel: "Orn", status: "late", note: "Waits for the mail boat, +10 min" },
  { time: "17:30", route: "Lantern Isle", to: "Fyrholm", vessel: "Linnea", status: "ok", note: "Evening crossing" },
];
