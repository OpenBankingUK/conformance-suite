import { mutationTypes as types, mutations } from './index';

describe('Config', () => {
  describe('mutations', () => {
    let state;
    const discoveryModel = { b: 1 };

    beforeEach(() => {
      state = {
        discoveryModel: [],
      };
    });

    it(`${types.SET_DISCOVERY_MODEL}`, () => {
      const expectedState = {
        discoveryModel,
      };
      mutations[types.SET_DISCOVERY_MODEL](state, discoveryModel);
      expect(state).toEqual(expectedState);
    });
  });
});
