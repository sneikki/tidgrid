---
layout: default
title: Get started
---

{% comment %}
What the panels show. The tool's output is its real output in an empty project. In the row
samples, (( )) marks what is fine (see _includes/lines.html).
{% endcomment %}
<!-- tidgrid-ignore-start -->
{% capture install_run %}
$ go install github.com/sneikki/tidgrid/cmd/tidgrid@latest
{% endcapture %}
{% capture init_run %}
$ tidgrid init
((tidgrid init: wrote index.tg))
{% endcapture %}
{% capture build_run %}
$ tidgrid build
((wrote src/styles/tidgrid.css))
{% endcapture %}
{% capture check_run %}
$ tidgrid check
roles used:
  none
allow:
  none
$ echo $?
((0))
{% endcapture %}
{% capture agent_rules %}
Style only with the classes `tidgrid list` prints. No <style>, no style= for appearance, no own CSS.
Run `tidgrid check` before every commit; it must exit 0.
A new role goes in the same commit as the page that first uses it; `tidgrid check` fails on a role nothing uses.
{% endcapture %}
<!-- tidgrid-ignore-end -->

<section class="tg-pad-t(xl)">
  <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(md)">
    <h1 class="tg-text(page-title)">Get started</h1>
    <p class="tg-text(lead) tg-max-width(reading)">Five steps from nothing to a project whose agent can only use the classes you approved.</p>
  </div>
</section>

{% comment %}The five steps (story.js): each keeps its heading in view beside a terminal that prints row by row.{% endcomment %}
<div id="story" class="tg-col tg-pad-b(xxl)">
  <div id="progress" class="tg-surface(progress)" aria-hidden="true"></div>

  <section class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col lg:tg-row lg:tg-align(start) tg-gap(xl) tg-pad-y(xl)">
    <div class="tg-layout(pinned) tg-grow tg-col tg-y-gap(md)">
      <h2 class="tg-text(section-title)">1. Install</h2>
      <p class="tg-text(body)">One Go binary. Nothing is added to your project's dependencies.</p>
    </div>
    <div class="tg-grow tg-col tg-y-gap(xs)">
      <p class="tg-text(caption)">Your terminal</p>
      <div class="tg-surface(code)">
{% include lines.html text=install_run %}
      </div>
    </div>
  </section>

  <section class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col lg:tg-row lg:tg-align(start) tg-gap(xl) tg-pad-y(xl)">
    <div class="tg-layout(pinned) tg-grow tg-col tg-y-gap(md)">
      <h2 class="tg-text(section-title)">2. Write a starter vocabulary</h2>
      <p class="tg-text(body)">Run it in the root of your project. It writes a starter <span class="tg-text(keyword)">index.tg</span> and refuses to overwrite one that exists.</p>
      <p class="tg-text(body)">Open the file and set <span class="tg-text(keyword)">out</span> in <span class="tg-text(keyword)">project()</span> to where the stylesheet should go, and <span class="tg-text(keyword)">scan</span> to the folders that hold your pages. Then put your own colours, spacing and roles in.</p>
      <a class="tg-link(text) tg-self(start)" href="{{ site.baseurl }}/docs/vocabulary-file">Every call of the vocabulary file</a>
    </div>
    <div class="tg-grow tg-col tg-y-gap(xs)">
      <p class="tg-text(caption)">Your terminal</p>
      <div class="tg-surface(code)">
{% include lines.html text=init_run %}
      </div>
    </div>
  </section>

  <section class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col lg:tg-row lg:tg-align(start) tg-gap(xl) tg-pad-y(xl)">
    <div class="tg-layout(pinned) tg-grow tg-col tg-y-gap(md)">
      <h2 class="tg-text(section-title)">3. Build the stylesheet</h2>
      <p class="tg-text(body)">This validates the vocabulary and writes one static stylesheet to <span class="tg-text(keyword)">out</span>. Link it from your pages like any other stylesheet.</p>
      <p class="tg-text(body)">Build again whenever you change <span class="tg-text(keyword)">index.tg</span> or a file it loads.</p>
    </div>
    <div class="tg-grow tg-col tg-y-gap(xs)">
      <p class="tg-text(caption)">Your terminal</p>
      <div class="tg-surface(code)">
{% include lines.html text=build_run %}
      </div>
    </div>
  </section>

  <section class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col lg:tg-row lg:tg-align(start) tg-gap(xl) tg-pad-y(xl)">
    <div class="tg-layout(pinned) tg-grow tg-col tg-y-gap(md)">
      <h2 class="tg-text(section-title)">4. Check the pages</h2>
      <p class="tg-text(body)">Each violation prints as <span class="tg-text(keyword)">path:line: message</span>. After them come the roles each file uses and the exceptions you allowed. A role nothing uses is a violation too.</p>
      <p class="tg-text(body)">Exit code 0 means fine, 1 means violations, 2 means the vocabulary file itself is invalid or missing.</p>
      <a class="tg-link(text) tg-self(start)" href="{{ site.baseurl }}/docs/checker">How the checker works</a>
    </div>
    <div class="tg-grow tg-col tg-y-gap(xs)">
      <p class="tg-text(caption)">Your terminal, before any page uses a role</p>
      <div class="tg-surface(code)">
{% include lines.html text=check_run %}
      </div>
    </div>
  </section>

  <section class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col lg:tg-row lg:tg-align(start) tg-gap(xl) tg-pad-y(xl)">
    <div class="tg-layout(pinned) tg-grow tg-col tg-y-gap(md)">
      <h2 class="tg-text(section-title)">5. Tell your agent</h2>
      <p class="tg-text(body)">Put these three lines in the project's agent instructions.</p>
      <p class="tg-text(body)"><span class="tg-text(keyword)">tidgrid list</span> prints the vocabulary as compact text. It is what the agent reads before UI work. To see the look yourself, run <span class="tg-text(keyword)">tidgrid sheet</span>: it writes one HTML page that renders every token and role.</p>
      <div class="tg-col sm:tg-row tg-wrap tg-gap(sm)">
        <a class="tg-action(primary)" href="{{ site.baseurl }}/docs">Read the documentation</a>
        <a class="tg-action(secondary)" href="{{ site.baseurl }}/docs/pitfalls">See the pitfalls</a>
      </div>
    </div>
    <div class="tg-grow tg-col tg-y-gap(xs)">
      <p class="tg-text(caption)">Your agent's instructions</p>
      <div class="tg-surface(code)">
{% include lines.html text=agent_rules %}
      </div>
    </div>
  </section>
</div>

<script type="module" src="{{ site.baseurl }}/assets/js/story.js?v={{ site.time | date: "%s" }}"></script>
