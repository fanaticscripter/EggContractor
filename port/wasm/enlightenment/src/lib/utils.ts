// Trim trailing zeros, and possibly the decimal point.
export function trimTrailingZeros(s: string): string {
  s = s.replace(/0+$/, "");
  if (s.endsWith(".")) {
    s = s.substring(0, s.length - 1);
  }
  return s;
}
