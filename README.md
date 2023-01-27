<p align="center">
    <img src="assets/logo.svg">
</p>

# Tidgrid

Modern CSS framework to ease front-end development.

# Quick usage

## Building

Run `yarn build` to build Tidgrid from source. The build will be minimized,
and source map is generated automatically. To compile with more options,
please refer to the next section.

The compiled CSS files will be placed under the `build` directory.
To remove all generated build files and directories, run `yarn clean`.

Please note that the files under build are only updated within
releases. To get the unstable version, you need to build it from the source.

## Compiling

Run `yarn compile` to generate the CSS files. Note that the resulting files
are provided as-is, having comments in place and no minimizing. To strip
comments and minimize the result, run `yarn compile:minimized`.

To compile automatically after changes to source files, run `yarn watch`.
Run `yarn repl:scss` to run an interactive SassScript REPL session. 

## Checking code style

There is a linter configured for this project. To use it,
run `yarn lint:scss`.
