---
layout: doc-page
title: Grid
---

The grid module consists of two elements: 
<span class="doc__keyword">row</span> and
<span class="doc__keyword">cell</span>. Rows are used
to group cells together, and cells contain the actual
elements of the document. Rows can contain any number
of cells, but they may not contain any other elements.
Modifiers control the structure of the cells within
a single row. Row modifiers declare the default
structure for all child cells. Cell modifiers
affect only a single cell.

## Modes

### stacked

<span class="doc__keyword">stacked</span> is the default mode
when no other mode is explicitly used. Each cell takes up all
horizontal space, and thus each cell has its own line. Cells
that are in <span class="doc__keyword">stacked</span> mode
force the next cells to break a new line, no matter what mode
the other one has. <span class="doc__keyword">stacked</span> comes handy with
smaller viewports, where the elements would not fit if they
were arranged horizontally.

### auto
In <span class="doc__keyword">auto</span> mode,
cells take up all the space horizontally. If there are adjacent
cells with other size modifiers, cells with
<span class="doc__keyword">auto</span> mode
respect their size and fill the remaining space.

### thin

<span class="doc__keyword">thin</span> mode is different from
the previous two modes. Instead of expanding horizontally, it
shrinks the cells as much as possible. Cells take up only as
much space as their content requires. It may be thought like
the inline box in CSS.

### disabled

<span class="doc__keyword">disabled</span> mode removes
the cells from the page layout.

## Other sizing modifiers

The following modifiers are not independent and
not controlled with <span class="doc__keyword">tg-mode</span>.
However, they control the size of the cells, too, and may
be used alongside the modes presented above.

### span

<span class="doc__keyword">tg-span</span> sets the size
of a cell based on the column count of the row. For example,
within a row of 16 columns,
<span class="doc__keyword">tg-span(8)</span> would make
the cell span half of its horizontal space, and
<span class="doc__keyword">tg-span(4)</span> one quarter.

When applied to a row, the row must be decessor of an
element defining the column count using
<span class="doc__keyword">tg-size</span>. The row does
not need to be a direct children; there might be multiple
ancestors between. This way the column count can be set
for the whole document, or part of it, without needing
to write it over and over again.

### frac

<span class="doc__keyword">tg-frac</span> works the same
way as <span class="doc__keyword">tg-span</span>, but
uses fractions instead of columns. 
<span class="doc__keyword">tg-span(4)</span> would take up
1/3 of the space within a row with 12 columns, but only 1/4
within a row with 16 columns. 
<span class="doc__keyword">tg-frac(1/2)</span> always takes up
half of the horizontal space, regardless of the column count.

## Spacing

Spacing is relatively straightforward, and it is controlled
with just one modifier: <span class="doc__keyword">tg-space</span>,
that controls both horizontal and vertical spacing. There are also
two derivatives to control only one of the axes:
<span class="doc__keyword">tg-space-x</span> and
<span class="doc__keyword">tg-space-y</span>
The values are practically same as the ones used in Flexbox:

* <span class="doc__keyword">between</span>
* <span class="doc__keyword">around</span>
* <span class="doc__keyword">evenly</span>
* <span class="doc__keyword">none</span>

## Positioning

Like spacing, positioning has just one modifier and two derivatives:
<span class="doc__keyword">tg-pos</span>,
<span class="doc__keyword">tg-pos-x</span> and
<span class="doc__keyword">tg-pos-y</span>.

Supported values for this modifier are the following:

* <span class="doc__keyword">start</span>
* <span class="doc__keyword">center</span>
* <span class="doc__keyword">end</span>
* <span class="doc__keyword">none</span>

## Gaps

Gaps are declared for <span class="doc__keyword">row</span>
elements. The <span class="doc__keyword">tg-gap</span>
modifier, and its two derivatives, 
<span class="doc__keyword">tg-x-gap</span> and
<span class="doc__keyword">tg-y-gap</span>, control
the gaps for horizontal and vertical axis.
Gap sizes are configurable, and Tidgrid provides the
following default configuration:

* <span class="doc__keyword">none</span>: 0rem
* <span class="doc__keyword">xs</span>: 0.25rem
* <span class="doc__keyword">sm</span>: 0.5rem
* <span class="doc__keyword">md</span>: 0.75rem
* <span class="doc__keyword">lg</span>: 1rem
* <span class="doc__keyword">xl</span>: 2rem

## Order

The order modifier, <span class="doc__keyword">tg-order</span>,
is used for cells only. Two values are currently supported:

* <span class="doc__keyword">first</span>
* <span class="doc__keyword">last</span>
