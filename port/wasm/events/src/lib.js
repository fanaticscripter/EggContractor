import events from "@/events.json";
import legacyEvents from "@/legacy_events.json";

import dayjs from "dayjs";
import utc from "dayjs/plugin/utc";

dayjs.extend(utc);

// Version history extracted from https://apps.apple.com/us/app/egg-inc/id993492744:
//
// [
//   {
//     "versionDisplay": "1.20.9",
//     "releaseNotes": "- New Contracts UI and features",
//     "releaseDate": "2021-04-30"
//   },
//   {
//     "versionDisplay": "1.20.7",
//     "releaseNotes": "- Small improvements and bug fixes",
//     "releaseDate": "2021-03-15"
//   },
//   {
//     "versionDisplay": "1.20.6",
//     "releaseNotes": "- Small improvements and bug fixes",
//     "releaseDate": "2021-03-11"
//   },
//   {
//     "versionDisplay": "1.20.5",
//     "releaseNotes": "- Many small/medium improvements and bug fixes\n- Increased rewards for random gifts and video gifts\n- Slightly reduced rewards for some trophies and daily gifts",
//     "releaseDate": "2021-03-01"
//   },
//   {
//     "versionDisplay": "1.20.4",
//     "releaseNotes": "- Large increases to rewards for challenges, trophies, and daily gifts\n- Bug fixes",
//     "releaseDate": "2021-02-16"
//   },
//   {
//     "versionDisplay": "1.20.3",
//     "releaseNotes": "- Stability improvements and bug fixes",
//     "releaseDate": "2021-02-05"
//   },
//   {
//     "versionDisplay": "1.20.2",
//     "releaseNotes": "- Stability improvements and bug fixes",
//     "releaseDate": "2021-01-29"
//   },
//   {
//     "versionDisplay": "1.20.1",
//     "releaseNotes": "ARTIFACTS IS HERE!\nAt long last, the largest update ever to Egg, Inc!\nCommission deep space exploration missions to uncover over 100 artifacts from long lost civilizations around the galaxy\nEXPORE. CRAFT. CUSTOMIZE.\nUnlock the secrets of the universe.",
//     "releaseDate": "2021-01-21"
//   },
//   {
//     "versionDisplay": "1.20",
//     "releaseNotes": "ARTIFACTS IS HERE!\nAt long last, the largest update ever to Egg, Inc!\nCommission deep space exploration missions to uncover over 100 artifacts from long lost civilizations around the galaxy\nEXPORE. CRAFT. CUSTOMIZE.\nUnlock the secrets of the universe.",
//     "releaseDate": "2021-01-20"
//   },
//   {
//     "versionDisplay": "1.12.13",
//     "releaseNotes": "- Addresses ad bugs",
//     "releaseDate": "2020-11-11"
//   },
//   {
//     "versionDisplay": "1.12.12",
//     "releaseNotes": "- iPhone 12 display bug fix",
//     "releaseDate": "2020-10-27"
//   },
//   {
//     "versionDisplay": "1.12.11",
//     "releaseNotes": "- bug fixes",
//     "releaseDate": "2020-07-05"
//   },
//   {
//     "versionDisplay": "1.12.10",
//     "releaseNotes": "- small improvements and bug fixes",
//     "releaseDate": "2020-06-30"
//   },
//   {
//     "versionDisplay": "1.12.9",
//     "releaseNotes": "- small improvements and bug fixes",
//     "releaseDate": "2020-06-11"
//   },
//   {
//     "versionDisplay": "1.12.8",
//     "releaseNotes": "- Contract progress details now displayed on contract egg screen\n- Exponential contract progress bars to give a better feel on progress (tap to toggle)\n- Press and hold boost token gift buttons\n- Get a Boost token as a video reward or as a random delivery!\n- Elite Contract Goals & Rewards - Contract goals & rewards are now split into standard and elite. Beginning players will have more achievable goals, and expert players will have more fun and earn larger rewards, Egg of Prophecy rewards will be the same between both types!",
//     "releaseDate": "2020-03-17"
//   },
//   {
//     "versionDisplay": "1.12.7",
//     "releaseNotes": "- Exponential contract progress bars to give a better feel on progress (tap to toggle)\n- Press and hold boost token gift buttons\n- Get a Boost token as a video reward or as a random delivery!\n- Elite Contract Goals & Rewards - Contract goals & rewards are now split into standard and elite. Beginning players will have more achievable goals, and expert players will have more fun and earn larger rewards, Egg of Prophecy rewards will be the same between both types!",
//     "releaseDate": "2020-03-16"
//   },
//   {
//     "versionDisplay": "1.12.6",
//     "releaseNotes": "- See boost tokens of your coop mates!\n- Several minor improvements and bug fixes",
//     "releaseDate": "2020-02-21"
//   },
//   {
//     "versionDisplay": "1.12.5",
//     "releaseNotes": "- Minor improvements and bug fixes\n- Enhanced Privacy Controls\n",
//     "releaseDate": "2020-01-29"
//   },
//   {
//     "versionDisplay": "1.12.4",
//     "releaseNotes": "* Alleviates coop kick bug and age restriction bug *\n- CCPA Compliance and data consent is now expanded to ALL Players\n- Players who do not consent to data collection can now watch *untracked* ads for rewards!\n- Piggy grows faster\n- Bug fixes",
//     "releaseDate": "2020-01-11"
//   },
//   {
//     "versionDisplay": "1.12.3",
//     "releaseNotes": "- CCPA Compliance and data consent is now expanded to ALL Players\n- Players who do not consent to data collection can now watch *untracked* ads for rewards!\n- Piggy grows faster\n- Bug fixes",
//     "releaseDate": "2020-01-10"
//   },
//   {
//     "versionDisplay": "1.12.2",
//     "releaseNotes": "* crash fix, ads fix, instant cash boost bug fix, running chicken bonus increase\n- Missions update: Complete 3 missions simultaneously\n- Boost Tokens! A new fun way boosts are regulated in Contracts\n- New Boost: Soul Mirror - Mirror the Earnings Bonus of the top player in your co-op!\n- Balance changes: Some cash have increased BUT while the MAX running chicken bonus reduced, but it now increases farm value AND drone rewards!\n- Numerous improvements and bug fixes with more to come!",
//     "releaseDate": "2019-12-16"
//   },
//   {
//     "versionDisplay": "1.12.1",
//     "releaseNotes": "* fix for a restore backup crash\n- Missions update: Complete 3 missions simultaneously\n- Boost Tokens! A new fun way boosts are regulated in Contracts\n- New Boost: Soul Mirror - Mirror the Earnings Bonus of the top player in your co-op!\n- Balance changes: Some cash have increased BUT while the MAX running chicken bonus reduced, but it now increases farm value AND drone rewards!\n- Numerous improvements and bug fixes with more to come!",
//     "releaseDate": "2019-12-14"
//   },
//   {
//     "versionDisplay": "1.12",
//     "releaseNotes": "- Missions update: Complete 3 missions simultaneously\n- Boost Tokens! A new fun way boosts are regulated in Contracts\n- New Boost: Soul Mirror - Mirror the Earnings Bonus of the top player in your co-op!\n- Balance changes: Some cash have increased BUT while the MAX running chicken bonus reduced, but it now increases farm value AND drone rewards!\n- Numerous improvements and bug fixes with more to come!",
//     "releaseDate": "2019-12-10"
//   },
//   {
//     "versionDisplay": "1.11.6",
//     "releaseNotes": "- Numerous small improvements and bug fixes",
//     "releaseDate": "2019-10-24"
//   },
//   {
//     "versionDisplay": "1.11.5",
//     "releaseNotes": "- New Help screen!\n- Different purchases increment the piggy different amounts\n",
//     "releaseDate": "2019-10-22"
//   },
//   {
//     "versionDisplay": "1.11.4",
//     "releaseNotes": "- Contract onboarding improvements\n- Piggy cap starts lower but grows faster and ends up increasing the cap significantly\n* fixes contract screen crash when no contracts are available",
//     "releaseDate": "2019-10-10"
//   }
// ]

const utcDateToTimestamp = date => dayjs.utc(date).unix();

const eventsAndAppUpdates = [
  ...legacyEvents,
  ...events,
  // {
  //   type: "app-update",
  //   version: "1.11.4",
  //   startTimestamp: utcDateToTimestamp("2019-10-10"),
  //   releaseNotes:
  //     "- Contract onboarding improvements\n- Piggy cap starts lower but grows faster and ends up increasing the cap significantly\n* fixes contract screen crash when no contracts are available",
  // },
  // {
  //   type: "app-update",
  //   version: "1.11.5",
  //   startTimestamp: utcDateToTimestamp("2019-10-22"),
  //   releaseNotes:
  //     "- New Help screen!\n- Different purchases increment the piggy different amounts\n",
  // },
  // {
  //   type: "app-update",
  //   version: "1.11.6",
  //   startTimestamp: utcDateToTimestamp("2019-10-24"),
  //   releaseNotes: "- Numerous small improvements and bug fixes",
  // },
  // {
  //   type: "app-update",
  //   version: "1.12",
  //   startTimestamp: utcDateToTimestamp("2019-12-10"),
  //   releaseNotes:
  //     "- Missions update: Complete 3 missions simultaneously\n- Boost Tokens! A new fun way boosts are regulated in Contracts\n- New Boost: Soul Mirror - Mirror the Earnings Bonus of the top player in your co-op!\n- Balance changes: Some cash have increased BUT while the MAX running chicken bonus reduced, but it now increases farm value AND drone rewards!\n- Numerous improvements and bug fixes with more to come!",
  // },
  // {
  //   type: "app-update",
  //   version: "1.12.1",
  //   startTimestamp: utcDateToTimestamp("2019-12-14"),
  //   releaseNotes:
  //     "* fix for a restore backup crash\n- Missions update: Complete 3 missions simultaneously\n- Boost Tokens! A new fun way boosts are regulated in Contracts\n- New Boost: Soul Mirror - Mirror the Earnings Bonus of the top player in your co-op!\n- Balance changes: Some cash have increased BUT while the MAX running chicken bonus reduced, but it now increases farm value AND drone rewards!\n- Numerous improvements and bug fixes with more to come!",
  // },
  // {
  //   type: "app-update",
  //   version: "1.12.2",
  //   startTimestamp: utcDateToTimestamp("2019-12-16"),
  //   releaseNotes:
  //     "* crash fix, ads fix, instant cash boost bug fix, running chicken bonus increase\n- Missions update: Complete 3 missions simultaneously\n- Boost Tokens! A new fun way boosts are regulated in Contracts\n- New Boost: Soul Mirror - Mirror the Earnings Bonus of the top player in your co-op!\n- Balance changes: Some cash have increased BUT while the MAX running chicken bonus reduced, but it now increases farm value AND drone rewards!\n- Numerous improvements and bug fixes with more to come!",
  // },
  {
    type: "app-update",
    version: "1.12.3",
    startTimestamp: utcDateToTimestamp("2020-01-10"),
    releaseNotes:
      "- CCPA Compliance and data consent is now expanded to ALL Players\n- Players who do not consent to data collection can now watch *untracked* ads for rewards!\n- Piggy grows faster\n- Bug fixes",
  },
  {
    type: "app-update",
    version: "1.12.4",
    startTimestamp: utcDateToTimestamp("2020-01-11"),
    releaseNotes:
      "* Alleviates coop kick bug and age restriction bug *\n- CCPA Compliance and data consent is now expanded to ALL Players\n- Players who do not consent to data collection can now watch *untracked* ads for rewards!\n- Piggy grows faster\n- Bug fixes",
  },
  {
    type: "app-update",
    version: "1.12.5",
    startTimestamp: utcDateToTimestamp("2020-01-29"),
    releaseNotes: "- Minor improvements and bug fixes\n- Enhanced Privacy Controls\n",
  },
  {
    type: "app-update",
    version: "1.12.6",
    startTimestamp: utcDateToTimestamp("2020-02-21"),
    releaseNotes:
      "- See boost tokens of your coop mates!\n- Several minor improvements and bug fixes",
  },
  {
    type: "app-update",
    version: "1.12.7",
    startTimestamp: utcDateToTimestamp("2020-03-16"),
    releaseNotes:
      "- Exponential contract progress bars to give a better feel on progress (tap to toggle)\n- Press and hold boost token gift buttons\n- Get a Boost token as a video reward or as a random delivery!\n- Elite Contract Goals & Rewards - Contract goals & rewards are now split into standard and elite. Beginning players will have more achievable goals, and expert players will have more fun and earn larger rewards, Egg of Prophecy rewards will be the same between both types!",
  },
  {
    type: "app-update",
    version: "1.12.8",
    startTimestamp: utcDateToTimestamp("2020-03-17"),
    releaseNotes:
      "- Contract progress details now displayed on contract egg screen\n- Exponential contract progress bars to give a better feel on progress (tap to toggle)\n- Press and hold boost token gift buttons\n- Get a Boost token as a video reward or as a random delivery!\n- Elite Contract Goals & Rewards - Contract goals & rewards are now split into standard and elite. Beginning players will have more achievable goals, and expert players will have more fun and earn larger rewards, Egg of Prophecy rewards will be the same between both types!",
  },
  {
    type: "app-update",
    version: "1.12.9",
    startTimestamp: utcDateToTimestamp("2020-06-11"),
    releaseNotes: "- small improvements and bug fixes",
  },
  {
    type: "app-update",
    version: "1.12.10",
    startTimestamp: utcDateToTimestamp("2020-06-30"),
    releaseNotes: "- small improvements and bug fixes",
  },
  {
    type: "app-update",
    version: "1.12.11",
    startTimestamp: utcDateToTimestamp("2020-07-05"),
    releaseNotes: "- bug fixes",
  },
  {
    type: "app-update",
    version: "1.12.12",
    startTimestamp: utcDateToTimestamp("2020-10-27"),
    releaseNotes: "- iPhone 12 display bug fix",
  },
  {
    type: "app-update",
    version: "1.12.13",
    startTimestamp: utcDateToTimestamp("2020-11-11"),
    releaseNotes: "- Addresses ad bugs",
  },
  {
    type: "app-update",
    version: "1.20",
    startTimestamp: utcDateToTimestamp("2021-01-20"),
    releaseNotes:
      "ARTIFACTS IS HERE!\nAt long last, the largest update ever to Egg, Inc!\nCommission deep space exploration missions to uncover over 100 artifacts from long lost civilizations around the galaxy\nEXPORE. CRAFT. CUSTOMIZE.\nUnlock the secrets of the universe.",
  },
  {
    type: "app-update",
    version: "1.20.1",
    startTimestamp: utcDateToTimestamp("2021-01-21"),
    releaseNotes:
      "ARTIFACTS IS HERE!\nAt long last, the largest update ever to Egg, Inc!\nCommission deep space exploration missions to uncover over 100 artifacts from long lost civilizations around the galaxy\nEXPORE. CRAFT. CUSTOMIZE.\nUnlock the secrets of the universe.",
  },
  {
    type: "app-update",
    version: "1.20.2",
    startTimestamp: utcDateToTimestamp("2021-01-29"),
    releaseNotes: "- Stability improvements and bug fixes",
  },
  {
    type: "app-update",
    version: "1.20.3",
    startTimestamp: utcDateToTimestamp("2021-02-05"),
    releaseNotes: "- Stability improvements and bug fixes",
  },
  {
    type: "app-update",
    version: "1.20.4",
    startTimestamp: utcDateToTimestamp("2021-02-16"),
    releaseNotes:
      "- Large increases to rewards for challenges, trophies, and daily gifts\n- Bug fixes",
  },
  {
    type: "app-update",
    version: "1.20.5",
    startTimestamp: utcDateToTimestamp("2021-03-01"),
    releaseNotes:
      "- Many small/medium improvements and bug fixes\n- Increased rewards for random gifts and video gifts\n- Slightly reduced rewards for some trophies and daily gifts",
  },
  {
    type: "app-update",
    version: "1.20.6",
    startTimestamp: utcDateToTimestamp("2021-03-11"),
    releaseNotes: "- Small improvements and bug fixes",
  },
  {
    type: "app-update",
    version: "1.20.7",
    startTimestamp: utcDateToTimestamp("2021-03-15"),
    releaseNotes: "- Small improvements and bug fixes",
  },
  {
    type: "app-update",
    version: "1.20.9",
    startTimestamp: utcDateToTimestamp("2021-04-30"),
    releaseNotes: "- New Contracts UI and features",
  },
].sort((e1, e2) => e1.startTimestamp - e2.startTimestamp);

const existingEventTypes = new Set(eventsAndAppUpdates.map(e => e.type));

const eventTypes = [
  ["app-update", "APP UPDATE"],
  // epic-research-sale is hoisted to the top (from before research-sale) in order to
  // (a) Make the badge distinguishable from research-sale (both use the same icon);
  // (b) Not waste a good color (red) on the useless even that is piggy-boost.
  ["epic-research-sale", "EPIC RESEARCH SALE"],
  ["piggy-boost", "PIGGY GROWTH"],
  ["piggy-cap-boost", "UNLIMITED PIGGY"],
  ["prestige-boost", "PRESTIGE BOOST"],
  ["earnings-boost", "CASH BOOST"],
  ["gift-boost", "GENEROUS GIFTS"],
  ["drone-boost", "GENEROUS DRONES"],
  ["research-sale", "RESEARCH SALE"],
  ["hab-sale", "HEN HOUSE SALE"],
  ["vehicle-sale", "VEHICLE SALE"],
  ["boost-sale", "BOOST SALE"],
  ["boost-duration", "BOOST TIME+"],
  ["crafting-sale", "CRAFTING SALE"],
  ["mission-fuel", "MISSION FUEL BOOST"],
  ["mission-capacity", "MISSION CAPACITY BOOST"],
  ["mission-duration", "MISSION DURATION CUT"],
].filter(([type, name]) => existingEventTypes.has(type));

const baseColors = [
  "gray",
  "red",
  "orange",
  "amber",
  "yellow",
  "lime",
  "green",
  "emerald",
  "teal",
  "cyan",
  "light-blue",
  "blue",
  "indigo",
  "violet",
  "purple",
  "fuchsia",
  "pink",
];

// A list of classes used, for purgecss.
[
  "text-gray-500",
  "text-red-500",
  "text-orange-500",
  "text-amber-500",
  "text-yellow-500",
  "text-lime-500",
  "text-green-500",
  "text-emerald-500",
  "text-teal-500",
  "text-cyan-500",
  "text-light-blue-500",
  "text-blue-500",
  "text-indigo-500",
  "text-violet-500",
  "text-purple-500",
  "text-fuchsia-500",
  "text-pink-500",

  "text-gray-300",
  "text-red-300",
  "text-orange-300",
  "text-amber-300",
  "text-yellow-300",
  "text-lime-300",
  "text-green-300",
  "text-emerald-300",
  "text-teal-300",
  "text-cyan-300",
  "text-light-blue-300",
  "text-blue-300",
  "text-indigo-300",
  "text-violet-300",
  "text-purple-300",
  "text-fuchsia-300",
  "text-pink-300",

  "bg-gray-500",
  "bg-red-500",
  "bg-orange-500",
  "bg-amber-500",
  "bg-yellow-500",
  "bg-lime-500",
  "bg-green-500",
  "bg-emerald-500",
  "bg-teal-500",
  "bg-cyan-500",
  "bg-light-blue-500",
  "bg-blue-500",
  "bg-indigo-500",
  "bg-violet-500",
  "bg-purple-500",
  "bg-fuchsia-500",
  "bg-pink-500",
];

const eventIcon = eventId => {
  switch (eventId) {
    case "piggy-boost":
      return "egginc/icon_piggy.png";
    case "piggy-cap-boost":
      return "egginc/icon_piggy.png";
    case "prestige-boost":
      return "egginc-extras/icon_prestige_boost.png";
    case "earnings-boost":
      return "egginc-extras/icon_earnings_boost.png";
    case "gift-boost":
      return "egginc-extras/icon_gift_boost.png";
    case "drone-boost":
      return "egginc-extras/icon_drone_boost.png";
    case "epic-research-sale":
      return "egginc-extras/icon_research_sale.png";
    case "research-sale":
      return "egginc-extras/icon_research_sale.png";
    case "hab-sale":
      return "egginc-extras/icon_hab_sale.png";
    case "vehicle-sale":
      return "egginc-extras/icon_vehicle_sale.png";
    case "boost-sale":
      return "egginc-extras/icon_lightning.png";
    case "boost-duration":
      return "egginc-extras/icon_boost_duration.png";
    case "crafting-sale":
      return "egginc/icon_afx_craft.png";
    case "mission-fuel":
      return "egginc/icon_afx_mission.png";
    default:
      return "";
  }
};

const eventCaption = (eventId, multiplier) => {
  switch (eventId) {
    case "prestige-boost":
    case "earnings-boost":
    case "gift-boost":
    case "drone-boost":
    case "boost-duration":
    case "mission-fuel":
      return `${multiplier}x`;
    case "epic-research-sale":
    case "research-sale":
    case "hab-sale":
    case "vehicle-sale":
    case "boost-sale":
    case "crafting-sale":
      return `${Math.round((1 - multiplier) * 100)}% OFF`;
    case "piggy-boost":
      return `+${multiplier}`;
    case "piggy-cap-boost":
      return "NO CAP";
    default:
      return "";
  }
};

const eventBaseColor = eventId => {
  for (let i = 0; i < eventTypes.length; i++) {
    if (eventId === eventTypes[i][0]) {
      return baseColors[i];
    }
  }
  return "gray";
};

const eventFgClass = eventId => `text-${eventBaseColor(eventId)}-500`;

const eventBrightFgClass = eventId => `text-${eventBaseColor(eventId)}-300`;

const eventBgClass = eventId => `bg-${eventBaseColor(eventId)}-500`;

export {
  eventsAndAppUpdates as events,
  eventTypes,
  eventIcon,
  eventCaption,
  eventFgClass,
  eventBrightFgClass,
  eventBgClass,
};
