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

const taken = new Set();
let size = 0; // one cell in px
let travelled = 0; // px the grid has moved up

function follow() {
  size = parseFloat(getComputedStyle(grid).backgroundSize);
  travelled = scrollY * SHARE;
  const within = travelled % size;
  grid.style.transform = `translateY(${-within}px)`;
  lit.style.transform = `translateY(${within - travelled}px)`;
}

function pick(list) {
  return list[Math.floor(Math.random() * list.length)];
}

// The header, what the first screen says, and everything after it: no cell lights up behind these.
const KEEP_CLEAR = "header, [data-clear] > *, section ~ section, footer";

// Where a cell would sit on screen, or null when something to keep clear or another cell is there.
function freePlace(busy) {
  const column = Math.floor(Math.random() * Math.ceil(innerWidth / size));
  const row = Math.floor((travelled + Math.random() * innerHeight) / size);
  const key = `${column}:${row}`;
  if (taken.has(key)) return null;

  const left = column * size;
  const top = row * size - travelled;
  const covered = busy.some((box) => box.left < left + size && box.right > left && box.top < top + size && box.bottom > top);
  return covered ? null : { key, left, top: row * size };
}

function light() {
  if (document.hidden || scrollY > innerHeight) return;
  const busy = [...document.querySelectorAll(KEEP_CLEAR)].map((element) => element.getBoundingClientRect());

  for (let attempt = 0; attempt < TRIES; attempt++) {
    const place = freePlace(busy);
    if (!place) continue;

    const foreign = Math.random() < FOREIGN;
    const cell = document.createElement("div");
    cell.className = "tg-surface(cell)";
    cell.classList.toggle("tg-state(foreign)", foreign);
    cell.textContent = pick(foreign ? OUTSIDE : OWN);
    cell.style.left = `${place.left}px`;
    cell.style.top = `${place.top}px`;
    cell.addEventListener("animationend", () => {
      cell.remove();
      taken.delete(place.key);
    });
    taken.add(place.key);
    lit.append(cell);
    return;
  }
}

if (!matchMedia("(prefers-reduced-motion: reduce)").matches) {
  follow();
  addEventListener("scroll", follow, { passive: true });
  addEventListener("resize", follow);
  if (grid.hasAttribute("data-cells")) setInterval(light, EVERY);
}
