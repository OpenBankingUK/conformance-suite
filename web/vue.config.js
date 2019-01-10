/**
 * https://cli.vuejs.org/config/#vue-config-js
 */

// vue.config.js
module.exports = {
  // options...

  // https://cli.vuejs.org/config/#devserver
  // https://webpack.js.org/configuration/dev-server/
  //
  // These settings control the server that is started when you do `yarn serve`.
  devServer: {
    clientLogLevel: 'info',
    compress: true,
    overlay: {
      warnings: false,
      errors: true,
    },
    open: 'Google Chrome',
  },

  // https://cli.vuejs.org/config/#runtimecompiler
  runtimeCompiler: true,

  // https://cli.vuejs.org/config/#chainwebpack
  chainWebpack: (config) => {
    // Inline fonts and images so we don't do another fetch for them.
    // If we set `limit` to zero, all the fonts and images are inlined.
    // Explanation can be found in:
    // * https://cli.vuejs.org/guide/webpack.html#replacing-loaders-of-a-rule
    // * https://github.com/vuejs/vue-cli/issues/3215

    config.module.rule('fonts').use('url-loader').tap((opts) => {
      const options = Object.assign(opts, { limit: 0 });
      return options;
    });

    config.module.rule('images').use('url-loader').tap((opts) => {
      const options = Object.assign(opts, { limit: 0 });
      return options;
    });
  },

  // https://cli.vuejs.org/config/#css-sourcemap
  css: {
    sourceMap: true,
  },
};
