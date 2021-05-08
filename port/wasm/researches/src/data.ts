import initSqlJs, { Database, ParamsObject, SqlValue } from "sql.js";
import sqljsWasmURL from "sql.js/dist/sql-wasm.wasm?url";

import researchesData from "./researches.json";

type ResearchType = "common" | "epic";
type ResearchCategory =
  | "egg_laying_rate"
  | "egg_value"
  | "fleet_size"
  | "hab_capacity"
  | "hatchery_capacity"
  | "hatchery_refill_rate"
  | "internal_hatchery_rate"
  | "running_chicken_bonus"
  | "shipping_capacity";
type ResearchCategories = ResearchCategory | "" | "egg_laying_rate,egg_value";
type ResearchEffectType = "additive" | "multiplicative";
type ResearchCompoundType = "additive" | "multiplicative";
type Research = {
  serial_id: number;
  id: string;
  name: string;
  type: ResearchType;
  tier?: number;
  categories: ResearchCategories;
  description: string;
  effect_type: ResearchEffectType;
  levels: number;
  per_level: number;
  levels_compound: ResearchCompoundType;
  prices: number[];
};

const researches = researchesData as Research[];

export { researches };

let db: Database;

export const schema = `CREATE TABLE research(
  serial_id INTEGER PRIMARY KEY,
  id TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  -- Research type, either 'common' or 'epic'.
  type TEXT NOT NULL,
  -- Tier number for common research, or NULL for epic research.
  tier INTEGER,
  -- Comma-delimited categories. Possible catogries:
  -- * 'egg_laying_rate'
  -- * 'egg_value'
  -- * 'fleet_size'
  -- * 'hab_capacity'
  -- * 'hatchery_capacity'
  -- * 'hatchery_refill_rate'
  -- * 'internal_hatchery_rate'
  -- * 'running_chicken_bonus'
  -- * 'shipping_capacity';
  -- One research could belong to more than one categories;
  -- e.g. 'improved_genetics' has categories 'egg_laying_rate,egg_value'.
  -- One of a kind epic researches have categories set to the empty string.
  categories TEXT NOT NULL,
  description TEXT NOT NULL,
  -- Whether this research adds to its base effect, or multiplies it.
  -- Either 'additive' or 'multiplicative'.
  effect_type TEXT NOT NULL,
  -- Total number of levels.
  levels INTEGER NOT NULL,
  -- Bonus value per level.
  per_level REAL NOT NULL,
  -- How levels of this research are compounded with each other.
  -- Either 'additive' or 'multiplicative'.
  levels_compound TEXT NOT NULL,
  -- The prices (before discount) for each research level, encoded in a JSON array.
  -- The "Lab Upgrade" epic research is considered a discount.
  -- For common research, these are the bock prices; for epic research, these are GE prices.
  -- You can use SQLite JSON1 extension functions to query this field, and the result will be JSON-decoded as appropriate.
  -- E.g. json_extract(prices, '$[#-1]') extracts the price of the final level.
  prices TEXT NOT NULL
);`;

function undefinedToNull(x: SqlValue | undefined): SqlValue {
  return x !== undefined ? x : null;
}

// Because top-level await support is currently shit, and I don't want to create
// yet another Chrome-only tool, I just have to remember to call initDatabase
// before querying.
//
// https://github.com/tc39/proposal-top-level-await
// https://caniuse.com/mdn-javascript_operators_await_top_level
export async function initDatabase() {
  const SQL = await initSqlJs({ locateFile: file => sqljsWasmURL });
  db = new SQL.Database();
  db.run(schema);
  db.run(
    `INSERT INTO research(
      serial_id,
      id,
      name,
      type,
      tier,
      categories,
      description,
      effect_type,
      levels,
      per_level,
      levels_compound,
      prices
    ) VALUES ` + Array(researches.length).fill("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)").join(", "),
    researches
      .map(r => [
        r.serial_id,
        r.id,
        r.name,
        r.type,
        undefinedToNull(r.tier),
        r.categories,
        r.description,
        r.effect_type,
        r.levels,
        r.per_level,
        r.levels_compound,
        JSON.stringify(r.prices || []),
      ])
      .flat()
  );
}

export function executeQuery(query: string) {
  const stmt = db.prepare(query);
  const results = [];
  while (stmt.step()) {
    results.push(transformResult(stmt.getAsObject()));
  }
  return results;
}

function transformResult(x: ParamsObject): { [key: string]: any } {
  const result: { [key: string]: any } = {};
  for (const [key, val] of Object.entries(x)) {
    // Skip null fields so that they don't get serialized.
    if (val === null) {
      continue;
    }
    // Try to decode JSON objects and arrays.
    if (typeof val === "string" && val.length > 0 && val.match(/^(\{.*\}|\[.*\])$/)) {
      result[key] = JSON.parse(val);
    } else {
      result[key] = val;
    }
  }
  return result;
}
