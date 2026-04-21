// wait for the page to fully load
document.addEventListener("DOMContentLoaded", main);

function main() {
  configureSort();
  configureFilter();

  sortTable("file", "asc");
  updateHeaderClasses(
    Array.from(document.querySelectorAll("th[data-col]")),
    document.querySelector("th[data-col='file']"),
    "asc",
  );
}

function configureSort() {
  let headers = Array.from(document.querySelectorAll("th[data-col]"));

  headers.forEach((header) => {
    header.addEventListener("click", () => {
      let col = header.getAttribute("data-col");
      let direction = header.classList.contains("asc") ? "desc" : "asc";

      sortTable(col, direction);
      updateHeaderClasses(headers, header, direction);
    });
  });
}

function configureFilter() {
  let filterInput = document.getElementById("filter-input");
  filterInput.addEventListener("input", () => {
    filterTable(filterInput.value);
  });

  document.addEventListener("keydown", function (event) {
    const active = document.activeElement;

    const isTyping =
      active.tagName === "INPUT" ||
      active.tagName === "TEXTAREA" ||
      active.isContentEditable;

    if (event.key === "/" && !isTyping) {
      event.preventDefault();
      filterInput?.focus();
    }
  });
}

function sortTable(col, direction) {
  let table = document.querySelector("table.coverage-summary");
  let rows = Array.from(table.querySelectorAll("tr")).slice(1); // skip header row

  rows.sort((a, b) => {
    let aValue = a.querySelector(`td[data-col="${col}"]`).textContent.trim();
    let bValue = b.querySelector(`td[data-col="${col}"]`).textContent.trim();

    if (col === "file") {
      return direction === "asc"
        ? aValue.localeCompare(bValue)
        : bValue.localeCompare(aValue);
    } else {
      return direction === "asc"
        ? parseFloat(aValue) - parseFloat(bValue)
        : parseFloat(bValue) - parseFloat(aValue);
    }
  });

  rows.forEach((row) => table.appendChild(row)); // re-append sorted rows
}

function filterTable(query) {
  let table = document.querySelector("table.coverage-summary");
  let rows = Array.from(table.querySelectorAll("tr")).slice(1); // skip header row

  rows.forEach((row) => {
    let fileName = row.querySelector(`td[data-col="file"]`).textContent.trim();
    if (fileName.toLowerCase().includes(query.toLowerCase())) {
      row.style.display = "";
    } else {
      row.style.display = "none";
    }
  });
}

function updateHeaderClasses(headers, activeHeader, direction) {
  headers.forEach((header) => {
    header.classList.remove("asc", "desc");
    if (header === activeHeader) {
      header.classList.add(direction);
    }
  });
}
