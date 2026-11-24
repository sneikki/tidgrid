---
layout: default
title: Home
---

{% comment %}The samples the story shows. They show classes, they do not use them.{% endcomment %}
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
{% capture page %}
<main class="tg-col tg-gap(md) tg-pad(md)">
  <h1 class="tg-text(title)">Hello</h1>
  <section class="tg-surface(card) tg-pad(sm)">
    One card.
  </section>
</main>
{% endcapture %}
{% capture listing %}
# roles
## surface
tg-surface(card)  A panel that lifts off the ground. One per group of related content.
## text
tg-text(title)  The one title of a page.
{% endcapture %}
{% capture slip %}
slip.html:1: unknown class tg-gap(2); tg-gap takes: sm md
slip.html:2: appearance in style attribute: color; use a role (`tidgrid list`)
slip.html:3: unknown class border-l-4: not in the vocabulary; run `tidgrid list`
slip.html:7: own CSS: <style> block; make it a role in the vocabulary
{% endcapture %}
<!-- tidgrid-ignore-end -->

{% comment %}The first screen. data-clear: backdrop.js lights grid cells only where these children leave room.{% endcomment %}
<section class="tg-min-height(hero) tg-col tg-justify(center) tg-pad-y(xl)">
  <div data-clear class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(lg)">
    <h1 class="tg-text(wordmark-hero) tg-self(start)">tidgrid<span class="tg-text(wordmark-accent)">.dev</span></h1>
    <p class="tg-text(tagline) tg-self(start)">A closed CSS vocabulary for coding agents.</p>
    <p class="tg-text(body) tg-max-width(reading)">One file per project defines every class a page may use. A build turns it into one static stylesheet. A check fails on anything outside it.</p>
    <div class="tg-col sm:tg-row sm:tg-self(start) tg-gap(sm)">
      <a class="tg-action(primary)" href="{{ site.baseurl }}/get-started">Get started</a>
      <a class="tg-action(secondary)" href="{{ site.baseurl }}/docs">Read the documentation</a>
    </div>
  </div>
</section>

{% comment %}
The story (story.js). Steps count from 0 in page order; a part of the drawing carries the
steps where it exists (data-show) and where it is the subject (data-hot). Insert a step and
every number after it moves.
  0 without it   1 one file   2 build   3 list   4 a page   5 atoms and roles
  6 check, a slip   7 check, the rest   8 your turn
{% endcomment %}
<div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col lg:tg-row lg:tg-align(start) tg-gap(xl)">
  <div id="progress" class="tg-surface(progress)" aria-hidden="true"></div>

  <aside id="stage" class="tg-surface(stage) tg-grow" aria-hidden="true">
    <p id="badge" class="tg-text(badge) tg-pad-b(xs)">Without it</p>
    <svg class="tg-drawing(canvas)" viewBox="0 0 400 316" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <marker id="head" viewBox="0 0 8 8" refX="7" refY="4" markerWidth="7" markerHeight="7" orient="auto"><path class="tg-drawing(fill)" d="M0,0 L8,4 L0,8 z"/></marker>
        <marker id="head-in" viewBox="0 0 8 8" refX="7" refY="4" markerWidth="7" markerHeight="7" orient="auto"><path class="tg-drawing(fill) tg-drawing(in)" d="M0,0 L8,4 L0,8 z"/></marker>
        <marker id="head-out" viewBox="0 0 8 8" refX="7" refY="4" markerWidth="7" markerHeight="7" orient="auto"><path class="tg-drawing(fill) tg-drawing(out)" d="M0,0 L8,4 L0,8 z"/></marker>
      </defs>

      <rect class="tg-drawing(box)" data-hot="0,3" x="20" y="14" width="130" height="46" rx="5"/>
      <text class="tg-drawing(title)" x="32" y="33">coding agent</text>
      <text class="tg-drawing(note) tg-drawing(part)" data-show="0-2" x="32" y="48">writes any CSS it likes</text>
      <text class="tg-drawing(note) tg-drawing(part)" data-show="3-8" x="32" y="48">reads the list first</text>

      <g class="tg-drawing(part)" data-show="1-8">
        <rect class="tg-drawing(box)" data-hot="1" x="250" y="14" width="130" height="46" rx="5"/>
        <text class="tg-drawing(title)" x="262" y="33">tidgrid.toml</text>
        <text class="tg-drawing(note)" x="262" y="48">every allowed class</text>
      </g>
      <path class="tg-drawing(wire) tg-drawing(in) tg-drawing(part)" data-show="3-8" d="M250,37 L156,37" marker-end="url(#head-in)"/>
      <text class="tg-drawing(mono) tg-drawing(in) tg-drawing(part)" data-show="3-8" x="203" y="30" text-anchor="middle">tidgrid list</text>

      <path class="tg-drawing(wire) tg-drawing(out) tg-drawing(part)" data-show="0-3" d="M85,60 L85,112" marker-end="url(#head-out)"/>
      <path class="tg-drawing(wire) tg-drawing(in) tg-drawing(part)" data-show="4-8" d="M85,60 L85,112" marker-end="url(#head-in)"/>
      <text class="tg-drawing(note)" x="93" y="91">writes</text>

      <rect class="tg-drawing(box)" data-hot="4-5" x="20" y="118" width="130" height="84" rx="5"/>
      <text class="tg-drawing(title)" x="32" y="137">a page</text>
      <g class="tg-drawing(chip) tg-drawing(out) tg-drawing(part)" data-show="0-3"><rect x="32" y="146" width="52" height="18" rx="3"/><text x="58" y="158.5" text-anchor="middle">.btn-new</text></g>
      <g class="tg-drawing(chip) tg-drawing(out) tg-drawing(part)" data-show="0-3"><rect x="90" y="146" width="50" height="18" rx="3"/><text x="115" y="158.5" text-anchor="middle">#ff0066</text></g>
      <g class="tg-drawing(chip) tg-drawing(out) tg-drawing(part)" data-show="0-3"><rect x="32" y="170" width="50" height="18" rx="3"/><text x="57" y="182.5" text-anchor="middle">&lt;style&gt;</text></g>
      <g class="tg-drawing(chip) tg-drawing(in) tg-drawing(part)" data-show="4-8" data-hot="4-5"><rect x="32" y="146" width="44" height="18" rx="3"/><text x="54" y="158.5" text-anchor="middle">tg-col</text></g>
      <g class="tg-drawing(chip) tg-drawing(in) tg-drawing(part)" data-show="4-5,8" data-hot="4-5"><rect x="32" y="170" width="98" height="18" rx="3"/><text x="81" y="182.5" text-anchor="middle">tg-surface(card)</text></g>
      <g class="tg-drawing(chip) tg-drawing(out) tg-drawing(part)" data-show="6-7"><rect x="32" y="170" width="66" height="18" rx="3"/><text x="65" y="182.5" text-anchor="middle">border-l-4</text></g>

      <g class="tg-drawing(part)" data-show="2-8">
        <rect class="tg-drawing(box)" data-hot="2" x="250" y="118" width="130" height="46" rx="5"/>
        <text class="tg-drawing(title)" x="262" y="137">site.css</text>
        <text class="tg-drawing(note)" x="262" y="152">one static stylesheet</text>
        <path class="tg-drawing(wire) tg-drawing(in)" d="M315,60 L315,112" marker-end="url(#head-in)"/>
        <text class="tg-drawing(mono) tg-drawing(in)" x="322" y="91">tidgrid build</text>
      </g>
      <g class="tg-drawing(part)" data-show="4-8">
        <path class="tg-drawing(wire)" d="M250,141 L156,141" marker-end="url(#head)"/>
        <text class="tg-drawing(note)" x="203" y="134" text-anchor="middle">styles it</text>
      </g>

      <g class="tg-drawing(part)" data-show="6-8">
        <rect class="tg-drawing(box)" data-hot="6-8" x="20" y="236" width="360" height="46" rx="5"/>
        <text class="tg-drawing(title)" x="32" y="256">tidgrid check</text>
        <text class="tg-drawing(note)" x="32" y="271">every class against the vocabulary</text>
        <path class="tg-drawing(wire)" d="M85,202 L85,230" marker-end="url(#head)"/>
        <path class="tg-drawing(wire)" d="M315,164 L315,230" marker-end="url(#head)"/>
        <text class="tg-drawing(note)" x="308" y="203" text-anchor="end">still what build writes?</text>
      </g>
      <text class="tg-drawing(mono) tg-drawing(out) tg-drawing(part)" data-show="6-7" x="368" y="263" text-anchor="end">exit 1</text>
      <text class="tg-drawing(mono) tg-drawing(in) tg-drawing(part)" data-show="8" x="368" y="263" text-anchor="end">exit 0</text>

      <rect class="tg-drawing(fill) tg-drawing(in)" x="20" y="299" width="8" height="8" rx="1"/>
      <text class="tg-drawing(note)" x="33" y="307">in the vocabulary</text>
      <rect class="tg-drawing(fill) tg-drawing(out)" x="130" y="299" width="8" height="8" rx="1"/>
      <text class="tg-drawing(note)" x="143" y="307">outside it: refused</text>
    </svg>
  </aside>

  <div class="tg-grow tg-col tg-pad-b(xxl)">

    <section data-step data-badge="Without it" class="tg-surface(step) tg-col tg-y-gap(md) tg-pad-y(xl) lg:tg-min-height(step)">
      <h2 class="tg-text(section-title)">What goes wrong without it?</h2>
      <p data-rise class="tg-motion(rise) tg-text(body)">Left alone, a coding agent writes whatever CSS comes to mind: a new class name here, a hex colour there, a <span class="tg-text(keyword)">&lt;style&gt;</span> block when it is in a hurry.</p>
      <p data-rise class="tg-motion(rise) tg-text(body)">Every page is a little different, and nobody reviews CSS line by line.</p>
    </section>

    <section data-step data-badge="One file" class="tg-surface(step) tg-col tg-y-gap(md) tg-pad-y(xl) lg:tg-min-height(step)">
      <h2 class="tg-text(section-title)">Where is the list of classes?</h2>
      <p data-rise class="tg-motion(rise) tg-text(body)">In one file per project, <span class="tg-text(keyword)">tidgrid.toml</span>. It defines every class a page may use. What the project looks like is decided once, in that file, and reviewed there.</p>
      <div data-rise class="tg-motion(rise)">
{% include code.html code=vocabulary %}
      </div>
    </section>

    <section data-step data-badge="Build" class="tg-surface(step) tg-col tg-y-gap(md) tg-pad-y(xl) lg:tg-min-height(step)">
      <h2 class="tg-text(section-title)">What turns the file into CSS?</h2>
      <p data-rise class="tg-motion(rise) tg-text(body)"><span class="tg-text(keyword)">tidgrid build</span> reads the file and writes one static stylesheet, the one named in <span class="tg-text(keyword)">out</span>: here <span class="tg-text(keyword)">site.css</span>.</p>
      <p data-rise class="tg-motion(rise) tg-text(body)">Nobody edits that stylesheet. When it differs from what <span class="tg-text(keyword)">build</span> would write, the check fails.</p>
    </section>

    <section data-step data-badge="List" class="tg-surface(step) tg-col tg-y-gap(md) tg-pad-y(xl) lg:tg-min-height(step)">
      <h2 class="tg-text(section-title)">How does the agent know what exists?</h2>
      <p data-rise class="tg-motion(rise) tg-text(body)">It runs <span class="tg-text(keyword)">tidgrid list</span> before any UI work and styles pages only with what it finds there. Every role comes with a sentence that says when to use it.</p>
      <div data-rise class="tg-motion(rise) tg-col tg-y-gap(sm)">
        <p class="tg-text(body)"><span class="tg-text(keyword)">tidgrid list</span>, the part about roles</p>
{% include code.html code=listing plain=true %}
      </div>
    </section>

    <section data-step data-badge="A page" class="tg-surface(step) tg-col tg-y-gap(md) tg-pad-y(xl) lg:tg-min-height(step)">
      <h2 class="tg-text(section-title)">What does a page look like?</h2>
      <div data-rise class="tg-motion(rise) tg-col tg-y-gap(sm)">
        <p class="tg-text(body)"><span class="tg-text(keyword)">index.html</span></p>
{% include code.html code=page %}
      </div>
      <p data-rise class="tg-motion(rise) tg-text(body)">Layout comes from atoms (<span class="tg-text(keyword)">tg-col</span>, <span class="tg-text(keyword)">tg-gap(md)</span>) that take their values from the space tokens. Colour and border come only from the roles. <span class="tg-text(keyword)">tidgrid check</span> passes on this page.</p>
    </section>

    <section data-step data-badge="Atoms and roles" class="tg-surface(step) tg-col tg-y-gap(md) tg-pad-y(xl) lg:tg-min-height(step)">
      <h2 class="tg-text(section-title)">What is free and what is fixed?</h2>
      <a data-rise class="tg-motion(rise) tg-link(card) tg-col" href="{{ site.baseurl }}/docs/layout-atoms">
        <div class="tg-surface(card) tg-col tg-gap(sm)">
          <h3 class="tg-text(label)">Layout is atoms</h3>
          <p class="tg-text(body)">Rows, gaps, padding and widths are free to combine. They always run on the project's token scales, so a gap is <span class="tg-text(keyword)">tg-gap(md)</span> and never an arbitrary number.</p>
        </div>
      </a>
      <a data-rise class="tg-motion(rise) tg-link(card) tg-col" href="{{ site.baseurl }}/docs/roles">
        <div class="tg-surface(card) tg-col tg-gap(sm)">
          <h3 class="tg-text(label)">Appearance is roles</h3>
          <p class="tg-text(body)">Colour, border, shadow and type are reachable only through named roles such as <span class="tg-text(keyword)">tg-surface(card)</span>. What a role looks like is decided once, in the vocabulary file. There is no border atom and no colour atom.</p>
        </div>
      </a>
    </section>

    <section data-step data-badge="Check" class="tg-surface(step) tg-col tg-y-gap(md) tg-pad-y(xl) lg:tg-min-height(step)">
      <h2 class="tg-text(section-title)">What happens on a slip?</h2>
      <div data-rise class="tg-motion(rise) tg-col tg-y-gap(sm)">
        <p class="tg-text(body)"><span class="tg-text(keyword)">tidgrid check</span>, after a slip</p>
{% include code.html code=slip plain=true %}
      </div>
      <p data-rise class="tg-motion(rise) tg-text(body)">An unknown atom value, an inline colour, a class from nowhere and a <span class="tg-text(keyword)">&lt;style&gt;</span> block each fail, and the exit code is 1.</p>
    </section>

    <section data-step data-badge="Check" class="tg-surface(step) tg-col tg-y-gap(md) tg-pad-y(xl) lg:tg-min-height(step)">
      <h2 class="tg-text(section-title)">What else does check catch?</h2>
      <div data-rise class="tg-motion(rise) tg-text(prose)" markdown="1">

- A class that is not in the vocabulary, such as `border-l-4`.
- A class name built from pieces, such as `is-${state}`.
- Appearance in a `style` attribute. Geometry is allowed.
- Own CSS: a style block, or a stylesheet that is not allowed by name.
- A stylesheet that differs from what `build` would write.
- An ignore comment that is never closed.

</div>
      <p data-rise class="tg-motion(rise) tg-text(body)">A message names the fix where there is one, for example <span class="tg-text(keyword)">src/pages/index.astro:14: unknown class tg-gap(2); tg-gap takes: sm md</span>.</p>
      <a data-rise class="tg-motion(rise) tg-link(text) tg-self(start)" href="{{ site.baseurl }}/docs/checker">How the checker works</a>
    </section>

    <section data-step data-badge="Your turn" class="tg-surface(step) tg-col tg-y-gap(md) tg-pad-y(xl) lg:tg-min-height(step)">
      <h2 class="tg-text(section-title)">How do I start?</h2>
      <p data-rise class="tg-motion(rise) tg-text(body)">Install the tool, write a starter vocabulary, and give your agent three lines of instructions.</p>
      <div data-rise class="tg-motion(rise) tg-col sm:tg-row tg-wrap tg-gap(sm)">
        <a class="tg-action(primary)" href="{{ site.baseurl }}/get-started">Get started</a>
        <a class="tg-action(secondary)" href="{{ site.baseurl }}/docs">Read the documentation</a>
        <a class="tg-action(secondary)" href="https://github.com/sneikki/tidgrid">Source on GitHub</a>
      </div>
    </section>

  </div>
</div>

{% comment %}A module: its names stay its own (theme.js and backdrop.js share the page's global names).{% endcomment %}
<script type="module" src="{{ site.baseurl }}/assets/js/story.js?v={{ site.time | date: "%s" }}"></script>
