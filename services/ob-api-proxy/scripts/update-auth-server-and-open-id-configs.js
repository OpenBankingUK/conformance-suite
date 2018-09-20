const log = require('debug')('log');
const { updateOpenIdConfigs } = require('../app/authorisation-servers');
const { fetchOBAccountPaymentServiceProviders } = require('../app/ob-directory');

const cacheLatestConfigs = async () => {
  log('Running fetchOBAccountPaymentServiceProviders');
  await fetchOBAccountPaymentServiceProviders();

  log('Running updateOpenIdConfigs');
  await updateOpenIdConfigs();
};

cacheLatestConfigs().then(() => {
  if (process.env.NODE_ENV !== 'test') {
    process.exit();
  }
});

exports.cacheLatestConfigs = cacheLatestConfigs;
