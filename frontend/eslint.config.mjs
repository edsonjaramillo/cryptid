import { antfu } from '@antfu/eslint-config';

export default antfu({
  type: 'app',
  typescript: true,
  react: true,
  ignores: ['**/routeTree.gen.ts'],
  stylistic: {
    semi: true,
    quotes: 'single',
    indent: 2,
    overrides: {
      'style/max-len': ['error', { code: 100, ignoreStrings: true }],
    },
  },
  rules: {
    'node/prefer-global/process': 'off',
    'react-refresh/only-export-components': 'off',
    'style/eol-last': 'off',
  },
});
