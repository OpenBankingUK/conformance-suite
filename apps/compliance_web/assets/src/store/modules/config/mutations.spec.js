import mutations from './mutations';
import * as types from './mutation-types';

describe('Config', () => {
  describe('mutations', () => {
    let state;
    const config = '{"a": 1}';
    const payload = '{"b": 1}';

    beforeEach(() => {
      state = {
        payload: {},
      };
    });

    it(`${types.SET_CONFIG}`, () => {
      const expectedState = {
        raw: config,
        payload: {},
      };
      mutations[types.SET_CONFIG](state, config);
      expect(state).toEqual(expectedState);
    });

    it(`${types.SET_PAYLOAD}`, () => {
      const expectedState = {
        payload: {
          raw: payload,
        },
      };
      mutations[types.SET_PAYLOAD](state, payload);
      expect(state).toEqual(expectedState);
    });

    it(`${types.SUBMIT_CONFIG}`, () => {
      const initialState = {
        raw: config,
        payload: {
          raw: payload,
        },
      };
      const expectedState = {
        raw: config,
        parsed: JSON.parse(config),
        payload: {
          parsed: JSON.parse(payload),
          raw: payload,
        },
      };
      mutations[types.SUBMIT_CONFIG](initialState);
      expect(initialState).toEqual(expectedState);
    });
  });
});
