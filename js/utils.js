/**
 * @param {!string} key
 * @returns {string|undefined}
 */
function getLocalStorage(key) {
  try {
    return localStorage[`${window.location.pathname}_${key}`];
  } catch (err) {
    console.error(err);
    return undefined;
  }
}

/**
 * @param {!string} key
 * @param {any} val
 */
function setLocalStorage(key, val) {
  try {
    localStorage[`${window.location.pathname}_${key}`] = val;
  } catch (err) {
    console.error(err);
  }
}

/**
 * @param {!string} key
 * @returns {string|undefined}
 */
function getSessionStorage(key) {
  try {
    return sessionStorage[`${window.location.pathname}_${key}`];
  } catch (err) {
    console.error(err);
    return undefined;
  }
}

/**
 * @param {!string} key
 * @param {any} val
 */
function setSessionStorage(key, val) {
  try {
    sessionStorage[`${window.location.pathname}_${key}`] = val;
  } catch (err) {
    console.error(err);
  }
}

/**
 * @param {!string} key
 * @returns {string}
 */
function getCurrentQueryParam(key) {
  return new URL(location.href).searchParams.get(key);
}

/**
 * url is modified in place. If val is null or undefined, deletes key.
 * @param {!URL} url
 * @param {!string} key
 * @param {string} val
 */
function setURLQueryParam(url, key, val) {
  if (val === null || val === undefined) {
    url.searchParams.delete(key);
  } else {
    url.searchParams.set(key, val);
  }
}

/**
 * history.replaceState with query param set/replaced or deleted (if val is null
 * or undefined) in the URL.
 * @param {!string} key
 * @param {string} val
 */
function replaceStateSetQueryParam(key, val) {
  const url = new URL(location.href);
  setURLQueryParam(url, key, val);
  history.replaceState(history.state, "", url.toString());
}

window.replaceStateSetQueryParam = replaceStateSetQueryParam;

export {
  setLocalStorage,
  getLocalStorage,
  setSessionStorage,
  getSessionStorage,
  getCurrentQueryParam,
  replaceStateSetQueryParam,
};
