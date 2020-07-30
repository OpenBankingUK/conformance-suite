import { createLocalVue, mount } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';
import merge from 'lodash/merge';
import VueRouter from 'vue-router';
import Vuex from 'vuex';
import TheWizardFooter from '../../components/Wizard/TheWizardFooter.vue';
import TheStore from '../../store';
import WizardExport from './WizardExport.vue';

const discoveryModel = {
  discoveryModel: {
    discoveryItems: [
      {
        apiSpecification: {
          name: 'Specification A',
          version: 'v3.1.5',
        },
      },
      {
        apiSpecification: {
          name: 'Specification A',
          version: 'v3.1.4',
        },
      },
      {
        apiSpecification: {
          name: 'Specification B',
          version: 'v3.1.5',
        },
      },
    ],
  },
};

describe('WizardExport', () => {
  const mountOptions = () => {
    const localVue = createLocalVue();
    // `VueRouter` is mounted as `TheWizardFooter` requires it, as we are doing a `mount` instead of a `shallowMount`.
    const router = new VueRouter();
    localVue.use(Vuex);
    localVue.use(BootstrapVue);
    localVue.use(VueRouter);

    // discoveryModel is used for testing displayed specification versions
    TheStore.commit('config/SET_DISCOVERY_MODEL', discoveryModel);
    return {
      localVue,
      store: TheStore,
      router,
    };
  };

  const createComponent = () => {
    const options = mountOptions();
    // Use `mount` instead of `shallowMount` as we want the `bootstrap-vue` components
    // to be rendered so we can test that interacting with elements causing updates to the store.
    // See official `bootstrap-vue` tests for examples:
    // * https://github.com/bootstrap-vue/bootstrap-vue/blob/dev/src/components/form-input/form-input.spec.js
    // * https://github.com/bootstrap-vue/bootstrap-vue/blob/dev/src/components/form-checkbox/form-checkbox.spec.js
    const wrapper = mount(WizardExport, options);

    return { wrapper, options };
  };

  /**
   * Remember to reset the store in test case if you mutate it.
   * We only have to do this because we are initialzing a global `vuex` store in `web/src/store/index.js`,
   * if we weren't doing this we wouldn't need to reset the state after the tests.
   * We basically need to refactor `web/src/store/index.js` and `web/src/main.js` if we want to avoid having
   * to call `resetStore` after tests.
   * @param {*} store The store to reset.
   */
  const resetStore = (store) => {
    store.commit('exporter/SET_IMPLEMENTER', '');
    store.commit('exporter/SET_AUTHORISED_BY', '');
    store.commit('exporter/SET_JOB_TITLE', '');
    store.commit('exporter/SET_HAS_AGREED', false);
    store.commit('exporter/SET_ADD_DIGITAL_SIGNATURE', false);
    store.commit('exporter/SET_EXPORT_RESULTS_BLOB', null);
    store.commit('exporter/SET_EXPORT_RESULTS_FILENAME', null);
  };

  test('all elements are present', () => {
    const { wrapper } = createComponent();

    // name of the registered component
    expect(wrapper.name()).toBe('WizardExport');

    // component heading
    expect(wrapper.find('.panel-heading').text()).toBe('Export');

    // tested versions
    expect(wrapper.contains('#versions')).toBe(true);

    // all items are listed for tested versions
    discoveryModel.discoveryModel.discoveryItems.forEach((item) => {
      const id = (item.apiSpecification.version + item.apiSpecification.name).replace(/[^a-zA-Z0-9-]/g, '_');
      expect(wrapper.contains(`#${id}`)).toBe(true);
    });

    // form
    expect(wrapper.contains('#implementer')).toBe(true);
    expect(wrapper.contains('#authorised_by')).toBe(true);
    expect(wrapper.contains('#job_title')).toBe(true);
    expect(wrapper.contains('#has_agreed')).toBe(true);
    expect(wrapper.contains('#add_digital_signature')).toBe(true);

    // footer
    // could remove a lot of redundant checks but leaving them for now.
    expect(wrapper.contains({ name: 'TheWizardFooter' })).toBe(true);
    expect(wrapper.find({ name: 'TheWizardFooter' }).is(TheWizardFooter)).toBe(true);
    expect(wrapper.contains(TheWizardFooter)).toBe(true);

    wrapper.destroy();
  });

  test('form elements are bound to store.exporter.state', () => {
    const { wrapper, options: { store } } = createComponent();
    const state = {
      environment: '',
      implementer: '',
      authorised_by: '',
      job_title: '',
      products: [],
      has_agreed: false,
      add_digital_signature: false,
      export_results_blob: null,
      export_results_filename: '',
    };

    // assert initial state - everything should be empty
    expect(store.state.exporter).toEqual(state);

    // set value on implementer input, then assert value has been updated in the store
    const implementer = 'implementer';
    wrapper.find('#implementer').setValue(implementer);
    expect(store.state.exporter).toEqual(merge(state, { implementer }));

    // set value on authorised_by input, then assert value has been updated in the store
    const authorised_by = 'authorised_by';
    wrapper.find('#authorised_by').setValue(authorised_by);
    expect(store.state.exporter).toEqual(merge(state, { authorised_by, implementer }));

    // set value on job_title input, then assert value has been updated in the store
    const job_title = 'job_title';
    wrapper.find('#job_title').setValue(job_title);
    expect(store.state.exporter).toEqual(merge(state, { job_title, authorised_by, implementer }));

    // set value on has_agreed checkbox, then assert value has been updated in the store
    const has_agreed = true;
    wrapper.find('#has_agreed').setChecked(has_agreed);
    expect(store.state.exporter).toEqual(merge(state, {
      has_agreed, job_title, authorised_by, implementer,
    }));

    // set value on has_agreed checkbox, then assert value has been updated in the store
    const add_digital_signature = true;
    wrapper.find('#add_digital_signature').setChecked(add_digital_signature);
    expect(store.state.exporter).toEqual(merge(state, {
      add_digital_signature, has_agreed, job_title, authorised_by, implementer,
    }));

    resetStore(store);
    wrapper.destroy();
  });

  test('"Export Conformance Report" is enabled when form is valid', () => {
    const { wrapper, options: { store } } = createComponent();

    wrapper.find('#implementer').setValue('implementer');
    wrapper.find('#authorised_by').setValue('authorised_by');
    wrapper.find('#job_title').setValue('job_title');
    wrapper.find('#products').setValue('Business');
    wrapper.find('#has_agreed').setChecked(true);
    wrapper.find('#add_digital_signature').setChecked(true);

    const footer = wrapper.find({ name: 'TheWizardFooter' });
    const next = footer.find('#next');
    expect(footer.props('isNextEnabled')).toBe(true);
    expect(next.attributes('disabled')).toBeFalsy();

    resetStore(store);
    wrapper.destroy();
  });

  test('"Export Conformance Report" is disabled when form is invalid', () => {
    const { wrapper } = createComponent();

    const footer = wrapper.find({ name: 'TheWizardFooter' });
    const next = footer.find('#next');
    expect(footer.props('isNextEnabled')).toBe(false);
    expect(next.attributes('disabled')).toBe('disabled');
    expect(next.text()).toBe('Export Conformance Report');

    wrapper.destroy();
  });

  test('"Export Conformance Report" is rendered as the next button on TheWizardFooter', () => {
    const { wrapper } = createComponent();

    const footer = wrapper.find({ name: 'TheWizardFooter' });
    const next = footer.find('#next');
    expect(next.text()).toBe('Export Conformance Report');

    wrapper.destroy();
  });

  test('download report.zip is not rendered until after clicking "Export Conformance Report"', () => {
    const { wrapper } = createComponent();

    expect(wrapper.find('.download-report-link').exists()).toBe(false);
    // TODO(mbana): Probably rely on e2e for this one instead of unit-testing it.

    wrapper.destroy();
  });

  test('TheErrorStatus not rendered when there are no errors', () => {
    const { wrapper } = createComponent();

    expect(wrapper.find({ name: 'TheErrorStatus' }).exists()).toBe(false);

    wrapper.destroy();
  });
});
