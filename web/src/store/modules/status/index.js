import state from './state';
import * as mutationTypes from './mutation-types';
import actions from './actions';
import mutations from './mutations';
import getters from './getters';

export default {
  namespaced: true,
  state,
  actions,
  mutations,
  getters,
  mutationTypes,
};
