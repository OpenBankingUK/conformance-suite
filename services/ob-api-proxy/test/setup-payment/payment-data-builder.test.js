const { buildPaymentsData } = require('../../app/setup-payment/payment-data-builder');
const assert = require('assert');

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

describe('buildPaymentstData and then postPayments', () => {
  const instructionIdentification = 'ghghg';
  const endToEndIdentification = 'XXXgHTg';
  const reference = 'Things';
  const unstructured = 'XXX';

  const opts = {
    instructionIdentification,
    endToEndIdentification,
    reference,
    unstructured,
  };

  const risk = {
    foo: 'bar',
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

  describe(' For the /payments endpoint', () => {
    it('returns a body payload of the correct shape', async () => {
      const paymentsPayload = buildPaymentsData(opts, risk, creditorAccount, instructedAmount);
      const expectedPayload = {
        Data: paymentData,
        Risk: risk,
      };
      assert.deepEqual(paymentsPayload, expectedPayload);
    });
  });
});

describe('buildPaymentstData with optionality', () => {
  const instructionIdentification = 'ttttt';
  const endToEndIdentification = 'RRR';
  const reference = 'Ref2';

  const opts = {
    instructionIdentification,
    endToEndIdentification,
    reference,
  };

  const risk = {
    foo: 'bar',
  };

  const data = {
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
      },
    },
  };

  it('returns a body payload of the correct shape: with missing unstructured field', () => {
    const paymentsPayload = buildPaymentsData(opts, risk, creditorAccount, instructedAmount);
    const expectedPayload = {
      Data: {
        Initiation: data.Initiation,
      },
      Risk: risk,
    };
    assert.deepEqual(paymentsPayload, expectedPayload);
  });

  it('returns a body payload of the correct shape: with missing reference field', () => {
    opts.unstructured = 'blah';
    delete opts.reference;
    data.Initiation.RemittanceInformation = {
      Unstructured: opts.unstructured,
    };
    const paymentsPayload = buildPaymentsData(opts, risk, creditorAccount, instructedAmount);
    const expectedPayload = {
      Data: data,
      Risk: risk,
    };
    assert.deepEqual(paymentsPayload, expectedPayload);
  });

  it('returns a body payload of the correct shape: with missing reference AND unstructured fields', () => {
    delete opts.reference;
    delete opts.unstructured;
    delete data.Initiation.RemittanceInformation;
    const paymentsPayload = buildPaymentsData(opts, risk, creditorAccount, instructedAmount);
    const expectedPayload = {
      Data: data,
      Risk: risk,
    };
    assert.deepEqual(paymentsPayload, expectedPayload);
  });
});
