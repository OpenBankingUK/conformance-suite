const uuidv4 = require('uuid/v4');

const allowedCurrencies = ['GBP', 'EUR']; // TODO - refactor out of here

// For detailed spec see
// https://openbanking.atlassian.net/wiki/spaces/WOR/pages/23266217/Payment+Initiation+API+Specification+-+v1.1.1#PaymentInitiationAPISpecification-v1.1.1-POST/paymentsrequest

exports.buildPaymentsData = (opts, risk, creditorAccount, instructedAmount) => {
  if (!instructedAmount.Amount) throw new Error('InstructedAmount Amount missing');
  if (!instructedAmount.Currency) throw new Error('InstructedAmount Currency missing');
  if (!creditorAccount.SchemeName) throw new Error('CreditorAccount SchemeName missing');
  if (!creditorAccount.Identification) throw new Error('CreditorAccount Identification missing');
  if (!creditorAccount.Name) throw new Error('CreditorAccount Name missing');
  const {
    instructionIdentification,
    endToEndIdentification,
    reference,
    unstructured,
  } = opts;
  const currency = instructedAmount.Currency;

  if (allowedCurrencies.indexOf(currency) === -1) throw new Error('Disallowed currency');
  const payload = {
    Data: {
      Initiation: {
        InstructionIdentification: instructionIdentification || uuidv4().slice(0, 34),
        EndToEndIdentification: endToEndIdentification || uuidv4().slice(0, 34),
        InstructedAmount: instructedAmount,
        CreditorAccount: creditorAccount,
      },
    },
    Risk: risk || {},
  };

  // Optional Fields
  let remittanceInformation;
  if (reference || unstructured) {
    remittanceInformation = {};
    if (reference) remittanceInformation.Reference = reference;
    if (unstructured) remittanceInformation.Unstructured = unstructured;
  }
  if (remittanceInformation) payload.Data.Initiation.RemittanceInformation = remittanceInformation;

  return payload;
};
