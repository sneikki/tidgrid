---
layout: doc-page
title: Mixins
---

A mixin is a closed list of named sets, each a `css` dict, that roles share. A role includes sets by name. It renders as if the set's lines were written first, then its own. A mixin puts nothing else into the stylesheet and is never a class.

{% capture c %}
load("tokens.tg", "color", "space", "font", "sizes", "tint")

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
{% endcapture %}
{% include code.html code=c caption="roles.tg" %}

`type.heading` and `type["heading"]` are both a set. `tidgrid list` prints the mixins with the roles that include each set:

{% capture c %}
# mixins
type  A text style. Pick one; never set family, size or weight by hand.
  heading: font-family: {font.title}; font-size: {sizes.h2}; font-weight: 700
    included by: tg-text(section-title) tg-text(card-title)
  label: font-family: {font.sans}; font-size: {sizes.label}; font-weight: 600
    included by: tg-text(caption) tg-text(hint)
shape  Shared shapes of controls.
  button: padding: {space.sm} {space.md}; cursor: pointer
    included by: tg-action(primary) tg-action(secondary)
{% endcapture %}
{% include code.html code=c plain=true caption="Part of what tidgrid list prints" %}

## The rules

- A set takes no parameters. A set is a step of a scale: `type.heading`, `type.label`.
- `include` takes only mixin sets. Anything else is an error.
- Nested rules with the same key, in a set and in the role (`"&:hover"`), merge into one rule: the set's lines first, then the role's.
- The same property given twice is an error, at any depth: a role's own property that an included set already gives, or two included sets that give the same one, also inside a nested rule. The message names the path.
- A set needs at least two roles that include it. Only roles that exist count, and an unused role fails `check`, so a role invented to give a set its second user does not survive.
- A set needs at least two lines in total, nested ones counted. A single shared value belongs in a token.
- A role whose own top-level lines contain every top-level line of a set (same property, same value as written) without including it fails `build`, and the message gives the `include` to write. This only applies to sets with two or more top-level lines. A role may still write by hand what a mixin does not cover.
- `mixin` needs `use`. A mixin name declared twice is an error.

Each rule fails `build` with exit 2. The messages, from real runs:

{% capture c %}
index.tg:6: role.action.a.color: given twice: by set m.a and by the role's own css
index.tg:5: mixin.m.a: included only by role action.a; a set needs at least two
index.tg:5: mixin.m.one: a set needs at least two lines; a single shared value belongs in a token
index.tg:8: role.action.c: own lines repeat set m.a; write include = [m.a] and drop them
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

## Variants of a role

A role that comes in variants (buttons that differ in colour) needs no feature of its own. The shared, fixed lines are a mixin set, `shape.button` above. The lines that differ come from a plain function, `filled` above, which may put the colour in as many places as it likes. Class names stay as they are (`tg-action(primary)`, `tg-action(secondary)`). No private custom property carries the colour, so `check` sees every value.

## Repeated lines

A line is a top-level (property, value as written) pair of a role's own `css`. Lines that come from an included set do not count. `build` fails when the same lines appear in several roles: 3 lines in 2 roles, or 2 lines in 4 roles. Each group is reported once, with the roles and the lines they share, and the message says the way out: make the lines a mixin set and include it.

{% capture c %}
index.tg:5: role.action.a: own lines also in role action.b: cursor: pointer; padding: 0; border: none; make them a mixin set and include it
index.tg: 1 group of repeated lines; fix the first and build again: one mixin set often clears several
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

The first line is at the first role of the group. The last line is the total. At most 20 groups are shown, the ones with the most roles first.

The numbers are a setting, a list of (lines, roles) pairs. `project(repeats=[(3, 2), (2, 4)])` is the default, and `repeats = []` turns the rule off. `tidgrid convert` writes `repeats = []`, because an old project repeats lines freely. Cleaning it up starts by removing that argument.

Two reported groups may share a line. Put the shared line in one set only. If a role includes two sets that both give it, `build` says `given twice`.

## Next

- [Roles](/docs/roles)
- [Values](/docs/values)
- [Commands](/docs/commands): `tidgrid convert`
