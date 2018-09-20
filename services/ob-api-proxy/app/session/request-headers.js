const uuidv4 = require('uuid/v4');
const { getUsername } = require('./session');
const { base64DecodeJSON } = require('../ob-util');

const parseMultiValueHeader = (headers, header) => {
  const SPLIT_REGEX = ' ';

  if (headers[header] && headers[header].length > 0) {
    return headers[header].split(SPLIT_REGEX);
  }

  return [];
};

exports.extractHeaders = async (headers) => {
  const sessionId = headers['authorization'];
  const authorisationServerId = headers['x-authorization-server-id'];
  const interactionId = headers['x-fapi-interaction-id'] || uuidv4();
  const username = await getUsername(sessionId);
  const validationRunId = headers['x-validation-run-id'];
  const swaggerUris = parseMultiValueHeader(headers, 'x-swagger-uris');
  const permissions = headers['x-permissions'] && headers['x-permissions'].split(' ');
  const config = headers['x-config'] && base64DecodeJSON(headers['x-config']);
  const fapiFinancialId = config.fapi_financial_id;

  const extracted = {
    authorisationServerId,
    fapiFinancialId,
    interactionId,
    sessionId,
    username,
    validationRunId,
    swaggerUris,
    permissions,
    config,
  };

  return extracted;
};
