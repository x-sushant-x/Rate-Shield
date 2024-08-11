/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        'global-bg': '#D7D6DF',
        'sidebar-bg': '#131719'
      }
    },
  },
  plugins: [],
}

