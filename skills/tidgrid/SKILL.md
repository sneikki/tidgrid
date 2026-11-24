---
name: tidgrid
description: Use when building or changing UI in a project that has an index.tg or other .tg vocabulary files, or an old tidgrid.toml (convert it with `tidgrid convert`), with classes such as tg-row, tg-gap(md), tg-surface(card); when `tidgrid check` fails; or when adding a colour, a spacing value, a role, a mixin, a state or dark mode to such a project. Also when setting tidgrid up in a new project.
---

# Styling with Tidgrid

Tidgrid is a closed CSS vocabulary. A Starlark program (`.tg` files, entry `index.tg`) defines every class a page may use. `tidgrid build` writes the stylesheet from it and `tidgrid check` fails on anything outside it. You never write CSS; you pick classes, and when none fits you change the vocabulary.

Commands take the entry file as an argument and default to `index.tg` in the current directory: `tidgrid build [file]`. There is no `-c` flag. Install (needs Go 1.25): `go install github.com/sneikki/tidgrid/cmd/tidgrid@latest`.

A project that has only a `tidgrid.toml` and no `index.tg`: `build`, `check`, `list`, `plan` and `sheet` exit 2 and say to run `tidgrid convert` (`init` and `convert` work). Run it once; it writes `index.tg`, `tokens.tg` and `roles.tg` and builds the same stylesheet. Commit the three files and delete the TOML. The converted `project()` has `repeats = []`; remove it and fix what `build` then reports (see Mixins).

## Before any UI work

Run `tidgrid list` and read it. It prints the tokens, the atoms with the values they take, every role with the sentence that says when to use it, the mixins, and the exceptions. Do not guess a class from memory or from another project: role names differ per project.

## Picking a class

Ask in this order and stop at the first yes.

1. **Is it layout?** Direction, gap, padding, width, height, grid columns, order, hiding. Use atoms: `tg-col md:tg-row tg-gap(md) tg-pad(lg)`. Atoms take token names only (`tg-gap(md)`, never `tg-gap(12px)`). A breakpoint prefix means "from this breakpoint up".
2. **Is it a look that a role already has?** Read the `use` lines in `tidgrid list`. Use the role even when its name does not match your component's name; roles are named by purpose.
3. **Would a small change to an existing role cover it?** A hover, a state, a breakpoint: add a nested rule to that role.
4. **Do the lines you need already sit in a mixin set?** Read the mixins in `tidgrid list`. A new role that needs a text style or a control shape includes the set (`include = [type.heading]`) and writes only what differs.
5. **Only then: a new role.**

An element takes a property from an atom or from a role, never both. Atoms win.

## Adding a role

A role that one element uses once is CSS with a longer name. Before adding one:

- Name it by what it is for (`surface(note)`, `text(lead)`), not by where it sits (`surface(homepage-hero-left)`).
- Check that it holds a look. If every declaration is direction, gap, padding or size, it is layout: use atoms. A role is right for layout only where atoms cannot reach (an overlay, a fixed panel, grid placement).
- Write `use` for the next agent: when the role applies, and when it does not.
- Values are token references (`color.ink`, `space.md`, or `"{color.ink}"` in a string) or keywords, unitless numbers, `calc()` and `color-mix()` over tokens. A raw colour, length or time is refused: add a token first. So is a percentage that sizes a box, a spacing or a grid track (`"width": "40%"`, `"padding": "5%"`): it goes to a token. A percentage elsewhere (`border-radius: 50%`, `left: 50%`, a gradient stop) is fine.
- Hover, focus, children and breakpoints are nested dicts of the role (`"&:hover"`, `"& a"`, `"@md"`), not new roles.
- Put the role and the page that first uses it in the same commit. `check` fails on a role nothing uses, so a role in a commit of its own fails until the page arrives. A change to an existing role or to tokens can be its own commit. Never add a role "for later".

```python
role("text", "lead", use = "Opening paragraph of a page. Once per page.", css = {"color": color.ink, "font-size": sizes.h2})
```

## Mixins

Lines that two or more roles share go to a mixin set, not into both roles. A mixin is a closed list of named sets; a role includes one by name:

```python
type = mixin("type", use = "A text style. Pick one; never set family, size or weight by hand.", sets = {
    "heading": {"font-family": font.title, "font-size": sizes.h2, "font-weight": 700},
    "label": {"font-family": font.sans, "font-size": sizes.label, "font-weight": 600},
})
role("text", "section-title", use = "...", include = [type.heading], css = {"color": color.ink})
```

- A set needs at least two roles that include it and at least two lines (nested ones count); `build` fails otherwise. A single shared value belongs in a token. So make a set when the second user exists, not before.
- A role may not give a property that an included set gives, and two included sets may not overlap; this holds inside nested rules too. A nested rule with the same key in the set and the role merges, the set's lines first.
- A role whose own top-level lines repeat a whole set of two or more top-level lines fails `build` and the message says which `include` to write.
- Repeated lines: when the same lines sit in 3 lines x 2 roles or 2 lines x 4 roles (`project(repeats=[(3, 2), (2, 4)])`, `[]` is off), `build` fails with `role.action.a: own lines also in role action.b: cursor: pointer; padding: 0; border: none; make them a mixin set and include it`. At most 20 groups show, most roles first. When two groups share a line, put it in one set only, or the role gets it twice (`given twice: by set m.s1 and by set m.s2`).
- A set takes no parameters. Variants that differ in a colour (two buttons): a set for the fixed lines plus a plain function for the rest, `css = filled(color.accent, color.ground)`. No custom property that one role sets and another reads; the checker cannot follow that.
- Mixins are Starlark only. A `.toml` vocabulary has none.

## Writing the vocabulary

- Tokens: `color = tokens("color", {"ink": "#e6e9eb"})`. `color.ink` and `color["sand-ink"]` are references; a name the group lacks is an error at that line. A name with a hyphen needs the brackets: `color.sand-ink` is a subtraction and fails with `undefined: ink`.
- A file runs top to bottom: a token group or a mixin is defined above the first role that uses it.
- The refusal of raw lengths and times covers every rule: element rules, nested and pseudo-element rules, offsets, outline widths, shadow sizes. Small one-off sizes go in a token group of their own.
- Values are strings. A `def` that returns a string names a repeated expression (`tint(c, pct)`); a top-level `for` makes several roles. No `while`, no recursion.
- `load("tokens.tg", "color")` imports names; `source("roles.tg")` runs a file for what it declares and imports no names (the sourced file may `load`). A file runs once however often it is named. Paths stay inside the entry file's folder.
- Kind and name of a role are letters, digits, `-` and `_`.
- `project(out=..., scan=[...])`, `breakpoints(...)`, `element(sel, css=...)`, `keyframes(name, {...})`, `allow_class(name, reason)`, `allow_css(path, reason)`.
- `tidgrid plan` shows what the stylesheet would gain, lose or change, without writing. Read it before `build` after a vocabulary edit.

## Class names in code

The checker reads class attributes and string literals in class expressions. It cannot see a name that travels through a variable.

- Write whole names. `` `tg-kind(${kind})` `` is a violation: a built name cannot be checked.
- A lookup table of names is a variable too. Pass each name through a helper so it is read:
  `{ primary: cn("tg-action(primary)") }`. Helpers are `cn clsx cva classNames twMerge` by default; `helpers` in `project()` replaces that list, so repeat the ones in use.
- `check` reports a role as not used although the page shows it: the checker could not see the name. Fix the route, not the report. A role a script sets counts when the script names it whole, as a string. A name that appears only in another role's or element's selector does not count.

## State

A class a script sets is a role of kind `state` with only a `use` line; other roles react to it in a nested rule (`"&.tg-state(open)"`, `".tg-state(ready) &"`). The role is unused for `check` until the script writes the name whole: `card.classList.toggle("tg-state(open)")`. ARIA and data attributes work the same way and need no state role: `"&[aria-current]"`, `"&[data-state=retired]"`.

## Dark mode

Colours stay in tokens, as pairs: `"ink": "light-dark(#32495d, #d7e2ec)"`, and `element("html", css = {"color-scheme": "light dark"})`. The page follows the system with no second stylesheet. For a manual switch: declare `role("state", "dark", use = "Set on html to force the dark scheme.")`, add `"&.tg-state(dark)": {"color-scheme": "dark"}` to the html element's css, and have the script write the name whole, `document.documentElement.classList.add("tg-state(dark)")`. Without that line `check` fails on the unused role.

## Before every commit

```
tidgrid build && tidgrid check
```

Commit the stylesheet `build` wrote. Exit code 0 is fine, 1 means violations (`check`), 2 means the vocabulary does not load.

| Message | Fix |
|---|---|
| `unknown class tg-gap(2); tg-gap takes: sm md` | Use one of the listed token names. |
| `unknown class card-accent: not in the vocabulary; run `tidgrid list`` | Pick a class from `tidgrid list`, or add a role. |
| ``appearance in style attribute: color; use a role (`tidgrid list`)`` | Move the look to a role. Geometry in `style` is allowed. |
| `own CSS: <style> block; make it a role in the vocabulary` | Make it a role. |
| `role tg-text(lead) is not used in the scanned sources: use it, remove it, ...` | Use it in the page of the same commit, or remove it. |
| `stylesheet differs from the vocabulary; run `tidgrid build`` | Run `tidgrid build` and commit the result. |
| `role.text.pad.padding: raw length or time 12px (move it to tokens)` / `raw colour #fff` / `raw size 40% ...` | Add a token, reference it. |
| `role.text.plain: own lines repeat set type.heading; write include = [type.heading] and drop them` | Include the set. |
| `mixin.solo.a: included only by role text.s1; a set needs at least two` | Inline the set into its one role, or find the second user. |
| `role.text.dup.font-weight: given twice: by set type.label and by the role's own css` | Drop the line from the role, or pick another set. |
| `index.tg not found, but tidgrid.toml is here: run `tidgrid convert` ...` | Run `tidgrid convert`. |

A vocabulary message (exit 2) starts with `file:line:` of the call that declared the thing. A page violation (`check`, exit 1) starts with the page and line, such as `src/b.html:1:`. Messages marked `...` are cut.

## Never

- Write a stylesheet, a `<style>` block or a CSS module. The only way in is `allow_css(path, reason)`, for what the tool cannot write (an `@font-face` file).
- Edit the generated stylesheet.
- Add to `allow_class` to silence the checker. It is for classes a third-party script needs, each with its reason.
- Wrap markup in `tidgrid-ignore-start` to get past a violation. It is for code samples that show classes.
- Add a role, or a set with one user, that nothing needs yet.

## Setting up a project

`tidgrid init [file]` writes a starter `index.tg` with grey colour tokens, the reset `tidiness("full")` (nothing has a margin, a padding, a list marker or a link colour until an element rule or a role gives one; `"minimal"` keeps the browser's looks) and the starter roles commented out. Set `out` and `scan` in `project()`, put the project's colours and scales in the token groups, uncomment a role when a page uses it, then run the check in CI: `tidgrid build`, fail if the stylesheet changed, `tidgrid check`. `tidgrid sheet` writes one HTML page that shows every token and role, for a person to look at.

Reference: https://tidgrid.dev/docs
