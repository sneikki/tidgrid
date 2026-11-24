// A story: steps read by scrolling beside a pinned drawing. The step whose top has passed
// 45 % of the free screen is the one being read. The drawing shows the parts that exist at
// that step (data-show="1-4,7", steps counted from 0) and marks the ones it is about
// (data-hot). Blocks marked data-rise come in once. Without this script the page reads as
// plain sections: the drawing stays hidden and nothing is dimmed.
const root = document.documentElement;
const steps = [...document.querySelectorAll("[data-step]")];
const stage = document.getElementById("stage");
const badge = document.getElementById("badge");
const progress = document.getElementById("progress");
const parts = [...stage.querySelectorAll("[data-show], [data-hot]")].map((element) => ({
  element,
  show: element.hasAttribute("data-show") ? ranges(element.dataset.show) : null,
  hot: ranges(element.dataset.hot),
}));
const wide = matchMedia("(min-width: 75em)"); // breakpoint lg: the drawing sits beside the steps

const READ_AT = 0.45; // share of the free screen a step's top must pass
let rising = [...document.querySelectorAll("[data-rise]")];
let current = -1;
let queued = false;

// "1-4,7" -> [[1, 4], [7, 7]]
function ranges(text) {
  if (!text) return [];
  return text.split(",").map((part) => {
    const [from, to = from] = part.split("-").map(Number);
    return [from, to];
  });
}

function within(list, step) {
  return list.some(([from, to]) => step >= from && step <= to);
}

function show(step) {
  if (step === current) return;
  current = step;
  for (const { element, show, hot } of parts) {
    if (show) element.classList.toggle("tg-state(on)", within(show, step));
    element.classList.toggle("tg-state(hot)", within(hot, step));
  }
  steps.forEach((element, index) => element.classList.toggle("tg-state(active)", index === step));
  badge.textContent = steps[step].dataset.badge;
}

function stepBeingRead() {
  const pinnedAbove = wide.matches ? 0 : stage.getBoundingClientRect().height;
  const line = pinnedAbove + (innerHeight - pinnedAbove) * READ_AT;
  let step = 0;
  steps.forEach((element, index) => {
    if (element.getBoundingClientRect().top <= line) step = index;
  });
  return step;
}

// How far the story has been read, 0 to 1: from its first step reaching the line to the page's end.
function shareRead() {
  const start = steps[0].getBoundingClientRect().top + scrollY - innerHeight * READ_AT;
  const end = root.scrollHeight - innerHeight;
  return end > start ? Math.min(1, Math.max(0, (scrollY - start) / (end - start))) : 0;
}

function frame() {
  queued = false;
  const share = shareRead();
  show(share > 0.995 ? steps.length - 1 : stepBeingRead());
  progress.style.transform = `scaleX(${share})`;
  rising = rising.filter((element) => {
    if (element.getBoundingClientRect().top > innerHeight * 0.9) return true;
    element.classList.add("tg-state(shown)");
    return false;
  });
}

function request() {
  if (queued) return;
  queued = true;
  requestAnimationFrame(frame);
}

if (matchMedia("(prefers-reduced-motion: reduce)").matches) {
  rising.forEach((element) => element.classList.add("tg-state(shown)"));
  rising = [];
}
root.classList.add("tg-state(scripted)");
addEventListener("scroll", request, { passive: true });
addEventListener("resize", request);
frame();
