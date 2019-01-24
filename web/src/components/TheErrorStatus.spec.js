import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import BootstrapVue from 'bootstrap-vue';
import status from '../store/modules/status';
import TheErrorStatus from './TheErrorStatus.vue';

const localVue = createLocalVue();

localVue.use(Vuex);
localVue.use(BootstrapVue);

describe('TheErrorStatus.vue', () => {
  let state;

  const stateWithErrors = {
    errors: [
      new Error('Error message'),
      'text message',
    ],
  };

  const stateNoErrors = {
    errors: [],
  };

  const mockStore = (errors) => {
    state = {
      errors,
    };

    return new Vuex.Store({
      modules: {
        status: {
          namespaced: true,
          state,
          getters: status.getters,
        },
      },
    });
  };

  const component = ({ errors }) => {
    const store = mockStore(errors);
    return shallowMount(TheErrorStatus, {
      store,
      localVue,
    });
  };

  it('does not render when no errors', () => {
    const wrapper = component(stateNoErrors);
    expect(wrapper.find('.error-status').exists()).toBe(false);
  });

  it('renders when errors', () => {
    const wrapper = component(stateWithErrors);
    expect(wrapper.find('.error-status').exists()).toBe(true);
  });

  it('displays error messages when errors', () => {
    const wrapper = component(stateWithErrors);
    expect(wrapper.text()).toMatch(/Error message/);
    expect(wrapper.text()).toMatch(/text message/);
  });
});
