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

export { getLocalStorage, setLocalStorage, iconURL };
