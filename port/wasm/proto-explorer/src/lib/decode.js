import { ei } from "./proto";
import * as $protobuf from "protobufjs/minimal";

/**
 * @typedef {Object} decodeResult
 * @property {Object} payload
 * @property {String} code
 * @property {String} error
 */

/**
 * @param {!String} messageName - Name of the protobuf message.
 * @param {!String} encoded - base64-encoded protobuf payload.
 * @param {Boolean} authenticated - Whether to decode as AuthenticatedMessage.
 * @returns {decodeResult}
 */
function decodeMessage(messageName, encoded, authenticated) {
  authenticated ||= false;
  const message = ei[messageName];
  if (message === undefined) {
    return {
      error: `Message ${messageName} does not exist.`,
    };
  }

  if (authenticated) {
    const wrapperResult = decodeMessage("AuthenticatedMessage", encoded, false);
    if (wrapperResult.error !== undefined) {
      return {
        error: wrapperResult.error,
      };
    }
    const innerResult = decodeMessage(messageName, wrapperResult.payload.message, false);
    innerResult.code = wrapperResult.payload.code;
    return innerResult;
  }

  let binary;
  try {
    binary = atob(encoded);
  } catch (e) {
    return {
      error: "Error decoding input as base64.",
    };
  }
  const buf = new Uint8Array(new ArrayBuffer(binary.length));
  for (let i = 0; i < binary.length; i++) {
    buf[i] = binary.charCodeAt(i);
  }

  try {
    const decoded = message.decode(buf);
    return {
      payload: decoded.toJSON(),
    };
  } catch (e) {
    if (e instanceof $protobuf.util.ProtocolError) {
      const partiallyDecoded = e.instance;
      return {
        payload: partiallyDecoded.toJSON(),
        error: `Partially decoded due to error: ${e}.`,
      };
    } else {
      return {
        error: `Decoding failed with error: ${e}.`,
      };
    }
  }
}

export { decodeMessage };
