import { Build, Config } from "../models";
import { aggregateEffect } from "./common";

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function researchPriceMultiplier(build, config) {
  return aggregateEffect(
    build,
    config,
    [proto.ArtifactSpec.Name.PUZZLE_CUBE],
    (aggregate, effect) =>
      effect.multiplier <= 1
        ? (1 + effect.delta * effect.multiplier) * aggregate
        : ((1 + effect.delta) / effect.multiplier) * aggregate,
    1
  );
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function researchPriceDiscount(build, config) {
  return researchPriceMultiplier(build, config) - 1;
}

export { researchPriceMultiplier, researchPriceDiscount };
