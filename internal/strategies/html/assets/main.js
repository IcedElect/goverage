document.addEventListener("DOMContentLoaded", main);

function main() {
  configureDarkMode();
}

function configureDarkMode() {
  let darkMode = localStorage.getItem("dark");

  if (darkMode === "true") {
    setTheme("dark");
  } else if (darkMode === "false") {
    setTheme("light");
  } else {
    initDarkMode();
  }

  let switcher = document.querySelector("#theme-switcher");
  if (switcher) {
    switcher.addEventListener("click", () => {
      console.log("click");
      toggleDarkMode();
    });
  }
}

function initDarkMode() {
  let systemDarkMode = window.matchMedia(
    "(prefers-color-scheme: dark)"
  ).matches;
  if (systemDarkMode) {
    setTheme("dark");
  }

  window
    .matchMedia("(prefers-color-scheme: dark)")
    .addEventListener("change", (event) => {
      setTheme(event.matches ? "dark" : "light");
    });
}

function toggleDarkMode() {
  if (document.documentElement.classList.contains("dark")) {
    localStorage.setItem("dark", "false");
    setTheme("light");
  } else {
    localStorage.setItem("dark", "true");
    setTheme("dark");
  }
}

function setTheme(mode) {
  console.log("setTheme", mode);
  let lightStyle = document.querySelector('link[href*="github.min.css"]');
  let darkStyle = document.querySelector('link[href*="github-dark.min.css"]');

  if (mode === "dark") {
    document.documentElement.classList.add("dark");
    lightStyle.setAttribute("disabled", "disabled");
    darkStyle.removeAttribute("disabled");
  } else {
    document.documentElement.classList.remove("dark");
    lightStyle.removeAttribute("disabled");
    darkStyle.setAttribute("disabled", "disabled");
  }
}
