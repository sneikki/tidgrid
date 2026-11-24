---
layout: default
title: Documentation
---

<section class="tg-pad-y(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(lg)">
    <h1 class="tg-text(page-title)">Documentation</h1>
    <p class="tg-text(body)">Start with the vocabulary file. Everything else is a closer look at one part of it, or at a command.</p>

    <a href="{{ site.baseurl }}/docs/vocabulary-file" class="tg-link(card) tg-col">
      <div class="tg-surface(card) tg-col tg-gap(sm)">
        <h2 class="tg-text(section-title)">The vocabulary file</h2>
        <p class="tg-text(body)">One tidgrid.toml defines every class a page may use. A whole example, table by table.</p>
        <div class="tg-row tg-wrap tg-gap(xs)">
          <span class="tg-text(keyword)">[project]</span>
          <span class="tg-text(keyword)">[breakpoints]</span>
          <span class="tg-text(keyword)">[tokens]</span>
          <span class="tg-text(keyword)">[element]</span>
          <span class="tg-text(keyword)">[role]</span>
          <span class="tg-text(keyword)">[keyframes]</span>
          <span class="tg-text(keyword)">[allow]</span>
        </div>
      </div>
    </a>

    <div class="tg-col">
      {% include doc-entry.html href="/docs/tokens" sample='ink = "#32495d"' title="Tokens" description="The only place for raw colours, lengths and times. Each becomes a custom property." %}
      {% include doc-entry.html href="/docs/layout-atoms" sample="tg-row md:tg-gap(md)" title="Layout atoms" description="Rows, gaps, padding, widths and grids on the project's scales, per breakpoint." %}
      {% include doc-entry.html href="/docs/roles" sample="tg-surface(card)" title="Roles" description="Colour, border, shadow and type, reachable only through named roles." %}
      {% include doc-entry.html href="/docs/values" sample="calc({space.md} * 2)" title="Values" description="What a value in the vocabulary file may and may not contain." %}
      {% include doc-entry.html href="/docs/checker" sample="index.html:3: unknown class" title="The checker and how classes are found" description="What check reports, and what the scanner reads and does not read." %}
      {% include doc-entry.html href="/docs/commands" sample="tidgrid build" title="Commands" description="init, build, check, list and sheet, with their flags and exit codes." %}
      {% include doc-entry.html href="/docs/pitfalls" sample="Not like this" title="Pitfalls" description="Mistakes agents made in real projects, each beside its fix. Half of them pass the check." %}
    </div>
  </div>
</section>
