---
layout: doc-page
title: Values
---

Outside `[tokens]` and `[breakpoints]`, a value in the vocabulary file may not contain a raw colour, a raw length, a raw time or a percentage used as a size. Those live in tokens. `build` refuses the file if one slips through.

## Not allowed

| Not allowed | Examples |
|---|---|
| A raw colour | `#hex`, `rgb()`, `hsl()`, `oklch()` and the like, or a CSS colour name such as `red` |
| A raw length | a number with any CSS length unit: `px`, `rem`, `vw`, `dvh`, `cqw` |
| A raw time | a number with `ms` or `s` |
| A percentage that sizes a box, a spacing or a grid track | `width = "40%"`, `padding = "5%"`, `grid-template-columns = "50% 50%"` |

The percentage rule covers `width`, `height`, `inline-size`, `block-size`, their `min-` and `max-` forms, `flex-basis`, `flex`, every `margin` and `padding` property, the three `gap` properties, and the `grid-template` and `grid-auto-` properties. Give the size a name instead: `half = "50%"` under `[tokens.size]`, then `width = "{size.half}"`.

## Allowed

| Allowed | Examples |
|---|---|
| Keywords | `transparent`, `currentColor` |
| Zero and unitless numbers | `0`, `1.75` |
| Percentages anywhere else | `border-radius = "50%"`, `left = "50%"`, a gradient stop, `translate(-50%)` |
| `fr` and angles | `1fr`, `45deg` |
| Functions over tokens | shown below |

{% capture c %}
padding    = "calc({space.md} * 2)"
background = "color-mix(in srgb, {color.accent} 18%, transparent)"
{% endcapture %}
{% include code.html code=c caption="tidgrid.toml" %}

Text in quotes and inside `url()` is not looked at.

## At-rule keys

The key of an at-rule sub-table follows the same rule. A reference in it is replaced with the token's written value, because a media query cannot read a custom property:

{% capture c %}
[role.layout.panel."@media (max-height: {size.short})"]
display = "none"
{% endcapture %}
{% include code.html code=c caption="tidgrid.toml" %}

## Next

- [Tokens](/docs/tokens)
- [Roles](/docs/roles)
