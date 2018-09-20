const { requireAuthorization } = require('./authorization');
const { session, getUsername } = require('./session');
const { login } = require('./login');
const { extractHeaders } = require('./request-headers');

module.exports = {
  login,
  requireAuthorization,
  session,
  getUsername,
  extractHeaders,
};
