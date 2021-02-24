import { Build, Artifact, Stone } from "./models";

/**
 * @param {!Artifact} artifact
 * @param {!Stone} stone
 * @returns {Number}
 */
function stoneSettingCost(artifact, stone) {
  return Math.floor(artifact.base_crafting_price * 0.05 + stone.base_crafting_price * 0.1);
}

/**
 * @param {!Build} build
 * @returns {Number}
 */
function aggregateStoneSettingCost(build) {
  let sum = 0;
  for (const artifact of build.artifacts) {
    if (artifact.isEmpty()) {
      continue;
    }
    for (const stone of artifact.activeStones) {
      if (stone === null) {
        continue;
      }
      sum += stoneSettingCost(artifact, stone);
    }
  }
  return sum;
}

export { stoneSettingCost, aggregateStoneSettingCost };
