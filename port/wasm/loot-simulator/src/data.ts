import lunr from "lunr";

import data from "./app-data.json";
import { Item, ItemId, Mission, MissionId } from "./types";

export const items = data.items as Item[];
export const itemIds = items.map(item => item.id);
export const itemIdToItem = new Map(items.map(item => [item.id, item]));
export const missions = [...data.missions].reverse() as Mission[];
export const missionIds = missions.map(mission => mission.id);
export const missionIdToMission = new Map(missions.map(mission => [mission.id, mission]));

const itemsSearchIndex = lunr(function () {
  this.ref("id");
  this.field("display");
  this.field("tier_name");
  for (const item of items) {
    this.add(item);
  }
});

const missionsSearchIndex = lunr(function () {
  this.ref("id");
  this.field("display");
  for (const mission of missions) {
    this.add(mission);
  }
});

// These words or prefix of words aren't indexed, and would cause zero matches
// if otherwise queried as required.
const searchTermIgnoreList = ["a", "i", "in", "o", "of", "t", "th", "the"];

// Since lunr's wildcard doesn't match the empty string (e.g. "Simple Demeters
// necklace" is matched by "+necklace" but not "+necklace*"), we have to search
// twice, once with wildcard and once without, in order to return search results
// as the user types.
//
// See https://github.com/olivernn/lunr.js/issues/370
function search<T>(index: lunr.Index, userQuery: string, refToItem: (ref: string) => T): T[] {
  let terms = userQuery
    .replace(/[^A-Za-z0-9\s]/g, " ")
    .split(/\s+/)
    .map(term => term.toLowerCase());
  // As long as the query doesn't end in whitespace, the final term should be
  // treated as partial (user is in the middle of typing the term).
  const partialFinal = terms[terms.length - 1] !== "";
  terms = terms.filter(term => term !== "");
  const fullMatches = index.query(query => {
    terms.forEach(term => {
      if (!searchTermIgnoreList.includes(term)) {
        query.term(term, { presence: lunr.Query.presence.REQUIRED });
      }
    });
  });
  if (!partialFinal) {
    return fullMatches.map(result => refToItem(result.ref));
  }
  const partialFinalMatches = index.query(query => {
    terms.forEach((term, index) => {
      if (index === terms.length - 1) {
        // Add a trailing wildcard to the final term, which is being typed out.
        query.term(term, {
          wildcard: lunr.Query.wildcard.TRAILING,
          presence: lunr.Query.presence.REQUIRED,
        });
      } else {
        if (!searchTermIgnoreList.includes(term)) {
          query.term(term, { presence: lunr.Query.presence.REQUIRED });
        }
      }
    });
  });
  const matches = fullMatches
    .concat(partialFinalMatches)
    .sort((result1, result2) => result2.score - result1.score);
  // Deduplicate and keep the highest score entry of each result.
  const matchRefs = new Set(matches.map(result => result.ref));
  return [...matchRefs.entries()].map(([ref, _]) => refToItem(ref));
}

export function searchItems(userQuery: string): Item[] {
  return search(itemsSearchIndex, userQuery, ref => itemIdToItem.get(ref as ItemId)!);
}

export function searchMissions(userQuery: string): Mission[] {
  return search(missionsSearchIndex, userQuery, ref => missionIdToMission.get(ref as MissionId)!);
}
