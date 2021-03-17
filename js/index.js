import {
  getCurrentQueryParam,
  replaceStateSetQueryParam,
  getLocalStorage,
  setLocalStorage,
} from "./utils";

const CONTRACT_FILTER_QUERY_PARAM_KEY = "c";
const CONTRACT_FILTER_LOCALSTORAGE_KEY = "contract_filter";

(() => {
  const contractFilterSelect = document.getElementById("contract_filter");
  if (contractFilterSelect) {
    const contractFilterOnChange = () => {
      const filter = contractFilterSelect.value;
      replaceStateSetQueryParam(CONTRACT_FILTER_QUERY_PARAM_KEY, filter === "" ? null : filter);
      for (const el of document.querySelectorAll("[data-contract]")) {
        const contractId = el.dataset.contract;
        const show = filter === "" ? true : contractId === filter;
        el.style.display = show ? "" : "none";
      }
      setLocalStorage(CONTRACT_FILTER_LOCALSTORAGE_KEY, filter);
    };
    contractFilterSelect.addEventListener("change", contractFilterOnChange);

    // Set initial filter.
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
    contractFilterOnChange();
  }
})();
