---
layout: doc-page
title: Grid overview
---

{% include note.html type="note" content="This page is incomplete." %}

## Introduction

The grid module is the heart of Tidgrid. It offers a variety
of options to control the structure of a web page. Tidgrid provides
two basic elements for controlling the grid:

* Rows
* Cells

Row is the top-level element in a grid. It acts as a container
for cells, that in turn contain the actual elements of the web page.
A row defines the overall structure for the cells it contains.
A cell may define more specific styles that apply only to that cell.

## Core

<div class="tg-size(12)">
    <div class="tg-row tg-gap(xl) md:tg-span(4)">
        {% include doc-card.html title="Row"
            href="/docs/grid/row"
            description="Lorem ipsum dolor sit amet" %}

        {% include doc-card.html title="Cell"
            href="/docs/grid/cell"
            description="Lorem ipsum dolor sit amet" %}

        {% include doc-card.html title="Modes"
            href="/docs/grid/modes"
            description="Lorem ipsum dolor sit amet" %}
    </div>
</div>

## Spacing, positioning & gaps

<div class="tg-size(12)">
    <div class="tg-row tg-gap(xl) md:tg-span(4)">
        {% include doc-card.html title="Spacing"
            href="/docs/grid/spacing"
            description="Lorem ipsum dolor sit amet" %}

        {% include doc-card.html title="Positioning"
            href="/docs/grid/positioning"
            description="Lorem ipsum dolor sit amet" %}

        {% include doc-card.html title="Gaps"
            href="/docs/grid/gaps"
            description="Lorem ipsum dolor sit amet" %}
    </div>
</div>
