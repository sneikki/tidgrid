"use strict";
// Builds the package contents next to this file: the wasm binary, Go's JS
// glue with its licence, and the repo's README and licence.
const { execFileSync } = require("child_process");
const fs = require("fs");
const os = require("os");
const path = require("path");

const here = __dirname;
const repo = path.join(here, "..");
const run = (cmd, args, opts) =>
  execFileSync(cmd, args, { encoding: "utf8", stdio: ["ignore", "pipe", "inherit"], ...opts }).trim();

// A commit that is exactly a release tag must carry that tag's number.
function checkVersion() {
  const version = require("./package.json").version;
  let tag;
  try {
    tag = run("git", ["describe", "--exact-match", "--tags", "HEAD"], { cwd: repo, stdio: ["ignore", "pipe", "ignore"] });
  } catch {
    return;
  }
  if (tag !== "v" + version) throw new Error(`HEAD is tagged ${tag}, npm/package.json says ${version}`);
}

function buildWasm() {
  const out = path.join(here, "tidgrid.wasm");
  run("go", ["build", "-trimpath", "-ldflags=-s -w", "-o", out, "./cmd/tidgrid"], {
    cwd: repo,
    env: { ...process.env, GOOS: "js", GOARCH: "wasm", CGO_ENABLED: "0" },
  });
  // The wasm is published: it must not carry this machine's paths or user.
  const wasm = fs.readFileSync(out);
  const home = os.homedir();
  const user = os.userInfo().username;
  const leaks = [home, home.replace(/\\/g, "/"), `/${user}/`, `\\${user}\\`].filter((s) => s.length > 3 && wasm.includes(s));
  if (leaks.length) throw new Error(`tidgrid.wasm contains builder paths: ${leaks.join(", ")}`);
}

function copyGoGlue() {
  const goroot = run("go", ["env", "GOROOT"]);
  fs.copyFileSync(path.join(goroot, "lib", "wasm", "wasm_exec.js"), path.join(here, "wasm_exec.js"));
  // Arch's package keeps the licence outside GOROOT.
  const licence = [path.join(goroot, "LICENSE"), "/usr/share/licenses/go/LICENSE"].find(fs.existsSync);
  if (!licence) throw new Error("Go LICENSE not found");
  fs.copyFileSync(licence, path.join(here, "GO-LICENSE"));
}

checkVersion();
buildWasm();
copyGoGlue();
fs.copyFileSync(path.join(repo, "README.md"), path.join(here, "README.md"));
fs.copyFileSync(path.join(repo, "LICENSE.md"), path.join(here, "LICENSE.md"));
