module.exports = {
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: 'module',
    project: './tsconfig.json',
    ecmaFeatures: {
      jsx: true,
    },
  },
  settings: {
    react: {
      version: 'detect', // Tells eslint-plugin-react to automatically detect the version of React to use
    },
  },
  plugins: ['lodash', 'react-hooks'],
  extends: [
    'plugin:react-hooks/recommended',
    'plugin:@typescript-eslint/recommended', // Uses the recommended rules from the @typescript-eslint/eslint-plugin
    'prettier', // Uses eslint-config-prettier to disable ESLint rules from @typescript-eslint/eslint-plugin that would conflict with prettier
    'next',
    'plugin:compat/recommended',
    'plugin:prettier/recommended', // Enables eslint-plugin-prettier and eslint-config-prettier. This will display prettier errors as ESLint errors. Make sure this is always the last configuration in the extends array.
  ],
  env: {
    browser: true,
  },
  rules: {
    '@typescript-eslint/explicit-module-boundary-types': 'off',
    'react/react-in-jsx-scope': 'off',
    '@typescript-eslint/restrict-template-expressions': 'error',
    '@typescript-eslint/switch-exhaustiveness-check': 'error',
    'lodash/import-scope': 'error',
    'react/prop-types': 'off',
    '@next/next/no-img-element': 'off',
    'jsx-a11y/alt-text': 'off',
    'no-debugger': 'error',
    'no-warning-comments': ['error', { terms: ['fixme'] }],
    'no-console': ['error', { allow: ['warn', 'error'] }],
    'no-warning-comments': ['error', { terms: ['fixme'] }],
    '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
    'no-extra-boolean-cast': 'error',
    '@typescript-eslint/strict-boolean-expressions': 'error',
  },
}
