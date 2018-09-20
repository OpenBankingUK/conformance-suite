import actions from './actions';
import getters from './getters';
import router from '../../../router';

describe('Config', () => {
  describe('actions', () => {
    let dispatch;
    let commit;
    let routerSpy;

    beforeEach(() => {
      dispatch = jest.fn();
      commit = jest.fn();
      routerSpy = jest.spyOn(router, 'push');
    });

    it('setPayload', () => {
      const payload = '{"a": 1}';
      actions.setPayload({ commit }, payload);
      expect(commit).toHaveBeenCalledWith('SET_PAYLOAD', payload);
    });

    it('setConfig', () => {
      const config = '{"a": 1}';
      actions.setConfig({ commit }, config);
      expect(commit).toHaveBeenCalledWith('SET_CONFIG', config);
    });

    it('resetValidationsRun', () => {
      actions.resetValidationsRun({ commit });
      expect(commit).toHaveBeenCalledWith('reporter/SET_WEBSOCKET_LAST_UPDATE', null, { root: true });
      expect(commit).toHaveBeenCalledWith('validations/SET_VALIDATION_PAYLOAD', null, { root: true });
    });

    it('submitConfig', () => {
      actions.submitConfig({ commit });
      expect(commit).toHaveBeenCalledWith('SUBMIT_CONFIG');
    });

    it('startValidation', () => {
      actions.startValidation({ dispatch, getters });
      expect(dispatch).toHaveBeenNthCalledWith(1, 'resetValidationsRun');
      expect(dispatch).toHaveBeenNthCalledWith(2, 'submitConfig');
      expect(dispatch).toHaveBeenNthCalledWith(
        3,
        'validations/validate', {
          payload: getters.getPayload,
          config: getters.getConfig,
        },
        { root: true },
      );
      expect(routerSpy).toHaveBeenCalledWith('/reports');
    });
  });
});
