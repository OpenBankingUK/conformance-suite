/* eslint camelcase: 0 */
const { verifyHeaders, postPayments } = require('../setup-payment/payments');
const { retrievePaymentDetails } = require('../setup-payment/persistence');

const makePayment = async ({ resource_endpoint, api_version }, headers, paymentData) => {
  const response = await postPayments(
    resource_endpoint,
    `/open-banking/v${api_version}/payment-submissions`,
    headers,
    paymentData,
  );

  if (response && response.Data && response.Data.Status !== 'Rejected') {
    return response.Data.PaymentSubmissionId;
  }
  const error = new Error('Payment submission failed');
  error.status = 500;
  throw error;
};

exports.submitPayment = async (authorisationServerId, headers) => {
  const { config } = headers;
  verifyHeaders(headers);
  const paymentData = await retrievePaymentDetails(headers.interactionId);
  return makePayment(config, headers, paymentData);
};
