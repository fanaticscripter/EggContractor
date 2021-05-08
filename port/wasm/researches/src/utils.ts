export function isApplePlatform(): boolean {
  return !!navigator.platform.match(/(Mac|iPhone|iPod|iPad)/i);
}
