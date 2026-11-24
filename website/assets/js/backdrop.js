// The page grid. It travels slower than the page, and on the front page its cells
// light up with class names: the vocabulary's own, and now and then one from outside it.
// With reduced motion asked for, the grid stays still and empty.
const grid = document.getElementById("grid");
const lit = grid.firstElementChild;

const SHARE = 0.35; // of the page's scroll that the grid travels
const EVERY = 450; // ms between two cells lighting up
const FOREIGN = 0.3; // share of cells that show a name from outside the vocabulary
const TRIES = 8; // random places tried before a round is skipped

// Shown as text in the cells, never set as classes.
const OWN = ["tg-col", "tg-row", "tg-gap(md)", "tg-pad(sm)", "tg-grow", "tg-center", "tg-wrap", "tg-cols(3)", "tg-hide"];
const OUTSIDE = ["mt-[13px]", "border-l-4", ".btn-new", "#ff0066", "is-active", "p-[7px]"];

// The header, what the first screen says, and everything after it: no cell shows behind these.
const KEEP_CLEAR = "header, [data-clear] > *, section ~ section, footer";

const cells = new Map(); // "column:row" -> { element, left, top } in the lit layer's own coordinates
let size = 0; // one cell in px
let travelled = 0; // px the grid has moved up
let queued = false;

function measure() {
  size = parseFloat(getComputedStyle(grid).backgroundSize);
  follow();
}

function keptClear() {
  return [...document.querySelectorAll(KEEP_CLEAR)].map((element) => element.getBoundingClientRect());
}

// Whether a cell at this place of the lit layer is on screen behind one of the boxes.
function covered(boxes, left, top) {
  const onScreen = top - travelled;
  return boxes.some((box) => box.left < left + size && box.right > left && box.top < onScreen + size && box.bottom > onScreen);
}

function put(key) {
  cells.get(key)?.element.remove();
  cells.delete(key);
}

// The page moves faster than the grid, so text catches up with a lit cell: such a cell goes out.
function follow() {
  queued = false;
  const boxes = cells.size ? keptClear() : [];
  travelled = scrollY * SHARE;
  const within = travelled % size;
  grid.style.transform = `translateY(${-within}px)`;
  lit.style.transform = `translateY(${within - travelled}px)`;
  for (const [key, { left, top }] of cells) if (covered(boxes, left, top)) put(key);
}

function onScroll() {
  if (queued) return;
  queued = true;
  requestAnimationFrame(follow);
}

function pick(list) {
  return list[Math.floor(Math.random() * list.length)];
}

function freePlace(boxes) {
  const column = Math.floor(Math.random() * Math.ceil(innerWidth / size));
  const row = Math.floor((travelled + Math.random() * innerHeight) / size);
  const key = `${column}:${row}`;
  const left = column * size;
  const top = row * size;
  return cells.has(key) || covered(boxes, left, top) ? null : { key, left, top };
}

function light() {
  if (document.hidden || scrollY > innerHeight) return;
  const boxes = keptClear();

  for (let attempt = 0; attempt < TRIES; attempt++) {
    const place = freePlace(boxes);
    if (!place) continue;

    const foreign = Math.random() < FOREIGN;
    const element = document.createElement("div");
    element.className = "tg-surface(cell)";
    element.classList.toggle("tg-state(foreign)", foreign);
    element.textContent = pick(foreign ? OUTSIDE : OWN);
    element.style.left = `${place.left}px`;
    element.style.top = `${place.top}px`;
    element.addEventListener("animationend", () => put(place.key));
    cells.set(place.key, { element, left: place.left, top: place.top });
    lit.append(element);
    return;
  }
}

if (!matchMedia("(prefers-reduced-motion: reduce)").matches) {
  measure();
  addEventListener("scroll", onScroll, { passive: true });
  addEventListener("resize", measure);
  if (grid.hasAttribute("data-cells")) setInterval(light, EVERY);
}
