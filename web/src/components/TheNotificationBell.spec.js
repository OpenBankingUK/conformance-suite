import { mount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import BootstrapVue from 'bootstrap-vue';
import status from '../store/modules/status';
import TheNotificationBell from './TheNotificationBell.vue';

const localVue = createLocalVue();

localVue.use(Vuex);
localVue.use(BootstrapVue);

describe('TheNotificationBell.vue', () => {
  let state;

  const stateWithNotifications = {
    notifications: [
      {
        message: 'example message',
        extURL: 'https://www.example.com',
      },
      {
        message: 'another example message',
        extURL: 'https://www.examplecorp.com',
        infoUrl: '/example-info-alert',
      },
    ],
  };

  const stateNoNotifications = {
    notifications: [],
  };

  const mockStore = (notifications) => {
    state = {
      notifications,
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

  const component = ({ notifications }) => {
    const store = mockStore(notifications);
    return mount(TheNotificationBell, {
      store,
      localVue,
    });
  };

  it('when no notifications, does not render counter', () => {
    const wrapper = component(stateNoNotifications);
    expect(wrapper.find('.odometer-value').exists()).toBe(false);
  });

  it('when notifications, does render counter', () => {
    const wrapper = component(stateWithNotifications);
    expect(wrapper.find('.odometer-value').exists()).toBe(true);
  });

  it('when notifications, does render correct counter value', () => {
    const wrapper = component(stateWithNotifications);
    expect(wrapper.find('.odometer-value').text()).toBe(stateWithNotifications.notifications.length.toString());
  });
});
