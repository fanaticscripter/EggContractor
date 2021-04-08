const API_ROOT = "https://wasmegg.zw.workers.dev/?url=https://afx-2-dot-auxbrainhome.appspot.com";
const TIMEOUT = 10000;

class APIClient {
  constructor() {
    /**
     * @type {!Array<!AbortController>}
     */
    this._controllers = [];
  }

  /**
   * @async
   * @param {!String} endpoint
   * @param {!String} encodedPayload
   * @param {Boolean} flush - Whether to flush (abort) existing requests.
   * @returns {Promise<!String>}
   * @throws Will throw an error if request fails with a non-2XX response or times out.
   */
  async request(endpoint, encodedPayload, flush) {
    if (flush) {
      this._flush();
    }

    const controller = new AbortController();
    this._controllers.push(controller);
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
    } finally {
      this._popController(controller);
    }
  }

  destroy() {
    this._flush();
  }

  _popController(controller) {
    for (let i = 0; i < this._controllers.length; i++) {
      if (this._controllers[i] === controller) {
        this._controllers.splice(i, 1);
        break;
      }
    }
  }

  _flush() {
    while (this._controllers.length > 0) {
      const controller = this._controllers.pop();
      controller.abort();
    }
  }
}

export { APIClient };
