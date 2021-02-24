import { Build, Config } from "../models";
import { additiveEffect } from "./common";

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function maxRunningChickenBonusMultiplier(build, config) {
  const maxRCBBase = baseMaxRunningChickenBonus(config);
  const maxRCB = maxRunningChickenBonus(build, config);
  return maxRCB / maxRCBBase;
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function maxRunningChickenBonus(build, config) {
  return (
    baseMaxRunningChickenBonus(config) +
    additiveEffect(build, config, [
      proto.ArtifactSpec.Name.VIAL_MARTIAN_DUST,
      proto.ArtifactSpec.Name.TERRA_STONE,
    ])
  );
}

/**
 * @param {!Config} config
 * @returns {!Number}
 */
function baseMaxRunningChickenBonus(config) {
  // Assume max common and epic research.
  return 540;
}

export { maxRunningChickenBonus, maxRunningChickenBonusMultiplier };
