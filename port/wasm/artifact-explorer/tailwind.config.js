module.exports = {
  purge: ["src/*.html", "src/**/*.vue"],
  darkMode: false,
  theme: {
    extend: {
      maxWidth: {
        column: "12rem",
      },
    },
  },
  variants: {
    extend: {
      opacity: ["disabled"],
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
