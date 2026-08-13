import { shallowMount, createLocalVue } from '@vue/test-utils';
import cloneDeep from 'lodash/cloneDeep';
import Vuex from 'vuex';
import { getters, mutations, state as configState } from '../../store/modules/config';
import TheConfigurationForm from './TheConfigurationForm.vue';

jest.mock('../../api/apiUtil', () => ({
  get: jest.fn(() => Promise.resolve({ json: () => Promise.resolve([]) })),
}));

const localVue = createLocalVue();
localVue.use(Vuex);

const discoveryModel = version => ({
  discoveryModel: {
    discoveryItems: [
      {
        apiSpecification: {
          name: 'Payment Initiation API',
          version,
          schemaVersion: `https://example.com/${version}/payment-initiation-openapi.json`,
        },
      },
    ],
  },
});

const storeForDiscovery = (model) => {
  const state = cloneDeep(configState);
  state.discoveryModel = model;

  return new Vuex.Store({
    modules: {
      config: {
        namespaced: true,
        state,
        getters,
        mutations,
      },
    },
  });
};

const component = model => shallowMount(TheConfigurationForm, {
  localVue,
  store: storeForDiscovery(model),
  stubs: [
    'b-form',
    'b-card',
    'b-form-group',
    'b-input-group',
    'b-input-group-prepend',
    'b-input-group-append',
    'b-button',
    'b-form-input',
    'b-form-select',
    'b-form-checkbox',
    'b-container',
    'b-row',
    'b-col',
  ],
});

describe('TheConfigurationForm frequency controls', () => {
  it('shows only legacy frequency for v3-only payment discovery', () => {
    const wrapper = component(discoveryModel('v3.1.6'));

    expect(wrapper.find('paymentfrequency-stub').exists()).toEqual(true);
    expect(wrapper.find('v4standingorderfrequency-stub').exists()).toEqual(false);
  });

  it('shows only v4 standing-order frequency for v4-only payment discovery', () => {
    const wrapper = component(discoveryModel('v4.0.1'));

    expect(wrapper.find('paymentfrequency-stub').exists()).toEqual(false);
    expect(wrapper.find('v4standingorderfrequency-stub').exists()).toEqual(true);
  });

  it('shows both controls for mixed payment discovery', () => {
    const model = {
      discoveryModel: {
        discoveryItems: [
          discoveryModel('v3.1.6').discoveryModel.discoveryItems[0],
          discoveryModel('v4.0.1').discoveryModel.discoveryItems[0],
        ],
      },
    };

    const wrapper = component(model);

    expect(wrapper.find('paymentfrequency-stub').exists()).toEqual(true);
    expect(wrapper.find('v4standingorderfrequency-stub').exists()).toEqual(true);
  });

  it('shows both controls when discovery is unavailable', () => {
    const wrapper = component(null);

    expect(wrapper.find('paymentfrequency-stub').exists()).toEqual(true);
    expect(wrapper.find('v4standingorderfrequency-stub').exists()).toEqual(true);
  });
});
