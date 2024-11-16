---
layout: default
title: Home
---

<section class="tg-section jumbo">
    <div class="tg-container">
        <div class="tg-row tg-y-gap(xl) tg-pos-y(center)">

            <div class="tg-cell lg:tg-span(8)">
                <div class="jumbo__text-logo">
                    <p>tidgrid</p>
                    <p class="jumbo__text-logo--accent">.dev</p>
                </div>
            </div>

            <div class="tg-cell lg:tg-span(8)">
                <div class="tg-row jumbo__buttons tg-y-gap(lg) tg-x-gap(lg) tg-x-gap(xl)">
                    <div class="tg-cell tg-span(16) sm:tg-mode(auto)">
                        <a href="{{ site.baseurl }}/get-started" class="jumbo__button jumbo__button--primary">
                            Get started
                        </a>
                    </div>

                    <div class="tg-cell xs:tg-span(16) sm:tg-mode(auto)">
                        <a href="{{ site.baseurl }}/docs" class="jumbo__button jumbo__button--secondary">
                            Documentation
                        </a>
                    </div>
                </div>
            </div>

            <div class="tg-cell">
                <p class="jumbo__description">
                    Tiny and expressive CSS library for responsive web development
                </p>
            </div>
        </div>
    </div>
</section>

{% capture source %}
{% highlight html %}
<div class="tg-row">
    <div class="tg-cell tg-span(4)">
        <div class="box"></div>
    </div>

    <div class="tg-cell tg-span(12)">
        <div class="box"></div>
    </div>

    <div class="tg-cell tg-span(8)">
        <div class="box"></div>
    </div>

    <div class="tg-cell tg-span(8)">
        <div class="box"></div>
    </div>

    <div class="tg-cell tg-span(12)">
        <div class="box"></div>
    </div>

    <div class="tg-cell tg-span(4)">
        <div class="box"></div>
    </div>
</div>
{% endhighlight %}
{% endcapture %}

{% capture result %}
    <div class="tg-row tg-gap(lg)">
        <div class="tg-cell tg-span(4)">
            <div class="features__box"></div>
        </div>

        <div class="tg-cell tg-span(12)">
            <div class="features__box"></div>
        </div>

        <div class="tg-cell tg-span(8)">
            <div class="features__box"></div>
        </div>

        <div class="tg-cell tg-span(8)">
            <div class="features__box"></div>
        </div>

        <div class="tg-cell tg-span(12)">
            <div class="features__box"></div>
        </div>

        <div class="tg-cell tg-span(4)">
            <div class="features__box"></div>
        </div>
    </div>
{% endcapture %}

{% include feature.html
    title="Powerful yet simple grid system"
    source=source
    result=result %}



{% capture source %}
{% highlight html %}
<div class="tg-row">
    <div class="tg-cell sm:tg-span(8) md:tg-mode(auto)">
        <div class="box"></div>
    </div>

    <div class="tg-cell sm:tg-span(8) md:tg-mode(auto)"> 
        <div class="box"></div>
    </div>

    <div class="tg-cell sm:tg-span(8) md:tg-mode(auto)">
        <div class="box"></div>
    </div>

    <div class="tg-cell sm:tg-span(8) md:tg-mode(auto)">
        <div class="box"></div>
    </div>

    <div class="tg-cell sm:tg-span(16) lg:tg-span(8)">
        <div class="box"></div>
    </div>

    <div class="tg-cell sm:tg-span(16) lg:tg-span(8)">
        <div class="box"></div>
    </div>
</div>
{% endhighlight %}
{% endcapture %}

{% capture result %}
    <div class="tg-row tg-gap(lg)">
        <div class="tg-cell sm:tg-span(8) md:tg-mode(auto)">
            <div class="features__box"></div>
        </div>

        <div class="tg-cell sm:tg-span(8) md:tg-mode(auto)"> 
            <div class="features__box"></div>
        </div>

        <div class="tg-cell sm:tg-span(8) md:tg-mode(auto)">
            <div class="features__box"></div>
        </div>

        <div class="tg-cell sm:tg-span(8) md:tg-mode(auto)">
            <div class="features__box"></div>
        </div>

        <div class="tg-cell sm:tg-span(16) lg:tg-span(8)">
            <div class="features__box features__box--dark"></div>
        </div>

        <div class="tg-cell sm:tg-span(16) lg:tg-span(8)">
            <div class="features__box features__box--dark"></div>
        </div>
    </div>
{% endcapture %}

{% include feature.html
    title="Responsivity baked in"
    reverse=true
    source=source
    result=result %}




{% capture source %}
{% highlight html %}
<div class="tg-row sm:tg-mode(auto) md:tg-mode(thin)
    lg:tg-mode(stacked) md:tg-space-x(between)">
    <div class="tg-cell">
        <div class="box">Box</div>
    </div>

    <div class="tg-cell">
        <div class="box">Box</div>
    </div>

    <div class="tg-cell">
        <div class="box">This box has some content</div>
    </div>

    <div class="tg-cell">
        <div class="box">Box</div>
    </div>

    <div class="tg-cell">
        <div class="box">Box</div>
    </div>

    <div class="tg-cell">
        <div class="box">This one too</div>
    </div>
</div>
{% endhighlight %}
{% endcapture %}

{% capture result %}
    <div class="tg-row sm:tg-mode(auto) md:tg-mode(thin) lg:tg-mode(stacked) md:tg-space-x(between) tg-gap(lg)">
        <div class="tg-cell">
            <div class="features__box">Box</div>
        </div>

        <div class="tg-cell">
            <div class="features__box">Box</div>
        </div>

        <div class="tg-cell">
            <div class="features__box">This box has some content</div>
        </div>

        <div class="tg-cell">
            <div class="features__box">Box</div>
        </div>

        <div class="tg-cell">
            <div class="features__box">Box</div>
        </div>

        <div class="tg-cell">
            <div class="features__box">This one too</div>
        </div>
    </div>
{% endcapture %}

{% include feature.html
    title="Easy-to-use modes"
    source=source
    result=result %}





{% capture source %}
{% highlight html %}
<div class="tg-row tg-y-gap(xl)">
    <div class="tg-cell">
        <div class="tg-row tg-x-gap(md)">
            <div class="tg-cell tg-span(5)">
                <div class="box"></div>
            </div>

            <div class="tg-cell tg-mode(fill)">
                <div class="box"></div>
            </div>

            <div class="tg-cell tg-span(3)">
                <div class="box"></div>
            </div>
        </div>
    </div>

    <div class="tg-cell">
        <div class="tg-row tg-mode(auto) tg-x-gap(xl)">
            <div class="tg-cell">
                <div class="box"></div>
            </div>

            <div class="tg-cell">
                <div class="box"></div>
            </div>
        </div>
    </div>
</div>
{% endhighlight %}
{% endcapture %}

{% capture result %}
    <div class="tg-row tg-y-gap(xl)">
        <div class="tg-cell">
            <div class="tg-row tg-x-gap(md)">
                <div class="tg-cell tg-span(5)">
                    <div class="features__box"></div>
                </div>

                <div class="tg-cell tg-mode(fill)">
                    <div class="features__box"></div>
                </div>

                <div class="tg-cell tg-span(3)">
                    <div class="features__box"></div>
                </div>
            </div>
        </div>

        <div class="tg-cell">
            <div class="tg-row tg-mode(auto) tg-x-gap(xl)">
                <div class="tg-cell">
                    <div class="features__box"></div>
                </div>

                <div class="tg-cell">
                    <div class="features__box"></div>
                </div>
            </div>
        </div>
    </div>
{% endcapture %}

{% include feature.html
    title="Gaps"
    reverse=true
    source=source
    result=result %}
