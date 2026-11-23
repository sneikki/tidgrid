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
go install github.com/sneikki/tidgrid/cmd/tidgrid@latest
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
`tg-<name>`. Group and token names are made of letters, digits, `-` and `_`.

`project.scan` and `project.helpers` are lists. `check` reports a scan path
that does not exist.

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
| `tg-cols(1…6\|12)` | `display: grid; grid-template-columns: repeat(N, minmax(0, 1fr))` |
| `tg-span(1…12)` | `grid-column: span N` |
| `tg-order(first\|last)` | `order: -1`, `order: 1` |
| `tg-hide` | `display: none` |
| `tg-hide-below(B)` | `display: none` under breakpoint B (B = a breakpoint name) |

Every atom except `tg-hide-below` also exists per breakpoint, mobile first:
`md:tg-gap(md)` applies from `md` up; `tg-hide-below(md)` hides below `md`,
and nothing has to be restored above it. Layout the atoms do not
cover (an overlay, a fixed panel) is a role.

An element takes a property from an atom or from a role, never both: atoms
are written after roles and win.

The stylesheet is written in this order, so a later part wins over an
earlier one: tokens, elements, keyframes, roles, atoms, then the atoms of
each breakpoint in file order, then the `tg-hide-below` atoms, so they win
over a `md:tg-row` on the same element. `tg-hide` is the last atom of its
block.

### Roles

`[role.<kind>.<name>]` becomes the class `tg-<kind>(<name>)`. Kinds are free
words (`surface`, `text`, `emphasis`, `action`, `state`, …) but may not
collide with an atom name.

- `use` is required: one or two sentences that tell an agent when the role
  applies. It is not emitted as CSS, and it is an error inside a sub-table.
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
(`#hex`, `rgb()`, `hsl()`, `oklch()` and the like, or a CSS colour name such
as `red`), a raw length or a raw time (a number with any CSS length or time
unit: `px`, `rem`, `vw`, `dvh`, `cqw`, `ms` …). Those live in tokens.
So does a percentage that sizes a box: in `width`, `height`, `inline-size`,
`block-size`, their `min-` and `max-` forms and `flex-basis`, write
`{size.half}`, not `50%`.
Keywords, `transparent`, `currentColor`, `0`, unitless numbers, percentages
anywhere else (`border-radius: 50%`, a gradient stop, `translate(-50%)`),
`fr`, angles and functions over tokens (`calc({space.md} * 2)`,
`color-mix(in srgb, {color.accent} 18%, transparent)`) are fine. Text in
quotes and inside `url()` is not looked at.

The key of an at-rule sub-table follows the same rule. A reference in it is
replaced with the token's written value, because a media query cannot read a
custom property: `"@media (max-height: {size.short})"`.

## What `check` reports

| Violation | Example |
|---|---|
| A class that is not in the vocabulary or under `allow.class` | `border-l-4`, `card-accent` |
| A class name built from pieces | `` `is-${state}` `` |
| Appearance in a `style` attribute (geometry is allowed) | `style="color: red"` |
| Own CSS: a `<style>` block, or a stylesheet not under `allow.css` | `Card.module.css` |
| The stylesheet at `project.out` differs from what `build` would write | vocabulary edited, not rebuilt |
| A `tidgrid-ignore-start` that is never closed | the rest of the file would go unchecked |

A message names the fix where there is one:
`src/pages/index.astro:14: unknown class tg-gap(2); tg-gap takes: sm md`.

After the violations: per file, the roles it uses; the roles nothing uses;
the `allow` entries with their reasons.

### How classes are found

A scanner, not a parser, so it works the same in HTML, Astro, JSX, Vue and
Svelte. It reads:

- `class="…"` and `className="…"` attributes, in any letter case;
- every string literal inside `class={…}`, `className={…}`, `class:list={…}`,
  `:class="…"`;
- every string literal on the right of a `className = …` assignment in
  script;
- every string literal passed to `cn`, `clsx`, `cva`, `classNames`,
  `twMerge` and to `classList.add`, `remove`, `toggle`, `contains`.
  `project.helpers` replaces that list of five, so repeat the ones still in
  use;
- the name in a `class:name` directive (Svelte); `style:prop` counts as a
  style attribute.

A string literal that is only compared (`kind === "primary"`) is not a class
and is not read.

Inside a class value, a template tag `{% … %}` is skipped and the names
around it are read. An interpolation (`{{ … }}`, `${…}`, `{…}`) glued to a
name is a built name; standing alone between spaces it is skipped.

Not read: fenced code blocks in Markdown, inline code in `.md` files (not
`.mdx`, where a backtick can build a class name), `{% highlight %}` and
`{% raw %}` blocks, and anything between the comments `tidgrid-ignore-start`
and `tidgrid-ignore-end` (for code samples that show classes).

Files: `.html .htm .astro .jsx .tsx .js .mjs .cjs .ts .vue .svelte .md .mdx
.liquid .njk` under `project.scan`. Stylesheets: `.css .scss .sass .less .styl`. Never entered:
`node_modules .git dist build _site .astro vendor`. A symlinked directory is
not entered either; `check` reports it unless its target is under
`project.scan`.

Known limit: a class kept in a variable far from these places is not seen.
A name built from pieces cannot be checked at all, which is why it is a
violation: write the whole names out in a lookup table. The table is itself
a variable, so pass each name through a helper to have it checked:
`{ primary: cn("tg-action(primary)") }`.

## Working with an agent

Three lines for the project's agent instructions:

```
Style only with the classes `tidgrid list` prints. No <style>, no style= for appearance, no own CSS.
Run `tidgrid check` before every commit; it must exit 0.
A new role is a change to tidgrid.toml in its own commit, never mixed into page work.
```

For Claude Code there is a skill in `skills/tidgrid`: how to pick a class, when
a new role is right, how names in code stay visible to the checker, dark
mode. It loads by itself in a project that has a `tidgrid.toml`. Install it
for yourself:

```
git clone https://github.com/sneikki/tidgrid && ln -s "$PWD/tidgrid/skills/tidgrid" ~/.claude/skills/tidgrid
```

## License

MIT
