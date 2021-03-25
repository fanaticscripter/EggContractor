const colors = require("tailwindcss/colors");

module.exports = {
  purge: ["src/*.html", "src/**/*.js", "src/**/*.vue"],
  darkMode: false,
  theme: {
    extend: {
      colors: {
        ...colors,
        "light-blue": colors.lightBlue,
      },
      screens: {
        "2col": "1512px",
        "3col": "2280px",
      },
      width: {
        6.5: "1.625rem",
      },
      maxWidth: {
        "2.5xl": "45rem",
      },
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
