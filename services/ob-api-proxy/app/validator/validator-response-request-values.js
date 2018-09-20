const { diff } = require('deep-diff');
const _ = require('lodash');

/**
 * Map of the request path and request method to the list of
 * properties that should be equal in the request and response for
 * request path and the request method combination. E.g., in the request
 * for /open-banking/v1.1/payments POST and response we expect that
 * ['Data.Initiation', 'Risk'] in the response should match the
 * ['Data.Initiation', 'Risk'] in the request.
 */
const ENDPOINT_TO_PROPERTIES_PATH = new Map([
  // See https://openbanking.atlassian.net/wiki/spaces/DZ/pages/5786479/Payment+Initiation+API+Specification+-+v1.1.0#PaymentInitiationAPISpecification-v1.1.0-Endpoints
  // for list sof payments endpoints for v1.1.
  //
  // https://openbanking.atlassian.net/wiki/spaces/DZ/pages/5786479/Payment+Initiation+API+Specification+-+v1.1.0#PaymentInitiationAPISpecification-v1.1.0-POST/payments
  //
  // Validate that `OBPaymentSetup1` and `OBPaymentSetupResponse1` in
  // POST /open-banking/v1.1/payments match:
  // * OBPaymentSetup1/Data (OBPaymentDataSetup1) ==
  //    OBPaymentSetupResponse1/Data (OBPaymentDataSetupResponse1)
  // * OBPaymentSetup1/Risk (OBRisk1) ==
  //    OBPaymentSetupResponse1/Risk (OBRisk1)
  [
    '/open-banking/v1.1/payments_POST',
    [
      ['Data', 'Initiation'],
      ['Risk'],
    ],
  ],
  // https://openbanking.atlassian.net/wiki/spaces/DZ/pages/5785171/Account+and+Transaction+API+Specification+-+v1.1.0#AccountandTransactionAPISpecification-v1.1.0-Endpoints
  //
  // Validate that `OBReadRequest1` and `OBReadResponse1` in POST /account-requests match:
  // * OBReadRequest1/Data/Permissions (OBExternalPermissions1Code) ==
  //    OBReadResponse1/Data/Permissions (OBExternalPermissions1Code)
  // * OBReadRequest1/Data/ExpirationDateTime (ISODateTime) ==
  //    OBReadResponse1/Data/ExpirationDateTime (ISODateTime)
  // * OBReadRequest1/Data/TransactionFromDateTime (ISODateTime) ==
  //    OBReadResponse1/Data/TransactionFromDateTime (ISODateTime)
  // * OBReadRequest1/Data/TransactionToDateTime (ISODateTime) ==
  //    OBReadResponse1/Data/TransactionToDateTime (ISODateTime)
  // * OBReadRequest1/Risk (OBRisk2) ==
  //    OBReadResponse1/Risk (OBRisk2)
  [
    '/open-banking/v1.1/account-requests_POST',
    [
      ['Data', 'Permissions'],
      ['Data', 'ExpirationDateTime'],
      ['Data', 'TransactionFromDateTime'],
      ['Data', 'TransactionToDateTime'],
      ['Risk'],
    ],
  ],
  // https://openbanking.atlassian.net/wiki/spaces/DZ/pages/127009546/Account+and+Transaction+API+Specification+-+v2.0.0#AccountandTransactionAPISpecification-v2.0.0-Endpoints
  // https://openbanking.atlassian.net/wiki/spaces/DZ/pages/129040562/Account+Requests+v2.0.0
  // https://openbanking.atlassian.net/wiki/spaces/DZ/pages/129040562/Account+Requests+v2.0.0#AccountRequestsv2.0.0-AccountRequests-Request
  // https://openbanking.atlassian.net/wiki/spaces/DZ/pages/129040562/Account+Requests+v2.0.0#AccountRequestsv2.0.0-AccountRequests-Response
  //
  // Validate that `OBReadRequest1` and `OBReadResponse1` in POST /account-requests match:
  // * OBReadRequest1/Data/Permissions (OBExternalPermissions1Code) ==
  //    OBReadResponse1/Data/Permissions (OBExternalPermissions1Code)
  // * OBReadRequest1/Data/ExpirationDateTime (ISODateTime) ==
  //    OBReadResponse1/Data/ExpirationDateTime (ISODateTime)
  // * OBReadRequest1/Data/TransactionFromDateTime (ISODateTime) ==
  //    OBReadResponse1/Data/TransactionFromDateTime (ISODateTime)
  // * OBReadRequest1/Data/TransactionToDateTime (ISODateTime) ==
  //    OBReadResponse1/Data/TransactionToDateTime (ISODateTime)
  // * OBReadRequest1/Risk (OBRisk2) ==
  //    OBReadResponse1/Risk (OBRisk2)
  [
    '/open-banking/v2.0/account-requests_POST',
    [
      ['Data', 'Permissions'],
      ['Data', 'ExpirationDateTime'],
      ['Data', 'TransactionFromDateTime'],
      ['Data', 'TransactionToDateTime'],
      ['Risk'],
    ],
  ],
]);

/**
 * Make key to index into the `ENDPOINT_TO_PROPERTIES_PATH` Map.
 *
 * @param {string} requestPath The request path, e.g., `/open-banking/v2.0/account-requests`.
 * @param {string} requestMethod The request method, e.g., `POST`.
 * @returns {string} Key to be used in `ENDPOINT_TO_PROPERTIES_PATH` map.
 */
const makeEndpointToPropertiesPathKey = (requestPath, requestMethod) => {
  const key = `${requestPath}_${requestMethod}`;
  return key;
};

/**
 * Validate that certain parts of the request's body matches certain
 * parts of the response's body.
 *
 * @param {object} requestBody
 *  The request that was sent to /open-banking/*.
 * @param {object} responseBody
 *  The response from a call to /open-banking/*.
 * @param {Array<Array<string>>} propertiesPath
 *  The paths to validate as a 2d array, e.g., propertiesPath=[['Data', 'Initiation'], ['Risk']]
 * @returns {object}
 *  Object indicating if it passed validation and list of errors, if any.
 */
const validatePropertiesMatch = (requestBody, responseBody, propertiesPath) => {
  const PROPERTY_PATH_SEPARATOR = '.';

  const errors = [];

  _.forEach(propertiesPath, (propertyPathArray) => {
    // e.g., propertyPathArray=['Data', 'Initiation'], propertyPath='Data.Initiation'
    const propertyPath = propertyPathArray.join(PROPERTY_PATH_SEPARATOR);

    // diff request against response
    //
    // Will contain something like:
    // [
    //   DiffNew {
    //     kind: 'N',
    //     path: ['EndToEndIdentification'],
    //     rhs: '8a30c4fe-a779-436f-b231-f21c05bd22'
    //   },
    //   DiffNew {
    //     kind: 'N',
    //     path: ['CreditorAccount'],
    //     rhs:
    //     {
    //       SchemeName: 'SortCodeAccountNumber',
    //       Name: 'Sam Morse',
    //       Identification: '11111112345678'
    //     }
    //   }
    // ]
    const differences = diff(
      _.get(requestBody, propertyPath, {}), /* lhs */
      _.get(responseBody, propertyPath, {}), /* rhs */
    );

    // push an object describing each difference into errors
    _.forEach(differences, (diffNew) => {
      const isArrayChange = diffNew.kind === 'A';

      // prefix to get the values of the left and right hand side
      let valuePathPrefix = '';
      // path to property modification as an array
      let diffNewPath;

      if (isArrayChange) {
        // deal with array modification differently. e.g.,
        // diffNew=DiffArray {
        //   kind: 'A',
        //   index: 8,
        //   item: DiffNew { kind: 'N', rhs: 'ReadTransactionsDetail' }
        // }
        valuePathPrefix = 'item.';
        diffNewPath = [diffNew.index];
      } else {
        diffNewPath = diffNew.path || [];
      }

      const path = _.concat(propertyPathArray, diffNewPath);
      const absolutePropertyPath = path.join(PROPERTY_PATH_SEPARATOR);
      const requestValue = _.get(diffNew, `${valuePathPrefix}lhs`, '');
      const responseValue = _.get(diffNew, `${valuePathPrefix}rhs`, '');
      const message = `request.${absolutePropertyPath}=${JSON.stringify(requestValue)} != response.${absolutePropertyPath}=${JSON.stringify(responseValue)}`;
      const code = 'OBJECTS_NOT_EQUAL';

      errors.push({
        path,
        message,
        code,
      });
    });
  });

  const isValid = _.isEmpty(errors);
  return {
    isValid,
    errors,
  };
};

module.exports = {
  /**
   * Set errors as the joined errors of `logObject` and `validateResponseResult`.
   * Concat errors, then set it as `logObject.report.results.errors`.
   *
   * @param {object} logObject Object returned by `logFormat`.
   * @param {object} validateResponseResult Object return by `validateResponse`.
   * @returns {object} The `logObject` parameter.
   */
  mergeErrors(logObject, validateResponseResult) {
    const errors = _.concat(
      _.get(logObject, 'report.results.errors', []),
      _.get(validateResponseResult, 'errors', []),
    );

    // set errors as the joined errors
    return _.set(logObject, 'report.results.errors', errors);
  },
  /**
   * Do additional validation such as checking that the response returns
   * the same data that was sent in the request for certain fields.
   *
   * @param {object} request A request to a /open-banking/* endpoint.
   * @param {object} response The response to the `request` parameter.
   * @returns {object} Object indicating if it passed validation and list of errors, if any.
   */
  validateResponse(request, response) {
    const { path: requestPath, method: requestMethod, body: requestBody } = request;
    const responseBody = JSON.parse(response.body);
    const key = makeEndpointToPropertiesPathKey(requestPath, requestMethod);

    if (ENDPOINT_TO_PROPERTIES_PATH.has(key)) {
      const propertiesPath = ENDPOINT_TO_PROPERTIES_PATH.get(key);
      return validatePropertiesMatch(
        requestBody,
        responseBody,
        propertiesPath,
      );
    }

    return {
      isValid: true,
      errors: [],
    };
  },
};
