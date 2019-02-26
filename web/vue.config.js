/**
 * https://cli.vuejs.org/config/#vue-config-js
 */

const path = require('path');
const _ = require('lodash');

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

    // Copy `ReDoc` (https://github.com/Rebilly/ReDoc) to `/static/redoc/bundles/redoc.standalone.js`
    // so that `pkg/server/swagger_handlers.go` can use it as the `RedocURL`.
    //
    // Uses:
    // * https://github.com/webpack-contrib/copy-webpack-plugin
    // * https://webpack.js.org/plugins/copy-webpack-plugin
    config.plugin('copy').tap((opts) => {
      const outputDir = config.output.get('path');
      const pattern = {
        from: path.resolve(__dirname, 'node_modules/redoc/bundles/redoc.standalone.js'),
        to: `${outputDir}/static/redoc/bundles/redoc.standalone.js`,
        toType: 'file',
      };
      const patterns = _.concat(pattern, opts[0]);

      return [
        patterns,
        {
          debug: process.env.NODE_ENV === 'production' ? 'warn' : 'warn',
        },
      ];
    });
  },

  // https://cli.vuejs.org/config/#css-sourcemap
  css: {
    sourceMap: true,
  },
};
