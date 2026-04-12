import tseslint from "@typescript-eslint/eslint-plugin"
import tsParser from "@typescript-eslint/parser"
import prettier from "eslint-config-prettier"
import simpleImportSort from "eslint-plugin-simple-import-sort"
import vue from "eslint-plugin-vue"
import vueParser from "vue-eslint-parser"

/** @type {import("eslint").Linter.FlatConfig[]} */
export default [
  ...vue.configs["flat/recommended"],
  {
    files: ["**/*.{ts,js,vue,mjs}"],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: tsParser,
        ecmaVersion: 2020,
        sourceType: "module"
      }
    },
    plugins: {
      "@typescript-eslint": tseslint,
      "simple-import-sort": simpleImportSort
    },
    rules: {
      "vue/multi-word-component-names": "off",
      "@typescript-eslint/no-unused-vars": "warn",
      "simple-import-sort/imports": "error",
      "simple-import-sort/exports": "error"
    },
    ignores: ["dist/**", "node_modules/**"]
  },
  prettier
]
