(() => {
  const navMenuToggle = document.getElementById("nav__menu-toggle");
  const navMenuToggleClosedIcon = document.getElementById("nav__menu-toggle-closed-icon");
  const navMenuToggleOpenIcon = document.getElementById("nav__menu-toggle-open-icon");
  const navMenu = document.getElementById("nav__menu");
  navMenuToggle?.addEventListener("click", () => {
    navMenuToggleClosedIcon?.classList.toggle("hidden");
    navMenuToggleOpenIcon?.classList.toggle("hidden");
    navMenu?.classList.toggle("hidden");
  });
})();
