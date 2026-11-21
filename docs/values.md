---
layout: doc-page
title: Values
---

Outside `[tokens]` and `[breakpoints]`, a value in the vocabulary file may not contain a raw colour, a raw length or a raw time. Those live in tokens. `build` refuses the file if one slips through.

## Not allowed

| Not allowed | Examples |
|---|---|
| A raw colour | `#hex`, `rgb()`, `hsl()`, `oklch()` and the like, or a CSS colour name such as `red` |
| A raw length | a number with any CSS length unit: `px`, `rem`, `vw`, `dvh`, `cqw` |
| A raw time | a number with `ms` or `s` |

## Allowed

| Allowed | Examples |
|---|---|
| Keywords | `transparent`, `currentColor` |
| Zero and unitless numbers | `0`, `1.75` |
| Percentages, `fr` and angles | `50%`, `1fr`, `45deg` |
| Functions over tokens | shown below |

{% capture c %}
padding    = "calc({space.md} * 2)"
background = "color-mix(in srgb, {color.accent} 18%, transparent)"
{% endcapture %}
{% include code.html code=c %}

Text in quotes and inside `url()` is not looked at.

## At-rule keys

The key of an at-rule sub-table follows the same rule. A reference in it is replaced with the token's written value, because a media query cannot read a custom property:

{% capture c %}
[role.layout.panel."@media (max-height: {size.short})"]
display = "none"
{% endcapture %}
{% include code.html code=c %}

## Next

- [Tokens](/docs/tokens)
- [Roles](/docs/roles)
