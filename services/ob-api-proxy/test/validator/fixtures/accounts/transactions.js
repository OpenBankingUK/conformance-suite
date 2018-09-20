exports.default = {
  valid() {
    return {
      request: {
        method: 'GET',
        url: 'http://localhost:8001/open-banking/v1.1/accounts/22292/transactions',
        path: '/open-banking/v1.1/accounts/22292/transactions',
        headers: {
          'Authorization': 'Bearer 2YotnFZFEjr1zCsicMWpAA',
          'Accept': 'application/json',
          'x-fapi-financial-id': 'aaax5nTR33811QyQfi',
          'x-fapi-interaction-id': '0f2253b5-30bb-40a2-93f6-0708b4e76325',
        },
      },

      basic: {
        statusCode: 200,
        headers: {
          'access-control-allow-origin': '*',
          'access-control-allow-methods': 'GET',
          'access-control-allow-headers': '',
          'access-control-allow-credentials': 'false',
          'access-control-max-age': '0',
          'content-type': 'application/json; charset=utf-8',
          'content-length': '621',
          'etag': 'W/"26d-/CEtMNK6kuJdSw//7SDW6kTgV90"',
          'date': 'Wed, 07 Feb 2018 11:58:01 GMT',
          'connection': 'close',
        },
        body: JSON.stringify({
          Data: {
            Transaction: [
              {
                AccountId: 'string',
                TransactionId: 'string',
                TransactionReference: 'string',
                Amount: {
                  Amount: '15000.00',
                  Currency: 'GBP',
                },
                CreditDebitIndicator: 'Credit',
                Status: 'Booked',
                BookingDateTime: '2018-07-03T14:40:18.155Z',
                ValueDateTime: '2018-07-03T14:40:18.155Z',
                AddressLine: 'string',
                BankTransactionCode: {
                  Code: 'string',
                  SubCode: 'string',
                },
                ProprietaryBankTransactionCode: {
                  Code: 'string',
                  Issuer: 'string',
                },
              },
            ],
          },
          Links: {
            Self: 'string',
            First: 'string',
            Prev: 'string',
            Next: 'string',
            Last: 'string',
          },
          Meta: {
            TotalPages: 0,
            FirstAvailableDateTime: '2018-07-03T14:40:18.155Z',
            LastAvailableDateTime: '2018-07-03T14:40:18.155Z',
          },
        }),
      },

      detail: {
        statusCode: 200,
        headers: {
          'access-control-allow-origin': '*',
          'access-control-allow-methods': 'GET',
          'access-control-allow-headers': '',
          'access-control-allow-credentials': 'false',
          'access-control-max-age': '0',
          'content-type': 'application/json; charset=utf-8',
          'content-length': '621',
          'etag': 'W/"26d-/CEtMNK6kuJdSw//7SDW6kTgV90"',
          'date': 'Wed, 07 Feb 2018 11:58:01 GMT',
          'connection': 'close',
        },
        body: JSON.stringify({
          Data: {
            Transaction: [
              {
                AccountId: 'string',
                TransactionId: 'string',
                TransactionReference: 'string',
                Amount: {
                  Amount: '15000.00',
                  Currency: 'GBP',
                },
                CreditDebitIndicator: 'Credit',
                Status: 'Booked',
                BookingDateTime: '2018-07-03T15:21:03.223Z',
                ValueDateTime: '2018-07-03T15:21:03.223Z',
                AddressLine: 'string',
                BankTransactionCode: {
                  Code: 'string',
                  SubCode: 'string',
                },
                ProprietaryBankTransactionCode: {
                  Code: 'string',
                  Issuer: 'string',
                },
                TransactionInformation: 'string',
                Balance: {
                  Amount: {
                    Amount: '15000.00',
                    Currency: 'GBP',
                  },
                  CreditDebitIndicator: 'Credit',
                  Type: 'ClosingAvailable',
                },
                MerchantDetails: {
                  MerchantName: 'string',
                  MerchantCategoryCode: '3000',
                },
              },
            ],
          },
          Links: {
            Self: 'string',
            First: 'string',
            Prev: 'string',
            Next: 'string',
            Last: 'string',
          },
          Meta: {
            TotalPages: 0,
            FirstAvailableDateTime: '2018-07-03T15:21:03.223Z',
            LastAvailableDateTime: '2018-07-03T15:21:03.223Z',
          },
        }),
      },
    };
  },
};
