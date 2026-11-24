---
layout: doc-page
title: The vocabulary file
---

A project's vocabulary is written in Starlark, a small Python-like language. Its entry file is `index.tg`. The vocabulary defines every class a page may use. Anything not in it, `check` rejects. The point is a closed vocabulary, decided once and checked.

Here is a small whole one. The calls below it are explained one by one.

{% capture c %}
project(out = "src/styles/tidgrid.css", scan = ["src"])
breakpoints(sm = "40em", md = "64em")

color = tokens("color", {         # tokens are the ONLY place for raw values
    "ground": "#0f1214",
    "ink": "#e6e9eb",
    "accent": "#5ed5dc",
    "link": "{color.accent}",     # a token may reference another
})
space = tokens("space", {"sm": "0.5rem", "md": "1rem"})        # gap and pad atoms run on this group
size = tokens("size", {"full": "100%", "measure": "80rem"})    # width and height atoms run on this group
border = tokens("border", {"hair": "1px"})                     # any other group name is free
motion = tokens("motion", {"slow": "1.1s"})

element("body", css = {"background": color.ground, "color": color.ink})

rise = keyframes("rise", {
    "0%": {"transform": "translateY(%s)" % space.sm},
    "100%": {"transform": "none"},
})

role("surface", "card",             # class: tg-surface(card)
    use = "A panel that lifts off the ground. One per group of related content.",
    css = {
        "background": color.ground,
        "border": "%s solid %s" % (border.hair, color.accent),
        "&:hover": {"border-color": color.ink},   # nested rule, native CSS nesting
        "@md": {"border-width": border.hair},     # from breakpoint md up
    })

role("state", "ready", use = "Set on body when the intro animation has finished.")

role("text", "title",
    use = "The one title of a page.",
    css = {
        "opacity": "0",
        ".tg-state(ready) &": {"opacity": "1", "animation": "%s %s ease" % (rise, motion.slow)},
    })

allow_class("has-js", "set by an inline script before first paint")
allow_css("src/styles/fonts.css", "self-hosted @font-face only")
{% endcapture %}
{% include code.html code=c caption="index.tg" %}

## The calls

The tool runs the entry file. Every call registers one thing in a single model, and the rules, `list`, `check`, `plan` and the stylesheet all run on that model. The calls are [builtins](/docs/builtins), available without importing.

- `project` says where the stylesheet goes (`out`, required), which folders `check` scans (`scan`, default `["."]`) and the class prefix (`prefix`, default `tg`). `helpers` sets the function names whose string arguments are class names. It replaces the default list (`cn`, `clsx`, `cva`, `classNames`, `twMerge`), so list those too if you still use them. `check` reports a scan path that does not exist. `repeats` sets the [repeated-lines rule](/docs/mixins#repeated-lines).
- `breakpoints` names the screen widths, smallest first. Every layout atom except `tg-hide-below` gets a prefixed variant per breakpoint. See [layout atoms](/docs/layout-atoms).
- `tokens` makes a group of raw values and returns it. See [tokens](/docs/tokens).
- `element` is a base rule for a bare element such as `body`, `a` or `"*"`. It is written like a role, and `use` is optional.
- `role` is a class with a look. See [roles](/docs/roles).
- `mixin` is a named set of lines that roles share. See [mixins](/docs/mixins).
- `keyframes` holds one entry per stop. It returns a reference that the build turns into the animation name, `tg-<name>`.
- `allow_class` and `allow_css` are visible exceptions. An allowed class or stylesheet is not reported, and its reason stays in the file and in the `check` summary.

Declaring `project`, a token group, a role or a mixin twice is an error that names both places.

## Order

Roles, declarations and nested rules are written in the order the calls ran and the dicts were written. The stylesheet as a whole is written in this order, so a later part wins over an earlier one: tokens, elements, keyframes, roles, atoms, then the atoms of each breakpoint in breakpoint order, then the `tg-hide-below` atoms, which therefore win over a `md:tg-row` on the same element. `tg-hide` is the last atom of its block.

## Starlark here

Top-level `for` and `if` work. So do `def`, list comprehensions and string formatting (`"%s %s" % (a, b)`). `while` and recursion do not, there is a step limit for the whole run, and nothing reads or writes outside the vocabulary's folder.

## Next

- [Files](/docs/files): a vocabulary over several files
- [Builtins](/docs/builtins)
- [Tokens](/docs/tokens)
- [Values](/docs/values)
