---
layout: doc-page
title: Values
---

Outside `tokens` and `breakpoints`, a value in the vocabulary may not contain a raw colour, a raw length, a raw time or a percentage used as a size. Those live in tokens. `build` refuses the vocabulary if one slips through, and prints the file, the line and the path to the value:

{% capture c %}
bad.tg:2: role.layout.p.width: raw size 40% (move it to tokens)
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

## Not allowed

| Not allowed | Examples |
|---|---|
| A raw colour | `#hex`, `rgb()`, `hsl()`, `oklch()` and the like, or a CSS colour name such as `red` |
| A raw length | a number with any CSS length unit: `px`, `rem`, `vw`, `dvh`, `cqw` |
| A raw time | a number with `ms` or `s` |
| A percentage that sizes a box, a spacing or a grid track | `"width": "40%"`, `"padding": "5%"`, `"grid-template-columns": "50% 50%"` |

The percentage rule covers `width`, `height`, `inline-size`, `block-size`, their `min-` and `max-` forms, `flex-basis`, `flex`, every `margin` and `padding` property, the three `gap` properties, and the `grid-template` and `grid-auto-` properties. Give the size a name instead: `"half": "50%"` in the `size` group, then `"width": size.half`.

## Allowed

| Allowed | Examples |
|---|---|
| Keywords | `transparent`, `currentColor` |
| Zero and unitless numbers | `0`, `1.75` |
| Percentages anywhere else | `"border-radius": "50%"`, `"left": "50%"`, a gradient stop, `translate(-50%)` |
| `fr` and angles | `1fr`, `45deg` |
| Functions over tokens | shown below |

Text in quotes and inside `url()` is not looked at.

## Functions and loops

A value is a string, so ordinary Starlark builds it. A `def` that returns a value names an expression. A loop makes several roles. A token group in a string is a reference, so `%s` puts it in:

{% capture c %}
def tint(c, pct):
    return "color-mix(in srgb, %s %d%%, transparent)" % (c, pct)

for n in ["info", "warn"]:
    role("notice", n, use = "A notice that is %s." % n, css = {"color": color.ink})

role("layout", "panel", use = "A panel.", css = {
    "width": size.half,
    "padding": "calc({space.md} * 2)",
    "background": tint(color.accent, 18),
})
{% endcapture %}
{% include code.html code=c caption="index.tg" %}

Role names may be computed. `check` fails on a role no page uses, and every role still needs its own written `use`. The finished values are validated exactly as written by hand. `while` and recursion do not exist; a function that calls itself fails with `called recursively`.

## At-rule keys

The key of an at-rule nested rule follows the same rule. A reference in it is replaced with the token's written value, because a media query cannot read a custom property:

{% capture c %}
role("layout", "panel", use = "A panel.", css = {
    "@media (max-height: {size.short})": {"display": "none"},
})
{% endcapture %}
{% include code.html code=c caption="index.tg" %}

## Next

- [Tokens](/docs/tokens)
- [Roles](/docs/roles)
