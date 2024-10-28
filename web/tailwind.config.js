/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./tmpl/*.gohtml", "./tmpl/partials/*.gohtml"],
  darkMode: [
    "class",
    "@media (prefers-color-scheme: dark) { &:not(.light *) }",
  ],
};
