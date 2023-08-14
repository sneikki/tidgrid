---
layout: doc-page
title: Installation
---

## Installing with a package manager

### npm

{% highlight sh %}
npm install tidgrid
{% endhighlight %}

### yarn

{% highlight sh %}
yarn add tidgrid
{% endhighlight %}

Import Tidgrid in a JavaScript file:

{% highlight js %}
// CommonJS
require("tidgrid");

// ES6
import "tidgrid";
{% endhighlight %}

Or in a Sass file:

{% highlight css %}
@import "tidgrid";
{% endhighlight %}
## Using CDN

Access Tidgrid directly and fetch the latest build from [jsDelivr](https://cdn.jsdelivr.net/gh/sneikki/tidgrid@0.1.0-alpha.2/build/tidgrid.css). The hosted file contains the compiled CSS. Attach it to your project with zero dependencies to package managers and build tools.

{% highlight html %}
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/sneikki/tidgrid@0.1.0-alpha.2/build/tidgrid.css" integrity="sha384-QB3RSUytwAFXoMg/cxubzIPCegtfmMedCoMd6mS3GlMco/vDi75t379PxstKly8F" crossorigin="anonymous">
{% endhighlight %}

Keep in mind that the hosted version loads all of the compiled CSS. You almost always use just a tiny subset of all the styles provided,
and the rest are useless for your project. This may or may not be a problem, based on the circumstances. The minimized build takes up
66 KB space. If your target audience has poor bandwidth, consider installing Tidgrid with a package manager and use tree-shaking
to remove all unused styles.

## Compiling yourself (using npm)

Download [the latest release](https://github.com/sneikki/tidgrid/releases/tag/0.1.0-alpha.2) or clone the repository to get the source code.

### Install dependencies

{% highlight sh %}
npm install
{% endhighlight %}

### Build

Compile Tidgrid and output the unminimized CSS in the build directory.

{% highlight sh %}
npm run build
{% endhighlight %}

### Minimized build

Compile Tidgrid, minimize the CSS and output it in the build directory.

{% highlight sh %}
npm run compile:minimized
{% endhighlight %}
