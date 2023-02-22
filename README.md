<p align="center">
    <img src="assets/logo.svg">
</p>

# Tidgrid

Tiny and expressive CSS library built around Flexbox to simplify your layouts.

[![Testing](https://github.com/sneikki/tidgrid/actions/workflows/test.yml/badge.svg)](https://github.com/sneikki/tidgrid/actions/workflows/test.yml)
[![Deployment](https://github.com/sneikki/tidgrid/actions/workflows/deploy.yml/badge.svg)](https://github.com/sneikki/tidgrid/actions/workflows/deploy.yml)

> :warning: **This project is still on a very early stage on development.**
> There certainly are bugs, and some of the features may not function correctly under all circumstances.
> Should you use this library in your on project, do it at your own discretion.
> If you find a bug or have a suggestion, don't hesitate to collaborate!

# Quick usage

## Inline stylesheet

The easiest way to get going is to add link to [a hosted version](https://cdn.jsdelivr.net/gh/sneikki/tidgrid@v0.1.0-alpha/build/tidgrid.css) of Tidgrid. This is intended merely for development purposes, as it includes all of the generated CSS on your site. You'll usually use just a small subset of Tidgrid. Thus, it is wise to purge unused styles to reduce the amount of CSS deployed to production.

## Building

Install dependencies by running `npm install`. Then run `npm run build` or `make build` to build Tidgrid from source. The build will be minimized,
and source map is generated automatically. To compile with more options,
please refer to the next section.

The compiled CSS files will be placed under the `build` directory.
To remove all generated build files and directories, run `make clean`.

Please note that the files under build are only updated within
releases. To get the unstable version, you need to build it from the source.

## Compiling

Run `npm run compile` to generate the CSS files. Note that the resulting files
are provided as-is, having comments in place and no minimizing. To strip
comments and minimize the result, run `make compile-minimized`.

## Checking code style

There is a linter configured for this project. To use it,
run `make lint-scss`. To fix style issues that may be automatically
fixed, run `make fix-scss`.
