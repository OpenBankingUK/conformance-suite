import { shallowMount, createLocalVue } from '@vue/test-utils';
import Vuex from 'vuex';
import BootstrapVue from 'bootstrap-vue';
import config from '../../store/modules/config';
import ConfigurationFormFile from './ConfigurationFormFile.vue';

const localVue = createLocalVue();

localVue.use(Vuex);
localVue.use(BootstrapVue);

describe('ConfigurationFormFile.vue', () => {
  let actions;
  let state;

  const mockStore = (storeConfigValue) => {
    state = {
      configuration: {
        signing_private: storeConfigValue,
        signing_public: storeConfigValue,
        transport_private: storeConfigValue,
        transport_public: storeConfigValue,
      },
    };

    actions = {
      clearErrors: jest.fn(),
      setErrors: jest.fn(),
    };

    return new Vuex.Store({
      modules: {
        config: {
          namespaced: true,
          state,
          actions,
          // eslint-disable-next-line import/no-named-as-default-member
          getters: config.getters,
        },
      },
    });
  };

  const component = ({ storeConfigValue }, id, validExtension) => {
    const store = mockStore(storeConfigValue);
    return shallowMount(ConfigurationFormFile, {
      store,
      localVue,
      propsData:
      {
        id: id || 'signing_private',
        setterMethodNameSuffix: 'x',
        label: 'y',
        validExtension: validExtension || '.key',
      },
    });
  };

  it('has empty description when store config value blank', () => {
    const wrapper = component({ storeConfigValue: '' });
    const { description } = wrapper.vm;
    expect(description).toBe('');
  });

  it('has description with store config value size', () => {
    const wrapper = component({ storeConfigValue: 'testCert' });
    const { description } = wrapper.vm;
    expect(description).toBe('Size: 8 bytes');
    expect(description).not.toContain('Last modified');
  });

  it('has description with file size and mod date when store config value blank', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'signing_private', '.key');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'example.key', size: 99, lastModifiedDate: Date() },
      validFile: true,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Size: 99 bytes');
    expect(description).toContain('Last modified');
  });

  it('has description with file size when file does not have lastModifiedDate and store config value blank', () => {
    const wrapper = component({ storeConfigValue: '' }, 'signing_private', '.key');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'example.key', size: 99 },
      validFile: true,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Size: 99 bytes');
    expect(description).not.toContain('Last modified');
  });

  it('has description with file size and mod date when store config value matches', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'signing_private', '.key');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'example.key', size: 99, lastModifiedDate: Date() },
      validFile: true,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Size: 99 bytes');
    expect(description).toContain('Last modified');
  });

  it('has description with store value size when store config value does not match file data', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'signing_private', '.key');
    wrapper.setData({
      data: 'differentValue',
      file: { name: 'example.key', size: 99, lastModifiedDate: Date() },
      validFile: true,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Size: 8 bytes');
    expect(description).not.toContain('Last modified');
  });

  it('.key file provided for signing_private file selection - .key required', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'signing_private', '.key');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'file.key', size: 99, lastModifiedDate: Date() },
      validFile: true,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Size: ');
    expect(description).toContain('Last modified: ');
  });

  it('.key file provided for signing_public file selection - .pem required', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'signing_public', '.pem');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'file.key', size: 99, lastModifiedDate: Date() },
      validFile: false,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Require file with extension .pem');
  });

  it('.pem file provided for signing_public file selection - .pem required', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'signing_private', '.pem');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'file.pem', size: 99, lastModifiedDate: Date() },
      validFile: true,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Size: ');
    expect(description).toContain('Last modified: ');
  });

  it('.pem file provided for signing_private file selection - .key required', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'signing_private', '.key');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'file.pem', size: 99, lastModifiedDate: Date() },
      validFile: false,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Require file with extension .key');
    expect(description).toContain('Invalid file format');
  });

  it('.key file provided for transport_private file selection - .key required', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'transport_private', '.key');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'file.key', size: 99, lastModifiedDate: Date() },
      validFile: true,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Size: ');
    expect(description).toContain('Last modified: ');
  });

  it('.key file provided for transport_public file selection - .pem required', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'transport_private', '.pem');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'file.key', size: 99, lastModifiedDate: Date() },
      validFile: false,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Require file with extension .pem');
    expect(description).toContain('Invalid file format');
  });

  it('.pem file provided for transport_public file selection - .pem required', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'transport_public', '.pem');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'file.pem', size: 99, lastModifiedDate: Date() },
      validFile: true,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Size: ');
    expect(description).toContain('Last modified: ');
  });

  it('.pem file provided for transport_private file selection - .key required', () => {
    const wrapper = component({ storeConfigValue: 'testCert' }, 'transport_private', '.key');
    wrapper.setData({
      data: 'testCert',
      file: { name: 'file.pem', size: 99, lastModifiedDate: Date() },
      validFile: false,
    });
    const { description } = wrapper.vm;
    expect(description).toContain('Require file with extension .key');
    expect(description).toContain('Invalid file format');
  });
});
