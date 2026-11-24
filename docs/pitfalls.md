---
layout: doc-page
title: Pitfalls
---

These came up when coding agents styled real projects with Tidgrid. The first group passes `tidgrid check`, so only a reader catches it. The second group the tool stops; each is shown with the message it prints.

{% comment %}Every sample on this page shows classes, it does not use them.{% endcomment %}
<!-- tidgrid-ignore-start -->

## What check cannot see

### Layout written as a role

If every declaration in a role is direction, gap, padding or size, it is layout, and the atoms already do it. A role is right for layout only where atoms cannot reach: an overlay, a fixed panel, grid placement.

{% capture bad %}
role("layout", "page",
    use = "The centred column of a page.",
    css = {
        "margin-inline": "auto",
        "width": size.full,
        "max-width": size.measure,
    })
{% endcapture %}
{% capture good %}
<main class="tg-center tg-width(full) tg-max-width(measure)">
  …
</main>
{% endcapture %}
{% include versus.html bad=bad good=good bad_in="in index.tg" good_in="in the page" %}

### A role per element

A role that one element uses once is CSS with a longer name. In one real project about 15 of 22 roles were used once. Name a role by what it is for, not by where it sits, and read the `use` lines in `tidgrid list` before adding one. A few roles are single by nature, such as the footer band; a heading is not.

{% capture bad %}
role("text", "pricing-heading",
    use = "The heading of the pricing page.",
    css = {"font-family": font.title, "font-size": type.h2})

role("text", "careers-heading",
    use = "The heading of the careers page.",
    css = {"font-family": font.title, "font-size": type.h2})
{% endcapture %}
{% capture good %}
role("text", "section-title",
    use = "Heading of a section inside a page.",
    css = {"font-family": font.title, "font-size": type.h2})
{% endcapture %}
{% include versus.html bad=bad good=good bad_in="in index.tg" good_in="in index.tg" %}

### A colour passed through a custom property

One role sets a private custom property and another reads it. The build accepts it, and nothing checks that the property is set where it is read: it is a channel beside the tokens. Write one role per variant.

{% capture bad %}
role("kind", "note",
    use = "A note.",
    css = {"--k": color.accent})

role("surface", "badge",
    use = "A badge.",
    css = {"border": "%s solid var(--k)" % border.hair})
{% endcapture %}
{% capture good %}
role("surface", "badge-note",
    use = "The badge of a note.",
    css = {"border": "%s solid %s" % (border.hair, color.accent)})
{% endcapture %}
{% include versus.html bad=bad good=good bad_in="in index.tg" good_in="in index.tg" %}

### An atom and a role on the same property

Atoms are written after roles and win. Here the role hides the panel until a script opens it, and `tg-col` sets `display` again: the panel is visible from the start. Give the atom its own element.

{% capture bad %}
<!-- surface(panel) sets display: none -->
<nav class="tg-surface(panel) tg-col">
  …
</nav>
{% endcapture %}
{% capture good %}
<nav class="tg-surface(panel)">
  <div class="tg-col">
    …
  </div>
</nav>
{% endcapture %}
{% include versus.html bad=bad good=good bad_in="in the page" good_in="in the page" %}

### A class kept in front matter

Front matter is not read. The class reaches the page through an interpolation, which `check` skips, so the typo below passes. Write the name in the template, inside a template tag: the names around a tag are read.

{% capture bad %}
---
main_class: tg-justify(centre)
---
<main class="tg-col {% raw %}{{ page.main_class }}{% endraw %}">
{% endcapture %}
{% capture good %}
---
centred: true
---
<main class="tg-col{% raw %}{% if page.centred %}{% endraw %} tg-justify(center){% raw %}{% endif %}{% endraw %}">
{% endcapture %}
{% include versus.html bad=bad good=good bad_in="in a Jekyll page" good_in="in a Jekyll page" %}

### A stylesheet cached from before the deploy

The page is fetched fresh and the stylesheet comes from the browser's cache. A class added in this deploy has no rule there, and the layout breaks for returning readers until the cache expires. Put something that changes with every build in the address. `check` accepts a query on the link.

{% capture bad %}
<link rel="stylesheet" href="/site.css">
{% endcapture %}
{% capture good %}
<link rel="stylesheet" href="/site.css?v={% raw %}{{ site.time | date: "%s" }}{% endraw %}">
{% endcapture %}
{% include versus.html bad=bad good=good bad_in="in the layout" good_in="in the layout" %}

### Reduced motion on every element

A rule on `*` has no specificity, so the `transition` of any role wins over it and the motion stays. Put the at-rule inside the role that moves.

{% capture bad %}
element("*", css = {
    "@media (prefers-reduced-motion: reduce)": {
        "transition-duration": motion.none,
    },
})
{% endcapture %}
{% capture good %}
role("motion", "rise",
    use = "Fades a block in.",
    css = {
        "@media (prefers-reduced-motion: reduce)": {"transition": "none"},
    })
{% endcapture %}
{% include versus.html bad=bad good=good bad_in="in index.tg" good_in="in index.tg" %}

### An exception to silence the checker

`allow_class` is for a class a third-party script needs, with the reason. Used for a class of your own, it makes `check` pass and leaves the class without a rule and without a review. The same goes for ignore comments around real markup: they are for code samples.

{% capture bad %}
allow_class("btn-new", "used on the careers page")
{% endcapture %}
{% capture good %}
role("action", "apply",
    use = "The button that starts a job application.",
    css = {"color": color.accent})
{% endcapture %}
{% include versus.html bad=bad good=good bad_in="in index.tg" good_in="in index.tg" %}

## What the tool stops

### A class name built from pieces

A name put together at run time cannot be checked, so it is a violation. Write every whole name out, each through a helper.

{% capture bad %}
<div className={`tg-kind(${kind})`}>…</div>
{% endcapture %}
{% capture good %}
const kinds = {
  note: cn("tg-kind(note)"),
  warning: cn("tg-kind(warning)"),
};

<div className={kinds[kind]}>…</div>
{% endcapture %}
{% capture says %}
$ tidgrid check
[[src/Card.jsx:1:]] class name built from pieces: tg-kind(${kind}); write every whole name out in a lookup table
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in Card.jsx" good_in="in Card.jsx" %}

### Class names in a lookup table

Whole names in a table pass the rule against built names, but the table is a variable and `check` does not read it. The typo below is not named: the role is only reported as unused, and the obvious fix is to remove the role. Pass each name through a helper (`cn`, `clsx`, `cva`, `classNames`, `twMerge`, or the ones in `helpers` of `project()`) and the same typo fails as an unknown class.

{% capture bad %}
const kinds = {
  primary: "tg-action(primery)",
};
{% endcapture %}
{% capture good %}
const kinds = {
  primary: cn("tg-action(primary)"),
};
{% endcapture %}
{% capture says %}
$ tidgrid check
[[index.tg:9:]] role tg-action(primary) is not used in the scanned sources: use it, remove it, or, if a script or lookup table sets it, write the name whole inside a helper call such as cn()
…
$ echo $?
[[1]]
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in kinds.js" good_in="in kinds.js" %}

### A raw value in the vocabulary

A colour, a length, a time or a size in per cent belongs in `tokens`. A lone size, such as a 5 px dot, needs a token too. `build` exits 2 and writes nothing.

{% capture bad %}
role("surface", "badge",
    use = "A badge.",
    css = {
        "padding": "12px",
        "color": "#ff0066",
        "width": "40%",
    })
{% endcapture %}
{% capture good %}
size = tokens("size", {"badge": "40%"})

role("surface", "badge",
    use = "A badge.",
    css = {
        "padding": space.sm,
        "color": color.accent,
        "width": size.badge,
    })
{% endcapture %}
{% capture says %}
$ tidgrid build
[[index.tg:9:]] role.surface.badge.padding: raw length or time 12px (move it to tokens)
[[index.tg:9:]] role.surface.badge.color: raw colour #ff0066
[[index.tg:9:]] role.surface.badge.width: raw size 40% (move it to tokens)
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in index.tg" good_in="in index.tg" %}

### An atom with a value of its own

Atoms take token names, never values.

{% capture bad %}
<div class="tg-row tg-gap(12px)">
{% endcapture %}
{% capture good %}
<div class="tg-row tg-gap(md)">
{% endcapture %}
{% capture says %}
$ tidgrid check
[[src/index.html:1:]] unknown class tg-gap(12px); tg-gap takes: sm md
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in index.html" good_in="in index.html" %}

### A role name remembered from another project

The atoms are the same in every project. The role names are not. Run `tidgrid list` and read it before writing a class.

{% capture bad %}
<section class="tg-surface(panel)">
{% endcapture %}
{% capture good %}
<section class="tg-surface(card)">
{% endcapture %}
{% capture says %}
$ tidgrid check
[[src/index.html:1:]] unknown class tg-surface(panel); tg-surface takes: card
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in index.html" good_in="in index.html" %}

### A state class without a state role

A class that a script sets is a class like any other. Give it a role of kind `state` with only a `use` line, let another role react to it, and have the script set `tg-state(open)`. An ARIA or data attribute works the same way and needs no role: `"&[aria-expanded=true]"`.

{% capture bad %}
menu.classList.add("is-open");
{% endcapture %}
{% capture good %}
role("state", "open",
    use = "Set by a script on the menu while it is open.")

role("layout", "menu",
    use = "The menu.",
    css = {"&.tg-state(open)": {"display": "block"}})
{% endcapture %}
{% capture says %}
$ tidgrid check
[[src/menu.js:1:]] unknown class is-open: not in the vocabulary; run `tidgrid list`
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in menu.js" good_in="in index.tg" %}

### A state role named only in a selector

The selector `"&.tg-state(open)"` styles the class; it does not set it. A state role counts as used only when a script writes the name whole, as a string the scanner reads. Until then `check` fails on the role.

{% capture bad %}
menu.classList.toggle(open);
{% endcapture %}
{% capture good %}
menu.classList.toggle("tg-state(open)");
{% endcapture %}
{% capture says %}
$ tidgrid check
[[index.tg:9:]] role tg-state(open) is not used in the scanned sources: use it, remove it, or, if a script or lookup table sets it, write the name whole inside a helper call such as cn()
…
$ echo $?
[[1]]
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in menu.js" good_in="in menu.js" %}

### A role in a commit of its own

`check` fails on a role nothing uses, so a commit that adds only the role is red until the page arrives. Put a new role in the same commit as the first page that uses it. A change to an existing role or to a token can still be a commit of its own. When a role replaces others, the commit that moves the last page off an old role also removes it: `check` exits 0 at every commit.

{% capture bad %}
git add index.tg site.css
git commit -m "Add role surface(job)"
git add careers.html
git commit -m "Add careers page"
{% endcapture %}
{% capture good %}
git add index.tg site.css careers.html
git commit -m "Add careers page"
{% endcapture %}
{% capture says %}
$ tidgrid check
[[index.tg:9:]] role tg-surface(job) is not used in the scanned sources: use it, remove it, or, if a script or lookup table sets it, write the name whole inside a helper call such as cn()
…
$ echo $?
[[1]]
{% endcapture %}
{% include versus.html bad=bad good=good says=says plain=true %}

### The same lines in several roles

Roles that repeat the same lines are one mixin set not yet written. `build` fails when 3 lines appear in 2 roles, or 2 lines in 4 roles, and says how to clear it. Lines that come from an included set do not count.

{% capture bad %}
role("action", "a",
    use = "The main button.",
    css = {"cursor": "pointer", "padding": "0", "border": "none", "color": color.accent})

role("action", "b",
    use = "A quiet button.",
    css = {"cursor": "pointer", "padding": "0", "border": "none", "color": color.ground})
{% endcapture %}
{% capture good %}
shape = mixin("shape",
    use = "Shared shapes of controls.",
    sets = {"button": {"cursor": "pointer", "padding": "0", "border": "none"}})

role("action", "a",
    use = "The main button.",
    include = [shape.button],
    css = {"color": color.accent})

role("action", "b",
    use = "A quiet button.",
    include = [shape.button],
    css = {"color": color.ground})
{% endcapture %}
{% capture says %}
$ tidgrid build
[[index.tg:9:]] role.action.a: own lines also in role action.b: cursor: pointer; padding: 0; border: none; make them a mixin set and include it
index.tg: 1 group of repeated lines; fix the first and build again: one mixin set often clears several
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in index.tg" good_in="in index.tg" %}

### A mixin set that is not shared

A set needs at least two lines and at least two roles that include it. A set of one line is a token in disguise. A set that one role includes is that role's own lines under another name. Only roles that exist count, and an unused role fails `check`, so a role made up to give a set its second user does not last.

{% capture bad %}
shape = mixin("shape",
    use = "Shared shapes of controls.",
    sets = {"button": {"cursor": "pointer", "padding": "0"}})

role("action", "a",
    use = "The main button.",
    include = [shape.button],
    css = {"color": color.accent})

role("action", "b",
    use = "A quiet button.",
    css = {"color": color.ground})
{% endcapture %}
{% capture good %}
role("action", "a",
    use = "The main button.",
    css = {"cursor": "pointer", "padding": "0", "color": color.accent})

role("action", "b",
    use = "A quiet button.",
    css = {"color": color.ground})
{% endcapture %}
{% capture says %}
$ tidgrid build
[[index.tg:9:]] mixin.shape.button: included only by role action.a; a set needs at least two
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in index.tg" good_in="in index.tg" %}

With one line the message is `mixin.shape.pointer: a set needs at least two lines; a single shared value belongs in a token`.

### Two groups of repeated lines that share a line

Roles a and b repeat three lines, and roles c and d repeat three, one of them the same line as in the first group. Put the shared line in one set only. A role that includes two sets that both give it fails with `given twice: by set shape.button and by set shape.bare`.

{% capture bad %}
role("action", "a",
    use = "The main button.",
    css = {"cursor": "pointer", "padding": "0", "border": "none", "color": color.accent})

role("action", "b",
    use = "A quiet button.",
    css = {"cursor": "pointer", "padding": "0", "border": "none", "color": color.ground})

role("text", "c",
    use = "A link.",
    css = {"border": "none", "outline": "none", "margin": "0", "color": color.accent})

role("text", "d",
    use = "A quiet link.",
    css = {"border": "none", "outline": "none", "margin": "0", "color": color.ground})
{% endcapture %}
{% capture good %}
shape = mixin("shape",
    use = "Shared shapes of controls.",
    sets = {
        "button": {"cursor": "pointer", "padding": "0", "border": "none"},
        "bare": {"outline": "none", "margin": "0"},
    })

role("action", "a",
    use = "The main button.",
    include = [shape.button],
    css = {"color": color.accent})

role("action", "b",
    use = "A quiet button.",
    include = [shape.button],
    css = {"color": color.ground})

role("text", "c",
    use = "A link.",
    include = [shape.bare],
    css = {"border": "none", "color": color.accent})

role("text", "d",
    use = "A quiet link.",
    include = [shape.bare],
    css = {"border": "none", "color": color.ground})
{% endcapture %}
{% capture says %}
$ tidgrid build
[[index.tg:9:]] role.action.a: own lines also in role action.b: cursor: pointer; padding: 0; border: none; make them a mixin set and include it
[[index.tg:17:]] role.text.c: own lines also in role text.d: border: none; outline: none; margin: 0; make them a mixin set and include it
index.tg: 2 groups of repeated lines; fix the first and build again: one mixin set often clears several
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in index.tg" good_in="in index.tg" %}

### The generated stylesheet edited by hand

The stylesheet is output. Change the vocabulary, run `tidgrid build`, commit both.

{% capture bad %}
.tg-action\(primary\) {
  color: red;
}
{% endcapture %}
{% capture good %}
role("action", "primary",
    use = "The main button.",
    css = {"color": color.accent})
{% endcapture %}
{% capture says %}
$ tidgrid check
[[src/site.css:1:]] stylesheet differs from the vocabulary; run `tidgrid build`
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in site.css" good_in="in index.tg" %}

### A role without a use line

`use` is what an agent reads to decide whether the role applies. Without it the role cannot be picked, so the build refuses it.

{% capture bad %}
role("surface", "badge",
    css = {"padding": space.sm})
{% endcapture %}
{% capture good %}
role("surface", "badge",
    use = "A small label on a card: a status, a count.",
    css = {"padding": space.sm})
{% endcapture %}
{% capture says %}
$ tidgrid build
[[index.tg:9:]] role.surface.badge.use: required: when the role applies
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in index.tg" good_in="in index.tg" %}

### A code sample that shows classes

The scanner reads a class attribute wherever it finds one, also escaped inside a `pre`. Put the sample between the two ignore comments. A region that is never closed is reported as well.

{% capture bad %}
<pre>&lt;a class="btn-next"&gt;Apply&lt;/a&gt;</pre>
{% endcapture %}
{% capture good %}
<!-- tidgrid-ignore-start -->
<pre>&lt;a class="btn-next"&gt;Apply&lt;/a&gt;</pre>
<!-- tidgrid-ignore-end -->
{% endcapture %}
<!-- tidgrid-ignore-start -->
{% capture says %}
$ tidgrid check
[[src/guide.html:1:]] unknown class btn-next: not in the vocabulary; run `tidgrid list`
{% endcapture %}
{% include versus.html bad=bad good=good says=says bad_in="in guide.html" good_in="in guide.html" %}

<!-- tidgrid-ignore-end -->

## Next

- [The checker and how classes are found](/docs/checker)
- [Roles](/docs/roles): when to add one
