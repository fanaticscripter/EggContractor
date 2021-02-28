import { Build, Config } from "../models";
import { multiplicativeEffect } from "./common";
import { earningsWithMaxRunningChickenBonusMultiplier } from "./earnings";

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function soulEggsGainMultiplier(build, config) {
  const virtualEarningsMultiplier =
    earningsWithMaxRunningChickenBonusMultiplier(build, config) *
    multiplicativeEffect(build, config, [proto.ArtifactSpec.Name.PHOENIX_FEATHER]);
  return Math.pow(virtualEarningsMultiplier, 0.21);
}

export { soulEggsGainMultiplier };
