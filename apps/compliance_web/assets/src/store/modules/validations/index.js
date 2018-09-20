import getters from './getters';
import mutations from './mutations';
import actions from './actions';
import * as mutationTypes from './mutation-types';
import state from './state';

export default {
  namespaced: true,
  state,
  actions,
  mutations,
  getters,
  mutationTypes,
};
