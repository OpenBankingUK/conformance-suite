const merge = require('webpack-merge'); // eslint-disable-line
const prodEnv = require('./prod.env');

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
});
