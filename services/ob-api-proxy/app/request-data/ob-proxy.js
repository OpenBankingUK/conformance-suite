const superagent = require('superagent');
const { createRequest, obtainResult } = require('../ob-util');
const { consentAccessTokenAndPermissions } = require('../authorise');
const { extractHeaders } = require('../session');
const debug = require('debug')('debug');
const error = require('debug')('error');

const accessTokenAndPermissions = async (username, authorisationServerId,
  validationRunId, scope) => {
  let accessToken;
  let permissions;
  try {
    const consentKeys = {
      username, authorisationServerId, validationRunId, scope,
    };
    ({ accessToken, permissions } = await consentAccessTokenAndPermissions(consentKeys));
  } catch (err) {
    accessToken = null;
    permissions = null;
  }
  return { accessToken, permissions };
};

const scopeAndUrl = (reqPath, host) => {
  const path = `/open-banking${reqPath}`;
  const proxiedUrl = `${host}${path}`;
  const scope = path.split('/')[3].startsWith('payment') ? 'payments' : 'accounts';
  return { proxiedUrl, scope };
};

const resourceRequestHandler = async (req, res) => {
  try {
    res.setHeader('Access-Control-Allow-Origin', '*');
    const reqHeaders = await extractHeaders(req.headers);
    const { config } = reqHeaders;
    const { proxiedUrl, scope } = scopeAndUrl(req.path, config.resource_endpoint);
    const { accessToken, permissions } =
      await accessTokenAndPermissions(
        reqHeaders.username,
        reqHeaders.authorisationServerId,
        reqHeaders.validationRunId,
        scope,
      );
    const headers = Object.assign({ accessToken, permissions, scope }, reqHeaders);
    debug({
      proxiedUrl,
      scope,
      accessToken,
      fapiFinancialId: headers.fapiFinancialId,
      validationRunId: headers.validationRunId,
    });
    const request = createRequest(proxiedUrl, superagent.get(proxiedUrl), headers);

    let response;
    try {
      response = await request.send();
      debug('services/ob-api-proxy/app/request-data/ob-proxy.js:resourceRequestHandler -> response=%j', response);
    } catch (err) {
      error('services/ob-api-proxy/app/request-data/ob-proxy.js:resourceRequestHandler -> error getting proxiedUrl=%j, err=%j', proxiedUrl, err.message);
      throw err;
    }

    const result = await obtainResult(request, response, headers);
    debug('services/ob-api-proxy/app/request-data/ob-proxy.js:resourceRequestHandler -> result=%j', result);

    return res
      .status(response.status)
      .json(result);
  } catch (err) {
    const status = err.response ? err.response.status : 500;
    return res.status(status).send(err.message);
  }
};

module.exports = {
  resourceRequestHandler,
  scopeAndUrl,
};
