import { shallowMount, createLocalVue } from '@vue/test-utils';
import cloneDeep from 'lodash/cloneDeep';
import Vuex from 'vuex';
import { mutations, state as configState } from '../../store/modules/config';
import V4StandingOrderFrequency from './V4StandingOrderFrequency.vue';

const localVue = createLocalVue();
localVue.use(Vuex);

const component = (frequency) => {
  const state = cloneDeep(configState);
  state.configuration.v4_standing_order_frequency = frequency;

  const store = new Vuex.Store({
    modules: {
      config: {
        namespaced: true,
        state,
        mutations,
      },
    },
  });

  const wrapper = shallowMount(V4StandingOrderFrequency, {
    localVue,
    store,
    stubs: [
      'b-container',
      'b-row',
      'b-col',
      'b-form-group',
      'b-input-group',
      'b-form-input',
      'b-form-select',
    ],
  });

  return { wrapper, store };
};

describe('V4StandingOrderFrequency', () => {
  it('includes the latest OBFrequency6Code values', () => {
    const { wrapper } = component({ Type: 'WEEK', PointInTime: '03' });
    const values = wrapper.vm.type_options.map(option => option.value);

    expect(values).toEqual(expect.arrayContaining(['LWMH', 'LXMH', 'TWYR']));
  });

  it('treats OBFrequency6Code Type values as structured mode', () => {
    const { wrapper } = component({ Type: 'LXMH' });

    expect(wrapper.vm.structured_type).toEqual('LXMH');
    expect(wrapper.vm.legacy_type).toEqual('');
  });

  it('treats Frequency_1 Type values as legacy mode', () => {
    const { wrapper } = component({ Type: 'IntrvlWkDay:01:03' });

    expect(wrapper.vm.legacy_type).toEqual('IntrvlWkDay:01:03');
    expect(wrapper.vm.structured_type).toEqual(null);
  });

  it('clears structured fields when legacy mode is entered', () => {
    const { wrapper, store } = component({
      Type: 'WEEK',
      PointInTime: '03',
    });

    wrapper.vm.legacy_type = 'IntrvlWkDay:01:03';

    expect(store.state.config.configuration.v4_standing_order_frequency).toEqual({
      Type: 'IntrvlWkDay:01:03',
    });
  });

  it('clears legacy mode when structured mode is selected', () => {
    const { wrapper, store } = component({ Type: 'IntrvlWkDay:01:03' });

    wrapper.vm.structured_type = 'TWYR';

    expect(store.state.config.configuration.v4_standing_order_frequency).toEqual({
      Type: 'TWYR',
    });
  });

  it('clears PointInTime when CountPerPeriod is entered', () => {
    const { wrapper, store } = component({
      Type: 'WEEK',
      PointInTime: '03',
    });

    wrapper.vm.count_per_period = 3;

    expect(store.state.config.configuration.v4_standing_order_frequency).toEqual({
      Type: 'WEEK',
      CountPerPeriod: 3,
    });
  });

  it('validates spec-aligned PointInTime values', () => {
    ['1', '03', '99', '-1'].forEach((PointInTime) => {
      const { wrapper } = component({ Type: 'WEEK', PointInTime });

      expect(wrapper.vm.point_in_time_state).toEqual(true);
    });
  });

  it('rejects invalid PointInTime values', () => {
    ['100', '-10', 'AA'].forEach((PointInTime) => {
      const { wrapper } = component({ Type: 'WEEK', PointInTime });

      expect(wrapper.vm.point_in_time_state).toEqual(false);
    });
  });

  it('validates positive Int32 CountPerPeriod values', () => {
    [1, 2147483647].forEach((CountPerPeriod) => {
      const { wrapper } = component({ Type: 'WEEK', CountPerPeriod });

      expect(wrapper.vm.count_per_period_state).toEqual(true);
    });
  });

  it('rejects non-positive or out-of-range CountPerPeriod values', () => {
    [0, -1, 1.5, 2147483648].forEach((CountPerPeriod) => {
      const { wrapper } = component({ Type: 'WEEK', CountPerPeriod });

      expect(wrapper.vm.count_per_period_state).toEqual(false);
    });
  });
});
