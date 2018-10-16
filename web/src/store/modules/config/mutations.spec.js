import mutations from './mutations';
import * as types from './mutation-types';

describe('Config', () => {
  describe('mutations', () => {
    let state;
    const config = { a: 1 };
    const payload = { b: 1 };

    beforeEach(() => {
      state = {
        payload: [],
      };
    });

    it(`${types.SET_CONFIG}`, () => {
      const expectedState = {
        main: config,
        payload: [],
      };
      mutations[types.SET_CONFIG](state, config);
      expect(state).toEqual(expectedState);
    });

    it(`${types.SET_PAYLOAD}`, () => {
      const expectedState = {
        payload,
      };
      mutations[types.SET_PAYLOAD](state, payload);
      expect(state).toEqual(expectedState);
    });

    it(`${types.UPLOAD_PAYLOAD}`, () => {
      const expectedState = {
        payload: [payload],
      };
      mutations[types.UPDATE_PAYLOAD](state, payload);
      expect(state).toEqual(expectedState);
    });

    it(`${types.DELETE_PAYLOAD}`, () => {
      const initialState = {
        payload: [payload],
      };
      const expectedState = {
        payload: [],
      };
      mutations[types.DELETE_PAYLOAD](initialState, payload);
      expect(state).toEqual(expectedState);
    });
  });
});
