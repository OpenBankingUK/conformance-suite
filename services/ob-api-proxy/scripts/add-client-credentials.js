const error = require('debug')('error');
const { updateClientCredentials } = require('../app/authorisation-servers');

const args = process.argv.slice(2).reduce((acc, arg) => {
  const [k, v = true] = arg.split('=');
  acc[k] = v;
  return acc;
}, {});

const addClientCredentials = async () => {
  if (!args.authServerId || !args.clientId || !args.clientSecret) {
    throw new Error('authServerId, clientId, and clientSecret must ALL be present!');
  }
  try {
    await updateClientCredentials(args.authServerId, {
      clientId: args.clientId,
      clientSecret: args.clientSecret,
    });
  } catch (err) {
    error('addClientCredentials, err=%o', err);
  }
};

addClientCredentials().then(() => {
  if (process.env.NODE_ENV !== 'test') {
    process.exit();
  }
});

exports.addClientCredentials = addClientCredentials;
