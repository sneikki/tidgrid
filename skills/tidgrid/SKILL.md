---
name: tidgrid
description: Use when building or changing UI in a project that has a tidgrid.toml (classes such as tg-row, tg-gap(md), tg-surface(card)), when `tidgrid check` fails, or when adding a colour, a spacing value, a role, a state or dark mode to such a project. Also when setting tidgrid up in a new project.
---

# Styling with Tidgrid

Tidgrid is a closed CSS vocabulary. One file, `tidgrid.toml`, defines every class a page may use. `tidgrid build` writes the stylesheet from it and `tidgrid check` fails on anything outside it. You never write CSS; you pick classes, and when none fits you change the vocabulary.

Run every command in the directory that holds `tidgrid.toml`. Install: `go install github.com/sneikki/tidgrid/cmd/tidgrid@latest`.

## Before any UI work

Run `tidgrid list` and read it. It prints the tokens, the atoms with the values they take, every role with the sentence that says when to use it, and the exceptions. Do not guess a class from memory or from another project: role names differ per project.

## Picking a class

Ask in this order and stop at the first yes.

1. **Is it layout?** Direction, gap, padding, width, height, grid columns, order, hiding. Use atoms: `tg-col md:tg-row tg-gap(md) tg-pad(lg)`. Atoms take token names only (`tg-gap(md)`, never `tg-gap(12px)`). A breakpoint prefix means "from this breakpoint up".
2. **Is it a look that a role already has?** Read the `use` lines in `tidgrid list`. Use the role even when its name does not match your component's name; roles are named by purpose.
3. **Would a small change to an existing role cover it?** A hover, a state, a breakpoint: add a nested rule to that role.
4. **Only then: a new role.**

An element takes a property from an atom or from a role, never both. Atoms win.

## Adding a role

A role that one element uses once is CSS with a longer name. Before adding one:

- Name it by what it is for (`surface(note)`, `text(lead)`), not by where it sits (`surface(homepage-hero-left)`).
- Check that it holds a look. If every declaration is direction, gap, padding or size, it is layout: use atoms. A role is right for layout only where atoms cannot reach (an overlay, a fixed panel, grid placement).
- Write `use` for the next agent: when the role applies, and when it does not.
- Values are token references (`{color.ink}`, `{space.md}`) or keywords, unitless numbers, `calc()` and `color-mix()` over tokens. A raw colour, length or time is refused: add a token first. So is a percentage that sizes a box (`width = "40%"`): it goes to `tokens.size`. A percentage elsewhere (`border-radius: 50%`, a gradient stop) is fine.
- Hover, focus, children and breakpoints are nested sub-tables of the role (`"&:hover"`, `"& a"`, `"@md"`), not new roles.
- Commit the vocabulary change on its own (`tidgrid.toml` and the rebuilt stylesheet), before the page that uses it.

Variants that share a shape and differ in one colour: one role per variant is fine. Avoid passing the colour through a custom property that one role sets and another reads; the checker cannot follow that.

## Class names in code

The checker reads class attributes and string literals in class expressions. It cannot see a name that travels through a variable.

- Write whole names. `` `tg-kind(${kind})` `` is a violation: a built name cannot be checked.
- A lookup table of names is a variable too. Pass each name through a helper so it is read:
  `{ primary: cn("tg-action(primary)") }`. Helpers are `cn clsx cva classNames twMerge` by default; `project.helpers` replaces that list, so repeat the ones in use.
- A role that `check` lists under "roles unused" although the page shows it is a name the checker could not see. Fix the route, not the report.

## State

A class a script sets is a role of kind `state` with only a `use` line; other roles react to it in a nested rule (`"&.tg-state(open)"`, `".tg-state(ready) &"`). ARIA and data attributes work the same way and need no state role: `"&[aria-current]"`, `"&[data-state=retired]"`.

## Dark mode

Colours stay in tokens, as pairs: `ink = "light-dark(#32495d, #d7e2ec)"`, and `[element.html]` sets `color-scheme = "light dark"`. The page follows the system with no second stylesheet. For a manual switch, a script sets a state class on `html` and the vocabulary narrows the scheme: `[element.html."&.tg-state(dark)"] color-scheme = "dark"`.

## Before every commit

```
tidgrid build && tidgrid check
```

Commit the stylesheet `build` wrote. Exit code 0 is fine, 1 means violations, 2 means the vocabulary file is invalid or missing.

| Message | Fix |
|---|---|
| `unknown class tg-gap(2); tg-gap takes: sm md` | Use one of the listed token names. |
| `unknown class card-accent: not in the vocabulary` | Pick a class from `tidgrid list`, or add a role. |
| `appearance in style attribute: color` | Move the look to a role. Geometry in `style` is allowed. |
| `own CSS: <style> block` | Make it a role. |
| `stylesheet differs from the vocabulary` | Run `tidgrid build` and commit the result. |
| `raw colour …` / `raw length or time …` / `raw size …` | Add a token, reference it. |

## Never

- Write a stylesheet, a `<style>` block or a CSS module. The only way in is `[allow.css]` with a reason, for what the tool cannot write (an `@font-face` file).
- Edit the generated stylesheet.
- Add to `[allow.class]` to silence the checker. It is for classes a third-party script needs, each with its reason.
- Wrap markup in `tidgrid-ignore-start` to get past a violation. It is for code samples that show classes.

## Setting up a project

`tidgrid init` writes a starter `tidgrid.toml`. Set `project.out` and `project.scan`, put the project's colours and scales in `[tokens]`, then run the check in CI: `tidgrid build`, fail if the stylesheet changed, `tidgrid check`. `tidgrid sheet` writes one HTML page that shows every token and role, for a person to look at.

Reference: https://tidgrid.dev/docs
