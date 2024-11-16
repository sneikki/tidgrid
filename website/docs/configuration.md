---
layout: doc-page
title: Configuration
---

{% include note.html type="note" content="This page is incomplete." %}

Tidgrid indeed ships with its more or less opinionated default values.
However, they can never be suitable for all situations. Therefore, Tidgrid
is highly configurable and there is room for differing personal preferences.

Each module defines the default values for variables it needs;
the values can easily be overridden by simply assigning to that
variable in your own stylesheet. For example, if you prefer other
breakpoint sizes or names, just assign new mapping to the
<span class="doc__keyword">$default-breakpoints</span> variable:

{% highlight html %}
$default-breakpoints: (
  mobile: 36em,
  tablet: 48em,
  laptop: 64em,
  fullhd: 90em
);
{% endhighlight %}

If you don't want to repeatedly write the tg- namespace,
assign 
<span class="doc__keyword">$use-namespace: false;</span>.
Or if the parentheses syntax isn't your piece of cake,
switch to dashes: 
<span class="doc__keyword">$syntax: dashes</span>.
