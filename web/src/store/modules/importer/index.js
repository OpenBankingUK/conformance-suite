import api from '../../../api';

const mutationTypes = {
  SET_IS_REVIEW: 'SET_IS_REVIEW',
  SET_IS_RERUN: 'SET_IS_RERUN',
  SET_REPORT_ZIP_ARCHIVE: 'SET_REPORT_ZIP_ARCHIVE',
  SET_IMPORT_RESPONSE: 'SET_IMPORT_RESPONSE',
};

export default {
  namespaced: true,
  state: {
    is_review: false,
    is_rerun: false,
    report_zip_archive: '',
    import_response: '',
  },
  mutationTypes,
  mutations: {
    [mutationTypes.SET_IS_REVIEW](state, value) {
      state.is_review = value;
    },
    [mutationTypes.SET_IS_RERUN](state, value) {
      state.is_rerun = value;
    },
    [mutationTypes.SET_REPORT_ZIP_ARCHIVE](state, value) {
      state.report_zip_archive = value;
    },
    [mutationTypes.SET_IMPORT_RESPONSE](state, value) {
      state.import_response = value;
    },
  },
  actions: {
    async doImport({ commit, state }) {
      // eslint-disable-next-line no-console
      console.log('state.is_review=', state.is_review);
      // eslint-disable-next-line no-console
      console.log('state.is_rerun=', state.is_rerun);
      // eslint-disable-next-line no-console
      console.log('state.report_zip_archive=', state.report_zip_archive);

      if (state.is_review) {
        const results = await api.importReview(state.report_zip_archive);
        commit(mutationTypes.SET_IMPORT_RESPONSE, results);
        return Promise.resolve(results);
      }

      if (state.is_rerun) {
        const results = await api.importRerun(state.report_zip_archive);
        commit(mutationTypes.SET_IMPORT_RESPONSE, results);
        return Promise.resolve(results);
      }

      return Promise.resolve({});
    },
  },
};
