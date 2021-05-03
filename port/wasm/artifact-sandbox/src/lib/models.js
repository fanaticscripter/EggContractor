import {
  artifactFromId,
  artifactFromAfxIdLevelRarity,
  stoneFromId,
  stoneFromAfxIdLevel,
} from "./data";

class Builds {
  /**
   * @param {!Array<!Build>} builds
   * @param {!Config} config
   */
  constructor(builds, config) {
    this.builds = builds;
    this.config = config;
  }

  /**
   * @returns {!Builds}
   */
  static newDefaultBuilds() {
    return new Builds([Build.newEmptyBuild()], Config.newDefaultConfig());
  }

  /**
   * @param {!String} serialized
   * @returns {!Builds}
   * @throws Will throw an error if the serialized string cannot be parsed.
   */
  static deserialize(serialized) {
    const binary = atob(serialized);
    const buf = new Uint8Array(new ArrayBuffer(binary.length));
    for (let i = 0; i < binary.length; i++) {
      buf[i] = binary.charCodeAt(i);
    }
    const builds = proto.Builds.deserializeBinary(buf);
    return Builds.fromProto(builds);
  }

  /**
   * @param {!proto.Builds} builds
   * @returns {!Builds}
   */
  static fromProto(builds) {
    return new Builds(
      builds.getBuildsList().map(Build.fromProto),
      Config.fromProto(builds.getConfig())
    );
  }

  /**
   * @returns {!proto.Builds}
   */
  toProto() {
    const builds = new proto.Builds();
    builds.setBuildsList(this.builds.map(b => b.toProto()));
    builds.setConfig(this.config.toProto());
    return builds;
  }

  /**
   * @returns {String}
   */
  serialize() {
    return serialize(this.toProto());
  }
}

class Build {
  /**
   * @param {!Array<!Artifact>} artifacts
   */
  constructor(artifacts) {
    this.artifacts = artifacts;
  }

  /**
   * @param {!Object} props
   * @returns {!Build}
   */
  static fromBuildProps(props) {
    return new Build(props.map(Artifact.fromBuildProps));
  }

  /**
   * @returns {!Build}
   */
  static newEmptyBuild() {
    return new Build([
      Artifact.newEmptyArtifact(),
      Artifact.newEmptyArtifact(),
      Artifact.newEmptyArtifact(),
      Artifact.newEmptyArtifact(),
    ]);
  }

  /**
   * @param {!proto.Build} build
   * @returns {!Build}
   */
  static fromProto(build) {
    return new Build(build.getArtifactsList().map(Artifact.fromProto));
  }

  buildProps() {
    return this.artifacts.map(a => a.buildProps());
  }

  /**
   * @returns {!proto.Build}
   */
  toProto() {
    const build = new proto.Build();
    build.setArtifactsList(this.artifacts.map(a => a.toProto()));
    return build;
  }

  /**
   * @returns {!Boolean}
   */
  isEmpty() {
    for (const artifact of this.artifacts) {
      if (!artifact.isEmpty()) {
        return false;
      }
    }
    return true;
  }

  /**
   * @returns {!Boolean}
   */
  hasDuplicates() {
    const families = new Set();
    for (const artifact of this.artifacts) {
      if (artifact.isEmpty()) {
        continue;
      }
      if (families.has(artifact.family_id)) {
        return true;
      }
      families.add(artifact.family_id);
    }
  }
}

class Artifact {
  /**
   * @param {!Object} artifactProps
   * @param {!Array<Stone>} stones
   */
  constructor(artifactProps, stones) {
    if (artifactProps === null) {
      this.props = null;
      this.id = "";
    } else {
      this.props = artifactProps;
      for (const prop in artifactProps) {
        Object.defineProperty(this, prop, {
          value: artifactProps[prop],
          writable: false,
        });
      }
    }
    this.stones = stones;
  }

  /**
   * @param {!Object} props
   * @returns {!Artifact}
   */
  static fromBuildProps(props) {
    const artifactProps = artifactFromId(props.id);
    return new Artifact(artifactProps, props.stones.map(Stone.fromId));
  }

  /**
   * @returns {!Artifact}
   */
  static newEmptyArtifact() {
    return new Artifact(null, [null, null, null]);
  }

  /**
   * @param {!proto.Artifact} artifact
   * @returns {!Artifact}
   */
  static fromProto(artifact) {
    const artifactProps = artifact.getIsEmpty()
      ? null
      : artifactFromAfxIdLevelRarity(
          artifact.getAfxId(),
          artifact.getAfxLevel(),
          artifact.getAfxRarity()
        );
    const stones = artifact.getStonesList().map(Stone.fromProto);
    return new Artifact(artifactProps, stones);
  }

  /**
   * @returns {!Object}
   */
  buildProps() {
    return {
      id: this.id,
      stones: this.stones.map(s => (s === null ? "" : s.id)),
    };
  }

  /**
   * @returns {!proto.Artifact}
   */
  toProto() {
    const artifact = new proto.Artifact();
    if (this.isEmpty()) {
      artifact.setIsEmpty(true);
    } else {
      artifact.setIsEmpty(false);
      artifact.setAfxId(this.afx_id);
      artifact.setAfxLevel(this.afx_level);
      artifact.setAfxRarity(this.afx_rarity);
    }
    artifact.setStonesList(
      this.stones.map(s => {
        if (s === null) {
          const stone = new proto.Stone();
          stone.setIsEmpty(true);
          return stone;
        } else {
          return s.toProto();
        }
      })
    );
    return artifact;
  }

  /**
   * @returns {!Boolean}
   */
  isEmpty() {
    return this.id === "";
  }

  /**
   * @returns {!Array<Stone>}
   */
  get activeStones() {
    if (this.id === "") {
      return [];
    }
    return this.stones.slice(0, this.slots).filter(s => s !== null);
  }

  /**
   * @returns {!Number}
   */
  get clarityEffect() {
    if (this.isEmpty()) {
      return 0;
    }
    if (this.afx_id === proto.ArtifactSpec.Name.LIGHT_OF_EGGENDIL) {
      return 1;
    }
    let effect = 0;
    for (const stone of this.activeStones) {
      if (stone === null) {
        continue;
      }
      if (stone.afx_id === proto.ArtifactSpec.Name.CLARITY_STONE) {
        effect += stone.effect_delta;
      }
    }
    return effect;
  }

  /**
   * @returns {!Boolean}
   */
  hasClarityStones() {
    if (this.isEmpty()) {
      return false;
    }
    for (const stone of this.activeStones) {
      if (stone === null) {
        continue;
      }
      if (stone.afx_id === proto.ArtifactSpec.Name.CLARITY_STONE) {
        return true;
      }
    }
    return false;
  }

  /**
   * @returns {!Boolean}
   */
  isEffectiveOnEnlightenment() {
    if (this.isEmpty()) {
      return false;
    }
    return this.afx_id === proto.ArtifactSpec.Name.LIGHT_OF_EGGENDIL || this.clarityEffect > 0;
  }

  /**
   * @returns {!Boolean}
   */
  isEffectiveOnRegular() {
    if (this.isEmpty()) {
      return false;
    }
    return this.afx_id !== proto.ArtifactSpec.Name.LIGHT_OF_EGGENDIL;
  }
}

class Stone {
  /**
   * @param {!Object} stoneProps
   */
  constructor(stoneProps) {
    this.props = stoneProps;
    for (const prop in stoneProps) {
      Object.defineProperty(this, prop, {
        value: stoneProps[prop],
        writable: false,
      });
    }
  }

  /**
   * @param {!String} id
   * @returns {!Stone}
   */
  static fromId(id) {
    const stoneProps = stoneFromId(id);
    if (stoneProps === null) {
      return null;
    }
    return new Stone(stoneProps);
  }

  /**
   * @param {!proto.Stone} stone
   * @returns {!Stone}
   */
  static fromProto(stone) {
    if (stone.getIsEmpty()) {
      return null;
    }
    const stoneProps = stoneFromAfxIdLevel(stone.getAfxId(), stone.getAfxLevel());
    if (stoneProps === null) {
      return null;
    }
    return new Stone(stoneProps);
  }

  /**
   * @returns {!proto.Stone}
   */
  toProto() {
    const stone = new proto.Stone();
    stone.setIsEmpty(false);
    stone.setAfxId(this.afx_id);
    stone.setAfxLevel(this.afx_level);
    return stone;
  }
}

const maxSoulFood = 140;
const maxProphecyBonus = 5;

class Config {
  constructor() {
    this.prophecyEggs = 0;
    this.soulEggs = 0;
    this.soulEggsInput = "";
    this.isEnlightenment = false;
    this.soulFood = maxSoulFood;
    this.prophecyBonus = maxProphecyBonus;
    this.birdFeedActive = false;
    this.tachyonPrismActive = false;
    this.soulBeaconActive = false;
    this.boostBeaconActive = false;
    this.tachyonDeflectorBonus = 0;
  }

  /**
   * @returns {!Config}
   */
  static newDefaultConfig() {
    const self = new Config();
    self.prophecyEggs = 1;
    self.soulEggs = 250;
    self.soulEggsInput = "250";
    return self;
  }

  /**
   * @param {!proto.Config} config
   * @returns {!Config}
   */
  static fromProto(config) {
    const self = new Config();
    self.prophecyEggs = config.getProphecyEggs();
    self.soulEggs = config.getSoulEggs();
    self.soulEggsInput = config.getSoulEggsInput();
    self.isEnlightenment = config.getIsEnlightenment();
    self.soulFood = maxSoulFood - config.getMissingSoulFood();
    self.prophecyBonus = maxProphecyBonus - config.getMissingProphecyBonus();
    self.birdFeedActive = config.getBirdFeedActive();
    self.tachyonPrismActive = config.getTachyonPrismActive();
    self.soulBeaconActive = config.getSoulBeaconActive();
    self.boostBeaconActive = config.getBoostBeaconActive();
    self.tachyonDeflectorBonus = config.getTachyonDeflectorBonus();
    return self;
  }

  /**
   * @returns {!proto.Config}
   */
  toProto() {
    const config = new proto.Config();
    config.setProphecyEggs(this.prophecyEggs);
    config.setSoulEggs(this.soulEggs);
    config.setSoulEggsInput(this.soulEggsInput);
    config.setIsEnlightenment(this.isEnlightenment);
    config.setMissingSoulFood(maxSoulFood - this.soulFood);
    config.setMissingProphecyBonus(maxProphecyBonus - this.prophecyBonus);
    config.setBirdFeedActive(this.birdFeedActive);
    config.setTachyonPrismActive(this.tachyonPrismActive);
    config.setSoulBeaconActive(this.soulBeaconActive);
    config.setBoostBeaconActive(this.boostBeaconActive);
    config.setTachyonDeflectorBonus(this.tachyonDeflectorBonus);
    return config;
  }

  /**
   * @returns {!Boolean}
   */
  epicResearchMaxed() {
    return this.soulFood === maxSoulFood && this.prophecyBonus === maxProphecyBonus;
  }

  /**
   * @returns {!Boolean}
   */
  anyBoostActive() {
    return (
      this.birdFeedActive ||
      this.tachyonPrismActive ||
      this.soulBeaconActive ||
      this.boostBeaconActive
    );
  }
}

function serialize(msg) {
  return btoa(String.fromCharCode(...msg.serializeBinary()));
}

export { Builds, Build, Artifact, Stone, Config };
