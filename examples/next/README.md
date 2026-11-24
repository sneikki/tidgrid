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
  Never edit it.
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
- **Client state.** A class a script sets is a `state` role (`tg-state(active)`)
  that other roles react to. `components/DepartureBoard.tsx` is a client
  component; it writes the name whole inside `cn()`. Server and client files
  are scanned the same.
- **Fonts.** `next/font` writes the `@font-face` and exposes each family as a
  custom property (`variable: "--font-display"`); `tokens.tg` points the
  `font` tokens at those names. Nothing in the stylesheet or the checker knows
  the font's own hashed class, so use `variable`, not `className`.
