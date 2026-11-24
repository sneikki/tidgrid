---
layout: doc-page
title: The vocabulary file
---

A project has one vocabulary file, `tidgrid.toml`. It defines every class a page may use. Anything not in it, `check` rejects.

Here is a whole one. The tables below it are explained one by one.

{% capture c %}
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
{% endcapture %}
{% include code.html code=c caption="tidgrid.toml" %}

## The tables

- `[project]` says where the stylesheet goes (`out`, required), which folders `check` scans (`scan`, default `["."]`) and the class prefix (`prefix`, default `tg`). `project.helpers` sets the function names whose string arguments are class names. It replaces the default list (`cn`, `clsx`, `cva`, `classNames`, `twMerge`), so list those too if you still use them. `check` reports a scan path that does not exist.
- `[breakpoints]` names the screen widths, smallest first. Every layout atom except `tg-hide-below` gets a prefixed variant per breakpoint. See [layout atoms](/docs/layout-atoms).
- `[tokens.<group>]` holds the raw values. See [tokens](/docs/tokens).
- `[element.<selector>]` is a base rule for a bare element such as `body`, `a` or `"*"`. It is written like a role, and `use` is optional.
- `[role.<kind>.<name>]` is a class with a look. See [roles](/docs/roles).
- `[keyframes.<name>]` holds one sub-table per stop. A role refers to it as `{keyframes.<name>}`, which becomes the generated animation name `tg-<name>`.
- `[allow.class]` and `[allow.css]` are visible exceptions. An allowed class or stylesheet is not reported, and its reason stays in the file and in the `check` summary.

## Order

Roles, declarations and nested rules are written in file order. The stylesheet as a whole is written in this order, so a later part wins over an earlier one: tokens, elements, keyframes, roles, atoms, then the atoms of each breakpoint in file order, then the `tg-hide-below` atoms, which therefore win over a `md:tg-row` on the same element. `tg-hide` is the last atom of its block.

## Next

- [Tokens](/docs/tokens)
- [Values](/docs/values)
