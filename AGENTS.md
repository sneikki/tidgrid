# Site (`website/`, `docs/`)

The site tidgrid.dev: a Jekyll site in `website/`, its documentation pages in `docs/` (`website/docs` is a symlink to it). The rest of the repo is the tool itself.

- The live site is the branch `site` (GitHub Pages, `.github/workflows/site.yml`). It follows the default branch `master`: every push to master moves `site` forward and deploys. Never push to `site` yourself, with one exception, a hotfix: a pull request against `site` deploys when merged, and `site` is then merged into master with a merge commit, never squashed or rebased, before anything else lands there (the workflow fails on master until it is).
- A pull request that changes a file in `.github/workflows` cannot be carried to `site` by the workflow's own token: after it is merged the run on master fails at "Fast-forward site to master", and `git push origin master:site` by a person finishes it.
- A pull request runs the check and the build without deploying. The workflow builds the tool from the same commit, so a tool change that breaks the site's check fails the pull request.
- The site is styled with the tool it documents. Style only with the classes `tidgrid list` prints (run it in `website/`). No `<style>`, no `style=` for appearance, no own CSS.
- Run `cd website && go run ../cmd/tidgrid build && go run ../cmd/tidgrid check` before every commit; check must exit 0. Commit the stylesheet `build` wrote. That is the tool from this checkout, which is what the workflow uses; an installed `tidgrid` is a release and may differ.
- A new role goes in `website/index.tg` or a file it loads, in the same commit as the first page that uses it; `tidgrid check` fails on a role nothing uses.
- A code sample that shows classes goes between `tidgrid-ignore-start` and `tidgrid-ignore-end` comments, and into `{% include code.html %}`; pass `plain=true` for a shell command, tool output or prose.
- `version` in `website/_config.yml` is the released version the footer and the install line show; the workflow does not install it. After a tidgrid release, change that line; what the pages say about the tool must match that version.
