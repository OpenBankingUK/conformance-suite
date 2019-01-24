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
    const errorsResponse = [
      {
        key: 'DiscoveryModel.Version',
        error: "Field 'DiscoveryModel.Version' is required",
      },
      {
        key: 'DiscoveryModel.DiscoveryItems',
        error: "Field 'DiscoveryModel.DiscoveryItems' is required",
      },
      {
        key: 'DiscoveryModel.DiscoveryItems[0].APISpecification.Name',
        error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.Name' is required",
      },
      {
        key: 'DiscoveryModel.DiscoveryItems[0].APISpecification.URL',
        error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.URL' is required",
      },
    ];
    const expectedProblems = [
      {
        key: 'DiscoveryModel.Version',
        error: "Field 'discoveryModel.version' is required",
      },
      {
        key: 'DiscoveryModel.DiscoveryItems',
        error: "Field 'discoveryModel.discoveryItems' is required",
      },
      {
        key: 'DiscoveryModel.DiscoveryItems[0].APISpecification.Name',
        error: "Field 'discoveryModel.discoveryItems[0].apiSpecification.name' is required",
      },
      {
        key: 'DiscoveryModel.DiscoveryItems[0].APISpecification.URL',
        error: "Field 'discoveryModel.discoveryItems[0].apiSpecification.url' is required",
      },
    ];
    it('returns success false, and array of validation problem strings', async () => {
      fetch.mockResponseOnce(
        JSON.stringify({ error: errorsResponse }),
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
          error: "Field 'discoveryModel.discoveryVersion' is required",
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
