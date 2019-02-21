import testcases from './testcases';

describe('computeTestCases', () => {
  const EXPECTED_INPUT = '/api/test-cases';
  const EXPECTED_INIT = {
    method: 'GET',
    headers:
  {
    Accept: 'application/json; charset=UTF-8',
    'Content-Type': 'application/json; charset=UTF-8',
  },
  };

  let setShowLoading;

  beforeEach(() => {
    setShowLoading = jest.fn();
    fetch.resetMocks();
  });

  it('status 200 returns test cases', async () => {
    const data = [
      {
        apiSpecification: {
          name: 'Account and Transaction API Specification',
          url: 'https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0',
          version: 'v3.0',
          schemaVersion: 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json',
        },
        testCases: [
          {
            '@id': '#t1000',
            name: 'Create Account Access Consents',
            input: {
              method: 'POST',
              endpoint: '/account-access-consents',
              contextGet: {},
            },
            expect: {
              'status-code': 201,
              'schema-validation': true,
              contextPut: {},
            },
          },
        ],
      },
    ];
    fetch.mockResponseOnce(JSON.stringify(data), { status: 200 });

    try {
      const response = await testcases.computeTestCases(setShowLoading);
      expect(response).toEqual(data);

      // assert on the times called and arguments given to fetch
      expect(fetch.mock.calls.length).toEqual(1);
      expect(fetch.mock.calls[0][0]).toEqual(EXPECTED_INPUT);
      expect(fetch.mock.calls[0][1]).toEqual(Object.assign(
        {},
        EXPECTED_INIT,
      ));
    } catch (err) {
    // Should not get here.
      expect(err).toBeFalsy();
    }
  });

  it('status 400 throws error', async () => {
    const data = {
      error: 'error generation test cases, discovery model not set',
    };
    fetch.mockResponseOnce(JSON.stringify(data), { status: 400 });

    try {
      const response = await testcases.computeTestCases(setShowLoading);
      // Should not get here.
      expect(response).toBeFalsy();
    } catch (err) {
      expect(err).toEqual(data);

      // assert on the times called and arguments given to fetch
      expect(fetch.mock.calls.length).toEqual(1);
      expect(fetch.mock.calls[0][0]).toEqual(EXPECTED_INPUT);
      expect(fetch.mock.calls[0][1]).toEqual(Object.assign(
        {},
        EXPECTED_INIT,
      ));
    }
  });
});
