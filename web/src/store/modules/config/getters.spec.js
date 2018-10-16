import getters from './getters';

describe('Config', () => {
  let state;
  const example = '{"a": 1}';

  beforeEach(() => {
    state = {
      raw: example,
      parsed: JSON.parse(example),
      payload: {
        raw: example,
        parsed: JSON.parse(example),
      },
    };
  });

  describe('getters', () => {
    it('getConfig', () => {
      expect(getters.getConfig(state)).toEqual(state.main);
    });

    it('getPayload', () => {
      expect(getters.getPayload(state)).toEqual(state.payload);
    });
  });
});
