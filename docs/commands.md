---
layout: doc-page
title: Commands
---

Install with Go (needs Go 1.25):

{% capture c %}
go install github.com/sneikki/tidgrid/cmd/tidgrid@latest
{% endcapture %}
{% include code.html code=c plain=true caption="Your terminal" %}

## The seven commands

| Command | Does |
|---|---|
| `tidgrid init [file]` | Writes a starter vocabulary (default `index.tg`). Refuses to overwrite one. |
| `tidgrid build [file]` | Validates the vocabulary and writes the stylesheet to the `out` of `project()`. |
| `tidgrid check [file]` | Scans the sources. Prints `path:line: message` per violation, then the usage summary. See [the checker](/docs/checker). |
| `tidgrid list [file]` | Prints the vocabulary as compact text: what an agent reads before UI work. |
| `tidgrid plan [file]` | Prints what `build` would change in the stylesheet on disk. Writes nothing. |
| `tidgrid sheet [-o out] [file]` | Writes one HTML page that renders every token and role: what a person reads to see the look. The default page is `tidgrid-sheet.html`. |
| `tidgrid convert [tidgrid.toml]` | Writes a `.tg` vocabulary from an old TOML one. See below. |

## The file argument

The commands that read a vocabulary take the entry file as an argument, like a compiler. With no argument they read `index.tg` in the current directory. There is no `-c` flag. See [files](/docs/files).

## init

`tidgrid init` writes a starter vocabulary with tokens, two element rules and a set of roles. The roles are written commented out, so a fresh project has no unused role and `check` passes. Uncomment a role when a page uses it.

{% capture c %}
# Roles. Uncomment a role when a page uses it: `tidgrid check` fails on a role nothing uses.

# Surfaces: a panel that sits on the page. Raised is the normal panel, sunken is a well.
# role("surface", "raised",
#     use = "A panel holding one group of related content. Use for most boxes; at most a few per view.",
#     css = {
#         "background": color.raised,
#         "border-radius": radius.sm,
#     })
{% endcapture %}
{% include code.html code=c caption="Part of the index.tg init writes" %}

With an `index.tg` already there, `init` exits 2 and writes nothing:

{% capture c %}
tidgrid init: index.tg exists; not overwritten
{% endcapture %}
{% include code.html code=c plain=true caption="What init prints" %}

## plan

`tidgrid plan` renders the vocabulary in memory and compares it with the stylesheet at the output path: classes added and removed, the declarations that changed in a class, changed tokens with the classes that reach them, element rules and keyframes. Pages that use a removed class are not listed; `check` reports those as unknown classes once the stylesheet is built. It exits 0 whether or not there is a difference.

{% capture c %}
changed token color.ink: #e6e9eb -> #eeeeee
  reaches 3 classes: tg-text(section-title) tg-text(caption) tg-text(lead)
{% endcapture %}
{% include code.html code=c plain=true caption="What plan prints after a token changed" %}

With no stylesheet yet, everything is listed as added. With no difference it prints `no difference`.

## Exit codes

| Code | Means |
|---|---|
| 0 | Fine. |
| 1 | Violations (`check`). |
| 2 | The vocabulary itself is invalid or missing (`check` too), or the command was misused. |

## Converting an old project

The vocabulary used to be a `tidgrid.toml`. That form is still read when the file is named (`tidgrid check tidgrid.toml`), with its own rules unchanged, except that `check` now also fails on a role nothing uses. It has no mixins, functions or loops, and a file with the `.toml` extension is the only thing read as TOML.

`tidgrid convert` writes `index.tg`, `tokens.tg` and `roles.tg` beside the TOML file and refuses when any of the three exists. The converted project builds the same stylesheet byte for byte. Values stay strings as written, and TOML comments are not carried over. It writes `repeats = []` in `project`, because an old project repeats lines freely. See [repeated lines](/docs/mixins#repeated-lines).

{% capture c %}
wrote index.tg
wrote roles.tg
wrote tokens.tg
note: TOML comments are not carried over
{% endcapture %}
{% include code.html code=c plain=true caption="What convert prints" %}

In a directory that has a `tidgrid.toml` and no `index.tg`, `build`, `check`, `list`, `plan` and `sheet` exit 2 and say so. `init` and `convert` do not.

{% capture c %}
index.tg not found, but tidgrid.toml is here: run `tidgrid convert` to write index.tg
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

## In a project

Build after every change to the vocabulary, and check before every commit. The three lines that tell an agent so are on the [get started](/get-started) page.

## Next

- [Get started](/get-started)
- [The vocabulary file](/docs/vocabulary-file)
