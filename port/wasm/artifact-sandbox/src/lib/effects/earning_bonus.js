import { Build, Config } from "../models";
import { additiveEffect } from "./common";

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function earningBonus(build, config) {
  const peBonus = prophecyEggBonus(build, config);
  const peCount = config.prophecyEggs;
  const seBonus = soulEggBonus(build, config);
  const seCount = config.soulEggs;
  return seCount * seBonus * Math.pow(1 + peBonus, peCount);
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function earningBonusMultiplier(build, config) {
  const peBonusBase = baseProphecyEggBonus(config);
  const peBonus = prophecyEggBonus(build, config);
  const peCount = config.prophecyEggs;
  const seBonusBase = baseSoulEggBonus(config);
  const seBonus = soulEggBonus(build, config);
  return Math.pow((1 + peBonus) / (1 + peBonusBase), peCount) * (seBonus / seBonusBase);
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function prophecyEggBonus(build, config) {
  return (
    baseProphecyEggBonus(config) +
    additiveEffect(build, config, [
      proto.ArtifactSpec.Name.BOOK_OF_BASAN,
      proto.ArtifactSpec.Name.PROPHECY_STONE,
    ])
  );
}

/**
 * @param {!Config} config
 * @returns {!Number}
 */
function baseProphecyEggBonus(config) {
  return 0.05 + 0.01 * config.prophecyBonus;
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @returns {!Number}
 */
function soulEggBonus(build, config) {
  return (
    baseSoulEggBonus(config) + additiveEffect(build, config, [proto.ArtifactSpec.Name.SOUL_STONE])
  );
}

/**
 * @param {!Config} config
 * @returns {!Number}
 */
function baseSoulEggBonus(config) {
  return 0.1 + 0.01 * config.soulFood;
}

export { earningBonus, earningBonusMultiplier };
