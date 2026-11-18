import { mount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import BootstrapVue from 'bootstrap-vue';
import SpecificationHeader from './SpecificationHeader.vue';

jest.mock('axios', () => ({ get: jest.fn(), post: jest.fn() }));

const localVue = createLocalVue();
localVue.use(Vuex);
localVue.use(BootstrapVue);

const specName = 'Account and Transaction API Specification';
const consentUrl = 'https://example.com/auth?state=Token001';

const createStore = () => new Vuex.Store({
  modules: {
    config: {
      namespaced: true,
      getters: {
        tokenAcquisition: () => 'psu',
        callbackProxyUrl: () => '',
      },
    },
    testcases: {
      namespaced: true,
      state: { consentUrls: { [specName]: [consentUrl] } },
      getters: { tokenAcquired: () => () => false },
    },
    status: {
      namespaced: true,
      actions: { setErrors: () => {} },
    },
  },
});

describe('SpecificationHeader.vue', () => {
  it('opens the PSU consent popup without following href="#" (keeps scroll position)', () => {
    const openPopup = jest.fn();
    const wrapper = mount(SpecificationHeader, {
      localVue,
      store: createStore(),
      propsData: { apiSpecification: { name: specName, url: '', schemaVersion: '' } },
      methods: { openPopup },
    });

    const link = wrapper.find('.psu-consent-link');
    expect(link.exists()).toBe(true);

    const event = new MouseEvent('click', { bubbles: true, cancelable: true });
    link.element.dispatchEvent(event);

    expect(event.defaultPrevented).toBe(true);
    expect(openPopup).toHaveBeenCalledWith(consentUrl, 'PSU Consent', expect.any(Number), expect.any(Number));
  });
});
