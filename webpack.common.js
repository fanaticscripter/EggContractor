const webpack = require("webpack");
const { VueLoaderPlugin } = require("vue-loader");

module.exports = {
  entry: {
    app: "./js/app.js",
    coop: "./js/coop.js",
    events: "./js/events.js",
    peeker: "./js/peeker.js",
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: "vue-loader",
      },
    ],
  },
  plugins: [
    new VueLoaderPlugin(),
    new webpack.DefinePlugin({
      __VUE_OPTIONS_API__: true,
    }),
  ],
  cache: {
    type: "filesystem",
  },
};
