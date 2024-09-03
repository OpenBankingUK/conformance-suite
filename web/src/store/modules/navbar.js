export default {
  namespaced: true,
  state: {
    selectedVersion: 'v4.0.0',
  },
  mutations: {
    SET_SELECTED_VERSION(state, version) {
      state.selectedVersion = version;
    },
  },
  actions: {
    updateSelectedVersion({ commit }, version) {
      commit('SET_SELECTED_VERSION', version);
    },
  },
};
