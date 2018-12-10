import discovery from './discovery';

describe('validateDiscoveryConfig', () => {
  const discoveryModelStub = { };

  describe('when validation passes', () => {
    it('returns success true, and empty array of validation problem strings', async () => {
      fetch.mockResponseOnce(
        JSON.stringify(discoveryModelStub),
        { status: 200 },
      );
      const { success, problems } = await discovery.validateDiscoveryConfig(discoveryModelStub);
      expect(success).toBe(true);
      expect(problems).toEqual([]);
    });
  });

  describe('when validation fails', () => {
    const expectedProblems = [
      {
        key: 'DiscoveryModel.Version',
        error: 'Field validation for \'Version\' failed on the \'required\' tag',
      },
      {
        key: 'DiscoveryModel.DiscoveryItems',
        error: 'Field validation for \'DiscoveryItems\' failed on the \'required\' tag',
      },
    ];
    it('returns success false, and array of validation problem strings', async () => {
      fetch.mockResponseOnce(
        JSON.stringify({ error: expectedProblems }),
        { status: 400 },
      );
      const { success, problems } = await discovery.validateDiscoveryConfig(discoveryModelStub);
      expect(success).toBe(false);
      expect(problems).toEqual(expectedProblems);
    });
  });

  describe('when unexpected status code', () => {
    it('should throw an error', async () => {
      fetch.mockResponseOnce(
        JSON.stringify(discoveryModelStub),
        { status: 500 },
      );
      try {
        expect.assertions(1);
        await discovery.validateDiscoveryConfig(discoveryModelStub);
      } catch (e) {
        expect(e).toEqual(new Error('Expected 200 OK or 400 BadRequest Status.'));
      }
    });
  });
});
