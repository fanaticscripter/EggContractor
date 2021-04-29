import { ProtobufType } from "./types";
import { uint8ArrayToBinaryString } from "./utils";

/**
 * @param {!String} messageName - Name of the protobuf message.
 * @param {!Object} messageObj - Plain object of the payload to be encoded.
 * @param {Boolean} authenticated - Whether to encode as AuthenticatedMessage.
 * @returns {!String}
 * @throws Will throw an error if payload cannot be encoded.
 */

/**
 * Encode message object as base64-encoded protobuf.
 * @param message - Message type, e.g. ei.ContractCoopStatusRequest.
 * @param messageObj - Plain object of payload to be encoded.
 * @param authenticated - Whether to encode as AuthenticatedMessage.
 * @returns
 * @throws Throws on encoding failure.
 */
export function encodeMessage(message: ProtobufType, messageObj: object, authenticated = false) {
  if (authenticated) {
    throw new Error(`Authenticated encoding not implemented.`);
  }

  try {
    const buf = message.encode(messageObj).finish();
    return btoa(uint8ArrayToBinaryString(buf));
  } catch (e) {
    throw new Error(`Encoding ${JSON.stringify(messageObj)}: ${e}.`);
  }
}
