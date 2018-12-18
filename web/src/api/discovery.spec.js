import discovery from './discovery';

describe('validateDiscoveryConfig', () => {
  const discoveryModelStub = { };

  describe('when validation passes', () => {
    it('returns success true, and empty array of validation problem strings', async () => {
      fetch.mockResponseOnce(
        JSON.stringify(discoveryModelStub),
        { status: 201 },
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
        expect(e).toEqual(new Error('Expected 201 OK or 400 BadRequest Status.'));
      }
    });
  });
});

describe('annotationsAndMarkers', () => {
  describe('when no discoveryProblems', () => {
    it('returns empty array', async () => {
      const discoveryProblems = null;
      const { annotations, markers } = discovery.annotationsAndMarkers(discoveryProblems, '');
      expect(annotations).toEqual([]);
      expect(markers).toEqual([]);
    });
  });

  describe('when discoveryProblem at child location', () => {
    it('returns annotation at parent location when child not present', () => {
      const json = `{
        "discoveryModel": {
          "name": "example"
        }
      }`;
      const discoveryProblems = [
        {
          path: 'discoveryModel.discoveryVersion',
          parent: 'discoveryModel',
          error: 'Field validation for \'DiscoveryVersion\' failed on the \'required\' tag',
        },
      ];
      const { annotations, markers } = discovery.annotationsAndMarkers(discoveryProblems, json);
      expect(annotations).toEqual([
        {
          row: 1,
          column: 8,
          type: 'error',
          text: discoveryProblems[0].error,
        },
      ]);
      expect(markers).toEqual([
        {
          startRow: 1,
          startCol: 8,
          endRow: 1,
          endCol: 24,
          className: 'ace_error-marker',
          type: 'background',
        },
      ]);
    });
  });

  describe('when discoveryProblem at child location', () => {
    it('returns annotation at child location when child present', () => {
      const json = `{
        "discoveryModel": {
          "discoveryVersion": "9.9.9"
        }
      }`;
      const discoveryProblems = [
        {
          path: 'discoveryModel.discoveryVersion',
          parent: 'discoveryModel',
          error: 'Unsupported discoveryVersion',
        },
      ];
      const { annotations, markers } = discovery.annotationsAndMarkers(discoveryProblems, json);
      expect(annotations).toEqual([
        {
          row: 2,
          column: 10,
          type: 'error',
          text: discoveryProblems[0].error,
        },
      ]);
      expect(markers).toEqual([{
        className: 'ace_error-marker',
        endCol: 28,
        endRow: 2,
        startCol: 10,
        startRow: 2,
        type: 'background',
      }]);
    });
  });
});
