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
    'vue/attribute-hyphenation': 'off',
    'vue/name-property-casing': 'off',
    'vue/prop-name-casing': 'off',
    'max-len': 'off',
    /**
     * Do NOT turn these on as the IDE or `yarn lint`
     * will not warn of these errors if `process.env.NODE_ENV=production` is not set.
     * which is most likely the case.
     */
    // 'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    // 'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',

    // disallow reassignment of function parameters
    // disallow parameter object manipulation except for specific exclusions
    'no-param-reassign': ['error', {
      props: true,
      ignorePropertyModificationsFor: [
        'state', // for vuex state
        'acc', // for reduce accumulators
        'e', // for e.returnvalue
      ],
    }],
  },
  parserOptions: {
    parser: 'babel-eslint',
  },
};
