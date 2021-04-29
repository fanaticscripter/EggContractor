import { decodeMessage } from "./decode";
import { encodeMessage } from "./encode";
import { ei } from "./proto";
import { APP_VERSION, APP_BUILD, CLIENT_VERSION, PLATFORM, PLATFORM_STRING } from "./version";

const API_ROOT =
  import.meta.env.DEV && import.meta.env.VITE_APP_MOCK
    ? "/api"
    : "https://wasmegg.zw.workers.dev/?url=https://afx-2-dot-auxbrainhome.appspot.com";
const TIMEOUT = 5000;

/**
 * Makes an API request.
 * @param endpoint - Path of API endpoint, e.g. /ei/coop_status.
 * @param encodedPayload - base64-encoded request payload.
 * @returns base64-encoded response payload.
 * @throws Throws an error on network failure (including timeout) or non-2XX response.
 */
export async function request(endpoint: string, encodedPayload: string) {
  const controller = new AbortController();
  setTimeout(() => controller.abort(), TIMEOUT);
  const url = API_ROOT + endpoint;
  try {
    const resp = await fetch(url, {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: `data=${encodedPayload}`,
      signal: controller.signal,
    });
    const text = await resp.text();
    if (resp.status < 200 || resp.status >= 300) {
      throw new Error(`HTTP ${resp.status}: ${text}`);
    }
    return text;
  } catch (e) {
    if (e.name === "AbortError") {
      throw new Error(`POST ${url} data=${encodedPayload}: timeout after ${TIMEOUT}ms.`);
    } else {
      throw new Error(`POST ${url} data=${encodedPayload}: ${e}`);
    }
  }
}

/**
 * @param userId
 * @returns
 * @throws
 */
export async function requestFirstContact(userId: string) {
  const requestPayload: ei.IEggIncFirstContactRequest = {
    rinfo: basicRequestInfo(""),
    eiUserId: userId,
    clientVersion: CLIENT_VERSION,
    platform: PLATFORM,
  };
  const encodedRequestPayload = encodeMessage(ei.EggIncFirstContactRequest, requestPayload);
  const encodedResponsePayload = await request("/ei/first_contact", encodedRequestPayload);
  return decodeMessage(
    ei.EggIncFirstContactResponse,
    encodedResponsePayload,
    true
  ) as ei.IEggIncFirstContactResponse;
}

export function basicRequestInfo(userId: string): ei.IBasicRequestInfo {
  return {
    eiUserId: userId,
    clientVersion: CLIENT_VERSION,
    version: APP_VERSION,
    build: APP_BUILD,
    platform: PLATFORM_STRING,
  };
}
