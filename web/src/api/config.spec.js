import config from './config';

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
      const response = await config.validateConfiguration(data);
      expect(response).toEqual(data);

      // assert on the times called and arguments given to fetch
      expect(fetch.mock.calls.length).toEqual(1);
      expect(fetch.mock.calls[0][0]).toEqual(EXPECTED_INPUT);
      expect(fetch.mock.calls[0][1]).toEqual(Object.assign(
        {},
        EXPECTED_INIT,
        { body: JSON.stringify(data) },
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
      const response = await config.validateConfiguration(data);
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
