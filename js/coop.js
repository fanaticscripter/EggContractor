import { createApp } from "vue";
import CoopTable from "./CoopTable.vue";

for (const table of document.querySelectorAll(".CoopTable")) {
  createApp(CoopTable, {
    contractId: table.dataset.contract,
    code: table.dataset.code,
    members: JSON.parse(table.dataset.members),
  }).mount(table);
}
