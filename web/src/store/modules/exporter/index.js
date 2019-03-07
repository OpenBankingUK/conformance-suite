import pick from 'lodash/pick';
import api from '../../../api';

const mutationTypes = {
  SET_IMPLEMENTER: 'SET_IMPLEMENTER',
  SET_AUTHORISED_BY: 'SET_AUTHORISED_BY',
  SET_JOB_TITLE: 'SET_JOB_TITLE',
  SET_HAS_AGREED: 'SET_HAS_AGREED',
  SET_ADD_DIGITAL_SIGNATURE: 'SET_ADD_DIGITAL_SIGNATURE',
  SET_EXPORT_CONFORMANCE_REPORT: 'SET_EXPORT_CONFORMANCE_REPORT',
  SET_EXPORT_RESULTS: 'SET_EXPORT_RESULTS',
};

export default {
  namespaced: true,
  state: {
    implementer: '',
    authorised_by: '',
    job_title: '',
    has_agreed: false,
    add_digital_signature: false,
    export_results: null,
  },
  mutationTypes,
  mutations: {
    [mutationTypes.SET_IMPLEMENTER](state, value) {
      state.implementer = value;
    },
    [mutationTypes.SET_AUTHORISED_BY](state, value) {
      state.authorised_by = value;
    },
    [mutationTypes.SET_JOB_TITLE](state, value) {
      state.job_title = value;
    },
    [mutationTypes.SET_HAS_AGREED](state, value) {
      state.has_agreed = value;
    },
    [mutationTypes.SET_ADD_DIGITAL_SIGNATURE](state, value) {
      state.add_digital_signature = value;
    },
    [mutationTypes.SET_EXPORT_RESULTS](state, value) {
      state.export_results = value;
    },
  },
  actions: {
    async exportResults({ commit, state, dispatch }) {
      const payload = pick(state, [
        'implementer',
        'authorised_by',
        'job_title',
        'has_agreed',
        'add_digital_signature',
      ]);
      try {
        const results = await api.exportResults(payload);
        commit(mutationTypes.SET_EXPORT_RESULTS, results);
      } catch (err) {
        dispatch('status/setErrors', [err], { root: true });
      }
    },
  },
};
