/**
 * This creates a real store so avoid having to mock things.
 * This makes testing much easier.
 *
 * See the recommendation:
 * https://vue-test-utils.vuejs.org/guides/using-with-vuex.html#testing-a-running-store
 */
import { createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import cloneDeep from 'lodash/cloneDeep';
import moment from 'moment';

import exporter from './index';
import api from '../../../api';
// https://jestjs.io/docs/en/mock-functions#mocking-modules
jest.mock('../../../api');

describe('store/modules/exporter', () => {
  beforeEach(() => {
    jest.resetAllMocks();
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  /**
   * Creates a real store so we don't have to mock things out.
   */
  const createRealStore = () => {
    const localVue = createLocalVue();
    localVue.use(Vuex);
    const store = new Vuex.Store(cloneDeep(exporter));

    return store;
  };

  it('initial state', async () => {
    expect.assertions(1);
    const store = createRealStore();

    expect(store.state).toStrictEqual({
      implementer: '',
      authorised_by: '',
      job_title: '',
      has_agreed: false,
      add_digital_signature: false,
      export_results_blob: null,
      export_results_filename: '',
    });
  });

  describe('actions', () => {
    it('exportResults ok', async () => {
      expect.assertions(5);
      const store = createRealStore();
      const EXPORT_RESULTS = 'blob';

      expect(store.state.export_results_blob).toBeNull();

      api.exportResults.mockReturnValueOnce(EXPORT_RESULTS);
      await store.dispatch('exportResults');

      expect(store.state.export_results_blob).toBe(EXPORT_RESULTS);
      expect(store.state.export_results_filename).toMatch(/^report_/);
      expect(store.state.export_results_filename).toMatch(/\.zip$/);

      // remove prefix and post and check it is valid date
      const filename = store.state.export_results_filename.replace(/^report_/, '').replace(/\.zip$/, '');
      const date = moment(filename, 'report_YYYY-MM-DDTHH:mm:ssZ.zip');
      expect(date.isValid()).toBe(true);
    });
  });

  describe('mutations', () => {
    it('SET_IMPLEMENTER', async () => {
      expect.assertions(2);
      const store = createRealStore();
      const VALUE = 'Venom';

      expect(store.state.implementer).toBe('');
      store.commit(exporter.mutationTypes.SET_IMPLEMENTER, VALUE);
      expect(store.state.implementer).toBe(VALUE);
    });

    it('SET_AUTHORISED_BY', async () => {
      expect.assertions(2);
      const store = createRealStore();
      const VALUE = 'Venom';

      expect(store.state.authorised_by).toBe('');
      store.commit(exporter.mutationTypes.SET_AUTHORISED_BY, VALUE);
      expect(store.state.authorised_by).toBe(VALUE);
    });

    it('SET_JOB_TITLE', async () => {
      expect.assertions(2);
      const store = createRealStore();
      const VALUE = 'Venom';

      expect(store.state.job_title).toBe('');
      store.commit(exporter.mutationTypes.SET_JOB_TITLE, VALUE);
      expect(store.state.job_title).toBe(VALUE);
    });

    it('SET_HAS_AGREED', async () => {
      expect.assertions(2);
      const store = createRealStore();
      const VALUE = true;

      expect(store.state.has_agreed).toBe(false);
      store.commit(exporter.mutationTypes.SET_HAS_AGREED, VALUE);
      expect(store.state.has_agreed).toBe(VALUE);
    });

    it('SET_ADD_DIGITAL_SIGNATURE', async () => {
      expect.assertions(2);
      const store = createRealStore();
      const VALUE = true;

      expect(store.state.add_digital_signature).toBe(false);
      store.commit(exporter.mutationTypes.SET_ADD_DIGITAL_SIGNATURE, VALUE);
      expect(store.state.add_digital_signature).toBe(VALUE);
    });

    it('SET_EXPORT_RESULTS_BLOB', async () => {
      expect.assertions(2);
      const store = createRealStore();
      const VALUE = 'blob';

      expect(store.state.export_results_blob).toBeNull();
      store.commit(exporter.mutationTypes.SET_EXPORT_RESULTS_BLOB, VALUE);
      expect(store.state.export_results_blob).toBe(VALUE);
    });
  });
});
