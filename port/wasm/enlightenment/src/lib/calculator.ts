import { parseValueWithUnit, valueWithUnitRegExpGlobal } from "./units";

export {} from "./units";

// Map lower case version of a Math proper to the qualified property (with
// proper case), e.g. abs => Math.abs, pi => Math.PI.
//
// All symbols except Math.E are supported. E unfornately appears in scientific
// notations all the time.
const supportedMathSymbols = new Map(
  (Object.getOwnPropertyNames(Math)
    .filter(symbol => symbol !== "E")
    .map(symbol => [symbol.toLowerCase(), `Math.${symbol}`]) as [string, string][]).concat(
    // Additional supported aliases.
    [
      ["ln", "Math.log"],
      ["lg", "Math.log10"],
    ]
  )
);

const mathSymbolPattern = `\\b${[...supportedMathSymbols.keys()].join("|")}\\b`;
const mathSymbolRegExpGlobal = new RegExp(mathSymbolPattern, "gi");

// Returns null on error.
export function calculateWithOoMUnits(expr: string): number | null {
  expr = expr.replaceAll(valueWithUnitRegExpGlobal, value => parseValueWithUnit(value)!.toString());
  expr = expr.replaceAll(
    mathSymbolRegExpGlobal,
    value => supportedMathSymbols.get(value.toLowerCase())!
  );
  // Interpret ^ as exponentiation, because who the hell cares about xor.
  expr = expr.replaceAll(/\^/g, "**");
  try {
    const result = eval(expr);
    return typeof result === "number" ? result : null;
  } catch (e) {
    return null;
  }
}
