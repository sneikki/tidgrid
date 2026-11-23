# tidgrid-docs

The site tidgrid.dev: a Jekyll site in `website/`, its documentation pages in `docs/` (`website/docs` is a symlink to it). The tool itself lives in github.com/sneikki/tidgrid.

- The default branch is `master`; a push to it deploys the live site (GitHub Pages). A pull request runs the check and the build without deploying.
- The site is styled with the tool it documents. Style only with the classes `tidgrid list` prints (run it in `website/`). No `<style>`, no `style=` for appearance, no own CSS.
- Run `cd website && tidgrid build && tidgrid check` before every commit; check must exit 0. Commit the stylesheet `build` wrote.
- A new role is a change to `website/tidgrid.toml` in its own commit, never mixed into page work.
- A code sample that shows classes goes between `tidgrid-ignore-start` and `tidgrid-ignore-end` comments, and into `{% include code.html %}`; pass `plain=true` for a shell command, tool output or prose.
- What a page says about the tool must match the released version the workflow installs (`.github/workflows/deploy-website.yml`).
