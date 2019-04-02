import { getters } from './index';

describe('discoveryTemplates', () => {
  let state;
  const name = 'ob-v3.0-random';

  beforeEach(() => {
    state = {
      discoveryTemplates: [
        {
          model: {
            discoveryModel: { name },
          },
          image: 'imageData',
        },
      ],
    };
  });

  it('returns template with image set to PNG matching template name', async () => {
    const list = await getters.discoveryTemplates(state);
    expect(list[0].model).toEqual(state.discoveryTemplates[0].model);
    expect(list[0].image).toEqual('imageData');
  });
});

describe('discoveryProblems', () => {
  let state;

  beforeEach(() => {
    state = {
      problems: [
        {
          key: 'DiscoveryModel.DiscoveryItems[0].APISpecification.Name',
          error: 'Field validation for \'Name\' failed on the \'required\' tag',
        },
        {
          key: null,
          error: 'Unexpected token { in JSON at position 108',
        },
      ],
    };
  });

  it('returns object with JSON `path` and `parent` property for Key/Error problem string', () => {
    const list = getters.discoveryProblems(state);
    expect(list[0].path).toEqual('discoveryModel.discoveryItems[0].apiSpecification.name');
    expect(list[0].parent).toEqual('discoveryModel.discoveryItems[0].apiSpecification');
  });

  it('returns object with `error` property for Key/Error problem string', () => {
    const list = getters.discoveryProblems(state);
    expect(list[0].error).toEqual(`Field validation for 'Name' failed on the 'required' tag`); // eslint-disable-line
  });

  it('returns object with `error` property and null `path` for non Key/Error problem string', () => {
    const list = getters.discoveryProblems(state);
    expect(list[1].path).toEqual(null);
    expect(list[1].error).toEqual(`Unexpected token { in JSON at position 108`); // eslint-disable-line
  });

  it('returns null when null problems', () => {
    const list = getters.discoveryProblems({ problems: null });
    expect(list).toEqual(null);
  });
});

describe('Config', () => {
  let state;

  beforeEach(() => {
    state = {
      discoveryModel: {
        discoveryModel: {
          tokenAcquisition: 'headless',
        },
      },
    };
  });

  describe('getters', () => {
    it('discoveryModel', () => {
      expect(getters.discoveryModel(state)).toEqual(state.discoveryModel);
    });
    it('discoveryModelString', () => {
      expect(getters.discoveryModelString(state)).toEqual(JSON.stringify(state.discoveryModel, null, 2));
    });
    it('tokenAcquisition', () => {
      expect(getters.tokenAcquisition(state)).toEqual('headless');
    });
  });
});
