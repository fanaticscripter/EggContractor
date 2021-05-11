import { trimTrailingZeros } from "@/lib";

export function getLocalStorage(key: string, prefix?: string): string | undefined {
  if (prefix === undefined) {
    prefix = `${window.location.pathname}_`;
  }
  try {
    return localStorage[prefix + key];
  } catch (err) {
    console.error(err);
    return undefined;
  }
}

export function setLocalStorage(key: string, val: any, prefix?: string) {
  if (prefix === undefined) {
    prefix = `${window.location.pathname}_`;
  }
  try {
    localStorage[prefix + key] = val;
  } catch (err) {
    console.error(err);
  }
}

export function iconURL(relpath: string, size: number | string = "orig") {
  return `https://eggincassets.tcl.sh/${size}/${relpath}`;
}

export enum RoundingMode {
  Down = -1,
  Nearest = 0,
  Up = 1,
}

export function formatWithThousandSeparators(
  x: number,
  roundingMode = RoundingMode.Nearest
): string {
  let rounded: number;
  switch (roundingMode) {
    case RoundingMode.Down:
      rounded = Math.floor(x);
      break;
    case RoundingMode.Nearest:
      rounded = Math.round(x);
      break;
    case RoundingMode.Up:
      rounded = Math.ceil(x);
      break;
  }
  return rounded.toLocaleString("en-US");
}

export function formatPercentage(x: number, maxDecimals = 2): string {
  const s = (x * 100).toFixed(maxDecimals);
  return trimTrailingZeros(s) + "%";
}

export function formatDurationAuto(seconds: number): string {
  if (seconds < 0) {
    return "-" + formatDurationAuto(-seconds);
  }
  if (!isFinite(seconds)) {
    return "Forever";
  }
  let unit: string;
  let value: number;
  if (seconds < 60 * 59.5) {
    unit = "min";
    value = seconds / 60;
  } else if (seconds < 3600 * 23.5) {
    unit = "hr";
    value = seconds / 3600;
  } else {
    unit = "d";
    value = seconds / 86400;
  }
  return trimTrailingZeros(value.toFixed(1)) + unit;
}
