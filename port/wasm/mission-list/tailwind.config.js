module.exports = {
  purge: [
    "templates/**/*.html",
  ],
  darkMode: false,
  theme: {
    extend: {
      maxWidth: {
        "12xl": "120rem",
      }
    }
  },
  plugins: [
    require("@tailwindcss/forms"),
  ],
};
