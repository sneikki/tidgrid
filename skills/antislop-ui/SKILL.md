---
name: antislop-ui
description: Use before building any new UI (a site, a page, an app screen, a set of components), when asked to make a UI look good, better or less generic, and when someone says a UI looks generic, templated or AI-made. Decides the look once, with the person, as a Tidgrid vocabulary; pages are then built only from it, and `tidgrid check` fails on anything outside it.
---

# Anti-slop UI

A generated UI looks generated because every choice in it is the most probable one: the agent writes what it has seen most often. A prompt that says "be distinctive" is advice the agent can weigh and drop. This skill does two things instead. It moves the look into one decision, made on purpose for this subject and approved by the person before any page exists. Then it puts a command behind that decision that fails when a page steps outside it.

The command is Tidgrid, a closed CSS vocabulary: tokens, layout atoms on the token scales, and named roles. Colour, border, shadow and type exist only in roles and base element rules, never as atoms, so a page cannot put together a look of its own.

## What a command enforces and what it does not

| Step | Kept by |
|---|---|
| Pages use only classes from the vocabulary: no `<style>`, no own stylesheet, no appearance in `style=`, no class name built from pieces | `tidgrid check` fails |
| The stylesheet matches the vocabulary; no role sits unused | `tidgrid check` fails |
| Raw colours, lengths and times only in tokens; every role has a `use` line; the same lines are not copied across roles | `tidgrid build` fails |
| The vocabulary fits this subject and does not read as a default | your judgement (step 4), then the person |
| The person said yes to the look before pages were written | you; no command knows |
| Composition: what goes where, how large, how dense | your judgement; layout atoms are free, so a centred hero over three equal cards passes every check |

Never say a passing check means the design is good. It means the pages keep to the design.

## Running the tool

`npx tidgrid <command>` needs only Node (npm package `tidgrid`, 0.3.0 or later). In a project with a `package.json`, add it with `npm install -D tidgrid` so the version is fixed. A `tidgrid` binary on the PATH (`go install github.com/sneikki/tidgrid/cmd/tidgrid@latest`) works the same. Below, `tidgrid` means either.

Commands read `index.tg` in the current folder unless given a file. Exit 0 is fine, 1 means check violations, 2 means the vocabulary does not load.

How to write `.tg` files, pick a class, add a mixin, a state or dark mode, and keep class names in JSX visible to the checker is in the usage skill: `skills/tidgrid/SKILL.md` in https://github.com/sneikki/tidgrid, and https://tidgrid.dev/docs. Read it before writing the vocabulary. This skill does not repeat it.

## 1. Find out where the project stands

- **`index.tg` exists.** Run `tidgrid list`, and `tidgrid sheet -o <file>` to see it (step 5 says where the file may go). Run step 4 on what you see. If it reads as a default, tell the person before building on it; do not redesign it unasked. Otherwise go to step 6.
- **Only a `tidgrid.toml`.** Run `tidgrid convert`, then as above.
- **No vocabulary.** Steps 2 to 6. `tidgrid init` writes a starter `index.tg`. Its values are placeholders (system fonts, a grey ramp, one blue accent): that is the default look itself. Keep the structure, replace every value, and replace the starter roles with your own.
- **Styled some other way** (own CSS, a utility framework, a component library). Moving it to Tidgrid is a migration: propose it and wait for a yes. The current look is material for step 2, not a decision already made.

## 2. Write the brief of the look

Before any value, write a few sentences. They go to the person with the sheet.

- What the subject is, who looks at it, doing what, how often, on what screen.
- The visual world the subject already has outside software: its documents, instruments, places, printed matter, trade. A look that belongs to the subject comes from there more often than from other interfaces.
- The one quality the look must carry, and what it must not be mistaken for.
- Every kind of content the pages will hold: long text, figures, tables, forms, code, media, lists, links inside and outside prose, empty, loading and error states.

Every later decision must trace back to a sentence here. One that does not is a default.

## 3. Decide the vocabulary

Tokens first, then roles. Decide each one; the first value that comes to mind is the probable one.

- **Colour.** Start with the ground the pages sit on and the ink on it, then add colours one job at a time: one colour per job, and only as many accents as there are things that must stand out. Decide light and dark together (`light-dark()` pairs) if both are wanted. Check that every ink reaches WCAG AA contrast on every ground it sits on; no command checks this.
- **Type.** Choose by what is read most here: long prose, figures, labels or code. The display face, the reading face and the scale are three decisions; the distance between scale steps sets the tone as much as the face does. A face that is not on every machine needs a load path: a remote `<link>` passes the check, a local `@font-face` file needs `allow_css(path, reason)`, and in Next `next/font` with `variable`, the font token pointing at that custom property.
- **Space.** One scale, with steps and density chosen for the use: a tool someone works in all day packs tighter than a page someone reads once.
- **Shape and motion.** Corners, rules, borders, shadows, transitions: each a decision, and "none" is one.
- **Roles.** One per kind of content in the brief, named by purpose, each `use` line saying when it applies and when it does not. Give the subject's own things roles of their own (its records, measures, statuses, objects), not only levels of text. A vocabulary of title, heading, body, label, card and button is the starter again under new colours.
- **The unglamorous roles.** Tables, form fields with focus, disabled and error states, plain links, long words and code that must wrap, empty states. In our trial, pages written from a finished vocabulary had no generic patterns but came out dry: the vocabulary held text levels and a code panel and little else, so every page became heading, prose, panel.

`tidgrid build` must exit 0. `tidgrid check` lists every role as unused until pages use them; at this step that is expected.

## 4. Self-check: does it read as a default?

Run it on the vocabulary before showing it, and on every page before calling it done. For each decision:

1. **Swap.** Would this value, role or arrangement fit an unrelated product unchanged? Then it was not decided for this one.
2. **Origin.** Which sentence of the brief did it come from? "Clean", "modern" and "professional" are not in the brief.
3. **Neighbour.** If it replaced a default, is it just the next most common choice, another fashionable typeface or another gradient? That is the same default one step over.

The looks generated UIs fall into, so you recognise them in your own draft. They are symptoms, not a list to avoid: steering around them one at a time only finds the next one.

- **Colour:** white or near-black ground, a grey ramp, one saturated blue, indigo or violet accent on every link, button and highlight; a blue-to-violet gradient behind the hero or inside the headline.
- **Type:** the system stack, or the one geometric sans every template ships, for everything; headings that are the body face made bigger and bolder; small upper-case tracked labels above every heading.
- **Surfaces:** every group in a rounded card with a soft shadow or a hairline border; a coloured bar on a box's left edge; frosted glass; pill badges on everything.
- **Composition:** a centred hero with headline, subline and two buttons; three or four equal cards with an icon each; a row of big numbers; every section in the same centred column with the same padding; emoji as icons.
- **Sameness:** every page assembled from the same three parts.

When a check fails, go back to the brief, not to the opposite of the symptom.

## 5. Show it once and get a yes

Run `tidgrid build`, then `tidgrid sheet -o <file>` with the file outside the `scan` folders: inside them, `check` reads the sheet's own markup and fails. The sheet renders every token and role. If you can render HTML, look at it yourself first. Give the person the file, the brief, and one line per decision on why it fits the subject. Ask one question: yes, or what to change. Revise and show again until the answer is yes.

No page before the yes. If no one can answer, stop and say what is waiting; never approve your own look.

## 6. Build pages from the vocabulary only

- Read `tidgrid list` and pick classes in the usage skill's order: atoms for layout, an existing role for a look, a change to a role, a new role last.
- A page needs a look the vocabulary lacks: add the role, with its `use` line, in the same commit as the page (`check` fails on a role nothing uses). A role inside the approved look needs no new yes. A change to tokens or to the look of existing roles does: show `tidgrid plan` and a new sheet first.
- Composition is the part no command sees. Shape each page by what it is for, not by the last page, and run step 4 on it.
- **Plain HTML:** `project(out = ..., scan = [...])` with `scan` covering the pages, and the stylesheet linked from each page.
- **React or Next:** `out` inside the app (for Next, a file in `app/` imported once in the root layout); `scan` lists the source folders only, never the build output (`.next` is full of `style` attributes); add `"tg": "tidgrid build && tidgrid check"` to `package.json` and run it before `dev` and `build`; class names written whole, and tables of them through `cn()`.
- Never write own CSS, a `<style>` block or appearance in `style=`, never edit the generated stylesheet, and never use `allow_class` or ignore comments to get past the check.

Done means: `tidgrid build && tidgrid check` exit 0, the generated stylesheet committed, step 4 run on every page, and, if you can render, every page looked at at phone and desktop width.
