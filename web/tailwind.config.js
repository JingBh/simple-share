import defaultTheme from 'tailwindcss/defaultTheme'

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx,vue}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Noto Sans SC', ...defaultTheme.fontFamily.sans],
        name: ['Montserrat', 'Noto Sans SC', ...defaultTheme.fontFamily.sans],
        title: ['Kanit', 'Noto Sans SC', ...defaultTheme.fontFamily.sans],
      }
    },
  },
  plugins: [],
}
