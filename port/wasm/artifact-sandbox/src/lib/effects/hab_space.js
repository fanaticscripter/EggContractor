import { Build, Config } from "../models";
import { multiplicativeEffect } from "./common";

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function maxHabSpace(build, config) {
  return Math.floor(baseMaxHabSpace(config) * habSpaceMultiplier(build, config));
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function habSpaceMultiplier(build, config) {
  return multiplicativeEffect(build, config, [proto.ArtifactSpec.Name.ORNATE_GUSSET]);
}

/**
 * @param {!Config} config
 * @returns {!Number}
 */
function baseMaxHabSpace(config) {
  // Assume max common research.
  return 1.134e10;
}

export { habSpaceMultiplier, maxHabSpace, baseMaxHabSpace };
