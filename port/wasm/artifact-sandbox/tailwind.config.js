const defaultTheme = require("tailwindcss/defaultTheme");

module.exports = {
  purge: ["src/*.html", "src/**/*.vue"],
  darkMode: false,
  theme: {
    // Increase font size by 5% across the board since Always Together is a
    // small and hard to read font.
    fontSize: {
      xs: "calc(1.05 * .75rem)",
      sm: "calc(1.05 * .875rem)",
      base: "calc(1.05 * 1rem)",
      lg: "calc(1.05 * 1.125rem)",
      xl: "calc(1.05 * 1.25rem)",
      "2xl": "calc(1.05 * 1.5rem)",
      "3xl": "calc(1.05 * 1.875rem)",
      "4xl": "calc(1.05 * 2.25rem)",
      "5xl": "calc(1.05 * 3rem)",
      "6xl": "calc(1.05 * 4rem)",
      "7xl": "calc(1.05 * 5rem)",
    },
    extend: {
      fontFamily: {
        default: defaultTheme.fontFamily.sans,
        sans: ["Always Together", ...defaultTheme.fontFamily.sans],
      },
      colors: {
        dark: {
          20: "hsl(0, 0%, 20%)",
          21: "hsl(0, 0%, 21%)",
          22: "hsl(0, 0%, 22%)",
          23: "hsl(0, 0%, 23%)",
          24: "hsl(0, 0%, 24%)",
          25: "hsl(0, 0%, 25%)",
          30: "hsl(0, 0%, 30%)",
          40: "hsl(0, 0%, 40%)",
          45: "hsl(0, 0%, 45%)",
          50: "hsl(0, 0%, 50%)",
          60: "hsl(0, 0%, 60%)",
          70: "hsl(0, 0%, 70%)",
          80: "hsl(0, 0%, 80%)",
          90: "hsl(0, 0%, 90%)",
        },
      },
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
  variants: {
    extend: {
      opacity: ["disabled"],
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
