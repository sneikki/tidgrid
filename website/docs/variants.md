---
layout: doc-page
title: Variants
---

{% include note.html type="note" content="This page is incomplete." %}

Variants provide responsivity in Tidgrid. Almost all grid
modifiers are varying, which means that they can declared
to be active only for certain viewport sizes. Varying
modifiers can be prepended with breakpoints. Then,
that modifier becomes active only when the viewport is
wide enough to trigger the breakpoint. For example,
<span class="doc__keyword">md:tg-mode(auto)</span>
sets the auto mode, but only if viewport width is
at least 64em.

Breakpoints allow repeating the same modifier
for a single element. As long as the breakpoints
are different, an element could have class name
such as 
<span class="doc__keyword">tg-mode(auto) sm:tg-mode(stacked) xl:tg-mode(thin)</span>. Auto mode applies when viewport is narrower than the <span class="doc__keyword">sm</span> breakpoint, stacked mode until viewport width hits the <span class="doc__keyword">xl</span> breakpoint, and finally thin mode applies for viewports that trigger the <span class="doc__keyword">xl</span> breakpoint.

Tidgrid provides the following set of breakpoints by default:
* <span class="doc__keyword">xs</span>: 36em
* <span class="doc__keyword">sm</span>: 48em
* <span class="doc__keyword">md</span>: 64em
* <span class="doc__keyword">lg</span>: 75em
* <span class="doc__keyword">xl</span>: 90em
