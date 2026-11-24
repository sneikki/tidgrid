#!/usr/bin/env node
"use strict";
// Runs tidgrid.wasm under Go's JS glue. Arguments, cwd and environment pass
// through; the process exit code is the tool's own.
const fs = require("fs");
const os = require("os");
const path = require("path");

const root = path.join(__dirname, "..");

// Go's js port only knows "/" paths. On Windows, hand it "C:/dir" for the
// working directory and forward-slash, cwd-relative paths for path arguments.
function windowsPaths(args) {
  const cwd = process.cwd();
  process.cwd = () => cwd.replace(/\\/g, "/");
  return args.map((a) => {
    if (a.startsWith("-") || !/^[A-Za-z]:[\\/]|\\/.test(a)) return a;
    return path.relative(cwd, path.resolve(cwd, a)).replace(/\\/g, "/") || ".";
  });
}

function main() {
  globalThis.require = require;
  globalThis.fs = fs;
  globalThis.path = path;
  globalThis.TextEncoder = require("util").TextEncoder;
  globalThis.TextDecoder = require("util").TextDecoder;
  globalThis.performance ??= require("perf_hooks").performance;
  globalThis.crypto ??= require("crypto");
  require("../wasm_exec");

  const args = process.platform === "win32" ? windowsPaths(process.argv.slice(2)) : process.argv.slice(2);
  const go = new Go();
  go.argv = ["tidgrid", ...args];
  go.env = { TMPDIR: os.tmpdir(), ...process.env };
  // Not process.exit: let the event loop drain so stdout and stderr arrive whole.
  go.exit = (code) => {
    process.exitCode = code;
  };
  WebAssembly.instantiate(fs.readFileSync(path.join(root, "tidgrid.wasm")), go.importObject)
    .then((result) => {
      process.on("exit", (code) => {
        if (code === 0 && !go.exited) {
          // deadlock: make Go print its error and stack traces
          go._pendingEvent = { id: 0 };
          go._resume();
        }
      });
      return go.run(result.instance);
    })
    .catch((err) => {
      console.error(err);
      process.exitCode = 1;
    });
}

main();
