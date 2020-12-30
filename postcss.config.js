module.exports = api => ({
  plugins: [
    require("tailwindcss"),
    require("autoprefixer"),
    ...(api.env === "production"
      ? [
          require("cssnano")({
            preset: "default",
          }),
          require("postcss-hash")({
            algorithm: "sha256",
            trim: 8,
            manifest: "./static/manifest.postcss.json",
          }),
        ]
      : []),
  ],
});
