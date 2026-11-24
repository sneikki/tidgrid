---
layout: doc-page
title: Files and the entry file
---

A vocabulary may span several files, the way a program does. Splitting it changes nothing in the rules, `list`, `check`, `plan` or the stylesheet: they all run on the one model the calls build.

## The entry file

Every command that reads a vocabulary takes the entry file as an argument, like a compiler. With no argument it reads `index.tg` in the current directory.

{% capture c %}
tidgrid build
tidgrid build web/ui.tg
tidgrid check web/ui.tg
tidgrid sheet -o look.html web/ui.tg
{% endcapture %}
{% include code.html code=c plain=true caption="Your terminal" %}

The entry file can be anywhere and named anything. `out`, `scan` and `allow_css` paths resolve against its directory. One folder may hold several vocabularies. There is no `-c` flag. `sheet` takes its `-o <file>` before the file.

A missing file exits 2:

{% capture c %}
nope.tg: no such file; run `tidgrid init` to create a vocabulary, or give the file as an argument
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

## `load` and `source`

A small project in three files. `index.tg`:

{% capture c %}
project(out = "src/tidgrid.css", scan = ["src"])
breakpoints(sm = "40em", md = "64em")
load("tokens.tg", "color", "space")
source("roles.tg")
element("body", css = {"background": color.ground, "color": color.ink})
allow_class("has-js", "set by an inline script before first paint")
{% endcapture %}
{% include code.html code=c caption="index.tg" %}

`tokens.tg`:

{% capture c %}
color = tokens("color", {"ground": "#0f1214", "ink": "#e6e9eb", "sand-ink": "#d8cfc0", "accent": "#5ed5dc"})
space = tokens("space", {"sm": "0.5rem", "md": "1rem"})
font = tokens("font", {"title": "Georgia, serif", "sans": "system-ui, sans-serif"})
sizes = tokens("sizes", {"h2": "1.5rem", "label": "0.875rem"})

def tint(c, pct):
    return "color-mix(in srgb, %s %d%%, transparent)" % (c, pct)
{% endcapture %}
{% include code.html code=c caption="tokens.tg" %}

`roles.tg` starts with a `load` of what it needs and then declares roles and mixins. It is shown on the [mixins](/docs/mixins) page.

- `load("tokens.tg", "color", "space")` is Starlark's own `load`: it runs the file and imports the named variables. It must name at least one.
- `source("roles.tg")` runs a file for what it declares and imports no names into the caller. The sourced file may itself `load` what it needs.
- A path is relative to the file that names it and must stay inside the entry file's directory: no absolute path, no `..`, no symlink out.
- A file runs once, however many times `load` and `source` name it. A cycle is an error.
- The [builtins](/docs/builtins) are available in every file without importing.

Things go into the stylesheet in the order the calls ran, so the order of `load` and `source` lines is the order of the output.

## When a path is wrong

{% capture c %}
index.tg:2:7: "../x.tg": path leaves the directory of the entry file
index.tg:2:7: nope.tg: cannot read inside the entry file's directory: openat nope.tg: no such file or directory
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

## Next

- [Builtins](/docs/builtins)
- [Mixins](/docs/mixins)
