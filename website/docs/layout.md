---
layout: doc-page
title: Layout
---

{% include note.html type="note" content="This page is incomplete." %}

## Section

Sections split the page logically. Sections should independent from the other sections.
This is a universal practice, but Tidgrid aims to improve the layout. When a topic
changes, there is naturally some space between the previous.
The <span class="doc__keyword">section</span> adds vertical padding to create a
contrast between sections and make them visually distinguishable.


Adjust the amount of padding using <span class="doc__keyword">section-padding</span>.
Currently there are only two values:

* <span class="doc__keyword">none</span>
* <span class="doc__keyword">default</span>

### Footer

In addition to the general section element, Tidgrid provides a derivative footer section element.
The footer automatically adds extra padding to the bottom, which signifies the end of the web page.

## Container

Containers are what the cells are for rows; organizing the internal structure of sections.
Sections take care of vertical spacing, but anything placed inside a section would float
to the left side of the screen. Horizontal spacing is trickier compared to vertical, as
viewport width drastically changes by device. The container handles horizontal spacing
inside a section by taking care of margins
and centering. That is, containers are direct children of sections. A single section
might hold one or multiple adjacent containers.

By default, a <span class="doc__keyword">container</span>, when placed inside a section,
centerizes itself and leaveas equal amount of margin on both sides. Width of the container
is always the size of the currently active breakpoint. That is, if the viewport is 55em wide,
the container would be 48em wide and horizontally centered.

Containers grow until the last breakpoint is triggered. However, typically width of the
web page's content ceases to grow at a smaller width. To limit the width of a container,
give it a breakpoint as a value. For example, <span class="doc__keyword">tg-container(sm)</span>
would stop growing width when it is as wide as the sm breakpoint.

Consider a section containing text, but there should be a narrower off-text element
in the middle of the text. This is easily achieved by using containers of different sizes.

{% highlight html %}
<section class="tg-section">
    <div class="tg-container(lg)">
        <!-- Text -->
    </div>

    <div class="tg-container(md)">
        <!-- Narrower off-text element, e.g. image or advertisement -->
    </div>

    <div class="tg-container(lg)">
        <!-- More text -->
    </div>
</section>
{% endhighlight %}
