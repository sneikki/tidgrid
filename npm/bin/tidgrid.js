#!/usr/bin/env node
"use strict";
// Runs tidgrid.wasm under Go's JS glue. Arguments, cwd and environment pass
// through; the process exit code is the tool's own.
const fs = require("fs");
const os = require("os");
const path = require("path");

const root = path.join(__dirname, "..");

// Go's js port only knows "/" paths. On Windows, hand it "C:/dir" for the
// working directory and forward-slash, cwd-relative paths for path arguments,
// also the value of `--flag=value`.
// Limit: a path on another drive than the working directory stays "D:/x",
// which Go reads as relative to the working directory (not supported).
function windowsPaths(args) {
  const cwd = process.cwd();
  process.cwd = () => cwd.replace(/\\/g, "/");
  return args.map((a) => {
    const [, flag = "", v] = /^(-[^=]*=)?([^]*)$/.exec(a);
    if ((!flag && v.startsWith("-")) || !/^[A-Za-z]:[\\/]|\\/.test(v)) return a;
    return flag + (path.relative(cwd, path.resolve(cwd, v)).replace(/\\/g, "/") || ".");
  });
}

// Node on Windows has no O_DIRECTORY, so Go refuses to open any directory.
// Offer Go a bit for it and strip that bit again before Node sees the flags.
function windowsFs() {
  const O_DIRECTORY = 0x10000;
  return Object.create(fs, {
    constants: { value: { ...fs.constants, O_DIRECTORY } },
    open: { value: (p, flags, ...rest) => fs.open(p, flags & ~O_DIRECTORY, ...rest) },
  });
}

// Go's glue fits argv and environment in 8 KiB. npx and CI runners exceed
// that, so shed the npm_* variables first, then the largest, until it fits.
function fitEnv(env, args) {
  const cost = (s) => Buffer.byteLength(s) + 1 + 7 + 8;
  const entries = Object.entries(env).map(([k, v]) => [k, v, cost(`${k}=${v}`)]);
  let total = args.reduce((n, a) => n + cost(a), 0) + entries.reduce((n, e) => n + e[2], 0);
  const shed = [...entries].sort((a, b) => (b[0].startsWith("npm_") - a[0].startsWith("npm_")) || b[2] - a[2]);
  const drop = new Set();
  for (const e of shed) {
    if (total <= 7000) break;
    drop.add(e[0]);
    total -= e[2];
  }
  if (total > 7000) return null; // the arguments alone do not fit
  return Object.fromEntries(entries.filter((e) => !drop.has(e[0])).map((e) => e.slice(0, 2)));
}

function main() {
  globalThis.require = require;
  globalThis.fs = process.platform === "win32" ? windowsFs() : fs;
  globalThis.path = path;
  globalThis.TextEncoder = require("util").TextEncoder;
  globalThis.TextDecoder = require("util").TextDecoder;
  globalThis.performance ??= require("perf_hooks").performance;
  globalThis.crypto ??= require("crypto");
  require("../wasm_exec");

  const args = process.platform === "win32" ? windowsPaths(process.argv.slice(2)) : process.argv.slice(2);
  const go = new Go();
  go.argv = ["tidgrid", ...args];
  go.env = fitEnv({ TMPDIR: os.tmpdir(), ...process.env }, go.argv);
  if (!go.env) {
    process.stderr.write("tidgrid: arguments are too long (limit 8 KiB)\n");
    process.exitCode = 2;
    return;
  }
  // Not process.exit: let the event loop drain so stdout and stderr arrive whole.
  // A Go timer still pending would resume the exited program and throw.
  go.exit = (code) => {
    for (const t of go._scheduledTimeouts.values()) clearTimeout(t);
    go._scheduledTimeouts.clear();
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
