import lunr from "lunr";

import data from "./data.json";

export const artifacts = data.artifacts.map(artifact => ({
  ...artifact,
  display: `${artifact.name} (T${artifact.tier_number}), ${artifact.rarity}`,
  iconPath: `egginc/${artifact.icon_filename}`,
}));
export const stones = data.stones.map(stone => ({
  ...stone,
  display: `${stone.name} (T${stone.tier_number})`,
  iconPath: `egginc/${stone.icon_filename}`,
}));

export const artifactIdToArtifact = new Map(artifacts.map(artifact => [artifact.id, artifact]));
export const stoneIdToStone = new Map(stones.map(stone => [stone.id, stone]));

export function artifactFromId(id) {
  return artifactIdToArtifact.get(id) || null;
}

export function artifactFromAfxIdLevelRarity(afxId, afxLevel, afxRarity) {
  for (const artifact of artifacts) {
    if (
      artifact.afx_id === afxId &&
      artifact.afx_level === afxLevel &&
      artifact.afx_rarity === afxRarity
    ) {
      return artifact;
    }
  }
  return null;
}

export function stoneFromId(id) {
  return stoneIdToStone.get(id) || null;
}

export function stoneFromAfxIdLevel(afxId, afxLevel) {
  for (const stone of stones) {
    if (stone.afx_id === afxId && stone.afx_level === afxLevel) {
      return stone;
    }
  }
  return null;
}

const artifactsSearchIndex = lunr(function () {
  this.ref("id");
  this.field("display");
  for (const artifact of artifacts) {
    this.add(artifact);
  }
});

const stonesSearchIndex = lunr(function () {
  this.ref("id");
  this.field("display");
  for (const stone of stones) {
    this.add(stone);
  }
});

// Searching functionality is copied from data.ts of loot-simulator.

// These words or prefix of words aren't indexed, and would cause zero matches
// if otherwise queried as required.
const searchTermIgnoreList = ["a", "i", "in", "o", "of", "t", "th", "the"];

// Since lunr's wildcard doesn't match the empty string (e.g. "Simple Demeters
// necklace" is matched by "+necklace" but not "+necklace*"), we have to search
// twice, once with wildcard and once without, in order to return search results
// as the user types.
//
// See https://github.com/olivernn/lunr.js/issues/370
//
// Typescript signature:
// function search<T>(index: lunr.Index, userQuery: string, refToItem: (ref: string) => T): T[] {
function search(index, userQuery, refToItem) {
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

export function searchArtifacts(userQuery) {
  return search(artifactsSearchIndex, userQuery, ref => artifactIdToArtifact.get(ref));
}

export function searchStones(userQuery) {
  return search(stonesSearchIndex, userQuery, ref => stoneIdToStone.get(ref));
}
