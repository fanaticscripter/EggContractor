(() => {
  function formatCountdown(seconds) {
    if (seconds <= 0) {
      return "0:00:00";
    }
    const hh = Math.floor(seconds / 3600);
    seconds -= hh * 3600;
    const mm = Math.floor(seconds / 60);
    seconds -= mm * 60;
    const ss = Math.floor(seconds);
    return `${hh}:${mm.toString().padStart(2, "0")}:${ss.toString().padStart(2, "0")}`;
  }

  function updateCountdowns() {
    for (const row of document.querySelectorAll("tr.EventTable__event--active")) {
      const countdownCell = row.querySelector("td:last-child");
      const secondsLeft = parseFloat(row.dataset.expires) - new Date().getTime() / 1000;
      if (secondsLeft <= 0) {
        row.classList.add("EventTable__event--expired");
        row.classList.remove("EventTable__event--active");
        countdownCell.textContent = "-";
      } else {
        countdownCell.textContent = formatCountdown(secondsLeft);
      }
    }
  }

  updateCountdowns();
  setInterval(updateCountdowns, 1000);
})();
