---
layout: doc-page
title: Tokens
---

Tokens are the only place for raw colours, lengths and times. Everything else in the vocabulary file refers to them by name.

{% capture c %}
[tokens.color]
ground = "#0f1214"
accent = "#5ed5dc"
link   = "{color.accent}"

[tokens.space]
sm = "0.5rem"
md = "1rem"

[tokens.motion]
slow = "1.1s"
{% endcapture %}
{% include code.html code=c caption="tidgrid.toml" %}

`[tokens.<group>]` holds `name = "raw value"`. Each one becomes a custom property named `--tg-<group>-<name>` on `:root`, so scripts can read it. The example above gives `--tg-color-ground`, `--tg-space-md` and so on.

## Referring to a token

In any value, `{group.name}` is replaced with `var(--tg-group-name)`. A token may refer to another, as `link` does above.

Three things are replaced differently:

- Breakpoints are substituted literally, because a media query cannot read a custom property.
- `{keyframes.<name>}` becomes the generated animation name `tg-<name>`.
- A reference in the key of an at-rule sub-table is replaced with the token's written value, for the same reason as breakpoints.

Group and token names are made of letters, digits, `-` and `_`.

## Groups with a job

Two groups feed the layout atoms. Their names are fixed, their contents are yours.

- `tokens.space` is the scale for the gap and padding atoms: `tg-gap(md)` exists because `md` is a name in `tokens.space`.
- `tokens.size` is the scale for the width and height atoms: `tg-max-width(measure)` exists because `measure` is a name in `tokens.size`.

Every other group name is free: `color`, `border`, `motion`, `type`, `radius`, `shadow` or whatever the project needs. Roles can use all of them.

## Next

- [Layout atoms](/docs/layout-atoms)
- [Values](/docs/values)
