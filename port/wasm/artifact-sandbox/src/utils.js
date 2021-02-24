function iconURL(relpath, size) {
  const dir = !size ? "orig" : size.toString();
  return `https://eggincassets.tcl.sh/${dir}/${relpath}`;
}

function stringCmp(s1, s2) {
  return s1 < s2 ? -1 : s1 === s2 ? 0 : 1;
}

export { iconURL, stringCmp };
