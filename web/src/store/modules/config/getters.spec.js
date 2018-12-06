import getters from './getters';

describe('discoveryProblems', () => {
  let state;

  beforeEach(() => {
    state = {
      problems: [
        `Key: 'Model.DiscoveryModel.DiscoveryItems[0].APISpecification.Name' Error:Field validation for 'Name' failed on the 'required' tag`, // eslint-disable-line
        'Unexpected token { in JSON at position 108',
      ],
    };
  });

  it('returns object with JSON `path` property for Key/Error problem string', () => {
    const list = getters.discoveryProblems(state);
    expect(list[0].path).toEqual('discoveryModel.discoveryItems[0].apiSpecification.name');
  });

  it('returns object with `error` property for Key/Error problem string', () => {
    const list = getters.discoveryProblems(state);
    expect(list[0].error).toEqual(`Field validation for 'Name' failed on the 'required' tag`); // eslint-disable-line
  });

  it('returns object with `error` property and null `path` for non Key/Error problem string', () => {
    const list = getters.discoveryProblems(state);
    expect(list[1].path).toEqual(null); // eslint-disable-line
    expect(list[1].error).toEqual(`Unexpected token { in JSON at position 108`); // eslint-disable-line
  });
});

describe('Config', () => {
  let state;
  const example = '{"a": 1}';

  beforeEach(() => {
    state = {
      raw: example,
      parsed: JSON.parse(example),
      discoveryModel: {
        raw: example,
        parsed: JSON.parse(example),
      },
    };
  });

  describe('getters', () => {
    it('getConfig', () => {
      expect(getters.getConfig(state)).toEqual(state.main);
    });

    it('getDiscoveryModel', () => {
      expect(getters.getDiscoveryModel(state)).toEqual(state.discoveryModel);
    });
  });
});
