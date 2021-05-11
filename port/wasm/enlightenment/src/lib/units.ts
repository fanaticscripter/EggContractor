import { trimTrailingZeros } from "./utils";

// https://egg-inc.fandom.com/wiki/Order_of_Magnitude
const units = [
  { symbol: "K", oom: 3 },
  { symbol: "M", oom: 6 },
  { symbol: "B", oom: 9 },
  { symbol: "T", oom: 12 },
  { symbol: "q", oom: 15 },
  { symbol: "Q", oom: 18 },
  { symbol: "s", oom: 21 },
  { symbol: "S", oom: 24 },
  { symbol: "o", oom: 27 },
  { symbol: "N", oom: 30 },
  { symbol: "d", oom: 33 },
  { symbol: "U", oom: 36 },
  { symbol: "D", oom: 39 },
  { symbol: "Td", oom: 42 },
  { symbol: "qd", oom: 45 },
  { symbol: "Qd", oom: 48 },
  { symbol: "sd", oom: 51 },
  { symbol: "Sd", oom: 54 },
  { symbol: "Od", oom: 57 },
  { symbol: "Nd", oom: 60 },
  { symbol: "V", oom: 63 },
  { symbol: "uV", oom: 66 },
  { symbol: "dV", oom: 69 },
  { symbol: "tV", oom: 72 },
  { symbol: "qV", oom: 75 },
  { symbol: "QV", oom: 78 },
  { symbol: "sV", oom: 81 },
  { symbol: "SV", oom: 84 },
  { symbol: "OV", oom: 87 },
  { symbol: "NV", oom: 90 },
  { symbol: "tT", oom: 93 },
];

const oom2symbol = new Map(units.map(u => [u.oom, u.symbol]));
const minOom = units[0].oom;
const maxOom = units[units.length - 1].oom;

// When scientific on, the value is formatted as HTML.
export function formatEIValue(
  x: number,
  options?: { trim?: boolean; decimals?: number; scientific?: boolean }
): string {
  const trim = options?.trim === undefined ? false : options?.trim;
  const decimals = options?.decimals === undefined ? 3 : options?.decimals;
  const scientific = options?.scientific === undefined ? false : options?.scientific;
  if (isNaN(x)) {
    return "NaN";
  }
  if (x < 0) {
    return "-" + formatEIValue(-x, options);
  }
  if (!isFinite(x)) {
    return "infinity";
  }
  const oom = Math.log10(x);
  if (oom < minOom) {
    // Always round small number to an integer.
    return x.toFixed(0);
  }
  let oomFloor = Math.floor(oom);
  if (oom + 1e-9 >= oomFloor + 1) {
    // Fix problem of 1q being displayed as 1000T, 1N displayed as 1000o, etc,
    // where the floor is one integer down due to floating point imprecision.
    oomFloor++;
  }
  oomFloor -= oomFloor % 3;
  if (oomFloor > maxOom) {
    oomFloor = maxOom;
  }
  const principal = x / Math.pow(10, oomFloor);
  let numpart = principal.toFixed(decimals);
  if (trim) {
    numpart = trimTrailingZeros(numpart);
  }
  if (scientific) {
    return `${numpart}&times;10<sup>${oomFloor}</sup>`;
  } else {
    return numpart + oom2symbol.get(oomFloor)!;
  }
}
