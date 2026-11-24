---
layout: doc-page
title: Roles
---

Appearance lives in roles. Colour, border, shadow and type are reachable only through them: there is no colour atom and no border atom. What a role looks like is decided once, in the vocabulary file, and reviewed there.

{% capture c %}
[role.surface.card]
use        = "A panel that lifts off the ground. One per group of related content."
background = "{color.ground}"
border     = "{border.hair} solid {color.accent}"

[role.surface.card."&:hover"]
border-color = "{color.ink}"
{% endcapture %}
{% include code.html code=c caption="tidgrid.toml" %}

`[role.<kind>.<name>]` becomes the class `tg-<kind>(<name>)`. The role above is `tg-surface(card)`. Kinds are free words (`surface`, `text`, `emphasis`, `action`, `state` and so on), but a kind may not have the name of an atom.

## What goes in a role

- `use` is required. It is one or two sentences that tell an agent when the role applies, and `tidgrid list` prints it. It is not written to the stylesheet, and it is an error inside a sub-table.
- Every other string key is a CSS declaration.
- A sub-table is a nested rule. Its key is one of:
  - `@<breakpoint>`: from that breakpoint up.
  - `@below-<breakpoint>`: below that breakpoint.
  - any other `@…`: the at-rule as written, such as `@media (prefers-reduced-motion: reduce)` or `@supports …`.
  - a selector that contains `&`, such as `&:hover`, `&::before`, `& h1` or `.tg-state(ready) &`. Vocabulary classes in it are written plainly and must exist.
- Sub-tables nest to any depth.

Values follow the rules on the [values](/docs/values) page: colours, lengths and times come from tokens.

## State roles

A role may have no declarations at all. It is then a class that a script sets, and other roles react to:

{% capture c %}
[role.state.ready]
use = "Set on body when the intro animation has finished."

[role.text.title]
use     = "The one title of a page."
opacity = "0"

[role.text.title.".tg-state(ready) &"]
opacity = "1"
{% endcapture %}
{% include code.html code=c caption="tidgrid.toml" %}

## Elements and keyframes

`[element.<selector>]` works the same way for a bare element (`body`, `a`, `"*"`). `use` is optional there, and the element gets no class.

`[keyframes.<name>]` holds one sub-table per stop (`"0%"`, `"100%"`). A role refers to it as `{keyframes.<name>}`.

## When to add a role

Only when nothing in the vocabulary fits. A new role is a change to `tidgrid.toml`, in its own commit, never mixed into page work. That keeps the change to the look visible in review.

## Next

- [Values](/docs/values)
- [Commands](/docs/commands): `tidgrid sheet` renders every role
