import { Build, Config } from "../models";
import { multiplicativeEffect } from "./common";

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function boostMultiplier(build, config) {
  return multiplicativeEffect(build, config, [proto.ArtifactSpec.Name.DILITHIUM_MONOCLE]);
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function boostDurationMultiplier(build, config) {
  return multiplicativeEffect(build, config, [proto.ArtifactSpec.Name.DILITHIUM_STONE]);
}

export { boostMultiplier, boostDurationMultiplier };
