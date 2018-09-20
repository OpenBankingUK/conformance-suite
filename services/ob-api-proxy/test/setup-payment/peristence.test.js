const sinon = require('sinon');
const proxyquire = require('proxyquire');
const assert = require('assert');

const InstructedAmount = {
  Amount: '100.45',
  Currency: 'GBP',
};

const CreditorAccount = {
  SchemeName: 'SortCodeAccountNumber',
  Identification: '01122313235478',
  Name: 'Mr Kevin',
  SecondaryIdentification: '002',
};

const paymentData = {
  CreditorAccount,
  InstructedAmount,
  InstructionIdentification: 'testInstructionId',
  EndToEndIdentification: 'testEndToEndIdentification',
};

const PaymentId = '44673';
const interactionId = 'ABCD';

describe('persist payment details and retrieve it', () => {
  let setSpy;
  let getSpy;
  let persistence;

  beforeEach(() => {
    setSpy = sinon.spy();
    getSpy = sinon.spy();
    persistence = proxyquire('../../app/setup-payment/persistence', {
      '../storage': { set: setSpy, get: getSpy },
    });
  });
  it('verify valid payment details persistence', async () => {
    const { persistPaymentDetails } = persistence;
    const fullPaymentData = {
      Data: Object.assign({}, paymentData, { PaymentId }),
    };
    await persistPaymentDetails(interactionId, fullPaymentData);
    assert.ok(setSpy.called);
    assert.ok(setSpy.calledOnce);
    assert.ok(setSpy.calledWithExactly('payments', fullPaymentData, interactionId));
  });

  it('verify valid payment details retrieval', async () => {
    const { retrievePaymentDetails } = persistence;

    await retrievePaymentDetails(interactionId);
    assert.ok(getSpy.called);
    assert.ok(getSpy.calledOnce);
    assert.ok(getSpy.calledWithExactly('payments', interactionId, ['-id']));
  });

  it('verify error when PaymentId is not provided for payment details persistence', async () => {
    const { persistPaymentDetails } = persistence;
    const fullPaymentData = {
      Data: paymentData,
    };
    try {
      await persistPaymentDetails(interactionId, fullPaymentData);
    } catch (e) {
      assert.ok(e instanceof assert.AssertionError);
    }
    assert.ok(getSpy.notCalled);
  });

  it('verify error when interactionId is not provided for payment details persistence', async () => {
    const { persistPaymentDetails } = persistence;

    try {
      await persistPaymentDetails(null, paymentData);
    } catch (e) {
      assert.ok(e instanceof assert.AssertionError);
    }
    assert.ok(getSpy.notCalled);
  });

  it('verify error when interactionId not provided for payment details retrieval', async () => {
    const { retrievePaymentDetails } = persistence;

    try {
      await retrievePaymentDetails(null);
    } catch (e) {
      assert.ok(e instanceof assert.AssertionError);
    }
    assert.ok(getSpy.notCalled);
  });
});
