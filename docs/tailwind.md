---
layout: doc-page
title: Tidgrid and Tailwind
---

Tailwind scans the sources as plain text and generates CSS for every token that matches a utility. A token that matches none is thrown away: it gets no CSS and no message from the build, because to the scanner a mistyped class is one more word in the file.

Tidgrid generates the stylesheet too, and the checking is built in. It checks the pages against the vocabulary and the vocabulary against its own rules, and `build` or `check` fails when something is wrong.

| Stage | What it catches |
|---|---|
| `check`, from the pages | an unknown class, a class name built from pieces, appearance in `style=`, own CSS, a stale stylesheet, a role nothing uses |
| `build`, from the vocabulary | a raw value outside the tokens, the same lines in several roles, a mixin set that is not shared, a role without a `use` line |
| `plan` | what a change does to the stylesheet, before it is built |

## No appearance atoms

Tailwind has an atom for colour, border and shadow: `border-l-4`, `border-blue-500`, `shadow-md`. An agent composes a look on the element it is writing. A coloured stripe on the left edge of a box is two such atoms.

In Tidgrid colour, border, shadow and type are reachable only through a role, and a role is decided once, in the vocabulary. Layout stays free: atoms on token scales, because geometry does not pull a design off course. See [roles](/docs/roles) and [layout atoms](/docs/layout-atoms).

The same content, as an agent writes it with utility classes and with roles:

{% comment %}Every sample on this page shows classes, it does not use them.{% endcomment %}
<!-- tidgrid-ignore-start -->
<div class="tg-col tg-y-gap(md)">
  <div class="tg-col md:tg-row tg-gap(md)">
    <div class="tg-grow tg-col tg-y-gap(xs)">
{% capture c %}
<div class="rounded-lg border-l-4 border-blue-500 bg-white p-4 shadow-md">
  <h2 class="text-xl font-semibold">Plan</h2>
  <p class="text-gray-600">Renews on 1 March.</p>
</div>
{% endcapture %}
{% include code.html code=c caption="Utility classes" %}
    </div>
    <div class="tg-grow tg-col tg-y-gap(xs)">
{% capture c %}
<div class="tg-surface(raised) tg-col tg-y-gap(xs) tg-pad(md)">
  <h2 class="tg-text(heading)">Plan</h2>
  <p class="tg-text(label)">Renews on 1 March.</p>
</div>
{% endcapture %}
{% include code.html code=c caption="Roles" %}
    </div>
  </div>
{% capture c %}
role("surface", "raised",
    use = "A panel holding one group of related content. Use for most boxes; at most a few per view.",
    css = {
        "background": color.raised,
        "border": "%s solid %s" % (border.hair, color.line),
        "border-radius": radius.sm,
    })
{% endcapture %}
{% include code.html code=c caption="The role behind tg-surface(raised), in index.tg" %}
</div>
<!-- tidgrid-ignore-end -->

The agent that wants a stripe has two ways to write one, and `check` stops both. A utility class is not in the vocabulary, and the same look in `style=` is appearance on the element:

{% capture c %}
src/index.html:1: unknown class border-l-4: not in the vocabulary; run `tidgrid list`
src/index.html:1: appearance in style attribute: border-left; use a role (`tidgrid list`)
{% endcapture %}
{% include code.html code=c plain=true caption="What check prints" %}

The third way is a role of its own, and `build` holds that to the tokens:

{% capture c %}
index.tg:196: role.surface.stripe.use: required: when the role applies
index.tg:196: role.surface.stripe.border-left: raw colour #2456c4
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

A new role is a decision in the vocabulary, where it is reviewed, not a choice on one element.

## A value of your own

In Tailwind a value outside the theme is a class: `top-[117px]` and `bg-[#bada55]` are in its documentation. In Tidgrid a raw colour, length or time outside the tokens fails `build`, and an atom takes only token names. The value goes in a token first, see [values](/docs/values).

## Written from memory

An agent knows Tailwind's class names by heart, and writes them from memory: the same defaults any other project gets. A Tidgrid vocabulary has the names of its own project. The agent reads `tidgrid list`, the roles with the `use` line that says when each one applies, and what it reads is this project's decisions. A name from memory is an unknown class, and `check` says so.

## Row by row

| | Tailwind | Tidgrid |
|---|---|---|
| Colour, border, shadow | classes on the element: `border-blue-500`, `shadow-md` | only through a role |
| A value of your own | `w-[137px]` and `bg-[#bada55]` generate a class | `build` fails; the value belongs in a token |
| An unknown class | thrown away, no CSS | `check` fails |
| The same look in several places | the class list is repeated, or the project extracts it: a loop, a component or partial, a custom CSS class | one role, named once; lines that several roles share become a mixin set, or `build` fails |
| A class name built from pieces: `text-${tone}-600` | no CSS for the names that do not exist | `check` fails |
| A style nothing uses | not generated | `check` fails on a role nothing uses |

## The expected objection

Tailwind with a theme that holds only your own colours does much of this.

A theme with its own palette stops the wrong red: `--color-*: initial` removes the default palette, and your own colours take its place. The colour utilities and `border-l-4` are still there, and with them a stripe in the right colour on the left edge of a box. In Tidgrid the vocabulary has no such pair until a role makes it.

## Where Tailwind is the right choice

- A component library that arrives as Tailwind classes. With Tidgrid those looks are written in the vocabulary.
- A one-off page. Tidgrid pays back when the same look comes up again across pages.
