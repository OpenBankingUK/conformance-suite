const { postPayments } = require('../../app/setup-payment/payments');
const assert = require('assert');
const nock = require('nock');

const instructedAmount = {
  Amount: '100.45',
  Currency: 'GBP',
};

const creditorAccount = {
  SchemeName: 'SortCodeAccountNumber',
  Identification: '01122313235478',
  Name: 'Mr Kevin',
  SecondaryIdentification: '002',
};

const amount = instructedAmount.Amount;
const currency = instructedAmount.Currency;
const identification = creditorAccount.Identification;
const name = creditorAccount.Name;
const secondaryIdentification = creditorAccount.SecondaryIdentification;

const paymentId = '44673';
const paymentSubmissionId = '44673-001';

describe('postPayments request to remote payment endpoints', () => {
  const instructionIdentification = 'ghghg';
  const endToEndIdentification = 'XXXgHTg';
  const reference = 'Things';
  const unstructured = 'XXX';

  const risk = {
    foo: 'bar',
  };

  const accessToken = '2YotnFZFEjr1zCsicMWpAA';
  const fapiFinancialId = 'abc';
  const sessionId = 'testSessionId';
  const idempotencyKey = 'id-key-blah';
  const interactionId = 'xyz';
  const customerIp = '10.10.0.1';
  const customerLastLogged = 'Sun, 10 Sep 2017 19:43:31 UTC';
  const jwsSignature = 'testJwsSignature';

  const config = {
    transport_cert: '-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n',
    transport_key: '-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n',
  };

  const headers = {
    interactionId,
    customerIp,
    customerLastLogged,
    accessToken,
    fapiFinancialId,
    idempotencyKey,
    jwsSignature,
    sessionId,
    config,
  };

  const paymentData = {
    Initiation: {
      InstructionIdentification: instructionIdentification,
      EndToEndIdentification: endToEndIdentification,
      InstructedAmount: {
        Amount: amount,
        Currency: currency,
      },
      CreditorAccount: {
        SchemeName: 'SortCodeAccountNumber',
        Identification: identification,
        Name: name,
        SecondaryIdentification: secondaryIdentification,
      },
      RemittanceInformation: {
        Reference: reference,
        Unstructured: unstructured,
      },
    },
  };

  const paymentSubmissionData = Object.assign({}, paymentData);
  paymentSubmissionData.PaymentId = paymentId;

  const expectedPaymentResponse = {
    Data: {
      PaymentId: paymentId,
      Initiation: paymentData.Initiation,
    },
    Risk: risk,
    Links: {
      self: `/open-banking/v1.1/payments/${paymentId}`,
    },
    Meta: {
      'total-pages': 1,
    },
  };

  const expectedPaymentSubmissionResponse = {
    Data: {
      PaymentId: paymentId,
      PaymentSubmissionId: paymentSubmissionId,
    },
    Links: {
      self: `/open-banking/v1.1/payment-submissions/${paymentSubmissionId}`,
    },
    Meta: {},
  };


  // Request / response Mocks
  // Payment
  nock(/example\.com/)
    .post('/prefix/open-banking/v1.1/payments')
    .matchHeader('authorization', `Bearer ${accessToken}`) // required
    .matchHeader('x-fapi-financial-id', fapiFinancialId) // required
    .matchHeader('x-idempotency-key', idempotencyKey) // required
    .matchHeader('x-fapi-interaction-id', interactionId)
    .matchHeader('x-fapi-customer-ip-address', customerIp)
    .matchHeader('x-fapi-customer-last-logged-time', customerLastLogged)
    .matchHeader('x-jws-signature', jwsSignature) // required in v1.1.0 ( not v1.1.1 )
    .reply(201, expectedPaymentResponse);

  nock(/example\.com/)
    .post('/prefix/open-banking/v1.1/payment-submissions')
    .matchHeader('authorization', `Bearer ${accessToken}`) // required
    .matchHeader('x-fapi-financial-id', fapiFinancialId) // required
    .matchHeader('x-idempotency-key', idempotencyKey) // required
    .matchHeader('x-fapi-interaction-id', interactionId)
    .matchHeader('x-fapi-customer-ip-address', customerIp)
    .matchHeader('x-fapi-customer-last-logged-time', customerLastLogged)
    .matchHeader('x-jws-signature', jwsSignature) // required in v1.1.0 ( not v1.1.1 )
    .reply(201, expectedPaymentSubmissionResponse);

  nock(/example\.com/)
    .post('/prefix/open-banking/v1.1/non-exisits')
    .matchHeader('authorization', `Bearer ${accessToken}`) // required
    .matchHeader('x-fapi-financial-id', fapiFinancialId) // required
    .matchHeader('x-idempotency-key', idempotencyKey) // required
    .matchHeader('x-fapi-interaction-id', interactionId)
    .matchHeader('x-fapi-customer-ip-address', customerIp)
    .matchHeader('x-fapi-customer-last-logged-time', customerLastLogged)
    .matchHeader('x-jws-signature', jwsSignature) // required in v1.1.0 ( not v1.1.1 )
    .reply(404, {});

  nock(/example\.com/)
    .post('/prefix/open-banking/v1.1/bad-request')
    .matchHeader('authorization', `Bearer ${accessToken}`) // required
    .matchHeader('x-fapi-financial-id', fapiFinancialId) // required
    .matchHeader('x-idempotency-key', idempotencyKey) // required
    .matchHeader('x-fapi-interaction-id', interactionId)
    .matchHeader('x-fapi-customer-ip-address', customerIp)
    .matchHeader('x-fapi-customer-last-logged-time', customerLastLogged)
    .matchHeader('x-jws-signature', jwsSignature) // required in v1.1.0 ( not v1.1.1 )
    .reply(400, {});


  describe(' For requests to payments endpoints', () => {
    it('returns data when remote endpoint returns 201 OK', async () => {
      const resourceServerPath = 'http://example.com/prefix';
      const result = await postPayments(
        resourceServerPath,
        '/open-banking/v1.1/payments',
        headers,
        paymentData,
      );
      assert.deepEqual(result, expectedPaymentResponse);
    });
  });

  it('throws error when remote endpoints returns 404', async () => {
    const resourceServerPath = 'http://example.com/prefix';
    let error;
    try {
      await postPayments(
        resourceServerPath,
        '/open-banking/v1.1/non-exisits',
        headers,
        paymentData,
      );
    } catch (e) {
      error = e;
    }
    assert.equal('Not Found', error.message);
  });

  it('throws error when missing required request headers', async () => {
    const resourceServerPath = 'http://example.com/prefix';
    let error;
    try {
      await postPayments(
        resourceServerPath,
        '/open-banking/v1.1/bad-request',
        {},
        paymentData,
      );
    } catch (e) {
      error = e;
    }
    assert.equal('idempotencyKey missing from headers', error.message);
  });

  it('throws error when remote endpoing return 400 due to incorrect request format', async () => {
    const resourceServerPath = 'http://example.com/prefix';
    let error;
    try {
      await postPayments(
        resourceServerPath,
        '/open-banking/v1.1/bad-request',
        headers,
        paymentData,
      );
    } catch (e) {
      error = e;
    }
    assert.equal('Bad Request', error.message);
  });
});
