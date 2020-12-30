(() => {
  const peekerForm = document.querySelector("form#peeker");
  const contractIdInput = document.querySelector("input#contract-id");
  const coopCodeInput = document.querySelector("input#coop-code");
  const clearContractIdButton = document.querySelector("button#clear-contract-id");

  peekerForm?.addEventListener("submit", e => {
    e.preventDefault();
    let contractId = contractIdInput?.value.trim().toLowerCase();
    let coopCode = coopCodeInput?.value.trim().toLowerCase();
    contractId = contractId?.replace(/\(.*\)/g, "").trim();
    if (!contractId || !coopCode) {
      return;
    }
    const peekURL = `/peek/${encodeURIComponent(contractId)}/${encodeURIComponent(coopCode)}/`;
    window.open(peekURL, "_blank", "noopener");
  });

  clearContractIdButton?.addEventListener("click", e => {
    e.preventDefault();
    if (contractIdInput) {
      contractIdInput.value = "";
    }
  });
})();
