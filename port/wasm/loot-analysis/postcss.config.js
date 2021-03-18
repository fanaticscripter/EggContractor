const path = require("path");

module.exports = {
  plugins: [
    require("@tailwindcss/jit"),
    require("postcss-import"),
    require("tailwindcss"),
    require("autoprefixer"),
    require("cssnano")({
      preset: "default",
    }),
    require("postcss-hash")({
      algorithm: "sha256",
      trim: 8,
      manifest: "./dist/manifest.postcss.json",
    }),
  ],
};
