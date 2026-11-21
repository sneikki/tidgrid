---
layout: default
title: Home
---

<section class="tg-pad-y(xxl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(lg)">
    <h1 class="tg-text(wordmark-hero)">tidgrid<span class="tg-text(wordmark-accent)">.dev</span></h1>
    <p class="tg-text(tagline)">A closed CSS vocabulary for coding agents.</p>
    <p class="tg-text(body)">One file per project defines every class a page may use. A build turns it into one static stylesheet. A check fails on anything outside it.</p>
    <div class="tg-col sm:tg-row tg-gap(sm)">
      <a class="tg-action(primary)" href="{{ site.baseurl }}/get-started">Get started</a>
      <a class="tg-action(secondary)" href="{{ site.baseurl }}/docs">Read the documentation</a>
    </div>
  </div>
</section>

<section class="tg-pad-y(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(md)">
    <h2 class="tg-text(section-title)">Why</h2>
    <p class="tg-text(body)">Left alone, a coding agent writes whatever CSS comes to mind: a new class name here, a hex colour there, a <span class="tg-text(keyword)">&lt;style&gt;</span> block when it is in a hurry. Every page is a little different, and nobody reviews CSS line by line.</p>
    <p class="tg-text(body)">Tidgrid closes the list. The agent reads the vocabulary with <span class="tg-text(keyword)">tidgrid list</span> and styles pages only with what it finds there. <span class="tg-text(keyword)">tidgrid check</span> fails when it does not. What the project looks like is decided once, in one file, and reviewed there.</p>
  </div>
</section>

<section class="tg-pad-y(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(lg)">
    <h2 class="tg-text(section-title)">A vocabulary and a page that uses it</h2>

    <div class="tg-col lg:tg-row tg-gap(lg)">
      <div class="tg-grow tg-col tg-y-gap(sm)">
        <p class="tg-text(label)">tidgrid.toml</p>
<!-- tidgrid-ignore-start -->
{% capture vocabulary %}
[project]
out  = "site.css"
scan = ["."]

[tokens.color]
ground = "#0f1214"
ink    = "#e6e9eb"
accent = "#5ed5dc"

[tokens.space]
sm = "0.5rem"
md = "1rem"

[tokens.border]
hair = "1px"

[element.body]
background = "{color.ground}"
color      = "{color.ink}"

[role.surface.card]
use    = "A panel that lifts off the ground. One per group of related content."
border = "{border.hair} solid {color.accent}"

[role.text.title]
use   = "The one title of a page."
color = "{color.accent}"
{% endcapture %}
{% include code.html code=vocabulary %}
<!-- tidgrid-ignore-end -->
      </div>

      <div class="tg-grow tg-col tg-y-gap(sm)">
        <p class="tg-text(label)">index.html</p>
<!-- tidgrid-ignore-start -->
{% capture page %}
<main class="tg-col tg-gap(md) tg-pad(md)">
  <h1 class="tg-text(title)">Hello</h1>
  <section class="tg-surface(card) tg-pad(sm)">
    One card.
  </section>
</main>
{% endcapture %}
{% include code.html code=page %}
<!-- tidgrid-ignore-end -->
        <p class="tg-text(body)">Layout comes from atoms (<span class="tg-text(keyword)">tg-col</span>, <span class="tg-text(keyword)">tg-gap(md)</span>) that take their values from the space tokens. Colour and border come only from the roles. <span class="tg-text(keyword)">tidgrid check</span> passes on this page.</p>
        <p class="tg-text(label)">tidgrid list, what the agent reads</p>
<!-- tidgrid-ignore-start -->
{% capture listing %}
tg-surface(card)  A lifted panel.
tg-text(title)    The one title of a page.
tg-gap(sm|md)     gap on the space tokens
{% endcapture %}
{% include code.html code=listing %}
<!-- tidgrid-ignore-end -->
        <p class="tg-text(label)">tidgrid check, after a slip</p>
<!-- tidgrid-ignore-start -->
{% capture slip %}
index.html:3: unknown class tg-gap(2)
index.html:4: style attribute sets a colour
index.html:9: own CSS in a style block
{% endcapture %}
{% include code.html code=slip %}
<!-- tidgrid-ignore-end -->
        <p class="tg-text(body)">An unknown atom value, an inline colour and a <span class="tg-text(keyword)">&lt;style&gt;</span> block each fail, and the exit code is 1.</p>
      </div>
    </div>
  </div>
</section>

<section class="tg-pad-y(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(lg)">
    <h2 class="tg-text(section-title)">Atoms for layout, roles for looks</h2>

    <div class="tg-col md:tg-row tg-gap(lg)">
      <div class="tg-grow tg-col tg-y-gap(sm)">
        <p class="tg-text(label)">Layout is atoms</p>
        <p class="tg-text(body)">Rows, gaps, padding and widths are free to combine. They always run on the project's token scales, so a gap is <span class="tg-text(keyword)">tg-gap(md)</span> and never an arbitrary number.</p>
        <a class="tg-link(text) tg-self(start)" href="{{ site.baseurl }}/docs/layout-atoms">Layout atoms</a>
      </div>

      <div class="tg-grow tg-col tg-y-gap(sm)">
        <p class="tg-text(label)">Appearance is roles</p>
        <p class="tg-text(body)">Colour, border, shadow and type are reachable only through named roles such as <span class="tg-text(keyword)">tg-surface(card)</span>. What a role looks like is decided once, in the vocabulary file. There is no border atom and no colour atom.</p>
        <a class="tg-link(text) tg-self(start)" href="{{ site.baseurl }}/docs/roles">Roles</a>
      </div>
    </div>
  </div>
</section>

<section class="tg-pad-y(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(md)">
    <h2 class="tg-text(section-title)">What check catches</h2>
    <div class="tg-text(prose)" markdown="1">

- A class that is not in the vocabulary, such as `border-l-4`.
- A class name built from pieces, such as `is-${state}`.
- Appearance in a `style` attribute. Geometry is allowed.
- Own CSS: a style block, or a stylesheet that is not allowed by name.
- A stylesheet that differs from what `build` would write.
- An ignore comment that is never closed.

</div>
    <p class="tg-text(body)">A message names the fix where there is one, for example <span class="tg-text(keyword)">src/pages/index.astro:14: unknown class tg-gap(2); tg-gap takes: sm md</span>.</p>
    <a class="tg-link(text) tg-self(start)" href="{{ site.baseurl }}/docs/checker">How the checker works</a>
  </div>
</section>

<section class="tg-pad-y(xxl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(lg)">
    <h2 class="tg-text(section-title)">Try it</h2>
    <p class="tg-text(body)">Install the tool, write a starter vocabulary, and give your agent three lines of instructions.</p>
    <div class="tg-col sm:tg-row tg-gap(sm)">
      <a class="tg-action(primary)" href="{{ site.baseurl }}/get-started">Get started</a>
      <a class="tg-action(secondary)" href="{{ site.baseurl }}/docs">Read the documentation</a>
      <a class="tg-action(secondary)" href="https://github.com/sneikki/tidgrid">Source on GitHub</a>
    </div>
  </div>
</section>
