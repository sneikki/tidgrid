---
layout: default
title: Documentation
---

<section class="tg-pad-y(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(lg)">
    <h1 class="tg-text(page-title)">Documentation</h1>
    <p class="tg-text(body)">Start with the vocabulary file. Everything else is a closer look at one part of it, or at a command.</p>

    <div class="tg-layout(card-grid)">
      {% include doc-card.html href="/docs/vocabulary-file" title="The vocabulary file" description="One tidgrid.toml defines every class a page may use. A whole example, table by table." %}
      {% include doc-card.html href="/docs/tokens" title="Tokens" description="The only place for raw colours, lengths and times. Each becomes a custom property." %}
      {% include doc-card.html href="/docs/layout-atoms" title="Layout atoms" description="Rows, gaps, padding, widths and grids on the project's scales, per breakpoint." %}
      {% include doc-card.html href="/docs/roles" title="Roles" description="Colour, border, shadow and type, reachable only through named roles." %}
      {% include doc-card.html href="/docs/values" title="Values" description="What a value in the vocabulary file may and may not contain." %}
      {% include doc-card.html href="/docs/checker" title="The checker and how classes are found" description="What check reports, and what the scanner reads and does not read." %}
      {% include doc-card.html href="/docs/commands" title="Commands" description="init, build, check, list and sheet, with their flags and exit codes." %}
    </div>
  </div>
</section>
