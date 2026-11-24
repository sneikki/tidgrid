---
layout: default
title: Home
---

{% comment %}
The samples the scenes show. They show classes, they do not use them. In the row samples,
[[ ]] marks what is outside the vocabulary and (( )) what is fine (see _includes/lines.html).
The tool's output is its real output for these files.
{% endcomment %}
<!-- tidgrid-ignore-start -->
{% capture agent_page %}
<main>
  <h1>Careers</h1>
  <p [[style="color: #ff0066"]]>Join us.</p>
  <a class="[[btn-new]]">Apply</a>
</main>
[[<style>.btn-new { color: red }</style>]]
{% endcapture %}
{% capture live_page %}
<p [[style="color: #ff0066"]]>Join us.</p>
<a class="[[btn-new]]">Apply</a>
[[<style>.btn-new { color: red }</style>]]
{% endcapture %}
{% capture slip_page %}
1 <main class="tg-col [[tg-gap(2)]]">
2   <h1 class="tg-text(title)">Careers</h1>
3   <p [[style="color: #ff0066"]]>Join us.</p>
4   <a class="[[btn-new]]">Apply</a>
5 </main>
6 [[<style>.btn-new { color: red }</style>]]
{% endcapture %}
{% capture slip_check %}
$ tidgrid check
[[careers.html:1:]] unknown class tg-gap(2); tg-gap takes: sm md
[[careers.html:3:]] appearance in style attribute: color; use a role (`tidgrid list`)
[[careers.html:4:]] unknown class btn-new: not in the vocabulary; run `tidgrid list`
[[careers.html:6:]] own CSS: <style> block; make it a role in the vocabulary
…
$ echo $?
[[1]]
{% endcapture %}
{% capture tg_tokens %}
color = tokens("color", {"ground": "#0f1214", "ink": "#e6e9eb", "accent": "#5ed5dc"})
border = tokens("border", {"hair": "1px"})
{% endcapture %}
{% capture tg_card %}
role("surface", "card",
    use = "A panel that lifts off the ground. One per group of related content.",
    css = {"border": "%s solid %s" % (border.hair, color.accent)})
{% endcapture %}
{% capture tg_title %}
role("text", "title",
    use = "The one title of a page.",
    css = {"color": color.accent})
{% endcapture %}
{% capture list_run %}
$ tidgrid list
…
# roles
## surface
tg-surface(card)  A panel that lifts off the ground. One per group of related content.
## text
tg-text(title)  The one title of a page.
{% endcapture %}
{% capture start_run %}
$ go install github.com/sneikki/tidgrid/cmd/tidgrid@latest
$ tidgrid init
$ tidgrid build
$ tidgrid check
{% endcapture %}
<!-- tidgrid-ignore-end -->

{% comment %}The first screen. data-clear: backdrop.js lights grid cells only where these children leave room.{% endcomment %}
<section class="tg-min-height(hero) tg-col tg-justify(center) tg-pad-y(xl)">
  <div data-clear class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(lg)">
    <h1 class="tg-text(wordmark-hero) tg-self(start)">tidgrid<span class="tg-text(wordmark-accent)">.dev</span></h1>
    <p class="tg-text(tagline) tg-self(start)">Your agent has opinions about CSS. Revoke them.</p>
    <p class="tg-text(lead) tg-max-width(reading)">One vocabulary lists every class a page may use. Everything else fails the check.</p>
    <div class="tg-col sm:tg-row sm:tg-self(start) tg-gap(sm)">
      <a class="tg-action(primary)" href="{{ site.baseurl }}/get-started">Get started</a>
      <a class="tg-action(secondary)" href="{{ site.baseurl }}/docs">Read the documentation</a>
    </div>
  </div>
</section>

{% comment %}
The pitch (story.js): a few scenes, each one big statement with a picture of its own. It
sells; how the tool works is in the documentation. Every data-rise piece comes in once as
it scrolls into view. Two colours mean something: text(passed) = in the vocabulary,
text(refused) = outside it.
{% endcomment %}
<div id="story" class="tg-col tg-pad-b(xxl)">
  <div id="progress" class="tg-surface(progress)" aria-hidden="true"></div>

  {% comment %}Scene: what happens today. Agent, the flow of merges, the live site; the slip travels all the way.{% endcomment %}
  <section class="tg-col tg-y-gap(xl) tg-pad-y(xxl)">
    <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(md)">
      <h2 class="tg-text(punch)">It shipped <span class="tg-text(refused)">#ff0066</span>. Nobody asked.</h2>
      <p class="tg-text(lead) tg-max-width(reading)">A new class here, a hex colour there, a style block when it is in a hurry. Nobody reviews CSS, and your agent knows it.</p>
    </div>
    <div class="tg-center tg-width(full) tg-max-width(wide) tg-pad-x(md) tg-col lg:tg-row lg:tg-align(center) tg-gap(md)">
      <div class="tg-grow tg-col tg-y-gap(xs)">
        <p class="tg-text(caption)">Your agent writes careers.html</p>
        <div class="tg-surface(code)">
{% include lines.html text=agent_page %}
        </div>
      </div>
      <div data-rise class="tg-motion(rise) lg:tg-fixed lg:tg-width(aside) tg-col tg-y-gap(xs)">
        <p class="tg-text(caption)">Merged this week</p>
{% include stream.html %}
      </div>
      <div data-rise class="tg-motion(rise) tg-grow tg-col tg-y-gap(xs)">
        <p class="tg-text(caption)">Live on your site</p>
        <div class="tg-surface(code)" data-alarm>
{% include lines.html text=live_page %}
        </div>
      </div>
    </div>
  </section>

  {% comment %}Scene: the same page with the check in the way. It stops in the flow.{% endcomment %}
  <section class="tg-col tg-y-gap(xl) tg-pad-y(xxl)">
    <div class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(md)">
      <h2 class="tg-text(punch)">Same agent. Same crap. <span class="tg-text(refused)">Exit code 1.</span></h2>
      <p class="tg-text(lead) tg-max-width(reading)">Run before every commit, the check names the fix. The agent cleans up after itself. You never see it.</p>
    </div>
    <div class="tg-center tg-width(full) tg-max-width(wide) tg-pad-x(md) tg-col lg:tg-row lg:tg-align(center) tg-gap(md)">
      <div class="tg-grow tg-col tg-y-gap(xs)">
        <p class="tg-text(caption)">Your agent writes careers.html again</p>
        <div class="tg-surface(code)">
{% include lines.html text=slip_page %}
        </div>
      </div>
      <div data-rise class="tg-motion(rise) lg:tg-fixed lg:tg-width(aside) tg-col tg-y-gap(xs)">
        <p class="tg-text(caption)">Merged this week</p>
{% include stream.html stopped=true %}
      </div>
      <div data-rise class="tg-motion(rise) tg-grow tg-col tg-y-gap(xs)">
        <p class="tg-text(caption)">tidgrid check, before the commit</p>
        <div class="tg-surface(code)" data-alarm>
{% include lines.html text=slip_check %}
        </div>
      </div>
    </div>
  </section>

  {% comment %}Scene: the vocabulary, built up call by call.{% endcomment %}
  <section class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(xl) tg-pad-y(xxl)">
    <div class="tg-col tg-y-gap(md)">
      <h2 class="tg-text(punch)">The whole look. Decided once.</h2>
      <p class="tg-text(lead) tg-max-width(reading)">Colours, spacing and every named look live in the vocabulary. Nothing else gets to have an opinion.</p>
      <a class="tg-link(text) tg-self(start)" href="{{ site.baseurl }}/docs/vocabulary-file">The vocabulary file</a>
    </div>
    <div class="tg-max-width(reading) tg-col tg-y-gap(xs)">
      <p class="tg-text(caption)">index.tg</p>
      <div class="tg-surface(code) tg-col tg-y-gap(md)">
        <div data-rise class="tg-motion(rise)"><pre><code>{{ tg_tokens | strip | escape }}</code></pre></div>
        <div data-rise class="tg-motion(rise)"><pre><code>{{ tg_card | strip | escape }}</code></pre></div>
        <div data-rise class="tg-motion(rise)"><pre><code>{{ tg_title | strip | escape }}</code></pre></div>
      </div>
    </div>
  </section>

  {% comment %}Scene: list. What the agent reads, printed row by row.{% endcomment %}
  <section class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(xl) tg-pad-y(xxl)">
    <div class="tg-col tg-y-gap(md)">
      <h2 class="tg-text(punch)">It cannot guess. It has to read.</h2>
      <p class="tg-text(lead) tg-max-width(reading)">A guessed class fails the check. One command prints what exists and when to use it: your agent's entire CSS education.</p>
      <a class="tg-link(text) tg-self(start)" href="{{ site.baseurl }}/docs/roles">How roles work</a>
    </div>
    <div class="tg-max-width(reading) tg-col tg-y-gap(xs)">
      <p class="tg-text(caption)">The agent's terminal, the roles</p>
      <div class="tg-surface(code)">
{% include lines.html text=list_run %}
      </div>
    </div>
  </section>

  {% comment %}Scene: what it is not. The file types the checker reads, as chips.{% endcomment %}
  <section class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(lg) tg-pad-y(xxl)">
    <h2 class="tg-text(punch)">No runtime. No framework. No plugin.</h2>
    <p class="tg-text(lead) tg-max-width(reading)">One binary writes one static stylesheet. Your pages link it like it is 2005.</p>
    <div data-rise class="tg-motion(rise) tg-row tg-wrap tg-gap(xs)">
      <span class="tg-text(keyword)">.html</span>
      <span class="tg-text(keyword)">.astro</span>
      <span class="tg-text(keyword)">.jsx</span>
      <span class="tg-text(keyword)">.tsx</span>
      <span class="tg-text(keyword)">.vue</span>
      <span class="tg-text(keyword)">.svelte</span>
      <span class="tg-text(keyword)">.md</span>
      <span class="tg-text(keyword)">.mdx</span>
      <span class="tg-text(keyword)">.liquid</span>
      <span class="tg-text(keyword)">.njk</span>
    </div>
  </section>

  {% comment %}Scene: the reader's turn.{% endcomment %}
  <section class="tg-center tg-width(full) tg-max-width(measure) tg-pad-x(md) tg-col tg-y-gap(xl) tg-pad-y(xxl)">
    <div class="tg-col tg-y-gap(md)">
      <h2 class="tg-text(punch)">Four commands. Then go do something else.</h2>
      <p class="tg-text(lead) tg-max-width(reading)">Add three lines to your agent's instructions and stop reading CSS diffs.</p>
      <div class="tg-col sm:tg-row tg-wrap tg-gap(sm)">
        <a class="tg-action(primary)" href="{{ site.baseurl }}/get-started">Get started</a>
        <a class="tg-action(secondary)" href="{{ site.baseurl }}/docs">Read the documentation</a>
        <a class="tg-action(secondary)" href="https://github.com/sneikki/tidgrid">Source on GitHub</a>
      </div>
    </div>
    <div class="tg-max-width(reading) tg-col tg-y-gap(xs)">
      <p class="tg-text(caption)">Your terminal</p>
      <div class="tg-surface(code)">
{% include lines.html text=start_run %}
      </div>
    </div>
  </section>
</div>

{% comment %}A module: its names stay its own (theme.js and backdrop.js share the page's global names).{% endcomment %}
<script type="module" src="{{ site.baseurl }}/assets/js/story.js?v={{ site.time | date: "%s" }}"></script>
