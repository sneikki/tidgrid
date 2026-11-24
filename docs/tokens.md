---
layout: doc-page
title: Tokens
---

Tokens are the only place for raw colours, lengths and times. Everything else in the vocabulary refers to them by name.

{% capture c %}
color = tokens("color", {
    "ground": "#0f1214",
    "accent": "#5ed5dc",
    "link": "{color.accent}",
})
space = tokens("space", {"sm": "0.5rem", "md": "1rem"})
border = tokens("border", {"hair": "1px"})
motion = tokens("motion", {"slow": "1.1s"})
{% endcapture %}
{% include code.html code=c caption="index.tg" %}

`tokens("<group>", {...})` takes names and raw values. Each one becomes a custom property named `--tg-<group>-<name>` on `:root`, so scripts can read it. The example above gives `--tg-color-ground`, `--tg-space-md` and so on. `tokens(group, name=value, ...)` works too.

## Tokens as variables

A token group is also a Starlark value, which is why the call is assigned to a variable. `color.ground` is the reference `"{color.ground}"`, and `color["sand-ink"]` reaches a name that is not a plain identifier. Use them in a `css` dict, or in a string you build:

{% capture c %}
role("surface", "card",
    use = "A panel that lifts off the ground. One per group of related content.",
    css = {"background": color.ground, "padding": space.md,
           "border": "%s solid %s" % (border.hair, color.accent)})
{% endcapture %}
{% include code.html code=c caption="index.tg, below the tokens" %}

A name the group does not have is an error at that line:

{% capture c %}
index.tg:3:52: token group color has no token "nope"
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

Written as a plain string, `"{color.nope}"` is the same mistake, reported as `unknown reference`. A variable defined in one file reaches the others with [`load`](/docs/files).

## Referring to a token

In any value, `{group.name}` is replaced with `var(--tg-group-name)`. A token may refer to another, as `link` does above.

Three things are replaced differently:

- Breakpoints are substituted literally, because a media query cannot read a custom property.
- `{keyframes.<name>}` becomes the generated animation name `tg-<name>`.
- A reference in the key of an at-rule nested rule is replaced with the token's written value, for the same reason as breakpoints.

Group and token names are made of letters, digits, `-` and `_`. `keyframes` and `breakpoints` are not group names.

## Groups with a job

Two groups feed the layout atoms. Their names are fixed, their contents are yours.

- `space` is the scale for the gap and padding atoms: `tg-gap(md)` exists because `md` is a name in `space`.
- `size` is the scale for the width and height atoms: `tg-max-width(measure)` exists because `measure` is a name in `size`.

Every other group name is free: `color`, `border`, `motion`, `type`, `radius`, `shadow` or whatever the project needs. Roles can use all of them.

## Next

- [Layout atoms](/docs/layout-atoms)
- [Values](/docs/values)
