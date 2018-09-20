const merge = require('webpack-merge'); // eslint-disable-line
const devEnv = require('./dev.env');

module.exports = merge(devEnv, {
  NODE_ENV: '"test"',
});
