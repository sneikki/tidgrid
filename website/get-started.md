---
layout: page
title: Get started
---

Five steps from nothing to a project whose agent can only use the classes you approved.

## 1. Install

{% capture c %}
go install github.com/sneikki/tidgrid/cmd/tidgrid@latest
{% endcapture %}
{% include code.html code=c plain=true %}

## 2. Write a starter vocabulary

Run this in the root of your project:

{% capture c %}
tidgrid init
{% endcapture %}
{% include code.html code=c plain=true %}

It writes a starter `tidgrid.toml` and refuses to overwrite one that exists. Open it and set `project.out` to where the stylesheet should go, and `project.scan` to the folders that hold your pages. Then put your own colours, spacing and roles in. The [vocabulary file](/docs/vocabulary-file) page explains every table.

## 3. Build the stylesheet

{% capture c %}
tidgrid build
{% endcapture %}
{% include code.html code=c plain=true %}

This validates the vocabulary and writes one static stylesheet to `project.out`. Link it from your pages like any other stylesheet. Build again whenever you change `tidgrid.toml`.

## 4. Check the pages

{% capture c %}
tidgrid check
{% endcapture %}
{% include code.html code=c plain=true %}

Each violation prints as `path:line: message`. After them come the roles each file uses, the roles nothing uses, and the exceptions you allowed. Exit code 0 means fine, 1 means violations, 2 means the vocabulary file itself is invalid or missing. See [the checker](/docs/checker).

## 5. Tell your agent

Put these three lines in the project's agent instructions:

<!-- tidgrid-ignore-start -->
{% capture c %}
Style only with the classes `tidgrid list` prints. No <style>, no style= for appearance, no own CSS.
Run `tidgrid check` before every commit; it must exit 0.
A new role is a change to tidgrid.toml in its own commit, never mixed into page work.
{% endcapture %}
{% include code.html code=c plain=true %}
<!-- tidgrid-ignore-end -->

`tidgrid list` prints the vocabulary as compact text. It is what the agent reads before UI work. To see the look yourself, run `tidgrid sheet`: it writes one HTML page that renders every token and role.

## Next

- [Documentation index](/docs)
- [The vocabulary file](/docs/vocabulary-file)
- [Commands](/docs/commands)
