---
layout: doc-page
title: Builtins
---

A vocabulary is a Starlark program, and these are the functions it calls. They are available in every file without importing. Each call registers one thing in the model the tool builds.

| Call | Registers |
|---|---|
| `project(out, scan=["."], prefix="tg", helpers=[...], repeats=[(3, 2), (2, 4)])` | Where the stylesheet goes (`out`, required), which folders `check` scans, the class prefix, the helper functions `check` reads class names from, the repeated-lines setting. Once. |
| `breakpoints(sm="40em", ...)` or `breakpoints({...})` | Breakpoints, ascending in call order. Values are literal, because media queries cannot read custom properties. |
| `tokens(group, {name: value})` or `tokens(group, name=value, ...)` | A token group. Returns the group, so `color.ink` is a reference to a token. See [tokens](/docs/tokens). |
| `element(selector, css={...}, use="")` | A base rule for a bare element (`body`, `a`, `"*"`). `use` is optional. |
| `keyframes(name, {stop: css})` | An animation. Returns a reference that becomes the animation name `tg-<name>`. |
| `role(kind, name, use=..., css={...}, include=[...])` | The class `tg-<kind>(<name>)`. `use` is required. See [roles](/docs/roles). |
| `mixin(name, use=..., sets={...})` | Named sets of lines that roles include. See [mixins](/docs/mixins). |
| `allow_class(name, reason)` | A class `check` accepts although it is not in the vocabulary. |
| `allow_css(path, reason)` | A stylesheet `check` accepts. |
| `load(file, names...)`, `source(file)` | Run another file. See [files](/docs/files). |

A call of each:

{% capture c %}
project(out = "tidgrid.css", scan = ["."], repeats = [(3, 2), (2, 4)])
breakpoints(sm = "40em", md = "64em")
color = tokens("color", ink = "#e6e9eb", accent = "#5ed5dc")
motion = tokens("motion", {"slow": "1.1s"})
element("a", css = {"color": color.accent})
spin = keyframes("spin", {"to": {"transform": "rotate(360deg)"}})
bits = mixin("bits", use = "Two lines.", sets = {"pair": {"cursor": "pointer", "color": color.ink}})
role("action", "go", use = "Go.", include = [bits.pair], css = {"animation": "%s %s linear" % (spin, motion.slow)})
role("action", "stop", use = "Stop.", include = [bits.pair])
allow_class("has-js", "set by a script")
allow_css("fonts.css", "self-hosted @font-face only")
{% endcapture %}
{% include code.html code=c caption="index.tg" %}

Declaring `project`, a token group, a role or a mixin twice is an error that names both places:

{% capture c %}
index.tg:3:8: project: declared twice (index.tg:1 and index.tg:3)
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

A wrong argument is caught with a hint:

{% capture c %}
index.tg:2:5: role: unexpected keyword argument "usee" (did you mean use?)
{% endcapture %}
{% include code.html code=c plain=true caption="What build prints" %}

## Next

- [The vocabulary file](/docs/vocabulary-file)
- [Roles](/docs/roles)
