module.exports = {
  content: ["./src/views/**/*.templ"],
  mode: "jit",
  theme: {
    colors: {
      transparent: 'transparent',
      current: 'currentColor',
      'white': '#ffffff',
      'black': '#000000',
      'errorColor': "#dc2626",
      "grey": {
        50: "#f9fafb",
        100: "#f3f4f6",
        200: "#e5e7eb",
        300: "#d1d5db",
        400: "#9ca3af",
        500: "#6b7280",
        600: "#4b5563",
        700: "#374151",
        800: "#1f2937",
        900: "#111827",
        950: "#030712"
      }
    },
    extend: {
    },
  },
  plugins: [
    require('daisyui'),
    require('@tailwindcss/container-queries'),
  ],
  corePlugins: {
    preflight: true,
  },
  daisyui: {
    themes: [],
  },
};
