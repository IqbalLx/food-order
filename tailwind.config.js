/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/*/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography")],
};
