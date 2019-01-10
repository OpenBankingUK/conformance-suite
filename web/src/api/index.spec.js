/**
 * This test makes use of jest-fetch-mock.
 * See: https://github.com/jefflau/jest-fetch-mock#api
 */
import api from './';

describe('web/src/api', () => {
  describe('validateConfiguration', () => {
    const EXPECTED_INPUT = '/api/config/global';
    const EXPECTED_INIT = {
      method: 'POST',
      headers:
    {
      Accept: 'application/json; charset=UTF-8',
      'Content-Type': 'application/json; charset=UTF-8',
    },
    };

    beforeEach(() => {
      fetch.resetMocks();
    });

    it('status 201 returns original payload', async () => {
      const data = {
        signing_private: 'does_not_matter_what_the_value_is',
        signing_public: 'does_not_matter_what_the_value_is',
        transport_private: 'does_not_matter_what_the_value_is',
        transport_public: 'does_not_matter_what_the_value_is',
      };
      fetch.mockResponseOnce(JSON.stringify(data), { status: 201 });

      try {
        const response = await api.validateConfiguration(data);
        expect(response).toEqual(data);

        // assert on the times called and arguments given to fetch
        expect(fetch.mock.calls.length).toEqual(1);
        expect(fetch.mock.calls[0][0]).toEqual(EXPECTED_INPUT);
        expect(fetch.mock.calls[0][1]).toEqual(Object.assign(
          {}
          , EXPECTED_INIT
          , { body: JSON.stringify(data) },
        ));
      } catch (err) {
      // Should not get here.
        expect(err).toBeFalsy();
      }
    });

    it('status 400 throws error', async () => {
      const data = {
        error: "error with signing certificate: error with public key: asn1: structure error: tags don't match (16 vs {class:0 tag:2 length:1 isCompound:false}) {optional:false explicit:false application:false private:false defaultValue:\u003cnil\u003e tag:\u003cnil\u003e stringType:0 timeType:0 set:false omitEmpty:false} tbsCertificate @2",
      };
      fetch.mockResponseOnce(JSON.stringify(data), { status: 400 });

      try {
        const response = await api.validateConfiguration(data);
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
          { body: JSON.stringify(data) },
        ));
      }
    });
  });

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

    beforeEach(() => {
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
        const response = await api.computeTestCases(data);
        expect(response).toEqual(data);

        // assert on the times called and arguments given to fetch
        expect(fetch.mock.calls.length).toEqual(1);
        expect(fetch.mock.calls[0][0]).toEqual(EXPECTED_INPUT);
        expect(fetch.mock.calls[0][1]).toEqual(Object.assign(
          {}
          , EXPECTED_INIT,
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
        const response = await api.computeTestCases(data);
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
});
