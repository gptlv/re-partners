(() => {
  const flashEl = document.getElementById("app-flash");
  const calculateForm = document.getElementById("calculate-form");
  const calculateButton = document.getElementById("calculate-button");
  const orderAmountInput = document.getElementById("order-amount");
  const resultsEmptyEl = document.getElementById("results-empty");
  const resultsContainerEl = document.getElementById("results-container");
  const resultsBodyEl = document.getElementById("results-table-body");
  const resultsRequestedEl = document.getElementById("results-requested");
  const resultsPackedEl = document.getElementById("results-packed");
  const addPackForm = document.getElementById("add-pack-form");
  const addPackButton = addPackForm.querySelector("button[type='submit']");
  const newPackInput = document.getElementById("new-pack-size");
  const packsListEl = document.getElementById("packs-list");
  const packsEmptyEl = document.getElementById("packs-empty");

  const state = {
    packs: [],
    lastCalculation: null,
  };

  let flashTimer = null;

  const formatNumber = (value) =>
    new Intl.NumberFormat(undefined, {
      maximumFractionDigits: 0,
    }).format(value);

  function hideFlash() {
    if (!flashEl) {
      return;
    }
    flashEl.hidden = true;
    flashEl.className = "flash";
    flashEl.textContent = "";
    if (flashTimer) {
      clearTimeout(flashTimer);
      flashTimer = null;
    }
  }

  function showFlash(variant, message, options = {}) {
    if (!flashEl) {
      return;
    }
    const { persist = false } = options;
    if (flashTimer) {
      clearTimeout(flashTimer);
      flashTimer = null;
    }
    flashEl.textContent = message;
    flashEl.className = `flash show ${variant}`;
    flashEl.hidden = false;
    if (!persist) {
      flashTimer = setTimeout(hideFlash, 5000);
    }
  }

  function setButtonBusy(button, busy, busyLabel) {
    if (!button) {
      return;
    }
    if (busy) {
      if (!button.dataset.originalText) {
        button.dataset.originalText = button.textContent;
      }
      if (busyLabel) {
        button.textContent = busyLabel;
      }
      button.disabled = true;
      return;
    }

    if (button.dataset.originalText) {
      button.textContent = button.dataset.originalText;
      delete button.dataset.originalText;
    }
    button.disabled = false;
  }

  function renderPackList() {
    if (!packsListEl || !packsEmptyEl) {
      return;
    }

    packsListEl.innerHTML = "";

    if (!state.packs.length) {
      packsListEl.hidden = true;
      packsEmptyEl.hidden = false;
      return;
    }

    state.packs
      .slice()
      .sort((a, b) => a.size - b.size)
      .forEach((pack) => {
        const item = document.createElement("li");
        item.className = "list-item";

        const label = document.createElement("span");
        label.textContent = `${formatNumber(pack.size)} items`;

        const button = document.createElement("button");
        button.type = "button";
        button.className = "btn danger";
        button.dataset.action = "delete-pack";
        button.dataset.packId = String(pack.id);
        button.dataset.packSize = String(pack.size);
        button.textContent = "Delete";

        item.append(label, button);
        packsListEl.appendChild(item);
      });

    packsListEl.hidden = false;
    packsEmptyEl.hidden = true;
  }

  function renderResults(calculation) {
    if (!resultsBodyEl || !resultsContainerEl || !resultsEmptyEl) {
      return;
    }

    if (!calculation || !Array.isArray(calculation.packs) || !calculation.packs.length) {
      resultsBodyEl.innerHTML = "";
      resultsContainerEl.hidden = true;
      resultsEmptyEl.hidden = false;
      return;
    }

    resultsBodyEl.innerHTML = "";
    calculation.packs.forEach((pack) => {
      const row = document.createElement("tr");
      const sizeCell = document.createElement("td");
      sizeCell.textContent = `${formatNumber(pack.size)} items`;

      const countCell = document.createElement("td");
      countCell.textContent = formatNumber(pack.count);

      row.append(sizeCell, countCell);
      resultsBodyEl.appendChild(row);
    });

    const requested = Number(calculation.amount) || 0;
    const packed = calculation.packs.reduce(
      (total, pack) => total + Number(pack.size || 0) * Number(pack.count || 0),
      0,
    );

    if (resultsRequestedEl) {
      resultsRequestedEl.textContent = formatNumber(requested);
    }
    if (resultsPackedEl) {
      resultsPackedEl.textContent = formatNumber(packed);
    }

    resultsEmptyEl.hidden = true;
    resultsContainerEl.hidden = false;
  }

  async function fetchPacks() {
    try {
      const response = await fetch("/api/packs", {
        headers: {
          Accept: "application/json",
        },
      });
      if (!response.ok) {
        const message = (await response.text()).trim() || "Failed to load pack sizes.";
        throw new Error(message);
      }

      const data = await response.json();
      state.packs = Array.isArray(data.packs) ? data.packs : [];
      renderPackList();
    } catch (error) {
      console.error(error);
      showFlash("error", error.message || "Unable to load pack sizes. Please try again.");
    }
  }

  async function handleCalculate(event) {
    event.preventDefault();
    hideFlash();

    const amount = Number(orderAmountInput.value);
    if (!Number.isInteger(amount) || amount <= 0) {
      showFlash("error", "Enter an amount greater than zero.");
      orderAmountInput.focus();
      return;
    }

    setButtonBusy(calculateButton, true, "Calculating…");

    try {
      const response = await fetch("/api/calc", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({ amount }),
      });

      if (!response.ok) {
        const message = (await response.text()).trim() || "Calculation failed.";
        throw new Error(message);
      }

      const data = await response.json();
      state.lastCalculation = data;
      renderResults(data);
      showFlash("success", "Calculation complete.");
    } catch (error) {
      console.error(error);
      renderResults(null);
      showFlash("error", error.message || "Calculation failed. Please try again.");
    } finally {
      setButtonBusy(calculateButton, false);
    }
  }

  async function handleAddPack(event) {
    event.preventDefault();
    hideFlash();

    const size = Number(newPackInput.value);

    if (!Number.isInteger(size) || size <= 0) {
      showFlash("error", "Pack size must be a positive whole number.");
      newPackInput.focus();
      return;
    }

    setButtonBusy(addPackButton, true, "Adding…");

    try {
      const response = await fetch("/api/sizes", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify({ size }),
      });

      if (!response.ok) {
        const message = (await response.text()).trim() || "Failed to add pack size.";
        throw new Error(message);
      }

      const created = await response.json();
      state.packs.push(created);
      renderPackList();

      addPackForm.reset();
      newPackInput.focus();
      showFlash("success", `Added a ${formatNumber(size)} item pack.`);
    } catch (error) {
      console.error(error);
      showFlash("error", error.message || "Failed to add pack size. Please try again.");
    } finally {
      setButtonBusy(addPackButton, false);
    }
  }

  async function handleDeletePack(button) {
    const packId = Number(button.dataset.packId);
    const packSize = Number(button.dataset.packSize);

    if (!Number.isInteger(packId) || packId <= 0) {
      showFlash("error", "Invalid pack identifier.");
      return;
    }

    setButtonBusy(button, true, "Removing…");

    try {
      const response = await fetch(`/api/sizes/${packId}`, {
        method: "DELETE",
      });

      if (!response.ok) {
        const message = (await response.text()).trim() || "Failed to delete pack size.";
        throw new Error(message);
      }

      state.packs = state.packs.filter((pack) => pack.id !== packId);
      renderPackList();
      showFlash("success", `Removed the ${formatNumber(packSize)} item pack.`);
    } catch (error) {
      console.error(error);
      showFlash("error", error.message || "Failed to delete pack size.");
    } finally {
      setButtonBusy(button, false);
    }
  }

  function registerEventListeners() {
    if (calculateForm) {
      calculateForm.addEventListener("submit", handleCalculate);
    }

    if (addPackForm) {
      addPackForm.addEventListener("submit", handleAddPack);
    }

    if (packsListEl) {
      packsListEl.addEventListener("click", (event) => {
        const target = event.target;
        if (!(target instanceof HTMLElement)) {
          return;
        }
        const button = target.closest("[data-action='delete-pack']");
        if (button instanceof HTMLButtonElement) {
          handleDeletePack(button);
        }
      });
    }
  }

  function init() {
    registerEventListeners();
    fetchPacks();
  }

  document.addEventListener("DOMContentLoaded", init);
})();
