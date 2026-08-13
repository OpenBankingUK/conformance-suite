import { shallowMount, createLocalVue } from '@vue/test-utils';
import cloneDeep from 'lodash/cloneDeep';
import Vuex from 'vuex';
import { mutations, state as configState } from '../../store/modules/config';
import PaymentFrequency from './PaymentFrequency.vue';

const localVue = createLocalVue();
localVue.use(Vuex);

const component = (paymentFrequency) => {
  const state = cloneDeep(configState);
  state.configuration.payment_frequency = paymentFrequency;

  const store = new Vuex.Store({
    modules: {
      config: {
        namespaced: true,
        state,
        mutations,
      },
    },
  });

  return shallowMount(PaymentFrequency, {
    localVue,
    store,
    stubs: [
      'b-container',
      'b-row',
      'b-col',
      'b-form-group',
      'b-input-group',
      'b-form-input',
    ],
  });
};

describe('PaymentFrequency', () => {
  it('validates full Frequency_1 examples', () => {
    ['NotKnown', 'EvryDay', 'IntrvlDay:02', 'IntrvlDay:31', 'IntrvlWkDay:01:03', 'QtrDay:ENGLISH'].forEach((value) => {
      const wrapper = component(value);

      expect(wrapper.vm.payment_frequency_state).toEqual(true);
    });
  });

  it('rejects invalid Frequency_1 values', () => {
    ['IntrvlDay:01', 'IntrvlDay:32', 'WEEK'].forEach((value) => {
      const wrapper = component(value);

      expect(wrapper.vm.payment_frequency_state).toEqual(false);
    });
  });

  it('uses a text input instead of a schedule-code dropdown', () => {
    const wrapper = component('EvryDay');

    expect(wrapper.find('b-form-input-stub#payment_frequency').exists()).toEqual(true);
    expect(wrapper.find('b-form-select-stub').exists()).toEqual(false);
  });
});
