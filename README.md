<p align="center">
    <a href="https://tidgrid.dev">
        <img src="assets/logo.svg">
    </a>
</p>
<h3 align="center">
    <a href="https://tidgrid.dev">tidgrid.dev</a>
</h3>

# Tidgrid

Tiny and expressive CSS library built around Flexbox to simplify your layouts.

> :warning: **This project is still in a very early stage of development.**
> There certainly are bugs, and some of the features may not function correctly under all circumstances.
> Should you use this library in your project, do it at your discretion.
> If you find a bug or have a suggestion, don't hesitate to collaborate!

# Quick usage

## Inline stylesheet

The easiest way to get going is to add a link to your document to a [hosted version](https://cdn.jsdelivr.net/gh/sneikki/tidgrid@v0.1.0-alpha/build/tidgrid.css) of Tidgrid. This approach is intended merely for development purposes, as it includes all of the generated CSS on your site. You'll usually use just a tiny subset of Tidgrid. Thus, it is wise to purge unused styles to reduce the amount of CSS deployed to production.

## Building

Run 'npm install' to install dependencies. Then, to build Tidgrid from the source, run 'npm run build' or'make build'. The build will be minimized, and the source map will be automatically generated. Please see the following section for more options.

The CSS files that have been built will be put in the `build` directory.
Run `npm run clean` to delete all created build files and folders.

Please keep in mind that the files in the build directory are only changed between releases.
To obtain the unstable version, you must compile it from source.

## Compiling

To create the CSS files, use `npm run compile`.
Please keep in mind that the generated files are offered as-is, with no minimizing and no comments.
Run `npm run compile:minimized` to remove comments and reduce the size of the output.

## Checking code style

This project has a linter configured.
Run `npm run lint:scss` to utilize it.
Run `npm run fix:scss` to resolve auto-fixable style problems.
