module.exports = {
  purge: [
    "src/*.html",
    "src/**/*.vue",
  ],
  darkMode: false,
  theme: {
    extend: {
      screens: {
        "3xl": "1600px",
      },
      maxWidth: {
        "10xl": "104rem",
      },
    }
  },
  variants: {
    extend: {
      opacity: ["disabled"],
    }
  },
  plugins: [
    require("@tailwindcss/forms"),
  ],
};
