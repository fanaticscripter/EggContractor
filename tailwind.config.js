module.exports = {
  purge: [
    "templates/**/*.html",
    "js/**/*.vue",
  ],
  darkMode: false,
  theme: {
    extend: {
      cursor: {
        help: "help",
      },
      minHeight: {
        // min-h-stretch fixes 100vh != actual viewport height issue in iOS Safari.
        // autoprefixer should expand stretch to -webkit-fill-available for Chrome/Safari.
        // https://allthingssmitty.com/2020/05/11/css-fix-for-100vh-in-mobile-webkit/
        stretch: "stretch",
      },
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
  ],
};
