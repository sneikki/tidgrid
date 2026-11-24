---
layout: doc-page
title: Commands
---

Install with Go:

{% capture c %}
go install github.com/sneikki/tidgrid/cmd/tidgrid@latest
{% endcapture %}
{% include code.html code=c plain=true caption="Your terminal" %}

## The five commands

| Command | Does |
|---|---|
| `tidgrid init` | Writes a starter `tidgrid.toml`. Refuses to overwrite one. |
| `tidgrid build` | Validates the vocabulary and writes the stylesheet to `project.out`. |
| `tidgrid check` | Scans the sources. Prints `path:line: message` per violation, then the usage summary. See [the checker](/docs/checker). |
| `tidgrid list` | Prints the vocabulary as compact text: what an agent reads before UI work. |
| `tidgrid sheet` | Writes one HTML page that renders every token and role: what a person reads to see the look. |

## Flags

All commands take `-c <file>`, the vocabulary file. The default is `./tidgrid.toml`.

`sheet` also takes `-o <file>`, the page to write. The default is `tidgrid-sheet.html`.

## Exit codes

| Code | Means |
|---|---|
| 0 | Fine. |
| 1 | Violations. |
| 2 | The vocabulary file itself is invalid or missing. |

## In a project

Build after every change to the vocabulary file, and check before every commit. The three lines that tell an agent so are on the [get started](/get-started) page.

## Next

- [Get started](/get-started)
- [The vocabulary file](/docs/vocabulary-file)
