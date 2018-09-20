module.exports = {
  root: true,
  parserOptions: {
    sourceType: 'module',
    parser: 'babel-eslint'
  },
  env: {
    browser: true,
    mocha: true,
  },
  globals: {
    expect: true,
    gapi: true,
    jest: true,
  },
  plugins: [
    'vue'
  ],
  extends: [
    'eslint:recommended',
    'airbnb-base',
    'plugin:vue/strongly-recommended'
  ],
  settings: {
    'import/resolver': {
      node: {
        extensions: ['.js', '.jsx', '.vue']
      }
    }
  },
  rules: {
    'vue/name-property-casing': ['error', 'kebab-case'],
    semi: [
      'error',
      'always'
    ],
    // don't require .vue extension when importing
    'import/extensions': ['error', 'always', {
      js: 'never',
      vue: 'never'
    }],
    // disallow reassignment of function parameters
    // disallow parameter object manipulation except for specific exclusions
    'no-param-reassign': ['error', {
      props: true,
      ignorePropertyModificationsFor: [
        'state', // for vuex state
        'acc', // for reduce accumulators
        'e' // for e.returnvalue
      ]
    }],
    // allow optionalDependencies
    'import/no-extraneous-dependencies': ['error', {
      optionalDependencies: ['test/unit/index.js']
    }],
    // allow debugger during development
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'vue/attribute-hyphenation': [1, 'never']
  },
  overrides: [
    {
      files: [
        'webpack*.config.js',
        'vue.config.js',
        'rename-assets.js',
        'build/*.js',
        'config/*.js',
      ],
      env: {
        node: true
      }
    },
    {
      files: [
        'package.json'
      ],
      env: {
        node: false,
        browser: false
      }
    }
  ]
};
