const env = require('env-var');
const error = require('debug')('error');
const debug = require('debug')('debug');
const { getAll, get, set } = require('../storage');
const { getOpenIdConfig } = require('./openid-config');

const NO_SOFTWARE_STATEMENT_ID = 'local';
const ASPSP_AUTH_SERVERS_COLLECTION = 'aspspAuthorisationServers';

const isProvisionedForOpenBanking = () => env.get('OB_PROVISIONED').asString() === 'true';

const getSoftwareStatementId = () => (isProvisionedForOpenBanking() ?
  env.get('SOFTWARE_STATEMENT_ID').asString() : NO_SOFTWARE_STATEMENT_ID);

const sortByName = (list) => {
  list.sort((a, b) => {
    if (a.name > b.name) {
      return 1;
    } else if (a.name > b.name) {
      return -1;
    }
    return 0;
  });
  return list;
};

const decorateServerDataForClient = (data) => {
  const id = data.Id;
  const logoUri = data.CustomerFriendlyLogoUri;
  const name = data.CustomerFriendlyName;
  return {
    id,
    logoUri,
    name,
  };
};

const getAuthServerConfig = async id => get(ASPSP_AUTH_SERVERS_COLLECTION, id);

const setAuthServerConfig = async (id, authServer) =>
  set(ASPSP_AUTH_SERVERS_COLLECTION, authServer, id);

const requireAuthorisationServerId = async (req, res, next) => {
  const authServerId = req.headers['x-authorization-server-id'];
  if (!authServerId) {
    return res.status(400).send('request missing x-authorization-server-id header');
  }
  return next();
};

const storeAuthorisationServers = async (list) => {
  await Promise.all(list.map(async (item) => {
    const id = item.Id;
    const existing = await getAuthServerConfig(id);
    const authServer = existing || {};
    authServer.obDirectoryConfig = item;
    await setAuthServerConfig(id, authServer);
  }));
};

const allAuthorisationServers = async () => {
  try {
    const list = await getAll(ASPSP_AUTH_SERVERS_COLLECTION);
    if (!list) {
      return [];
    }
    return list;
  } catch (e) {
    error(e);
    return [];
  }
};

const updateOpenIdConfig = async (id, openidConfig) => {
  const authServer = await getAuthServerConfig(id);
  debug(`openidConfig: ${JSON.stringify(openidConfig)}`);
  authServer.openIdConfig = openidConfig;
  await setAuthServerConfig(id, authServer);
};

const fetchAndStoreOpenIdConfig = async (id, openidConfigUrl) => {
  try {
    if (openidConfigUrl === 'https://redirect.openbanking.org.uk') {
      return null; // ignore
    }
    const openidConfig = await getOpenIdConfig(openidConfigUrl);
    if (openidConfig) {
      await updateOpenIdConfig(id, openidConfig);
    } else {
      error(`OpenID config at ${openidConfigUrl} is blank`);
    }
  } catch (err) {
    error(`Error getting ${openidConfigUrl} : ${err.message}`);
  }
  return null;
};

const updateClientCredentials = async (id, newCredentials) => {
  const authServer = await getAuthServerConfig(id);
  const softwareStatementId = getSoftwareStatementId();
  if (!authServer) {
    throw new Error(`id=${id}, Auth Server Not Found`);
  }

  authServer.clientCredentials = Array.isArray(authServer.clientCredentials) ?
    authServer.clientCredentials : [];
  const found = authServer.clientCredentials.find(cred =>
    cred.softwareStatementId === softwareStatementId);
  const updated = Object.assign(found || { softwareStatementId }, newCredentials);
  if (!found) authServer.clientCredentials.push(updated);

  await setAuthServerConfig(id, authServer);
  return true;
};

const getRegisteredConfig = async (authServerId) => {
  const authServer = await getAuthServerConfig(authServerId);
  const softwareStatementId = getSoftwareStatementId();
  if (authServer
    && authServer.registeredConfigs
    && authServer.registeredConfigs.length > 0) {
    return authServer.registeredConfigs.find(config =>
      config.softwareStatementId === softwareStatementId);
  }

  const err = new Error(`Registered config not found for authServerId: [${authServerId}], softwareStatementId: [${softwareStatementId}]`);
  err.status = 500;
  throw err;
};

const updateRegisteredConfig = async (id, newConfig) => {
  const authServer = await getAuthServerConfig(id);
  const softwareStatementId = getSoftwareStatementId();
  if (!authServer) {
    throw new Error('Auth Server Not Found !');
  }
  authServer.registeredConfigs = authServer.registeredConfigs || [];
  const found = authServer.registeredConfigs.find(config =>
    config.softwareStatementId === softwareStatementId);
  const updated = Object.assign(found || { softwareStatementId }, newConfig);
  if (!found) authServer.registeredConfigs.push(updated);

  await setAuthServerConfig(id, authServer);
  return true;
};

const updateOpenIdConfigs = async () => {
  try {
    const list = await allAuthorisationServers();

    await Promise.all(list.map(async (authServer) => {
      try {
        const openidConfigUrl = authServer.obDirectoryConfig.OpenIDConfigEndPointUri;
        await fetchAndStoreOpenIdConfig(authServer.id, openidConfigUrl);
      } catch (err) {
        error(err);
      }
    }));
  } catch (err) {
    error(err);
  }
};

const authorisationServersForClient = async () => {
  try {
    const allServers = await allAuthorisationServers();
    const registeredServers = allServers.filter(s => s.clientCredentials);
    const servers = registeredServers.map(s => decorateServerDataForClient(s.obDirectoryConfig));
    return sortByName(servers);
  } catch (e) {
    error(e);
    return [];
  }
};

const openIdConfig = async (id) => {
  try {
    const config = await getAuthServerConfig(id);
    return (config && config.openIdConfig) ? config.openIdConfig : null;
  } catch (err) {
    error(err);
    return null;
  }
};

const openIdConfigValue = async (id, key) => {
  const config = await openIdConfig(id);
  const value = config ? config[key] : null;
  if (value === null) {
    const err = new Error(`${key} for auth server ${id} not found`);
    err.status = 500;
    throw err;
  }
  return value;
};

const requestObjectSigningAlgs = async (id) => {
  let registeredConfig;
  try {
    registeredConfig = await getRegisteredConfig(id);
  } catch (e) {
    registeredConfig = {};
  }

  if (registeredConfig.request_object_signing_alg) {
    return registeredConfig.request_object_signing_alg;
  }
  return openIdConfigValue(id, 'request_object_signing_alg_values_supported');
};

module.exports = {
  storeAuthorisationServers,
  allAuthorisationServers,
  authorisationServersForClient,
  updateOpenIdConfigs,
  updateOpenIdConfig,
  updateClientCredentials,
  setAuthServerConfig,
  requireAuthorisationServerId,
  requestObjectSigningAlgs,
  updateRegisteredConfig,
  getRegisteredConfig,
  ASPSP_AUTH_SERVERS_COLLECTION,
  NO_SOFTWARE_STATEMENT_ID,
};
