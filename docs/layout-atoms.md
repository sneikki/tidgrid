---
layout: doc-page
title: Layout atoms
---

Atoms are the classes for layout. Their grammar is fixed and the same in every project, so an agent that knows one Tidgrid project knows the atoms of all of them. Rows, gaps, padding and widths are free to combine, and always run on the project's token scales.

In the lists below, `S` is every name in `tokens.space` and `Z` is every name in `tokens.size`. Each line is a class, then the CSS it sets.

## Flex

{% capture c %}
tg-row                                              display: flex; flex-direction: row
tg-col                                              display: flex; flex-direction: column
tg-wrap                                             flex-wrap: wrap
tg-grow                                             flex: 1 1 0; min-width: 0
tg-fixed                                            flex: none
tg-center                                           margin-inline: auto
tg-align(start|center|end|baseline|stretch)         align-items
tg-justify(start|center|end|between|around|evenly)  justify-content
tg-self(start|center|end|stretch)                   align-self
tg-text-align(start|center|end)                     text-align
{% endcapture %}
{% include code.html code=c %}

## Space

{% capture c %}
tg-gap(S)  tg-x-gap(S)  tg-y-gap(S)                 gap, column-gap, row-gap
tg-pad(S)  tg-pad-x(S)  tg-pad-y(S)                 padding, padding-inline, padding-block
tg-pad-t(S)  tg-pad-r(S)  tg-pad-b(S)  tg-pad-l(S)  padding-top ... padding-left
{% endcapture %}
{% include code.html code=c %}

## Size

{% capture c %}
tg-width(Z)   tg-min-width(Z)   tg-max-width(Z)     width, min-width, max-width
tg-height(Z)  tg-min-height(Z)                      height, min-height
{% endcapture %}
{% include code.html code=c %}

## Grid and order

{% capture c %}
tg-cols(1...6|12)       display: grid; grid-template-columns: repeat(N, minmax(0, 1fr))
tg-span(1...12)         grid-column: span N
tg-order(first|last)    order: -1, order: 1
{% endcapture %}
{% include code.html code=c %}

## Display

{% capture c %}
tg-show    display: revert
tg-hide    display: none
{% endcapture %}
{% include code.html code=c %}

## Breakpoints

Every atom also exists per breakpoint, mobile first. A prefix such as `md:` applies the atom from that breakpoint up: `md:tg-gap(md)` sets the gap from `md` up, and `tg-hide md:tg-show` hides an element below `md`.

A small example, a column on a phone and a row from `md` up:

<!-- tidgrid-ignore-start -->
{% capture c %}
<div class="tg-col md:tg-row tg-gap(md) tg-pad(md)">
  <div class="tg-grow">One</div>
  <div class="tg-grow">Two</div>
</div>
{% endcapture %}
{% include code.html code=c %}
<!-- tidgrid-ignore-end -->

## Atom or role

An element takes a property from an atom or from a role, never both: atoms are written after roles and win.

Layout the atoms do not cover, such as an overlay or a fixed panel, is a role. Add it to the vocabulary file in its own commit.

## Next

- [Roles](/docs/roles)
- [Tokens](/docs/tokens)
