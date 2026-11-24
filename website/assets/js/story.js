// The front page story: scenes read by scrolling, each with a picture that builds up.
// A piece (data-rise) comes in once, when it scrolls into view; pieces that arrive together
// come one after another, so a terminal prints row by row. Without this script every piece
// is visible from the start. A module, so its names stay its own.
const root = document.documentElement;
const progress = document.getElementById("progress");
const story = document.getElementById("story");

const SHOW_AT = 0.86; // share of the screen height a piece's top must rise above
const APART = 110; // ms between two pieces that arrive together
const MOST = 8; // pieces after this many arrive with it: a fast scroll never waits long

let waiting = [...document.querySelectorAll("[data-rise]")];
let queued = false;

function show(element) {
  element.classList.add("tg-state(shown)");
}

// Pieces now in view, each a beat after the one that came into view with it.
function reveal() {
  const line = innerHeight * SHOW_AT;
  let together = 0;
  waiting = waiting.filter((element) => {
    if (element.getBoundingClientRect().top > line) return true;
    setTimeout(show, Math.min(together++, MOST) * APART, element);
    return false;
  });
}

// How far the story has been read, 0 to 1: from its top reaching mid-screen to the page's end.
function shareRead() {
  const start = story.getBoundingClientRect().top + scrollY - innerHeight / 2;
  const end = root.scrollHeight - innerHeight;
  return end > start ? Math.min(1, Math.max(0, (scrollY - start) / (end - start))) : 0;
}

function frame() {
  queued = false;
  reveal();
  progress.style.transform = `scaleX(${shareRead()})`;
}

function request() {
  if (queued) return;
  queued = true;
  requestAnimationFrame(frame);
}

if (matchMedia("(prefers-reduced-motion: reduce)").matches) {
  waiting.forEach(show);
  waiting = [];
}
root.classList.add("tg-state(scripted)");
addEventListener("scroll", request, { passive: true });
addEventListener("resize", request);
frame();
