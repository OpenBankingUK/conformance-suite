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

    it('setDiscoveryModel', () => {
      const discoveryModel = '{"a": 1}';
      actions.setDiscoveryModel({ commit }, discoveryModel);
      expect(commit).toHaveBeenCalledWith('SET_DISCOVERY_MODEL', discoveryModel);
    });

    it('setConfig', () => {
      const config = '{"a": 1}';
      actions.setConfig({ commit }, config);
      expect(commit).toHaveBeenCalledWith('SET_CONFIG', config);
    });

    it('resetValidationsRun', () => {
      actions.resetValidationsRun({ commit });
      expect(commit).toHaveBeenCalledWith('reporter/SET_WEBSOCKET_LAST_UPDATE', null, { root: true });
      expect(commit).toHaveBeenCalledWith('validations/SET_VALIDATION_DISCOVERY_MODEL', null, { root: true });
    });

    it('startValidation', () => {
      actions.startValidation({ dispatch, getters });
      expect(dispatch).toHaveBeenNthCalledWith(1, 'resetValidationsRun');
      expect(dispatch).toHaveBeenNthCalledWith(
        2,
        'validations/validate', {
          discoveryModel: getters.getDiscoveryModel,
          config: getters.getConfig,
        },
        { root: true },
      );
      expect(routerSpy).toHaveBeenCalledWith('/reports');
    });

    it('updateDiscoveryModel', () => {
      const discoveryModel = '{"a": 1}';
      actions.updateDiscoveryModel({ commit }, discoveryModel);
      expect(commit).toHaveBeenCalledWith('UPDATE_DISCOVERY_MODEL', discoveryModel);
    });

    it('deleteDiscoveryModel', () => {
      const discoveryModel = '{"a": 1}';
      actions.deleteDiscoveryModel({ commit }, discoveryModel);
      expect(commit).toHaveBeenCalledWith('DELETE_DISCOVERY_MODEL', discoveryModel);
    });
  });
});
