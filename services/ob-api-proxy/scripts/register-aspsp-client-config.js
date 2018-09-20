const debug = require('debug')('debug');
const error = require('debug')('error');
const { updateRegisteredConfig } = require('../app/authorisation-servers');

const parseArgs = rawArgs => rawArgs.reduce((acc, arg) => {
  const [k, v = true] = arg.split('=');
  acc[k] = v;
  return acc;
}, {});

const parseValue = (value) => {
  let result;
  try {
    result = JSON.parse(value);
    debug(`parseValue: result: ${result}`);
  } catch (e) {
    debug(`parseValue: error: ${e}`);
    result = value;
  }
  return result;
};

const registerAgreedConfig = async (args) => {
  const { authServerId, field, value } = parseArgs(args);
  if (!authServerId || !field || !value) {
    throw new Error('authServerId, field, and value must ALL be present!');
  }

  try {
    const config = {};
    config[field] = parseValue(value);
    debug(`config: ${JSON.stringify(config)}`);

    await updateRegisteredConfig(authServerId, config);
  } catch (e) {
    error(e);
  }
};

const exit = (env) => {
  if (env !== 'test') {
    process.exit();
  }
};

registerAgreedConfig(process.argv.slice(2))
  .then(() => exit(process.env.NODE_ENV))
  .catch((err) => {
    error(err);
    exit(process.env.NODE_ENV);
  });

exports.registerAgreedConfig = registerAgreedConfig;
