const { obtainClientCredentialsAccessToken } = require('../authorise');
const { postAccountRequests } = require('./account-requests');

const createRequest = async (resourcePath, headers) => {
  const response = await postAccountRequests(resourcePath, headers);

  if (response.Data) {
    const status = response.Data.Status;
    if (status === 'AwaitingAuthorisation' || status === 'Authorised') {
      if (response.Data.AccountRequestId && response.Data.Permissions) {
        return response;
        // return {
        //   accountRequestId: response.Data.AccountRequestId,
        //   permissions: response.Data.Permissions,
        // };
      }
    } else {
      const error = new Error(`Account request response status: "${status}"`);
      error.status = 500;
      throw error;
    }
  }

  const error = new Error('Account request response missing payload');
  error.status = 500;
  throw error;
};

exports.setupAccountRequest = async (headers) => {
  const { config } = headers;
  const accessToken = await obtainClientCredentialsAccessToken(config);
  const headersWithToken = Object.assign({ accessToken }, headers);
  const accountRequestIdAndPermissions = await createRequest(
    config.resource_endpoint,
    headersWithToken,
  );
  return accountRequestIdAndPermissions;
};
