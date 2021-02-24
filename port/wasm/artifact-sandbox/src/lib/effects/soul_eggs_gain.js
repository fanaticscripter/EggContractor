import { Build, Config } from "../models";
import { multiplicativeEffect } from "./common";
import { earningsWithMaxRunningChickenBonusMultipler } from "./earnings";

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function soulEggsGainMultipler(build, config) {
  const virtualEarningsMultiplier =
    earningsWithMaxRunningChickenBonusMultipler(build, config) *
    multiplicativeEffect(build, config, [proto.ArtifactSpec.Name.PHOENIX_FEATHER]);
  return Math.pow(virtualEarningsMultiplier, 0.21);
}

export { soulEggsGainMultipler };
