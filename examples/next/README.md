# Tidgrid in Next.js

A ferry timetable (Next.js App Router, React, TypeScript) styled only with
Tidgrid: no Tailwind, no CSS modules, no own CSS.

```
npm install
npm run dev      # or: npm run build && npm start
```

Needs the `tidgrid` binary (`devDependencies`), or one on your `PATH`. Both
`dev` and `build` first run `npm run tg` (`tidgrid build && tidgrid check`),
so a stale stylesheet or a class outside the vocabulary stops the run.

## Where things are

- `index.tg`: entry file, `project()`, breakpoints, base elements.
  `tokens.tg` holds colours (as `light-dark()` pairs, so dark mode needs no
  second stylesheet), scales and fonts; `roles.tg` the roles and mixins.
- `app/tidgrid.css`: generated, committed, imported once in `app/layout.tsx`.
  Never edit it. Committed so CI can fail when it differs from `tidgrid build`.
- `components/`, `lib/cn.ts`: the page parts and the few-line `cn` helper.

## What differs from a plain HTML project

- **Scan paths.** `scan = ["app", "components", "lib"]`, not `["."]`: the
  scanner does not skip `.next`, and Next's build output is full of `style`
  attributes that fail the check.
- **Helper list.** `helpers = ["cn"]`. A lookup table of class names is a
  variable the checker cannot follow, so each name goes through `cn()` where
  the table is written (`components/StatusBadge.tsx`). A name in any other
  prop (`tone="tg-text(label)"`) or a plain object is never read: not checked,
  and it does not count as use.
- **Client state.** `components/DepartureBoard.tsx` is a client component. Its
  filter buttons carry `aria-pressed`, which the `filter` role styles with
  `&[aria-pressed=true]`, so no state role is needed. A class a script sets
  needs a `state` role and the name written whole. Server and client files are
  scanned the same.
- **Fonts.** `next/font` writes the `@font-face` and exposes each family as a
  custom property (`variable: "--font-display"`); `tokens.tg` points the
  `font` tokens at those names. Nothing in the stylesheet or the checker knows
  the font's own hashed class, so use `variable`, not `className`.
- **Stylesheet edits.** A `.tg` edit needs `npm run tg` (or a restart of `dev`) before the page shows it.
- **CI.** A bare `next build` does not run the check, so CI calls `npm run build`.
- **Font network.** `next/font/google` needs network at build; `next/font/local` is the alternative.
- **Interpolations.** Put a space between two interpolations in `className`: `${a} ${b}`, never `${a}${b}`.
