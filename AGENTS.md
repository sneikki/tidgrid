# tidgrid-docs

The site tidgrid.dev: a Jekyll site in `website/`, its documentation pages in `docs/` (`website/docs` is a symlink to it). The tool itself lives in github.com/sneikki/tidgrid.

- The default branch is `master`; a push to it deploys the live site (GitHub Pages). A pull request runs the check and the build without deploying.
- The site is styled with the tool it documents. Style only with the classes `tidgrid list` prints (run it in `website/`). No `<style>`, no `style=` for appearance, no own CSS.
- Run `cd website && tidgrid build && tidgrid check` before every commit; check must exit 0. Commit the stylesheet `build` wrote.
- A new role goes in `website/index.tg` or a file it loads, in the same commit as the first page that uses it; `tidgrid check` fails on a role nothing uses.
- A code sample that shows classes goes between `tidgrid-ignore-start` and `tidgrid-ignore-end` comments, and into `{% include code.html %}`; pass `plain=true` for a shell command, tool output or prose.
- The tool version lives in one place: `version` in `website/_config.yml`. The footer shows it and the workflow installs it. After a tidgrid release, change that line; what the pages say about the tool must match that version.
