<p align="center">
    <img src="assets/logo.svg">
</p>

# Tidgrid

A closed CSS vocabulary for coding agents.

One file per project, `tidgrid.toml`, defines every class a page may use.
`tidgrid build` turns it into one static stylesheet. `tidgrid check` fails on
anything outside it: an unknown class, a `<style>` block, a raw colour.

The split that keeps pages from looking generic:

- **Layout is atoms.** Rows, gaps, padding, widths: free to combine, always on
  the project's token scales.
- **Appearance is roles.** Colour, border, shadow and type are reachable only
  through named roles (`tg-surface(card)`, `tg-emphasis(key)`). What a role
  looks like is decided once, in the vocabulary file, and reviewed there.
  There is no border or colour atom.

The versions up to `0.1.0-alpha.2` were an SCSS flexbox library; that code
lives in the git history.

## Install

```
go install github.com/sneikki/tidgrid@latest
```

## Commands

| Command | Does |
|---|---|
| `tidgrid init` | Writes a starter `tidgrid.toml`. Refuses to overwrite. |
| `tidgrid build` | Validates the vocabulary and writes the stylesheet to `project.out`. |
| `tidgrid check` | Scans the sources. Prints `path:line: message` per violation, then the usage summary. |
| `tidgrid list` | Prints the vocabulary as compact text: what an agent reads before UI work. |
| `tidgrid sheet` | Writes one HTML page that renders every token and role: what a person reads to see the look. |

All take `-c <file>` (default `./tidgrid.toml`). `sheet` takes `-o <file>`
(default `tidgrid-sheet.html`). Exit code 0: fine. 1: violations. 2: the
vocabulary file itself is invalid or missing.

## The vocabulary file

```toml
[project]
prefix = "tg"                       # default
out    = "src/styles/tidgrid.css"   # required
scan   = ["src"]                    # default ["."]

[breakpoints]                       # ascending, in file order
sm = "40em"
md = "64em"

[tokens.color]                      # tokens are the ONLY place for raw values
ground = "#0f1214"
ink    = "#e6e9eb"
accent = "#5ed5dc"
link   = "{color.accent}"           # a token may reference another

[tokens.space]                      # gap and pad atoms run on this group
sm = "0.5rem"
md = "1rem"

[tokens.size]                       # width and height atoms run on this group
full    = "100%"
measure = "80rem"

[tokens.border]                     # any other group name is free
hair = "1px"

[tokens.motion]
slow = "1.1s"

[element.body]                      # base rule for a bare element
background = "{color.ground}"
color      = "{color.ink}"

[role.surface.card]                 # class: tg-surface(card)
use        = "A panel that lifts off the ground. One per group of related content."
background = "{color.ground}"
border     = "{border.hair} solid {color.accent}"

[role.surface.card."&:hover"]       # nested rule, native CSS nesting
border-color = "{color.ink}"

[role.surface.card."@md"]           # from breakpoint md up
border-width = "{border.hair}"

[role.state.ready]                  # a class a script sets; other roles react to it
use = "Set on body when the intro animation has finished."

[role.text.title]
use     = "The one title of a page."
opacity = "0"

[role.text.title.".tg-state(ready) &"]
opacity   = "1"
animation = "{keyframes.rise} {motion.slow} ease"

[keyframes.rise."0%"]
transform = "translateY({space.sm})"

[keyframes.rise."100%"]
transform = "none"

[allow.class]                       # visible exceptions, each with its reason
"has-js" = "set by an inline script before first paint"

[allow.css]
"src/styles/fonts.css" = "self-hosted @font-face only"
```

### Tokens

`[tokens.<group>]` holds `name = "raw value"`. Each becomes a custom property
`--tg-<group>-<name>` on `:root`, so scripts can read it. In any value,
`{group.name}` is replaced with `var(--tg-group-name)`. Breakpoints are
substituted literally (media queries cannot read custom properties).
`{keyframes.<name>}` is replaced with the generated animation name
`tg-<name>`.

### Layout atoms

Fixed grammar, the same in every project. `S` is every name in
`tokens.space`, `Z` every name in `tokens.size`.

| Class | CSS |
|---|---|
| `tg-row` | `display: flex; flex-direction: row` |
| `tg-col` | `display: flex; flex-direction: column` |
| `tg-wrap` | `flex-wrap: wrap` |
| `tg-grow` | `flex: 1 1 0; min-width: 0` |
| `tg-fixed` | `flex: none` |
| `tg-center` | `margin-inline: auto` |
| `tg-align(start\|center\|end\|baseline\|stretch)` | `align-items` |
| `tg-justify(start\|center\|end\|between\|around\|evenly)` | `justify-content` |
| `tg-self(start\|center\|end\|stretch)` | `align-self` |
| `tg-text-align(start\|center\|end)` | `text-align` |
| `tg-gap(S)` `tg-x-gap(S)` `tg-y-gap(S)` | `gap`, `column-gap`, `row-gap` |
| `tg-pad(S)` `tg-pad-x(S)` `tg-pad-y(S)` | `padding`, `padding-inline`, `padding-block` |
| `tg-pad-t(S)` `tg-pad-r(S)` `tg-pad-b(S)` `tg-pad-l(S)` | `padding-top` … `padding-left` |
| `tg-width(Z)` `tg-min-width(Z)` `tg-max-width(Z)` | `width`, `min-width`, `max-width` |
| `tg-height(Z)` `tg-min-height(Z)` | `height`, `min-height` |
| `tg-hide` | `display: none` |

Every atom also exists per breakpoint, mobile first: `md:tg-gap(md)` applies
from `md` up. Layout the atoms do not cover (a grid, an overlay) is a role.

The stylesheet is written in this order, so a later part wins over an
earlier one: tokens, elements, keyframes, roles, atoms, then the atoms of
each breakpoint in file order. `tg-hide` is the last atom of its block.

### Roles

`[role.<kind>.<name>]` becomes the class `tg-<kind>(<name>)`. Kinds are free
words (`surface`, `text`, `emphasis`, `action`, `state`, …) but may not
collide with an atom name.

- `use` is required: one or two sentences that tell an agent when the role
  applies. It is not emitted as CSS.
- Every other string key is a CSS declaration.
- A sub-table is a nested rule. Its key is one of:
  - `@<breakpoint>`: from that breakpoint up.
  - `@below-<breakpoint>`: below that breakpoint.
  - any other `@…`: the at-rule as written (`@media (prefers-reduced-motion: reduce)`, `@supports …`).
  - a selector that contains `&` (`&:hover`, `&::before`, `& h1`,
    `.tg-state(ready) &`). Vocabulary classes in it are written plainly and
    must exist.
- Sub-tables nest to any depth.

`[element.<selector>]` works the same for a bare element (`body`, `a`,
`"*"`); `use` is optional there.

`[keyframes.<name>]` holds one sub-table per stop.

Roles, declarations and nested rules are written in file order.

### What a value may contain

Outside `[tokens]` and `[breakpoints]`, a value may not contain a raw colour
(`#hex`, `rgb()`, `hsl()`, `oklch()` and the like), a raw length or a raw
time (a number with a unit such as `px`, `rem`, `vw`, `ms`). Those live in
tokens. Keywords, `0`, unitless numbers, percentages and functions over
tokens (`calc({space.md} * 2)`,
`color-mix(in srgb, {color.accent} 18%, transparent)`) are fine.

## What `check` reports

| Violation | Example |
|---|---|
| A class that is not in the vocabulary or under `allow.class` | `border-l-4`, `card-accent` |
| A class name built from pieces | `` `is-${state}` `` |
| Appearance in a `style` attribute (geometry is allowed) | `style="color: red"` |
| Own CSS: a `<style>` block, or a stylesheet not under `allow.css` | `Card.module.css` |
| The stylesheet at `project.out` differs from what `build` would write | vocabulary edited, not rebuilt |

A message names the fix where there is one:
`src/pages/index.astro:14: unknown class tg-gap(2); tg-gap takes: sm md`.

After the violations: per file, the roles it uses; the roles nothing uses;
the `allow` entries with their reasons.

### How classes are found

A scanner, not a parser, so it works the same in HTML, Astro, JSX, Vue and
Svelte. It reads:

- `class="…"` and `className="…"`;
- every string literal inside `class={…}`, `className={…}`, `class:list={…}`;
- every string literal passed to `cn`, `clsx`, `cva`, `classNames`,
  `twMerge` (configurable: `project.helpers`) and to `classList.add`,
  `remove`, `toggle`, `contains`.

Known limit: a class kept in a variable far from these places is not seen.
A name built from pieces cannot be checked at all, which is why it is a
violation: write the whole names out in a lookup table.

## Working with an agent

Three lines for the project's agent instructions:

```
Style only with the classes `tidgrid list` prints. No <style>, no style= for appearance, no own CSS.
Run `tidgrid check` before every commit; it must exit 0.
A new role is a change to tidgrid.toml in its own commit, never mixed into page work.
```

## License

MIT
