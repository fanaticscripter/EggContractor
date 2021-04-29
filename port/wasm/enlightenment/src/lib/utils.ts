export function uint8ArrayToBinaryString(a: Uint8Array): string {
  return Array.from(a)
    .map(c => String.fromCharCode(c))
    .join("");
}

export function binaryStringToUint8Array(b: string): Uint8Array {
  const buf = new Uint8Array(new ArrayBuffer(b.length));
  for (let i = 0; i < b.length; i++) {
    buf[i] = b.charCodeAt(i);
  }
  return buf;
}

// Trim trailing zeros, and possibly the decimal point.
export function trimTrailingZeros(s: string): string {
  s = s.replace(/0+$/, "");
  if (s.endsWith(".")) {
    s = s.substring(0, s.length - 1);
  }
  return s;
}
