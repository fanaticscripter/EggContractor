function iconURL(relpath, size) {
  const dir = !size ? "orig" : size.toString();
  return `https://eggincassets.tcl.sh/${dir}/${relpath}`;
}

export { iconURL };
