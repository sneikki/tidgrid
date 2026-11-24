"use strict";
// Installs the packed tarball into a temp project and checks that every
// command prints, exits and writes the same as the native binary built from
// this commit. Run after `node build.js`.
const { spawnSync } = require("child_process");
const fs = require("fs");
const os = require("os");
const path = require("path");

const repo = path.join(__dirname, "..");
const tmp = fs.mkdtempSync(path.join(os.tmpdir(), "tidgrid-smoke-"));
const win = process.platform === "win32";
const exe = win ? ".exe" : "";
const native = path.join(tmp, "native" + exe);
const toml = fs.readFileSync(path.join(repo, "internal/vocab/vocabtest/example.toml"));

function sh(cmd, args, opts = {}) {
  const r = spawnSync(cmd, args, { encoding: "utf8", shell: process.platform === "win32" && cmd === "npm", ...opts });
  if (r.error) throw r.error;
  return r;
}

function setup() {
  const build = sh("go", ["build", "-o", native, "./cmd/tidgrid"], { cwd: repo });
  if (build.status !== 0) throw new Error(build.stderr);
  const pack = sh("npm", ["pack", "--pack-destination", tmp], { cwd: __dirname });
  if (pack.status !== 0) throw new Error(pack.stderr);
  const tarball = path.join(tmp, pack.stdout.trim().split("\n").pop());
  const proj = path.join(tmp, "install");
  fs.mkdirSync(proj);
  fs.writeFileSync(path.join(proj, "package.json"), "{}");
  const inst = sh("npm", ["install", "--no-audit", "--no-fund", tarball], { cwd: proj });
  if (inst.status !== 0) throw new Error(inst.stderr);
  return path.join(proj, "node_modules", "tidgrid", "bin", "tidgrid.js");
}

// Each scenario starts from `files`; a step is the argument list, with {dir}
// replaced by the project directory.
const badPage = { "src/page.html": '<div class="tg-nonsense">x</div>\n' };
const scenarios = [
  {
    name: "tg project",
    files: {},
    steps: [
      ["init"], ["init"], ["build"], ["check"], ["list"], ["plan"], ["sheet"], ["sheet", "-o", "s.html"],
      ["check", "{dir}/index.tg"], ["build", "./index.tg"], ["nonsense"], [], ["sheet", "-o={dir}/s2.html"],
    ],
    after: badPage,
  },
  { name: "no vocabulary", files: {}, steps: [["check"], ["build"], ["list"]] },
  // a long failing check leaves a Go timer pending when the tool exits
  {
    name: "many unknown classes",
    files: {},
    steps: [["init"], ["build"], ["build"], ["check"]],
    after: { "src/page.html": Array.from({ length: 10000 }, (_, i) => `<div class="tg-x${i}">x</div>\n`).join("") },
  },
  { name: "old toml", files: { "tidgrid.toml": toml }, steps: [["check"], ["convert"], ["convert"], ["build"], ["check"]] },
];

function snapshot(dir) {
  const out = {};
  (function walk(d) {
    for (const e of fs.readdirSync(d, { withFileTypes: true })) {
      const p = path.join(d, e.name);
      if (e.isDirectory()) walk(p);
      else out[path.relative(dir, p).replace(/\\/g, "/")] = fs.readFileSync(p, "utf8");
    }
  })(dir);
  return out;
}

// runs one scenario with `launch(args, dir)` and returns everything observable
function play(sc, name, launch) {
  const dir = path.join(tmp, name);
  fs.mkdirSync(dir);
  for (const [f, c] of Object.entries(sc.files)) fs.writeFileSync(path.join(dir, f), c);
  const norm = (s) => {
    s = s.split(dir).join("<dir>").split(dir.replace(/\\/g, "/")).join("<dir>").replace(/\r\n/g, "\n");
    // native Windows prints "\" in paths, wasm "/"; the launcher hands the wasm
    // paths relative to cwd, so it echoes "s2.html" where native echoes "<dir>/s2.html"
    return win ? s.replace(/\\/g, "/").replace("wrote <dir>/", "wrote ") : s;
  };
  const log = [];
  sc.steps.forEach((step, i) => {
    if (i === 3 && sc.after) {
      for (const [f, c] of Object.entries(sc.after)) {
        fs.mkdirSync(path.dirname(path.join(dir, f)), { recursive: true });
        fs.writeFileSync(path.join(dir, f), c);
      }
    }
    const r = launch(step.map((a) => a.replace("{dir}", dir)), dir);
    log.push({ step, code: r.status, out: norm(r.stdout), err: norm(r.stderr) });
  });
  return { log, files: snapshot(dir) };
}

// prints the first step and the first file that differ, native then npm
function report(a, b) {
  const cut = (s) => String(JSON.stringify(s)).slice(0, 400);
  const i = a.log.findIndex((l, n) => JSON.stringify(l) !== JSON.stringify(b.log[n]));
  if (i >= 0) console.log(`step ${i}\n go: ${cut(a.log[i])}\n npm: ${cut(b.log[i])}`);
  const f = Object.keys({ ...a.files, ...b.files }).find((k) => a.files[k] !== b.files[k]);
  if (f) console.log(`file ${f}\n go: ${cut(a.files[f])}\n npm: ${cut(b.files[f])}`);
}

function main() {
  const wasm = setup();
  let failed = 0;
  for (const sc of scenarios) {
    const a = play(sc, sc.name.replace(/\W/g, "") + "-go", (args, cwd) => sh(native, args, { cwd }));
    const b = play(sc, sc.name.replace(/\W/g, "") + "-npm", (args, cwd) => sh(process.execPath, [wasm, ...args], { cwd }));
    const same = JSON.stringify(a) === JSON.stringify(b);
    console.log(`${same ? "ok  " : "FAIL"} ${sc.name}: ${a.log.map((l) => l.code).join(" ")}`);
    if (!same) {
      failed++;
      report(a, b);
    }
  }
  // the real bin shim, as `npx tidgrid` runs it
  const proj = path.dirname(path.dirname(path.dirname(path.dirname(wasm))));
  const viaNpx = sh("npx", ["--no-install", "tidgrid", "check", "nowhere.tg"], { cwd: proj, shell: process.platform === "win32" });
  const npxOk = viaNpx.status === 2;
  console.log(`${npxOk ? "ok  " : "FAIL"} npx shim exits 2 (got ${viaNpx.status})`);
  if (!npxOk) {
    failed++;
    console.log(viaNpx.stdout + viaNpx.stderr);
  }
  // a reader that quits early must not make the launcher print anything
  if (!win) {
    const dir = path.join(tmp, "manyunknownclasses-npm");
    const head = sh("sh", ["-c", `node "${wasm}" check 2>err.txt | head -1 >/dev/null; cat err.txt`], { cwd: dir });
    const headOk = head.stdout === "" && head.stderr === "";
    console.log(`${headOk ? "ok  " : "FAIL"} no stderr into | head -1`);
    if (!headOk) {
      failed++;
      console.log(head.stdout + head.stderr);
    }
  }
  fs.rmSync(tmp, { recursive: true, force: true });
  process.exit(failed ? 1 : 0);
}

main();
