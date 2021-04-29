import * as $protobuf from "protobufjs/minimal";

import { ei } from "./proto";
import { ProtobufType } from "./types";
import { uint8ArrayToBinaryString, binaryStringToUint8Array } from "./utils";

/**
 * Decode base64-encoded message.
 * @param message - Message type, e.g. ei.ContractCoopStatusResponse.
 * @param encoded - base64-encoded message, or the corresponding base64-decoded Uint8Array.
 * @param authenticated - Whether to decode as an AuthenticatedMessage.
 * @returns Decoded plain object.
 * @throws Throws on decoding failure.
 */
export function decodeMessage(
  message: ProtobufType,
  encoded: string | Uint8Array,
  authenticated = false
): object {
  if (authenticated) {
    const wrapperPayload = decodeMessage(
      ei.AuthenticatedMessage,
      encoded,
      false
    ) as ei.IAuthenticatedMessage;
    if (wrapperPayload.message === null || wrapperPayload.message === undefined) {
      throw new Error("No message found behind wrap.");
    }
    return decodeMessage(message, wrapperPayload.message, false);
  }

  let binary;
  try {
    if (encoded instanceof Uint8Array) {
      // Note that protobuf.js's toObject automatically base64-decodes the
      // bytes, so we shouldn't double-decode. A misfeature if you ask me,
      // making an assumption by default and not offering an easy way out (using
      // conversion options { bytes: String } disables the base64-decode, or
      // rather, re-encodes, but then the type changes from Uint8Array to
      // String, now TypeScript is not happy).
      binary = uint8ArrayToBinaryString(encoded);
    } else {
      binary = atob(encoded);
    }
  } catch (e) {
    throw new Error(`Error decoding input as base64: ${e}`);
  }
  return message.toObject(message.decode(binaryStringToUint8Array(binary)));
}
