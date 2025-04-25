import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const tailwindStylesheetPath = resolve(__dirname, 'src/styles.css');

export default {
  printWidth: 100,
  quoteProps: 'preserve',
  singleQuote: true,
  semi: true,
  bracketSameLine: true,
  tailwindStylesheet: tailwindStylesheetPath,
  tailwindAttributes: ['cn', 'cva'],
  tailwindFunctions: ['cn', 'cva'],
  plugins: ['prettier-plugin-tailwindcss'],
};
