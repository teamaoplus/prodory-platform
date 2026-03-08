/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      colors: {
        primary: {
          50: '#E0F7FF',
          100: '#B3EBFF',
          200: '#80DEFF',
          300: '#4DD1FF',
          400: '#26C7FF',
          500: '#0B7DF0',
          600: '#0968C8',
          700: '#0753A0',
          800: '#053E78',
          900: '#032950',
        },
        accent: {
          50: '#E0FCFF',
          100: '#B3F7FF',
          200: '#80F2FF',
          300: '#4DEDFF',
          400: '#26E8FF',
          500: '#06B6D4',
          600: '#0598B2',
          700: '#047A8F',
          800: '#035C6C',
          900: '#023E49',
        },
        success: '#10B981',
        warning: '#F59E0B',
        danger: '#EF4444',
      },
    },
  },
  plugins: [],
}
