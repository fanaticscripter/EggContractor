function getLocalStorage(key) {
  try {
    return localStorage[`${window.location.pathname}_${key}`];
  } catch (err) {
    console.error(err);
    return undefined;
  }
}

function setLocalStorage(key, val) {
  try {
    localStorage[`${window.location.pathname}_${key}`] = val;
  } catch (err) {
    console.error(err);
  }
}

// Roughly based on https://github.com/tracker1/node-uuid4/blob/v2.0.2/index.js
function uuid4() {
  const rnd = new Uint8Array(16);
  crypto.getRandomValues(rnd);
  rnd[6] = (rnd[6] & 0x0f) | 0x40;
  rnd[8] = (rnd[8] & 0x3f) | 0x80;
  const segs = Array.prototype.map
    .call(rnd, x => x.toString(16).padStart(2, "0"))
    .join("")
    .match(/(.{8})(.{4})(.{4})(.{4})(.{12})/);
  segs.shift();
  return segs.join("-").toUpperCase();
}

export { getLocalStorage, setLocalStorage, uuid4 };
