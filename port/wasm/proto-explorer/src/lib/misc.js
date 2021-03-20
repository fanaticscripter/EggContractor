/**
 * @param {!number} min
 * @param {!number} max
 * @returns {!number}
 */
function randomRange(min, max) {
  return Math.floor(Math.random() * (max - min) + min);
}

/**
 * @returns {!string}
 */
function randomEiUserId() {
  return `EI${randomRange(1e15, 1e16)}`;
}

export { randomEiUserId };
