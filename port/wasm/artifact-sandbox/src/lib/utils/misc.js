/**
 *
 * @param {Number} x
 * @param {Number} [fractionDigits=3]
 * @returns {String}
 */
function formatFloat(x, fractionDigits) {
  if (fractionDigits === undefined) {
    fractionDigits = 3;
  }
  let s = x.toFixed(3);
  // Trim trailing zeros, and possibly the decimal point.
  s = s.replace(/0+$/, "");
  if (s.endsWith(".")) {
    s = s.substring(0, s.length - 1);
  }
  return s;
}

/**
 *
 * @param {Number} x
 * @param {Number} [fractionDigits=2]
 * @returns {String}
 */
function formatPercentage(x, fractionDigits) {
  return formatFloat(x * 100, fractionDigits) + "%";
}

export { formatFloat, formatPercentage };
