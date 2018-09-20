exports.default = {
  post: {
    request() {
      return {
        method: 'POST',
        url: 'http://localhost:8001/open-banking/v1.1/account-requests',
        path: '/open-banking/v1.1/account-requests',
        headers: {
          'Authorization': 'Bearer 2YotnFZFEjr1zCsicMWpAA',
          'Accept': 'application/json',
          'x-fapi-financial-id': 'aaax5nTR33811QyQfi',
          'x-fapi-interaction-id': '0f2253b5-30bb-40a2-93f6-0708b4e76325',
          'content-type': 'application/json; charset=utf-8',
        },
        body: {
          Data: {
            Permissions: [
              'ReadAccountsDetail',
              'ReadBalances',
              'ReadBeneficiariesDetail',
              'ReadDirectDebits',
              'ReadProducts',
              'ReadStandingOrdersDetail',
              'ReadTransactionsCredits',
              'ReadTransactionsDebits',
              'ReadTransactionsDetail',
            ],
            ExpirationDateTime: '2017-05-02T00:00:00+00:00',
            // ExpirationDateTime: '2017-05-02T00:00:00+00:00',
            TransactionFromDateTime: '2017-05-03T00:00:00+00:00',
            TransactionToDateTime: '2017-12-03T00:00:00+00:00',
          },
          Risk: {},
        },
      };
    },
    response() {
      return {
        statusCode: 200,
        url: 'http://localhost:8001/open-banking/v1.1/account-requests',
        path: '/open-banking/v1.1/account-requests',
        headers: {
          'Authorization': 'Bearer 2YotnFZFEjr1zCsicMWpAA',
          'Accept': 'application/json',
          'x-fapi-financial-id': 'aaax5nTR33811QyQfi',
          'x-fapi-interaction-id': '0f2253b5-30bb-40a2-93f6-0708b4e76325',
        },
        body: JSON.stringify({
          Data: {
            AccountRequestId: '88379',
            Status: 'AwaitingAuthorisation',
            CreationDateTime: '2017-05-02T00:00:00+00:00',
            Permissions: [
              'ReadAccountsDetail',
              'ReadBalances',
              'ReadBeneficiariesDetail',
              'ReadDirectDebits',
              'ReadProducts',
              'ReadStandingOrdersDetail',
              'ReadTransactionsCredits',
              'ReadTransactionsDebits',
              'ReadTransactionsDetail',
            ],
            ExpirationDateTime: '2017-05-02T00:00:00+00:00',
            // ExpirationDateTime: '2017-08-02T00:00:00+00:00',
            TransactionFromDateTime: '2017-05-03T00:00:00+00:00',
            TransactionToDateTime: '2017-12-03T00:00:00+00:00',
          },
          Risk: {},
          Links: {
            Self: '/account-requests/88379',
          },
          Meta: {
            TotalPages: 1,
          },
        }),
      };
    },
  },
};
