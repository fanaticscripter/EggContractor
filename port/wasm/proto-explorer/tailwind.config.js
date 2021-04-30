const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  purge: ["src/*.html", "src/**/*.vue"],
  darkMode: false,
  variants: {
    extend: {
      opacity: ["disabled"],
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
