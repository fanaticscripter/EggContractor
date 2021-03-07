import { ei } from "./proto";

/**
 * @param {!String} messageName - Name of the protobuf message.
 * @param {!Object} messageObj - Plain object of the payload to be encoded.
 * @param {Boolean} authenticated - Whether to encode as AuthenticatedMessage.
 * @returns {!String}
 * @throws Will throw an error if payload cannot be encoded.
 */
function encodeMessage(messageName, messageObj, authenticated) {
  authenticated ||= false;
  if (authenticated) {
    throw new Error(`Encoding ${messageName}: authenticated encoding not implemented.`);
  }

  const message = ei[messageName];
  if (message === undefined) {
    throw new Error(`Message ${messageName} does not exist.`);
  }

  try {
    const buf = message.encode(messageObj).finish();
    return btoa(String.fromCharCode(...buf));
  } catch (e) {
    throw new Error(`Encoding ${messageName} ${JSON.stringify(messageObj)}: ${e}.`);
  }
}

export { encodeMessage };
