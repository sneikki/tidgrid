import type { Metadata } from "next";
import { Barlow, Barlow_Condensed, IBM_Plex_Mono } from "next/font/google";
import "./tidgrid.css";

// next/font writes the @font-face and exposes the family as a custom property.
// The tokens in tokens.tg read these three names.
const display = Barlow_Condensed({ subsets: ["latin"], weight: ["700", "800"], variable: "--font-display" });
const body = Barlow({ subsets: ["latin"], weight: ["400", "600", "700"], variable: "--font-body" });
const mono = IBM_Plex_Mono({ subsets: ["latin"], weight: ["600"], variable: "--font-mono" });

export const metadata: Metadata = {
  title: "Skerry Line",
  description: "Ferry departures from Hallsund to the outer islands.",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={`${display.variable} ${body.variable} ${mono.variable}`}>
      <body>{children}</body>
    </html>
  );
}
