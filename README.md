<p align="center">
    <img src="assets/logo.svg">
</p>

# Tidgrid

A closed CSS vocabulary for coding agents.

A project's vocabulary, written in Starlark (`.tg` files, entry `index.tg`),
defines every class a page may use. `tidgrid build` turns it into one static
stylesheet. `tidgrid check` fails on anything outside it: an unknown class, a
`<style>` block, a raw colour, a role no page uses.

The point is a closed vocabulary, decided once and checked.

The split that keeps pages from looking generic:

- **Layout is atoms.** Rows, gaps, padding, widths: free to combine, always on
  the project's token scales.
- **Appearance is roles.** Colour, border, shadow and type are reachable only
  through named roles (`tg-surface(card)`, `tg-emphasis(key)`). What a role
  looks like is decided once, in the vocabulary, and reviewed there. There is
  no border or colour atom.

The versions up to `0.1.0-alpha.2` were an SCSS flexbox library; that code
lives in the git history.

## Install

Needs Go 1.25.

```
go install github.com/sneikki/tidgrid/cmd/tidgrid@latest
```

Without Go, Node 18 or newer runs the same tool as WebAssembly:
`npm install --save-dev tidgrid`, or `npx tidgrid check` once.

## Commands

| Command | Does |
|---|---|
| `tidgrid init [file]` | Writes a starter vocabulary (default `index.tg`). Refuses to overwrite. |
| `tidgrid build [file]` | Validates the vocabulary and writes the stylesheet to the `out` of `project()`. |
| `tidgrid check [file]` | Scans the sources. Prints `path:line: message` per violation, then the usage summary. |
| `tidgrid list [file]` | Prints the vocabulary as compact text: what an agent reads before UI work. |
| `tidgrid plan [file]` | Prints what `build` would change in the stylesheet on disk. Writes nothing. |
| `tidgrid sheet [file]` | Writes one HTML page that renders every token and role: what a person reads to see the look. |
| `tidgrid convert [tidgrid.toml]` | Writes a `.tg` vocabulary from an old TOML one. See [Converting an old project](#converting-an-old-project). |

The commands that read a vocabulary take the entry file as an argument, like a
compiler. With no argument they read `index.tg` in the current directory.
There is no `-c` flag. `sheet` also takes `-o <file>` (default
`tidgrid-sheet.html`).

Exit code 0: fine. 1: violations (`check`). 2: the vocabulary itself is
invalid or missing (`check` too), or the command was misused.

With no `index.tg` but a `tidgrid.toml` in the directory, `build`, `check`,
`list`, `plan` and `sheet` exit 2 and say to run `tidgrid convert`. `init` and
`convert` do not: `init` writes `index.tg` as usual. A file named with the `.toml` extension is
still read as TOML (see below); a `.tg` file is read as Starlark.

### plan

`tidgrid plan` renders the vocabulary in memory and compares it with the
stylesheet at the output path: classes added and removed, the declarations that
changed in a class, changed tokens with the classes that reach them, element
rules and keyframes. Pages that use a removed class are not listed; `check`
reports those as unknown classes once the stylesheet is built. Exit 0 whether
or not there is a difference, 2 when the vocabulary does not load.

```
$ tidgrid plan
changed class tg-text(caption): font-weight: 600 -> 500
changed class tg-text(hint): font-weight: 600 -> 500
added class tg-text(lead)
changed token color.ink: #e6e9eb -> #eeeeee
  reaches 6 classes: tg-text(section-title) tg-text(caption) tg-text(lead) tg-surface(card) tg-notice(info) tg-notice(warn)
```

With no stylesheet yet, everything is listed as added. With no difference it
prints `no difference`.

## The vocabulary

A vocabulary is a Starlark program. The tool runs the entry file, and every
builtin call registers one thing (a token group, a role, a mixin) in a single
model. Rules, `list`, `check`, `plan` and the stylesheet all run on that model,
so splitting the program over files changes nothing in them. Within each part
of the stylesheet, things are written in the order the calls ran.

Starlark here: top-level `for` and `if` work; `def`, list comprehensions and
string formatting (`"%s %s" % (a, b)`) work; `while` and recursion do not;
there is a step limit for the whole run; nothing reads or writes outside the
vocabulary's folder.

A small project, in three files. `index.tg`:

```python
project(out = "src/tidgrid.css", scan = ["src"])
breakpoints(sm = "40em", md = "64em")
load("tokens.tg", "color", "space")
source("roles.tg")
element("body", css = {"background": color.ground, "color": color.ink})
allow_class("has-js", "set by an inline script before first paint")
```

`tokens.tg`:

```python
color = tokens("color", {"ground": "#0f1214", "ink": "#e6e9eb", "sand-ink": "#d8cfc0", "accent": "#5ed5dc"})
space = tokens("space", {"sm": "0.5rem", "md": "1rem"})
font = tokens("font", {"title": "Georgia, serif", "sans": "system-ui, sans-serif"})
sizes = tokens("sizes", {"h2": "1.5rem", "label": "0.875rem"})
motion = tokens("motion", {"slow": "1.1s"})
border = tokens("border", {"hair": "1px"})

def tint(c, pct):
    return "color-mix(in srgb, %s %d%%, transparent)" % (c, pct)
```

`roles.tg`:

```python
load("tokens.tg", "border", "color", "space", "font", "sizes", "motion", "tint")

type = mixin("type",
    use = "A text style. Pick one; never set family, size or weight by hand.",
    sets = {
        "heading": {"font-family": font.title, "font-size": sizes.h2, "font-weight": 700},
        "label": {"font-family": font.sans, "font-size": sizes.label, "font-weight": 600},
    },
)

shape = mixin("shape", use = "Shared shapes of controls.", sets = {
    "button": {"padding": "%s %s" % (space.sm, space.md), "cursor": "pointer"},
})

def filled(fill, ink):
    return {"background": fill, "color": ink, "&:hover": {"background": tint(fill, 85)}}

role("action", "primary", use = "The one main button of a view.", include = [shape.button], css = filled(color.accent, color.ground))
role("action", "secondary", use = "A second button beside the primary one.", include = [shape.button], css = filled(color["sand-ink"], color.ground))

role("text", "section-title", use = "Heading of a section inside a page.", include = [type.heading], css = {"color": color.ink})
role("text", "card-title", use = "Title of a card.", include = [type.heading], css = {"color": color.accent})
role("text", "caption", use = "Small text under a figure.", include = [type.label], css = {"color": color.ink})
role("text", "hint", use = "A hint beside a field.", include = [type.label], css = {"color": color.accent})

rise = keyframes("rise", {
    "0%": {"transform": "translateY(%s)" % space.sm},
    "100%": {"transform": "none"},
})

role("state", "ready", use = "Set on body when the intro animation has finished.")

role("surface", "card",
    use = "A panel that lifts off the ground. One per group of related content.",
    css = {
        "background": color.ground,
        "border": "%s solid %s" % (border.hair, color.accent),
        "&:hover": {"border-color": color.ink},
        "@md": {"padding": space.md},
        "@media (prefers-reduced-motion: reduce)": {"animation": "none"},
        ".tg-state(ready) &": {"animation": "%s %s ease" % (rise, motion.slow)},
    })

for n in ["info", "warn"]:
    role("notice", n, use = "A notice that is %s." % n, css = {"color": color.ink})
```

`tidgrid list` prints the roles and mixins of that vocabulary like this:

```
# roles
## action
tg-action(primary)  The one main button of a view.
tg-action(secondary)  A second button beside the primary one.
## text
tg-text(section-title)  Heading of a section inside a page.
tg-text(card-title)  Title of a card.
tg-text(caption)  Small text under a figure.
tg-text(hint)  A hint beside a field.
## set by scripts
tg-state(ready)  Set on body when the intro animation has finished.
## surface
tg-surface(card)  A panel that lifts off the ground. One per group of related content.
## notice
tg-notice(info)  A notice that is info.
tg-notice(warn)  A notice that is warn.

# mixins
type  A text style. Pick one; never set family, size or weight by hand.
  heading: font-family: {font.title}; font-size: {sizes.h2}; font-weight: 700
    included by: tg-text(section-title) tg-text(card-title)
  label: font-family: {font.sans}; font-size: {sizes.label}; font-weight: 600
    included by: tg-text(caption) tg-text(hint)
shape  Shared shapes of controls.
  button: padding: {space.sm} {space.md}; cursor: pointer
    included by: tg-action(primary) tg-action(secondary)
```

### Files: `load` and `source`

- `load("tokens.tg", "color", "space")` is Starlark's own `load`: it runs the
  file and imports the named variables. It must name at least one.
- `source("roles.tg")` runs a file for what it declares and imports no names
  into the caller. The sourced file may itself `load` what it needs.
- A path is relative to the file that names it and must stay inside the entry
  file's directory: no absolute path, no `..`, no symlink out. A file runs
  once, however many times `load` and `source` name it. A cycle is an error.
- The builtins below are available in every file without importing.
- The entry file can be anywhere and named anything; `out`, `scan` and
  `allow_css` paths resolve against its directory. One folder may hold several
  vocabularies.

### Builtins

| Call | Registers |
|---|---|
| `project(out, scan=["."], prefix="tg", helpers=["cn","clsx","cva","classNames","twMerge"], repeats=[(3, 2), (2, 4)])` | Where the stylesheet goes (`out`, required), which folders `check` scans, the class prefix, the helper functions `check` reads class names from, the repeated-lines setting. Once. |
| `breakpoints({...})` or `breakpoints(sm="40em", ...)` | Breakpoints, ascending in call order. Values are literal (media queries cannot read custom properties). |
| `tokens(group, {name: value})` or `tokens(group, name=value, ...)` | A token group. Returns the group: `color.ink` and `color["sand-ink"]` give the reference `"{color.ink}"`. |
| `element(selector, css={...}, use="")` | A base rule for a bare element (`body`, `a`, `"*"`). `use` is optional. |
| `tidiness(level)` | A reset, written before the element rules. `"minimal"`: browser fixes only. `"full"`: also no margin, padding, border, heading size, list marker or link colour until an element rule or a role gives one. No call, no reset. |
| `keyframes(name, {stop: css})` | An animation. Returns `"{keyframes.name}"`. |
| `role(kind, name, use=..., css={...}, include=[...])` | The class `tg-<kind>(<name>)`. `use` is required. |
| `mixin(name, use=..., sets={...})` | Named sets of declarations that roles include. Returns the mixin: `type.heading` and `type["heading"]` give a set. |
| `allow_class(name, reason)` | A class `check` accepts although it is not in the vocabulary. |
| `allow_css(path, reason)` | A stylesheet `check` accepts. |
| `load(file, names...)`, `source(file)` | See above. |

Declaring a role, a token group, a mixin or `project` twice is an error that
names both places.

### The `css` dict

`css` is a dict of CSS declarations. Keys are strings.

- A string or number value is a declaration: `"font-weight": 700`,
  `"padding": "{space.md}"`.
- A dict value is a nested rule. Its key is one of:
  - `"@md"`: from breakpoint `md` up.
  - `"@below-md"`: below breakpoint `md`.
  - any other `@…`: the at-rule as written
    (`"@media (prefers-reduced-motion: reduce)"`, `"@supports …"`).
  - a selector that contains `&`: `"&:hover"`, `"&::before"`, `"& h1"`,
    `".tg-state(ready) &"`. Vocabulary classes in it are written plainly and
    must exist.
- Nested rules nest to any depth.
- `use` is not a key of `css`; it is an argument.

Declarations and nested rules are written in the order the dict has them.

### Tokens

`tokens("color", {...})` makes a group. Each token becomes a custom property
`--tg-<group>-<name>` on `:root`, so scripts can read it. Group and token names
are made of letters, digits, `-` and `_`; `keyframes` and `breakpoints` are not
group names.

A token group is also a Starlark value. `color.ink` and `color["sand-ink"]`
are the string `"{color.ink}"`, which a value may contain anywhere; the build
replaces `{group.name}` with `var(--tg-group-name)`. A name the group does not
have is an error at that line (`token group color has no token "nope"`). In a
plain string `"{color.nope}"` the same mistake is `unknown reference`.
A token's own value may reference another (`"{color.accent}"`).
Breakpoints are substituted literally. `{keyframes.<name>}` is the generated
animation name `tg-<name>`.

### Values: functions and loops

A value is a string, so ordinary Starlark builds it. A `def` that returns a
value names an expression; a loop makes several roles:

```python
def tint(c, pct):
    return "color-mix(in srgb, %s %d%%, transparent)" % (c, pct)

for n in ["info", "warn"]:
    role("notice", n, use = "A notice that is %s." % n, css = {"color": color.ink})
```

Role names may be computed. `check` fails on a role no page uses, and every
role still needs its own written `use`. The finished values are validated
exactly as written by hand (see
[What a value may contain](#what-a-value-may-contain)).

### Mixins

A mixin is a closed list of named sets, each a `css` dict. A role includes
sets by name:

```python
role("text", "section-title", use = "...", include = [type.heading], css = {"color": color.ink})
```

- The role renders as if the set's lines were written first, then its own.
  A mixin puts nothing else into the stylesheet and is never a class.
- A set takes no parameters. A set is a step of a scale: `type.heading`,
  `type.label`.
- `include` takes only mixin sets. Anything else is an error.
- Nested rules with the same key, in a set and in the role (`"&:hover"`), merge
  into one rule: the set's lines first, then the role's.
- The same property given twice is an error, at any depth: a role's own
  property that an included set already gives, or two included sets that give
  the same one, also inside a nested rule. The message names the path:
  `index.tg:7: role.action.b."&:hover".color: given twice: by set btn.base and by the role's own css`
- **A set needs at least two roles that include it.** A set with fewer fails
  `build`. Only roles that exist count, and an unused role fails `check`, so
  a role invented to give a set its second user does not survive.
- **A set needs at least two lines in total**, nested ones counted. A set of
  one line fails `build`: `mixin.m.one: a set needs at least two lines; a single shared value belongs in a token`.
- A role whose own top-level lines contain every top-level line of a set (same
  property, same value as written) without including it fails `build`, and the
  message gives the `include` to write. This only applies to sets with two or
  more top-level lines. A role may still write by hand what a mixin does not
  cover.
- `mixin` needs `use`. A mixin name declared twice is an error.
- Mixins exist only in `.tg` vocabularies.

### Variants of a role

A role that comes in variants (buttons that differ in colour) needs no
feature of its own. The shared, fixed lines are a mixin set. The lines that
differ come from a plain function that may put the colour in as many places
as it likes:

```python
shape = mixin("shape", use = "Shared shapes of controls.", sets = {
    "button": {"padding": "%s %s" % (space.sm, space.md), "cursor": "pointer"},
})

def filled(fill, ink):
    return {"background": fill, "color": ink, "&:hover": {"background": tint(fill, 85)}}

role("action", "primary", use = "The one main button of a view.", include = [shape.button], css = filled(color.accent, color.ground))
role("action", "secondary", use = "A second button beside the primary one.", include = [shape.button], css = filled(color["sand-ink"], color.ground))
```

Class names stay as they are (`tg-action(primary)`). No private custom
property carries the colour, so `check` sees every value.

### Repeated lines

A line is a top-level (property, value as written) pair of a role's own `css`.
Lines that come from an included set do not count. `build` fails (exit 2) when
the same lines appear in several roles: 3 lines in 2 roles, or 2 lines in 4
roles. Each group is reported once, with the roles and the lines they share,
and the message says the way out: make the lines a mixin set and include it.
Two roles that share lines and every other role that holds all of those lines
form a group. A report of two roles sharing three lines:

```
index.tg:3: role.action.a: own lines also in role action.b: cursor: pointer; padding: 0; border: none; make them a mixin set and include it
index.tg: 1 group of repeated lines; fix the first and build again: one mixin set often clears several
```

The first line is at the first role of the group and lists the shared lines;
the last line is the total. At most 20 groups are shown, the ones with the
most roles first (then the most lines); the last line then also says how many
are shown.

The numbers are a setting, a list of (lines, roles) pairs:
`project(repeats=[(3, 2), (2, 4)])` is the default. `repeats = []` turns the
rule off. `tidgrid convert` writes `repeats = []`, because an old project
repeats lines freely; cleaning it up starts by removing that argument.

Two reported groups may share a line. Put the shared line in one set only. If
a role includes two sets that both give it, `build` says so (`given twice: by
set m.s1 and by set m.s2`).

## Layout atoms

Fixed grammar, the same in every project. `S` is every name in the `space`
token group, `Z` every name in `size`.

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
each breakpoint in breakpoint order, then the `tg-hide-below` atoms, so they win
over a `md:tg-row` on the same element. `tg-hide` is the last atom of its
block.

## Roles

`role(kind, name, ...)` becomes the class `tg-<kind>(<name>)`. Kinds are free
words (`surface`, `text`, `emphasis`, `action`, `state`, …) but may not
collide with an atom name (`kind collides with atom tg-row`). Kind and name
are made of letters, digits, `-` and `_`.

- `use` is required: one or two sentences that tell an agent when the role
  applies. It is not emitted as CSS.
- `css` holds the declarations and nested rules (see above).
- `include` lists mixin sets.

A class a script sets is a role of kind `state` with only a `use`; other roles
react to it in a nested rule (`".tg-state(ready) &"`, `"&.tg-state(open)"`).
ARIA and data attributes need no role: `"&[aria-current]"`.

A name written only in another role's or element's selector does not count as
used. `check` fails on the state role until the script that sets the class
writes the name whole, as a string the scanner reads (see
[How classes are found](#how-classes-are-found)):

```js
card.classList.toggle("tg-state(open)");
```

Dark mode: keep colours in tokens as pairs, `"light-dark(#32495d, #d7e2ec)"`,
and set `color-scheme` on `html`:
`element("html", css = {"color-scheme": "light dark"})`. The page follows the
system with no second stylesheet. For a manual switch, declare a state role
for the class, narrow the scheme in a nested rule, and have the script write
the name whole:

```python
role("state", "dark", use = "Set on html to force the dark scheme.")
element("html", css = {"color-scheme": "light dark", "&.tg-state(dark)": {"color-scheme": "dark"}})
```

```js
document.documentElement.classList.add("tg-state(dark)");
```

## What a value may contain

Outside `tokens` and `breakpoints`, a value may not contain a raw colour
(`#hex`, `rgb()`, `hsl()`, `oklch()` and the like, or a CSS colour name such
as `red`), a raw length or a raw time (a number with any CSS length or time
unit: `px`, `rem`, `vw`, `dvh`, `cqw`, `ms` …). Those live in tokens.
So does a percentage that sizes a box, a spacing or a grid track: in `width`,
`height`, `inline-size`, `block-size`, their `min-` and `max-` forms,
`flex-basis`, `flex`, `margin…`, `padding…`, the three `gap` properties,
`grid-template…` and `grid-auto-…`, write `{size.half}`, not `50%`. Only
those properties are checked for it.
Keywords, `transparent`, `currentColor`, `0`, unitless numbers, percentages
anywhere else (`border-radius: 50%`, `left: 50%`, a gradient stop,
`translate(-50%)`),
`fr`, angles and functions over tokens (`calc({space.md} * 2)`,
`color-mix(in srgb, {color.accent} 18%, transparent)`) are fine. Text in
quotes and inside `url()` is not looked at.

The key of an at-rule nested rule follows the same rule. A reference in it is
replaced with the token's written value, because a media query cannot read a
custom property: `"@media (max-height: {size.short})"`.

## What `build` reports

`build`, `list`, `plan` and `sheet` load the vocabulary and print every
problem, one per line, then exit 2. A line names the file and line of the call
that declared the thing, then the path to it. Starlark's own errors keep their
file, line and column. These are the messages, from real runs:

| Problem | Message |
|---|---|
| A raw value | `roles.tg:46: role.text.pad.padding: raw length or time 12px (move it to tokens)` |
| A raw size | `roles.tg:46: role.text.wide.width: raw size 40% (move it to tokens)` |
| An unknown reference | `roles.tg:46: role.text.unk.color: unknown reference {color.nope}` |
| An unknown token through the variable | `roles.tg:46:47: token group color has no token "nope"` |
| A role without `use` | `roles.tg:46: role.text.nouse.use: required: when the role applies` |
| A mixin without `use` | `roles.tg:46: mixin.nouse.use: required: when the mixin applies` |
| A role declared twice | `roles.tg:46:5: role.text.caption: declared twice (roles.tg:23 and roles.tg:46)` |
| `project` declared twice | `index.tg:7:8: project: declared twice (index.tg:1 and index.tg:7)` |
| A kind that is an atom | `role.row: kind collides with atom tg-row` |
| An unknown breakpoint | `roles.tg:46: role.text.nn."@xl": unknown breakpoint "xl"` |
| A selector without `&` | `roles.tg:46: role.text.ss.hover: selector must contain &` |
| An unknown class in a selector | `roles.tg:46: role.text.cc."&.tg-nope": unknown class tg-nope` |
| `include` of a non-set | `roles.tg:46: role.text.ns.include: 1 is not a mixin set; write a set such as type.heading` |
| A property given by a set and by the role | `roles.tg:46: role.text.dup.font-weight: given twice: by set type.label and by the role's own css` |
| A property given by two sets | `roles.tg:46: role.text.y.font-weight: given twice: by set type.label and by set type.heading` |
| A set written out by hand | `roles.tg:46: role.text.plain: own lines repeat set type.heading; write include = [type.heading] and drop them` |
| A set with one user | `roles.tg:46: mixin.solo.a: included only by role text.s1; a set needs at least two` |
| A set with none | `roles.tg:46: mixin.solo.b: included by no role; a set needs at least two` |
| A set of one line | `roles.tg:46: mixin.m.one: a set needs at least two lines; a single shared value belongs in a token` |
| Repeated lines | `index.tg:3: role.action.a: own lines also in role action.b: cursor: pointer; padding: 0; border: none; make them a mixin set and include it` |
| A wrong argument | `roles.tg:46:5: role: unexpected keyword argument "usee" (did you mean use?)` |
| A path leaving the folder | `index.tg:7:1: cannot load ../x.tg: "../x.tg": path leaves the directory of the entry file` |
| An absolute path | `index.tg:7:7: "/etc/x.tg": path must be relative to the loading file` |
| A missing file | `index.tg:7:7: nope.tg: cannot read inside the entry file's directory: openat nope.tg: no such file or directory` |
| A load cycle | `roles.tg:46:7: load cycle: index.tg -> roles.tg -> index.tg` |
| Recursion | `roles.tg:47:12: function f called recursively`, then `called from` lines |

A repeated-lines failure (see above) is also a `build` error.

## What `check` reports

| Violation | Example |
|---|---|
| A class that is not in the vocabulary or under `allow_class` | `border-l-4`, `card-accent` |
| A class name built from pieces | `` `is-${state}` `` |
| Appearance in a `style` attribute (geometry is allowed) | `style="color: red"` |
| Own CSS: a `<style>` block, or a stylesheet not under `allow_css` | `Card.module.css` |
| A role no scanned source uses | a role added and never used |
| The stylesheet at `out` is missing, or differs from what `build` would write | vocabulary edited, not rebuilt |
| A `tidgrid-ignore-start` that is never closed | the rest of the file would go unchecked |
| A symlinked directory that is not under `scan` | |

Messages name the fix where there is one. From real runs:

```
src/b.html:1: appearance in style attribute: color; use a role (`tidgrid list`)
src/b.html:1: own CSS: <style> block; make it a role in the vocabulary
src/b.html:1: unknown class card-accent: not in the vocabulary; run `tidgrid list`
src/b.html:1: unknown class tg-gap(2); tg-gap takes: sm md
src/b.html:1: unknown class tg-text(nope); tg-text takes: section-title card-title caption hint lead
src/b.html:2: class name built from pieces: is-${state}; write every whole name out in a lookup table
roles.tg:46: role tg-text(lead) is not used in the scanned sources: use it, remove it, or, if a script or lookup table sets it, write the name whole inside a helper call such as cn()
src/tidgrid.css:1: stylesheet differs from the vocabulary; run `tidgrid build`
src/b.html:1: tidgrid-ignore-start is never closed; the rest of the file is not checked
src/styles/tidgrid.css:1: stylesheet missing; run `tidgrid build`
```

An unused role is reported at the line of the call that declared it. A role
that a script sets counts as used when the script names it whole, as a string
(`classList.add("tg-state(ready)")`, `cn("tg-state(ready)")`). A name that
appears only in another role's or element's selector (`"&.tg-state(open)"`)
does not count. A role made in a loop is reported once per role.

`check` loads the vocabulary first: when that fails it prints the problems
(see above) and exits 2, with no scan.

After the violations: per file, the roles it uses; the `allow` entries with
their reasons.

### Adding a role

A new role goes in the same commit as the first page that uses it. `check`
fails on a role nothing uses, so a role in a commit of its own fails until the
page arrives. A change to an existing role or to tokens can still be a commit
of its own. `tidgrid init` writes the starter roles commented out, so a fresh
project has no unused role and `check` passes; uncomment a role when a page
uses it.

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
  `helpers` in `project()` replaces that list of five, so repeat the ones
  still in use;
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
.liquid .njk` under `scan`. Stylesheets: `.css .scss .sass .less .styl`. Never entered:
`node_modules .git dist build _site .astro vendor`. A symlinked directory is
not entered either; `check` reports it unless its target is under
`scan`.

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
A new role goes in the same commit as the page that first uses it; `tidgrid check` fails on a role nothing uses.
```

For Claude Code there is a skill in `skills/tidgrid`: how to pick a class, when
a new role or a mixin is right, how names in code stay visible to the checker,
dark mode. It loads by itself in a project that has an `index.tg` (or a
`tidgrid.toml` to convert). Install it for yourself:

```
git clone https://github.com/sneikki/tidgrid && ln -s "$PWD/tidgrid/skills/tidgrid" ~/.claude/skills/tidgrid
```

Editors and GitHub do not know `.tg` as Starlark. A `.gitattributes` line
(`*.tg linguist-language=Starlark`) and an editor file association fix that.

## Converting an old project

The vocabulary used to be a `tidgrid.toml`. That form is still read when the
file is named (`tidgrid check tidgrid.toml`), with its own rules unchanged,
except that `check` now also fails on a role nothing uses. It has no mixins,
functions or loops.

`tidgrid convert [tidgrid.toml]` writes `index.tg`, `tokens.tg` and `roles.tg`
beside it and refuses when any of the three exists. The converted project
builds the same stylesheet byte for byte. Values stay strings as written, and
TOML comments are not carried over.

```
$ tidgrid convert
wrote index.tg
wrote roles.tg
wrote tokens.tg
note: TOML comments are not carried over
```

In a directory that has `tidgrid.toml` and no `index.tg`, `build`, `check`,
`list`, `plan` and `sheet` exit 2 (`init` and `convert` do not):

```
index.tg not found, but tidgrid.toml is here: run `tidgrid convert` to write index.tg
```

## License

MIT
