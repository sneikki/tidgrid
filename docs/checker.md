---
layout: doc-page
title: The checker and how classes are found
---

`tidgrid check` scans the project's sources and prints `path:line: message` for every violation. It exits 0 when there are none and 1 when there are. It loads the vocabulary first: when that fails it prints the problems and exits 2, with no scan. It takes the entry file as an argument, like the other commands (see [files](/docs/files)).

## What it reports

| Violation | Example |
|---|---|
| A class that is not in the vocabulary or under `allow_class` | `border-l-4`, `card-accent` |
| A class name built from pieces | `is-${state}` |
| Appearance in a `style` attribute (geometry is allowed) | `style="color: red"` |
| Own CSS: a style block, or a stylesheet not under `allow_css` | `Card.module.css` |
| A role no scanned source uses | a role added and never used |
| The stylesheet at `out` is missing, or differs from what `build` would write | vocabulary edited, not rebuilt |
| An ignore region (see below) that is never closed | the rest of the file would go unchecked |
| A symlinked directory that is not under `scan` | |

A message names the fix where there is one. From real runs:

<!-- tidgrid-ignore-start -->
{% capture c %}
src/b.html:1: appearance in style attribute: color; use a role (`tidgrid list`)
src/b.html:1: unknown class card-accent: not in the vocabulary; run `tidgrid list`
src/b.html:1: unknown class tg-gap(2); tg-gap takes: sm md
src/b.html:2: own CSS: <style> block; make it a role in the vocabulary
src/b.html:4: class name built from pieces: is-${state}; write every whole name out in a lookup table
src/c.html:1: tidgrid-ignore-start is never closed; the rest of the file is not checked
src/tidgrid.css:1: stylesheet differs from the vocabulary; run `tidgrid build`
src/tidgrid.css:1: stylesheet missing; run `tidgrid build`
{% endcapture %}
{% include code.html code=c plain=true caption="What check prints" %}
<!-- tidgrid-ignore-end -->

After the violations come two lists: per file, the roles it uses, and the `allow` entries with their reasons.

## Unused roles

A role no scanned source uses fails `check`, and the message is at the line of the call that declared it:

{% capture c %}
roles.tg:25: role tg-text(lead) is not used in the scanned sources: use it, remove it, or, if a script or lookup table sets it, write the name whole inside a helper call such as cn()
{% endcapture %}
{% include code.html code=c plain=true caption="What check prints" %}

A role that a script sets counts as used when the script names it whole, as a string: `classList.add("tg-state(ready)")`, `cn("tg-state(ready)")`. A name that appears only in another role's or element's selector does not count. A role made in a loop is reported once per role.

So a new role goes in the same commit as the first page that uses it. A role in a commit of its own fails until the page arrives.

## How classes are found

The checker is a scanner, not a parser, so it works the same in HTML, Astro, JSX, Vue and Svelte. It reads:

- `class="…"` and `className="…"` attributes, in any letter case.
- Every string literal inside `class={…}`, `className={…}`, `class:list={…}` and `:class="…"`.
- Every string literal on the right of a `className = …` assignment in script.
- Every string literal passed to `cn`, `clsx`, `cva`, `classNames` or `twMerge`, and to `classList.add`, `remove`, `toggle` and `contains`. `helpers` in `project()` replaces that list of five, so repeat the ones still in use.
- The name in a `class:name` directive (Svelte). `style:prop` counts as a style attribute.

A string literal that is only compared, such as `kind === "primary"`, is not a class and is not read.

Inside a class value, a template tag `{% raw %}{% … %}{% endraw %}` is skipped and the names around it are read. An interpolation (`{% raw %}{{ … }}{% endraw %}`, `${…}` or `{…}`) glued to a name is a built name, which is a violation. Standing alone between spaces, it is skipped.

## What it does not read

- Fenced code blocks in Markdown.
- Inline code in `.md` files. It is read in `.mdx`, where a backtick can build a class name.
- `{% raw %}{% highlight %}{% endraw %}` and `{% raw %}{% raw %}{% endraw %}` blocks.
- Anything between the comments `tidgrid-ignore-start` and `tidgrid-ignore-end`. Use them around code samples that show classes.

## Which files

Sources: `.html .htm .astro .jsx .tsx .js .mjs .cjs .ts .vue .svelte .md .mdx .liquid .njk` under `scan`. Stylesheets: `.css .scss .sass .less .styl`. These are never entered: `node_modules .git dist build _site .astro vendor`.

A symlinked directory is not entered either. `check` reports it unless its target is under `scan`, so add the target there.

## Known limit

A class kept in a variable far from the places above is not seen. A name built from pieces cannot be checked at all, which is why it is a violation. Write the whole names out in a lookup table instead:

<!-- tidgrid-ignore-start -->
{% capture c %}
const kinds = {
  primary: "tg-action(primary)",
  secondary: "tg-action(secondary)",
};
const cls = kinds[kind];
{% endcapture %}
{% include code.html code=c %}
<!-- tidgrid-ignore-end -->

The names in such a table are whole, but the table is still a variable: `check` does not read it, and a role used only there is listed as unused. To have them checked, pass each name through one of the helper functions (`cn("tg-action(primary)")`).

More mistakes that pass the check, each beside its fix, are on the [pitfalls](/docs/pitfalls) page.

## Next

- [Commands](/docs/commands)
- [Builtins](/docs/builtins): `allow_class` and `allow_css`
