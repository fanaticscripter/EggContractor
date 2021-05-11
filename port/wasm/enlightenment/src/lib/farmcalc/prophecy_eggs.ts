import { ei } from "../proto";

export function accountProphecyEggsCount(backup: ei.IBackup): number {
  return (
    trophiesProphecyEggsCount(backup) +
    contractsProphecyEggsCount(backup) +
    dailyGiftsProphecyEggsCount(backup)
  );
}

function trophiesProphecyEggsCount(backup: ei.IBackup): number {
  const trophyLevels = backup.game!.eggMedalLevel!;
  if (trophyLevels.length !== 19) {
    throw new Error(`expected trophy levels for 19 eggs, got ${trophyLevels.length}`);
  }
  let count = 0;
  for (let i = 0, egg = ei.Egg.EDIBLE; i < 19; i++, egg++) {
    const level = trophyLevels[i];
    if (egg === ei.Egg.ENLIGHTENMENT) {
      // Enlightenment egg.
      if (level >= 1) {
        count += 1;
      }
      if (level >= 2) {
        count += 2;
      }
      if (level >= 3) {
        count += 3;
      }
      if (level >= 4) {
        count += 5;
      }
      if (level >= 5) {
        count += 10;
      }
    } else {
      // All other eggs offer PE only at diamond, or none at all.
      if (level >= 5) {
        switch (egg) {
          case ei.Egg.EDIBLE:
            count += 5;
            break;
          case ei.Egg.SUPERFOOD:
            count += 4;
            break;
          case ei.Egg.MEDICAL:
            count += 3;
            break;
          case ei.Egg.ROCKET_FUEL:
            count += 2;
            break;
          case ei.Egg.SUPER_MATERIAL:
          case ei.Egg.FUSION:
          case ei.Egg.QUANTUM:
          case ei.Egg.IMMORTALITY:
          case ei.Egg.TACHYON:
            count += 1;
            break;
        }
      }
    }
  }
  return count;
}

function contractsProphecyEggsCount(backup: ei.IBackup): number {
  const contracts = (backup?.contracts?.contracts || []).concat(backup?.contracts?.archive || []);
  if (contracts.length === 0) {
    return 0;
  }
  let count = 0;
  for (const contract of contracts) {
    const props = contract.contract!;
    const league = contract.league || 0;
    let goals = props.goals;
    if (props.goalSets && props.goalSets.length > league) {
      goals = props.goalSets[league].goals;
    }
    if (!goals || goals.length === 0) {
      throw new Error(`no goals found for contract ${props.identifier!}`);
    }
    for (let i = 0; i < contract.numGoalsAchieved!; i++) {
      const goal = goals[i];
      if (goal.rewardType === ei.RewardType.EGGS_OF_PROPHECY) {
        count += Math.round(goal.rewardAmount!);
      }
    }
  }
  return count;
}

function dailyGiftsProphecyEggsCount(backup: ei.IBackup): number {
  return Math.min(Math.floor(backup.game!.numDailyGiftsCollected! / 28), 24);
}
