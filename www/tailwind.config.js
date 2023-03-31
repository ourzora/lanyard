/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx}',
    './components/**/*.{js,ts,jsx,tsx}',
    './utils/**/*.{js,ts,jsx,tsx}',
  ],
  // unused currently, but default is to use preferred-color-scheme so we ignore it entirely
  darkMode: 'class',
  theme: {
    screens: {
      sm: '640px',
      lg: '1280px',
    },
    extend: {
      colors: {
        brand: '#6DFA8C',
        'brand-light': '#F0FFEB',
        error: '#FA6D6D',
        'error-light': '#FFECEB',
      },
    },
  },
  plugins: [require('tailwind-aria')],
}
