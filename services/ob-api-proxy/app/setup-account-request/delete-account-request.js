const { obtainClientCredentialsAccessToken, consentAccountRequestId, deleteConsent } = require('../authorise');
const { deleteAccountRequest } = require('./account-requests');

exports.deleteRequest = async (headers) => {
  const {
    authorisationServerId, username, validationRunId, config,
  } = headers;
  const keys = {
    username, authorisationServerId, scope: 'accounts', validationRunId,
  };
  const accountRequestId = await consentAccountRequestId(keys);

  if (accountRequestId) {
    const accessToken = await obtainClientCredentialsAccessToken(config);
    const headersWithToken = Object.assign({ accessToken }, headers);
    const success = await deleteAccountRequest(
      accountRequestId,
      config.resource_endpoint, headersWithToken,
    );
    if (success) {
      await deleteConsent(keys);
      return 204;
    }
  }
  const error = new Error('Bad request - account request ID not found');
  error.status = 400;
  throw error;
};
