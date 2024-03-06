/** @type {import('tailwindcss').Config}*/
const config = {
  content: [
    "./src/**/*.{html,js,svelte,ts}",
    "./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}",
  ],

  plugins: [require("flowbite/plugin")],

  darkMode: "class",

  theme: {
    extend: {
      colors: {
        // flowbite-svelte
        primary: {
          50: "#f0f0fd",
          100: "#e3e4fc",
          200: "#cecdf8",
          300: "#b0aef3",
          400: "#978dec",
          500: "#8772e2",
          600: "#7857d4",
          700: "#6847bb",
          800: "#553c97",
          900: "#463778",
          950: "#322653",
        },
      },
    },
  },
};

module.exports = config;
