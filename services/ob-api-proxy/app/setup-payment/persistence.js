const { set, get } = require('../storage');
const assert = require('assert');

const persistPaymentDetails = async (interactionId, paymentData) => {
  assert(interactionId);
  assert(paymentData);
  assert(paymentData.Data.PaymentId);

  await set('payments', paymentData, interactionId);
};

const retrievePaymentDetails = async (interactionId) => {
  assert(interactionId);
  return get('payments', interactionId, ['-id']);
};

module.exports = {
  persistPaymentDetails,
  retrievePaymentDetails,
};
