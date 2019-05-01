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

      if (state.is_review) {
        const payload = {
          report: state.report_zip_archive,
        };
        const results = await api.importReview(payload);
        commit(mutationTypes.SET_IMPORT_RESPONSE, results);
        return Promise.resolve({});
      }

      if (state.is_rerun) {
        const payload = {
          report: state.report_zip_archive,
        };
        const results = await api.importRerun(payload);
        commit(mutationTypes.SET_IMPORT_RESPONSE, results);
        return Promise.resolve({});
      }

      return Promise.resolve({});
    },
  },
};
