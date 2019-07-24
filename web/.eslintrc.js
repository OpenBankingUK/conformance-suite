module.exports = {
  root: true,
  env: {
    node: true,
    jest: true,
  },
  extends: [
    'plugin:vue/recommended',
    '@vue/airbnb',
  ],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'warn',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'warn',

    'vue/attribute-hyphenation': 'off',
    'vue/name-property-casing': 'off',
    'vue/prop-name-casing': 'off',
    'max-len': 'off',

    // disallow reassignment of function parameters
    // disallow parameter object manipulation except for specific exclusions
    'no-param-reassign': [
      'error',
      {
        props: true,
        ignorePropertyModificationsFor: [
          'state', // for vuex state
          'acc', // for reduce accumulators
          'e', // for e.returnvalue
        ],
      },
    ],

    // Turn off this rather annoying warning.
    camelcase: 'off',
  },
  parserOptions: {
    parser: 'babel-eslint',
  },
};
