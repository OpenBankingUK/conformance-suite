import pick from 'lodash/pick';
import moment from 'moment';
import api from '../../../api';

/**
 * Example return value: `report_2019-03-25T11_41_05+00_00.zip`.
 * @param {*} prefix
 */
const generateFilename = function generateFilename(prefix) {
  const RFC3339 = 'YYYY-MM-DDTHH:mm:ssZ'; // "2006-01-02T15:04:05Z07:00"
  const datetime = moment(new Date()).format(RFC3339);
  const filename = `${prefix}report_${datetime}.zip`;

  return filename;
};

const mutationTypes = {
  SET_ENVIRONMENT: 'SET_ENVIRONMENT',
  SET_IMPLEMENTER: 'SET_IMPLEMENTER',
  SET_AUTHORISED_BY: 'SET_AUTHORISED_BY',
  SET_JOB_TITLE: 'SET_JOB_TITLE',
  SET_PRODUCTS: 'SET_PRODUCTS',
  SET_HAS_AGREED: 'SET_HAS_AGREED',
  SET_ADD_DIGITAL_SIGNATURE: 'SET_ADD_DIGITAL_SIGNATURE',
  SET_EXPORT_CONFORMANCE_REPORT: 'SET_EXPORT_CONFORMANCE_REPORT',
  SET_EXPORT_RESULTS_BLOB: 'SET_EXPORT_RESULTS_BLOB',
  SET_EXPORT_RESULTS_FILENAME: 'SET_EXPORT_RESULTS_FILENAME',
};

export default {
  namespaced: true,
  state: {
    environment: '',
    implementer: '',
    authorised_by: '',
    job_title: '',
    products: [],
    has_agreed: false,
    add_digital_signature: false,
    export_results_blob: null,
    export_results_filename: '',
  },
  mutationTypes,
  mutations: {
    [mutationTypes.SET_ENVIRONMENT](state, value) {
      state.environment = value;
    },
    [mutationTypes.SET_IMPLEMENTER](state, value) {
      state.implementer = value;
    },
    [mutationTypes.SET_AUTHORISED_BY](state, value) {
      state.authorised_by = value;
    },
    [mutationTypes.SET_JOB_TITLE](state, value) {
      state.job_title = value;
    },
    [mutationTypes.SET_PRODUCTS](state, value) {
      state.products = value;
    },
    [mutationTypes.SET_HAS_AGREED](state, value) {
      state.has_agreed = value;
    },
    [mutationTypes.SET_ADD_DIGITAL_SIGNATURE](state, value) {
      state.add_digital_signature = value;
    },
    [mutationTypes.SET_EXPORT_RESULTS_BLOB](state, value) {
      state.export_results_blob = value;
    },
    [mutationTypes.SET_EXPORT_RESULTS_FILENAME](state, value) {
      state.export_results_filename = value;
    },
  },
  actions: {
    async exportResults({ commit, state, dispatch }) {
      const payload = pick(state, [
        'environment',
        'implementer',
        'authorised_by',
        'job_title',
        'products',
        'has_agreed',
        'add_digital_signature',
      ]);
      try {
        commit(mutationTypes.SET_EXPORT_RESULTS_BLOB, null);
        commit(mutationTypes.SET_EXPORT_RESULTS_FILENAME, '');

        const results = await api.exportResults(payload);
        const filename = generateFilename(`${payload.implementer}_`);

        commit(mutationTypes.SET_EXPORT_RESULTS_BLOB, results);
        commit(mutationTypes.SET_EXPORT_RESULTS_FILENAME, filename);
      } catch (err) {
        dispatch('status/setErrors', [err], { root: true });
      }
    },
  },
  generateFilename,
};
