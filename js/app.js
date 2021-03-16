import tippy from "tippy.js";
import "tippy.js/dist/tippy.css";
import "tippy.js/themes/translucent.css";

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

  for (const el of document.querySelectorAll("[data-tooltip]")) {
    let content = el.dataset.tooltip;
    if (!content) {
      content = el.title;
    }
    if (!content) {
      continue;
    }
    tippy(el, {
      content,
      theme: "translucent",
    });
    el.removeAttribute("title");
  }

  const autoRefreshToggle = document.getElementById("AutoRefreshToggle");
  if (autoRefreshToggle) {
    const key = `auto-refresh-${window.location.pathname}`;
    let refreshTimeout;

    function toggleAutoRefresh() {
      const currentlyOn = autoRefreshToggle.classList.contains("AutoRefreshToggle--on");
      autoRefreshToggle.classList.toggle("AutoRefreshToggle--on");
      if (currentlyOn) {
        if (refreshTimeout !== undefined) {
          clearTimeout(refreshTimeout);
          refreshTimeout = undefined;
        }
      } else {
        refreshTimeout = setTimeout(() => {
          window.location.reload();
        }, 60000);
      }
      localStorage[key] = !currentlyOn;
    }

    const savedSetting = localStorage[key];
    if (savedSetting === "true") {
      toggleAutoRefresh();
    }
    autoRefreshToggle.addEventListener("click", toggleAutoRefresh);
  }
})();
