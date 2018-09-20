const utils = require('./utils');
const webpack = require('webpack'); // eslint-disable-line
const config = require('../config');
const merge = require('webpack-merge'); // eslint-disable-line
const baseWebpackConfig = require('./webpack.base.conf');
const HtmlWebpackPlugin = require('html-webpack-plugin'); // eslint-disable-line
const FriendlyErrorsPlugin = require('friendly-errors-webpack-plugin'); // eslint-disable-line
const ExtractTextPlugin = require('extract-text-webpack-plugin'); // eslint-disable-line

// add hot-reload related code to entry chunks
// Object.keys(baseWebpackConfig.entry).forEach((name) => {
//   baseWebpackConfig.entry[name] = ['./build/dev-client'].concat(baseWebpackConfig.entry[name]);
// });

module.exports = merge(baseWebpackConfig, {
  module: {
    // rules: utils.styleLoaders({ sourceMap: config.dev.cssSourceMap })
    rules: utils.styleLoaders({
      sourceMap: config.dev.cssSourceMap,
      extract: true,
    }),
  },
  // cheap-module-eval-source-map is faster for development
  devtool: '#cheap-module-eval-source-map',
  output: {
    path: config.build.assetsRoot,
    filename: utils.assetsPath('js/[name].js'),
    chunkFilename: utils.assetsPath('js/[id].js'),
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': config.dev.env,
    }),
    // https://github.com/glenjamin/webpack-hot-middleware#installation--usage
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin(),
    // https://github.com/ampedandwired/html-webpack-plugin
    // new HtmlWebpackPlugin({
    //   filename: 'index.html',
    //   template: 'index.html',
    //   inject: true
    // }),
    new FriendlyErrorsPlugin(),
    // extract css into its own file
    new ExtractTextPlugin({
      filename: utils.assetsPath('css/[name].css'),
    }),
  ],
});
