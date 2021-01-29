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

function iconURL(relpath, size) {
  const dir = !size ? "orig" : size.toString();
  return `https://eggincassets.tcl.sh/${dir}/${relpath}`;
}

function stringCmp(s1, s2) {
  return s1 < s2 ? -1 : s1 === s2 ? 0 : 1;
}

export { getLocalStorage, setLocalStorage, iconURL, stringCmp };
