export type ItemId =
  | "lunar-totem-1"
  | "lunar-totem-2"
  | "lunar-totem-3"
  | "lunar-totem-4"
  | "neodymium-medallion-1"
  | "neodymium-medallion-2"
  | "neodymium-medallion-3"
  | "neodymium-medallion-4"
  | "beak-of-midas-1"
  | "beak-of-midas-2"
  | "beak-of-midas-3"
  | "beak-of-midas-4"
  | "light-of-eggendil-1"
  | "light-of-eggendil-2"
  | "light-of-eggendil-3"
  | "light-of-eggendil-4"
  | "demeters-necklace-1"
  | "demeters-necklace-2"
  | "demeters-necklace-3"
  | "demeters-necklace-4"
  | "vial-of-martian-dust-1"
  | "vial-of-martian-dust-2"
  | "vial-of-martian-dust-3"
  | "vial-of-martian-dust-4"
  | "gusset-1"
  | "gusset-2"
  | "gusset-3"
  | "gusset-4"
  | "the-chalice-1"
  | "the-chalice-2"
  | "the-chalice-3"
  | "the-chalice-4"
  | "book-of-basan-1"
  | "book-of-basan-2"
  | "book-of-basan-3"
  | "book-of-basan-4"
  | "phoenix-feather-1"
  | "phoenix-feather-2"
  | "phoenix-feather-3"
  | "phoenix-feather-4"
  | "tungsten-ankh-1"
  | "tungsten-ankh-2"
  | "tungsten-ankh-3"
  | "tungsten-ankh-4"
  | "aurelian-brooch-1"
  | "aurelian-brooch-2"
  | "aurelian-brooch-3"
  | "aurelian-brooch-4"
  | "carved-rainstick-1"
  | "carved-rainstick-2"
  | "carved-rainstick-3"
  | "carved-rainstick-4"
  | "puzzle-cube-1"
  | "puzzle-cube-2"
  | "puzzle-cube-3"
  | "puzzle-cube-4"
  | "quantum-metronome-1"
  | "quantum-metronome-2"
  | "quantum-metronome-3"
  | "quantum-metronome-4"
  | "ship-in-a-bottle-1"
  | "ship-in-a-bottle-2"
  | "ship-in-a-bottle-3"
  | "ship-in-a-bottle-4"
  | "tachyon-deflector-1"
  | "tachyon-deflector-2"
  | "tachyon-deflector-3"
  | "tachyon-deflector-4"
  | "interstellar-compass-1"
  | "interstellar-compass-2"
  | "interstellar-compass-3"
  | "interstellar-compass-4"
  | "dilithium-monocle-1"
  | "dilithium-monocle-2"
  | "dilithium-monocle-3"
  | "dilithium-monocle-4"
  | "titanium-actuator-1"
  | "titanium-actuator-2"
  | "titanium-actuator-3"
  | "titanium-actuator-4"
  | "mercurys-lens-1"
  | "mercurys-lens-2"
  | "mercurys-lens-3"
  | "mercurys-lens-4"
  | "tachyon-stone-2"
  | "tachyon-stone-3"
  | "tachyon-stone-4"
  | "dilithium-stone-2"
  | "dilithium-stone-3"
  | "dilithium-stone-4"
  | "shell-stone-2"
  | "shell-stone-3"
  | "shell-stone-4"
  | "lunar-stone-2"
  | "lunar-stone-3"
  | "lunar-stone-4"
  | "soul-stone-2"
  | "soul-stone-3"
  | "soul-stone-4"
  | "prophecy-stone-2"
  | "prophecy-stone-3"
  | "prophecy-stone-4"
  | "quantum-stone-2"
  | "quantum-stone-3"
  | "quantum-stone-4"
  | "terra-stone-2"
  | "terra-stone-3"
  | "terra-stone-4"
  | "life-stone-2"
  | "life-stone-3"
  | "life-stone-4"
  | "clarity-stone-2"
  | "clarity-stone-3"
  | "clarity-stone-4"
  | "tachyon-stone-1"
  | "dilithium-stone-1"
  | "shell-stone-1"
  | "lunar-stone-1"
  | "soul-stone-1"
  | "prophecy-stone-1"
  | "quantum-stone-1"
  | "terra-stone-1"
  | "life-stone-1"
  | "clarity-stone-1"
  | "gold-meteorite-1"
  | "gold-meteorite-2"
  | "gold-meteorite-3"
  | "tau-ceti-geode-1"
  | "tau-ceti-geode-2"
  | "tau-ceti-geode-3"
  | "solar-titanium-1"
  | "solar-titanium-2"
  | "solar-titanium-3";

export type MissionId =
  | "chicken-one-short"
  | "chicken-one-standard"
  | "chicken-one-extended"
  | "chicken-nine-short"
  | "chicken-nine-standard"
  | "chicken-nine-extended"
  | "chicken-heavy-short"
  | "chicken-heavy-standard"
  | "chicken-heavy-extended"
  | "bcr-short"
  | "bcr-standard"
  | "bcr-extended"
  | "quintillion-chicken-short"
  | "quintillion-chicken-standard"
  | "quintillion-chicken-extended"
  | "cornish-hen-corvette-short"
  | "cornish-hen-corvette-standard"
  | "cornish-hen-corvette-extended"
  | "galeggtica-short"
  | "galeggtica-standard"
  | "galeggtica-extended"
  | "defihent-short"
  | "defihent-standard"
  | "defihent-extended"
  | "voyegger-short"
  | "voyegger-standard"
  | "voyegger-extended"
  | "henerprise-short"
  | "henerprise-standard"
  | "henerprise-extended";

export interface Item {
  id: ItemId;
  afx_id: number;
  afx_level: number;
  name: string;
  tier_number: number;
  tier_name: string;
  afx_type: number;
  type: string;
  icon_filename: string;
  display: string;
  iconPath: string;
  recipe:
    | {
        id: ItemId;
        count: number;
      }[]
    | null;
}

export interface Mission {
  id: MissionId;
  display: string;
  iconPath: string;
  capcity: number;
  durationSeconds: number;
  count: number;
  lootTotal: number;
  loot: {
    id: ItemId;
    count: number;
  }[];
}

export interface SimulationProgress {
  totalTrials: number;
  finishedTrials: number;
  successfulTrials: number;
  secondsElapsed: number;
  stopped: boolean;
}

export type ItemSpec = { id: ItemId; count: number };
export type ItemSelectSpec = { id: ItemId | null; count: number; rowid: string };
export type MissionSpec = { id: MissionId; count: number };
export type MissionSelectSpec = { id: MissionId | null; count: number; rowid: string };

export type SimulationReportCallback = (progress: SimulationProgress, missing?: ItemSpec[]) => void;

export interface SimulationWorkerInterface {
  // ping should return immediately, indicating that the worker is responsive.
  ping(): void;
  runSimulations(
    missions: MissionSpec[],
    targets: ItemSpec[],
    totalTrials: number,
    report: SimulationReportCallback,
    seed: string,
    signal?: AbortSignal
  ): Promise<void>;
}
