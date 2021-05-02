import { Build, Config } from "../models";
import { multiplicativeEffect } from "./common";
import { maxHabSpace, baseMaxHabSpace, habSpaceMultiplier } from "./hab_space";

/**
 * Laying rate multiplier per chicken.
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function layingRateMultiplier(build, config) {
  return (
    multiplicativeEffect(build, config, [
      proto.ArtifactSpec.Name.QUANTUM_METRONOME,
      proto.ArtifactSpec.Name.TACHYON_STONE,
    ]) *
    (1 + config.tachyonDeflectorBonus)
  );
}

/**
 * Multiplier of max laying rate, taking hab space increase into account.
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function maxLayingRateMultiplier(build, config) {
  return layingRateMultiplier(build, config) * habSpaceMultiplier(build, config);
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function maxHourlyLayingRate(build, config) {
  return (
    baseMaxHourlyLayingRatePerChicken(config) *
    layingRateMultiplier(build, config) *
    maxHabSpace(build, config)
  );
}

/**
 * @param {!Config} config
 * @returns {!Number}
 */
function baseMaxHourlyLayingRate(config) {
  return baseMaxHourlyLayingRatePerChicken(config) * baseMaxHabSpace(config);
}

/**
 * @param {!Config} config
 * @returns {!Number}
 */
function baseMaxHourlyLayingRatePerChicken(config) {
  // Base rate: 1 egg per 30 seconds
  // Affected by the following researches:
  // {id, perLevel, maxLevels}
  // {"comfy_nests", 0.10, 50},
  // {"hen_house_ac", 0.05, 50},
  // {"improved_genetics", 0.15, 30},
  // {"time_compress", 0.10, 20},
  // {"timeline_diversion", 0.02, 50},
  // {"relativity_optimization", 0.10, 10},
  // {"epic_egg_laying", 0.05, 20},
  return (
    (1 / 30) *
    (1 + 0.1 * 50) *
    (1 + 0.05 * 50) *
    (1 + 0.15 * 30) *
    (1 + 0.1 * 20) *
    (1 + 0.02 * 50) *
    (1 + 0.1 * 10) *
    (1 + 0.05 * 20) *
    3600
  );
}

export { layingRateMultiplier, maxLayingRateMultiplier, maxHourlyLayingRate };
