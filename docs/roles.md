---
layout: doc-page
title: Roles
---

Appearance lives in roles. Colour, border, shadow and type are reachable only through them: there is no colour atom and no border atom. What a role looks like is decided once, in the vocabulary, and reviewed there.

{% capture c %}
role("surface", "card",
    use = "A panel that lifts off the ground. One per group of related content.",
    css = {
        "background": color.ground,
        "border": "%s solid %s" % (border.hair, color.accent),
        "&:hover": {"border-color": color.ink},
    })
{% endcapture %}
{% include code.html code=c caption="index.tg" %}

`role(kind, name, ...)` becomes the class `tg-<kind>(<name>)`. The role above is `tg-surface(card)`. Kinds are free words (`surface`, `text`, `emphasis`, `action`, `state` and so on), but a kind may not have the name of an atom (`kind collides with atom tg-row`). Kind and name are made of letters, digits, `-` and `_`.

## What goes in a role

- `use` is required. It is one or two sentences that tell an agent when the role applies, and `tidgrid list` prints it. It is an argument, not a key of `css`, and it is not written to the stylesheet.
- `css` is a dict of declarations and nested rules, described below.
- `include` lists [mixin](/docs/mixins) sets.

## The `css` dict

Keys are strings.

- A string or number value is a declaration: `"font-weight": 700`, `"padding": "{space.md}"`.
- A dict value is a nested rule. Its key is one of:
  - `"@md"`: from breakpoint `md` up.
  - `"@below-md"`: below breakpoint `md`.
  - any other `@…`: the at-rule as written, such as `"@media (prefers-reduced-motion: reduce)"` or `"@supports …"`.
  - a selector that contains `&`, such as `"&:hover"`, `"&::before"`, `"& h1"` or `".tg-state(ready) &"`. Vocabulary classes in it are written plainly and must exist.
- Nested rules nest to any depth.

Declarations and nested rules are written in the order the dict has them. Values follow the rules on the [values](/docs/values) page: colours, lengths and times come from tokens.

## State roles

A role may have no `css` at all. It is then a class that a script sets, and other roles react to in a nested rule:

{% capture c %}
role("state", "ready", use = "Set on body when the intro animation has finished.")

role("text", "title",
    use = "The one title of a page.",
    css = {"opacity": "0", ".tg-state(ready) &": {"opacity": "1"}})
{% endcapture %}
{% include code.html code=c caption="index.tg" %}

ARIA and data attributes need no role: `"&[aria-current]"`.

A name written only in another role's or element's selector does not count as used. `check` fails on the state role until the script that sets the class writes the name whole, as a string the scanner reads:

<!-- tidgrid-ignore-start -->
{% capture c %}
body.classList.add("tg-state(ready)");
{% endcapture %}
{% include code.html code=c caption="A script" %}
<!-- tidgrid-ignore-end -->

## Elements and keyframes

`element(selector, css=...)` works the same way for a bare element (`body`, `a`, `"*"`). `use` is optional there, and the element gets no class.

`keyframes(name, {stop: css})` holds one entry per stop (`"0%"`, `"100%"`) and returns a reference to put in a value: `"%s {motion.slow} ease" % rise`.

## Roles in loops

Role names may be computed, so a loop can make several roles. Every role still needs its own written `use`. See [values](/docs/values).

## When to add a role

Only when nothing in the vocabulary fits. A new role goes in the same commit as the first page that uses it. `check` fails on a role no page uses, so a role in a commit of its own fails until the page arrives. A change to an existing role or to tokens can still be a commit of its own.

`tidgrid init` writes its starter roles commented out, so a fresh project has no unused role and `check` passes. Uncomment a role when a page uses it.

## Next

- [Mixins](/docs/mixins)
- [Values](/docs/values)
- [Commands](/docs/commands): `tidgrid sheet` renders every role
