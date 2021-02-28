import { Build, Config } from "../models";
import { multiplicativeEffect } from "./common";
import { boostMultiplier } from "./boosts";
import { earningBonusMultiplier } from "./earning_bonus";
import { layingRateMultiplier } from "./laying_rate";
import { maxRunningChickenBonusMultiplier } from "./rcb";

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function earningsMultiplier(build, config) {
  return (
    earningBonusMultiplier(build, config) *
    multiplicativeEffect(build, config, [
      proto.ArtifactSpec.Name.TUNGSTEN_ANKH,
      proto.ArtifactSpec.Name.DEMETERS_NECKLACE,
      proto.ArtifactSpec.Name.LIGHT_OF_EGGENDIL,
      proto.ArtifactSpec.Name.SHELL_STONE,
    ]) *
    layingRateMultiplier(build, config) *
    (config.birdFeedActive ? boostMultiplier(build, config) : 1)
  );
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function earningsWithMaxRunningChickenBonusMultiplier(build, config) {
  return earningsMultiplier(build, config) * maxRunningChickenBonusMultiplier(build, config);
}

export { earningsMultiplier, earningsWithMaxRunningChickenBonusMultiplier };
