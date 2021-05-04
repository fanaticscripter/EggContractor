import * as $protobuf from "protobufjs/minimal";
import { ei } from "./proto";

export class ProtobufMessage {
  toJSON(): { [k: string]: any };
}

export class ProtobufType {
  encode(message: object, writer?: $protobuf.Writer): $protobuf.Writer;

  /**
   * @throws {Error} If the payload is not a reader or valid buffer
   * @throws {$protobuf.util.ProtocolError} If required fields are missing
   */
  decode(reader: $protobuf.Reader | Uint8Array, length?: number): ProtobufMessage;

  toObject(message: ProtobufMessage, options?: $protobuf.IConversionOptions): { [k: string]: any };
}

export interface Item {
  key: string;
  afxId: ei.ArtifactSpec.Name;
  afxLevel: ei.ArtifactSpec.Level;
  afxRarity: ei.ArtifactSpec.Rarity;
  name: string;
  rarity: string;
  effectTarget: string;
  effectSize: string;
  effectDelta: number;
  slots: number;
  iconPath: string;
}

export type Stone = Item;

export interface Artifact extends Item {
  stones: Stone[];
  clarityEffect: number;
}

export interface Research {
  id: string;
  name: string;
  maxLevel: number;
  perLevel: number;
  epic?: boolean;
}

export interface ResearchInstance extends Research {
  level: number;
}
