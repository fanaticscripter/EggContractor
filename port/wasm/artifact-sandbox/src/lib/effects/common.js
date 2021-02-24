import { Build, Config } from "../models";

/**
 * @typedef effect
 * @property {!Number} delta
 * @property {!Number} multiplier
 */

/**
 * @param {!Build} build
 * @param {!Config} config
 * @param {!Array<!Number>} afxIds
 * @returns {!Array<!effect>}
 */
function gatherRelevantEffects(build, config, afxIds) {
  const deltas = [];
  for (const artifact of build.artifacts) {
    if (artifact.isEmpty()) {
      continue;
    }
    const effectMultiplier = config.isEnlightenment ? artifact.clarityEffect : 1;
    if (effectMultiplier === 0) {
      continue;
    }

    // Light of eggendil and slotted stones are ineffective on a non-enlightenment farm.
    if (!config.isEnlightenment && artifact.afx_id === proto.ArtifactSpec.Name.LIGHT_OF_EGGENDIL) {
      continue;
    }

    if (afxIds.includes(artifact.afx_id)) {
      deltas.push({
        delta: artifact.effect_delta,
        multiplier: effectMultiplier,
      });
    }

    for (const stone of artifact.activeStones) {
      if (stone === null) {
        continue;
      }
      if (afxIds.includes(stone.afx_id)) {
        deltas.push({
          delta: stone.effect_delta,
          multiplier: effectMultiplier,
        });
      }
    }
  }
  return deltas;
}

/**
 * @callback effectAggregator
 * @param {!Number} accumulator
 * @param {!effect} effect
 * @returns {!Number}
 */

/**
 * @param {!Build} build
 * @param {!Config} config
 * @param {!Array<!Number>} afxIds
 * @param {!effectAggregator} aggregator
 * @param {!Number} initialValue
 * @returns {!Number}
 */
function aggregateEffect(build, config, afxIds, aggregator, initialValue) {
  return gatherRelevantEffects(build, config, afxIds).reduce(aggregator, initialValue);
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @param {!Array<!Number>} afxIds
 * @returns {!Number}
 */
function additiveEffect(build, config, afxIds) {
  return aggregateEffect(
    build,
    config,
    afxIds,
    (aggregate, effect) => aggregate + effect.delta * effect.multiplier,
    0
  );
}

/**
 * @param {!Build} build
 * @param {!Config} config
 * @param {!Array<!Number>} afxIds
 * @returns {!Number}
 */
function multiplicativeEffect(build, config, afxIds) {
  return aggregateEffect(
    build,
    config,
    afxIds,
    (aggregate, effect) => aggregate * (1 + effect.delta * effect.multiplier),
    1
  );
}

export { additiveEffect, multiplicativeEffect, aggregateEffect };
