const headerNavigation = document.getElementById("header-navigation");
const headerNavigationToggler = document.getElementById("toggle-header-navigation");
const sidebarCloseButton = document.getElementById("close-sidebar");

document.body.addEventListener("click", ({ target }) => {
  if (!(
    headerNavigation.contains(target) ||
    headerNavigation.classList.contains("header__navigation--collapsed")
  )) {
    headerNavigation.classList.add("header__navigation--collapsed")
  } else if (headerNavigationToggler.contains(target)) {
    headerNavigation.classList.remove("header__navigation--collapsed")
  }
});

sidebarCloseButton.addEventListener("click", () =>
  headerNavigation.classList.add("header__navigation--collapsed")
);
