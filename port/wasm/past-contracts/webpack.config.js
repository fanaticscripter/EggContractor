const fs = require("fs");
const path = require("path");

const webpack = require("webpack");
const TerserPlugin = require("terser-webpack-plugin");
const { VueLoaderPlugin } = require("vue-loader");
const { WebpackManifestPlugin } = require("webpack-manifest-plugin");

const wasmManifest = JSON.parse(fs.readFileSync("dist/manifest.wasm.json"));
const wasmFile = wasmManifest["app.wasm"];
if (!wasmFile) {
  console.error("app.wasm mapping not found in dist/manifest.wasm.json");
  process.exit(1);
}

module.exports = {
  mode: "production",
  entry: {
    app: "./src/app.js",
  },
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "[name].[contenthash:8].js",
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: "vue-loader",
      },
      {
        test: /\.css$/,
        use: ["style-loader", "css-loader"],
      },
    ],
  },
  plugins: [
    new VueLoaderPlugin(),
    new webpack.DefinePlugin({
      __WASM_FILE__: JSON.stringify(wasmFile),
      __VUE_OPTIONS_API__: true,
    }),
    new WebpackManifestPlugin({
      fileName: "manifest.webpack.json",
      publicPath: "",
    }),
  ],
  externals: {
    Go: "Go",
  },
  optimization: {
    minimize: true,
    minimizer: [
      new TerserPlugin({
        extractComments: false,
        terserOptions: {
          format: {
            comments: false,
          },
        },
      }),
    ],
  },
};
