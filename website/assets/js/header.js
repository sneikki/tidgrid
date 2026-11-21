const headerNavigation = document.getElementById("header-navigation");
const headerNavigationToggler = document.getElementById("toggle-header-navigation");
const sidebarCloseButton = document.getElementById("close-sidebar");

document.body.addEventListener("click", ({ target }) => {
  const isOpen = headerNavigation.classList.contains("tg-state(open)");

  if (isOpen && !headerNavigation.contains(target)) {
    headerNavigation.classList.remove("tg-state(open)");
  } else if (headerNavigationToggler.contains(target)) {
    headerNavigation.classList.add("tg-state(open)");
  }
});

sidebarCloseButton.addEventListener("click", () =>
  headerNavigation.classList.remove("tg-state(open)")
);
