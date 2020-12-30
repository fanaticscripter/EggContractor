const path = require("path");

const { WebpackManifestPlugin } = require("webpack-manifest-plugin");

const { merge } = require("webpack-merge");
const common = require("./webpack.common.js");

module.exports = merge(common, {
  mode: "production",
  output: {
    path: path.resolve(__dirname, "static"),
    filename: "[name].[contenthash:8].js",
  },
  plugins: [
    new WebpackManifestPlugin({
      fileName: "manifest.webpack.json",
      publicPath: "",
    }),
  ],
});
