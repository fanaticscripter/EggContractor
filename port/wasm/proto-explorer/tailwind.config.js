const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  purge: [
    "src/*.html",
    "src/**/*.vue",
  ],
  darkMode: false,
  theme: {
    extend: {
      // h-stretch and min-h-stretch fix the 100vh != actual viewport height
      // issue in iOS Safari. autoprefixer should expand stretch to
      // -webkit-fill-available for Chrome/Safari.
      // https://allthingssmitty.com/2020/05/11/css-fix-for-100vh-in-mobile-webkit/
      height: {
        stretch: "stretch",
      },
      minHeight: {
        stretch: "stretch",
      },
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
  ],
};
