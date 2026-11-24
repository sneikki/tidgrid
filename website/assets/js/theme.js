// Loaded in <head> without defer, so a stored choice applies before the first paint.
// Without a stored choice the page follows the system setting (color-scheme: light dark).
const root = document.documentElement;

function applyScheme(scheme) {
  root.classList.toggle("tg-state(dark)", scheme === "dark");
  root.classList.toggle("tg-state(light)", scheme === "light");
}

function storedScheme() {
  try {
    return localStorage.getItem("color-scheme");
  } catch {
    return null; // storage blocked: the page still follows the system
  }
}

function currentScheme() {
  if (root.classList.contains("tg-state(dark)")) return "dark";
  if (root.classList.contains("tg-state(light)")) return "light";
  return matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
}

applyScheme(storedScheme());

document.addEventListener("click", ({ target }) => {
  if (!target.closest("[data-theme-toggle]")) return;
  const next = currentScheme() === "dark" ? "light" : "dark";
  applyScheme(next);
  try {
    localStorage.setItem("color-scheme", next);
  } catch {
    // choice lasts for this page only
  }
});
