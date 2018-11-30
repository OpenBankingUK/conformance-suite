import mutations from './mutations';
import * as types from './mutation-types';

describe('Config', () => {
  describe('mutations', () => {
    let state;
    const config = { a: 1 };
    const discoveryModel = { b: 1 };

    beforeEach(() => {
      state = {
        discoveryModel: [],
      };
    });

    it(`${types.SET_CONFIG}`, () => {
      const expectedState = {
        main: config,
        discoveryModel: [],
      };
      mutations[types.SET_CONFIG](state, config);
      expect(state).toEqual(expectedState);
    });

    it(`${types.SET_DISCOVERY_MODEL}`, () => {
      const expectedState = {
        discoveryModel,
      };
      mutations[types.SET_DISCOVERY_MODEL](state, discoveryModel);
      expect(state).toEqual(expectedState);
    });

    it(`${types.UPLOAD_DISCOVERY_MODEL}`, () => {
      const expectedState = {
        discoveryModel: [discoveryModel],
      };
      mutations[types.UPDATE_DISCOVERY_MODEL](state, discoveryModel);
      expect(state).toEqual(expectedState);
    });

    it(`${types.DELETE_DISCOVERY_MODEL}`, () => {
      const initialState = {
        discoveryModel: [discoveryModel],
      };
      const expectedState = {
        discoveryModel: [],
      };
      mutations[types.DELETE_DISCOVERY_MODEL](initialState, discoveryModel);
      expect(state).toEqual(expectedState);
    });
  });
});
