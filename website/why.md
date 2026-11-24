---
layout: default
title: Why Tidgrid
---

{% comment %}What the panels show. The text is the tool's real output (v0.3.0, rerun in an empty project).{% endcomment %}
<!-- tidgrid-ignore-start -->
{% capture unknown_run %}
src/page.html:1: unknown class tg-surface(rasied); tg-surface takes: raised sunken
{% endcapture %}
{% capture built_run %}
src/page.html:2: class name built from pieces: tg-text(${tone}); write every whole name out in a lookup table
{% endcapture %}
{% capture stripe_check %}
src/index.html:1: unknown class border-l-4: not in the vocabulary; run `tidgrid list`
src/index.html:1: appearance in style attribute: border-left; use a role (`tidgrid list`)
{% endcapture %}
{% capture stripe_build %}
index.tg:196: role.surface.stripe.use: required: when the role applies
index.tg:196: role.surface.stripe.border-left: raw colour #2456c4
{% endcapture %}
{% capture install_run %}
$ go install github.com/sneikki/tidgrid/cmd/tidgrid@latest
{% endcapture %}
<!-- tidgrid-ignore-end -->

<section class="tg-pad-t(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(md)">
    <h1 class="tg-text(page-title)">Why Tidgrid</h1>
    <p class="tg-text(lead) tg-max-width(reading)">What you get when your agent writes the UI. Every figure is measured on our own projects.</p>
  </div>
</section>

<section class="tg-pad-t(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(md)">
    <div class="tg-col md:tg-row tg-gap(md)">
      <figure class="tg-col tg-y-gap(xs) tg-grow">
        <a href="{{ site.baseurl }}/assets/img/landing-tailwind-full.png"><img class="tg-width(full) tg-surface(shot)" src="{{ site.baseurl }}/assets/img/landing-tailwind-top.png" width="1280" height="1500" alt="Top of a dark landing page for moonshine: a small navigation bar, a large white headline with a yellow button, a section with a heading on the left and text on the right, and the start of a grid of boxed items."></a>
        <figcaption class="tg-text(caption)">Styled with Tailwind CSS</figcaption>
      </figure>
      <figure class="tg-col tg-y-gap(xs) tg-grow">
        <a href="{{ site.baseurl }}/assets/img/landing-tidgrid-full.png"><img class="tg-width(full) tg-surface(shot)" src="{{ site.baseurl }}/assets/img/landing-tidgrid-top.png" width="1280" height="1500" alt="Top of a dark landing page for moonshine: a large monospace cream headline, two columns of intro text, and a dated list of one real case in four rows."></a>
        <figcaption class="tg-text(caption)">Written with the <a class="tg-link(text)" href="https://github.com/sneikki/tidgrid/tree/master/skills/antislop-ui">antislop-ui skill</a> and Tidgrid</figcaption>
      </figure>
    </div>
    <p class="tg-text(body) tg-max-width(reading)">The same request, "Build the landing page for moonshine. Everything the page may say is in content.md; the wordmark and favicon are in brand/. Plain HTML.", and the same content file, given to two agents on the same model. One was told to style it with Tailwind CSS. The other had the skill and Tidgrid: it first wrote a brief of the look and a vocabulary, a person approved that look, and only then did it build the page; <span class="tg-text(keyword)">tidgrid check</span> passes on it. One run each. The pictures show the top of each page; click one for the whole page.</p>
  </div>
</section>

<section class="tg-pad-y(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(xl)">
    <div class="tg-col tg-y-gap(xs)">
      <p class="tg-text(figure)">10 pages</p>
      <p class="tg-text(body) tg-max-width(reading)">A fresh agent wrote them from <span class="tg-text(keyword)">tidgrid list</span> alone. <span class="tg-text(keyword)">check</span> passed and no role was added.</p>
    </div>
    <div class="tg-col tg-y-gap(xs)">
      <p class="tg-text(figure)">24 groups, 7 mixin sets</p>
      <p class="tg-text(body) tg-max-width(reading)">This site's vocabulary had 24 groups of repeated lines in 57 roles; the largest was one title font and ink colour written by hand in 8 roles. The cleanup turned them into 7 mixin sets.</p>
    </div>
    <div class="tg-col tg-y-gap(xs)">
      <p class="tg-text(figure)">no difference</p>
      <p class="tg-text(body) tg-max-width(reading)">What <span class="tg-text(keyword)">tidgrid plan</span> printed against the live stylesheet after that cleanup.</p>
    </div>
    <div class="tg-col tg-y-gap(xs)">
      <p class="tg-text(figure)">0.1 s</p>
      <p class="tg-text(body) tg-max-width(reading)"><span class="tg-text(keyword)">tidgrid check</span> on this site, before every commit.</p>
    </div>
    <div class="tg-col tg-y-gap(xs)">
      <p class="tg-text(figure)">22 mistakes</p>
      <p class="tg-text(body) tg-max-width(reading)">Agents made them in real projects, and the pitfalls page lists each one.</p>
    </div>
  </div>
</section>

<section class="tg-pad-y(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(xl)">
    <h2 class="tg-text(section-title)">Next to Tailwind</h2>

    <div class="tg-col tg-y-gap(md)">
      <p class="tg-text(body)">In Tailwind a mistyped class, or one built from pieces, gets no CSS and no message from the build. Tidgrid exits 1.</p>
      {% include code.html code=unknown_run plain=true caption="A mistyped class" %}
      {% include code.html code=built_run plain=true caption="A class built from pieces" %}
    </div>

    <div class="tg-col tg-y-gap(md)">
      <p class="tg-text(body)">In Tailwind a coloured stripe on a box's left edge is two utility classes, <span class="tg-text(keyword)">border-l-4</span> and a border colour. Tidgrid stops all three ways to write it.</p>
      {% include code.html code=stripe_check plain=true caption="As a utility class or in style=: check" %}
      {% include code.html code=stripe_build plain=true caption="As a role of its own: build" %}
    </div>
  </div>
</section>

<section class="tg-pad-y(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(md)">
    {% include code.html code=install_run plain=true caption="Your terminal" %}
    <div class="tg-col sm:tg-row tg-wrap tg-gap(sm)">
      <a class="tg-action(primary)" href="{{ site.baseurl }}/get-started">Get started</a>
      <a class="tg-action(secondary)" href="{{ site.baseurl }}/docs/tailwind">The detailed comparison</a>
      <a class="tg-action(secondary)" href="{{ site.baseurl }}/docs/pitfalls">The pitfalls</a>
    </div>
  </div>
</section>
