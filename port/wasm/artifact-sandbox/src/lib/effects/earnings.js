import { Build, Config } from "../models";
import { multiplicativeEffect } from "./common";
import { earningBonusMultiplier } from "./earning_bonus";
import { layingRateMultiplier } from "./laying_rate";
import { maxRunningChickenBonusMultiplier } from "./rcb";

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function earningsMultipler(build, config) {
  return (
    earningBonusMultiplier(build, config) *
    multiplicativeEffect(build, config, [
      proto.ArtifactSpec.Name.TUNGSTEN_ANKH,
      proto.ArtifactSpec.Name.DEMETERS_NECKLACE,
      proto.ArtifactSpec.Name.LIGHT_OF_EGGENDIL,
      proto.ArtifactSpec.Name.SHELL_STONE,
    ]) *
    layingRateMultiplier(build, config)
  );
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function earningsWithMaxRunningChickenBonusMultipler(build, config) {
  return earningsMultipler(build, config) * maxRunningChickenBonusMultiplier(build, config);
}

export { earningsMultipler, earningsWithMaxRunningChickenBonusMultipler };
