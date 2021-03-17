import {
  getCurrentQueryParam,
  replaceStateSetQueryParam,
  getLocalStorage,
  setLocalStorage,
} from "./utils";

const CONTRACT_FILTER_QUERY_PARAM_KEY = "c";
const CONTRACT_FILTER_LOCALSTORAGE_KEY = "contract_filter";
const HIDE_SOLOS_QUERY_PARAM_KEY = "hide_solos";
const HIDE_SOLOS_LOCALSTORAGE_KEY = "hide_solos";

(() => {
  let contractFilter = "";
  let hideSolos = false;

  const showHideStatuses = () => {
    for (const el of document.querySelectorAll("[data-contract]")) {
      const contractId = el.dataset.contract;
      const type = el.dataset.type;
      const show =
        (contractFilter === "" ? true : contractId === contractFilter) &&
        (hideSolos ? type !== "solo" : true);
      el.style.display = show ? "" : "none";
    }
  };

  const contractFilterSelect = document.getElementById("contract_filter");
  if (contractFilterSelect) {
    const updateContractFilter = () => {
      contractFilter = contractFilterSelect.value;
      replaceStateSetQueryParam(
        CONTRACT_FILTER_QUERY_PARAM_KEY,
        contractFilter === "" ? null : contractFilter
      );
      setLocalStorage(CONTRACT_FILTER_LOCALSTORAGE_KEY, contractFilter);
    };

    contractFilterSelect.addEventListener("change", () => {
      updateContractFilter();
      showHideStatuses();
    });

    // Set initial filter.
    {
      let filter = getCurrentQueryParam(CONTRACT_FILTER_QUERY_PARAM_KEY);
      if (filter === null) {
        filter = getLocalStorage(CONTRACT_FILTER_LOCALSTORAGE_KEY);
        if (filter === undefined) {
          filter = "";
        }
        // Make sure the cached value is still valid, otherwise clear it.
        const optionValues = Array.prototype.map.call(
          contractFilterSelect.options,
          option => option.value
        );
        if (!optionValues.includes(filter)) {
          filter = "";
        }
      }
      contractFilterSelect.value = filter;
      updateContractFilter();
    }
  }

  const hideSolosCheckbox = document.getElementById("hide_solos");
  if (hideSolosCheckbox) {
    const updateHideSolos = () => {
      hideSolos = hideSolosCheckbox.checked;
      replaceStateSetQueryParam(HIDE_SOLOS_QUERY_PARAM_KEY, hideSolos ? 1 : null);
      setLocalStorage(HIDE_SOLOS_LOCALSTORAGE_KEY, hideSolos);
    };

    hideSolosCheckbox.addEventListener("change", () => {
      updateHideSolos();
      showHideStatuses();
    });

    // Set initial value.
    {
      let hide = getCurrentQueryParam(HIDE_SOLOS_QUERY_PARAM_KEY);
      if (hide === null) {
        hide = getLocalStorage(HIDE_SOLOS_LOCALSTORAGE_KEY) === "true";
      } else {
        hide = !!hide;
      }
      hideSolosCheckbox.checked = hide;
      updateHideSolos();
    }
  }

  showHideStatuses();
})();
