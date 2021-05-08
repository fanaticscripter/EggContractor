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
