---
layout: doc-page
title: Commands
---

Install with Go:

{% capture c %}
go install github.com/sneikki/tidgrid@latest
{% endcapture %}
{% include code.html code=c %}

## The five commands

- `tidgrid init` writes a starter `tidgrid.toml`. It refuses to overwrite one.
- `tidgrid build` validates the vocabulary and writes the stylesheet to `project.out`.
- `tidgrid check` scans the sources. It prints `path:line: message` per violation, then the usage summary. See [the checker](/docs/checker).
- `tidgrid list` prints the vocabulary as compact text. It is what an agent reads before UI work.
- `tidgrid sheet` writes one HTML page that renders every token and role. It is what a person reads to see the look.

## Flags

All commands take `-c <file>`, the vocabulary file. The default is `./tidgrid.toml`.

`sheet` also takes `-o <file>`, the page to write. The default is `tidgrid-sheet.html`.

## Exit codes

- 0: fine.
- 1: violations.
- 2: the vocabulary file itself is invalid or missing.

## In a project

Build after every change to the vocabulary file, and check before every commit. The three lines that tell an agent so are on the [get started](/get-started) page.

## Next

- [Get started](/get-started)
- [The vocabulary file](/docs/vocabulary-file)
