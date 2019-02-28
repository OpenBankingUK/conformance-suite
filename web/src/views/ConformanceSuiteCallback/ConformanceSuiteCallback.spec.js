import { shallowMount, createLocalVue } from '@vue/test-utils';
// https://vue-test-utils.vuejs.org/guides/testing-async-components.html#testing-asynchronous-behavior
import BootstrapVue from 'bootstrap-vue';
import flushPromises from 'flush-promises';
import ConformanceSuiteCallback from './ConformanceSuiteCallback.vue';


const localVue = createLocalVue();
localVue.use(BootstrapVue);

describe('ConformanceSuiteCallback', () => {
  const mount = ($route) => {
    const mocks = {
      $route,
    };
    return shallowMount(ConformanceSuiteCallback, {
      mocks,
      localVue,
      propsData:
      {
        autoCloseOnSuccess: false,
      },
    });
  };

  beforeEach(() => {
    fetch.resetMocks();
  });

  it('fragment', async () => {
    const $route = {
      fullPath: '/conformancesuite/callback#code=a052c795-742d-415a-843f-8a4939d740d1&scope=openid%20accounts&id_token=eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ&state=5a6b0d7832a9fb4f80f1170a',
    };

    const response = {
      code: 'a052c795-742d-415a-843f-8a4939d740d1',
      scope: 'openid accounts',
      id_token: 'eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ',
      state: '5a6b0d7832a9fb4f80f1170a',
    };
    fetch.mockResponseOnce(
      JSON.stringify(response),
      { status: 200 },
    );

    // render the component
    const wrapper = mount($route);

    // assert on the times called and arguments given to fetch
    expect(fetch.mock.calls.length).toEqual(1);
    expect(fetch.mock.calls[0][0]).toEqual('/api/redirect/fragment/ok');

    // assert element values
    expect(wrapper.find('.has-query').text()).toBe('false');
    expect(wrapper.find('.has-fragment').text()).toBe('true');
    expect(wrapper.find('.params').text()).toBe(JSON.stringify(response, null, 2));
    expect(wrapper.find('.is-error').text()).toBe('false');

    await flushPromises();
    expect(wrapper.find('.response').text()).toBe(JSON.stringify(response, null, 2));
  });

  it('query', async () => {
    const $route = {
      fullPath: '/conformancesuite/callback?code=a052c795-742d-415a-843f-8a4939d740d1&scope=openid%20accounts&id_token=eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ&state=5a6b0d7832a9fb4f80f1170a',
    };

    const response = {
      code: 'a052c795-742d-415a-843f-8a4939d740d1',
      scope: 'openid accounts',
      id_token: 'eyJ0eXAiOiJKV1QiLCJraWQiOiJGb2w3SXBkS2VMWm16S3RDRWdpMUxEaFNJek09IiwiYWxnIjoiRVMyNTYifQ.eyJzdWIiOiJtYmFuYSIsImF1ZGl0VHJhY2tpbmdJZCI6IjY5YzZkZmUzLWM4MDEtNGRkMi05Mjc1LTRjNWVhNzdjZWY1NS0xMDMzMDgyIiwiaXNzIjoiaHR0cHM6Ly9tYXRscy5hcy5hc3BzcC5vYi5mb3JnZXJvY2suZmluYW5jaWFsL29hdXRoMi9vcGVuYmFua2luZyIsInRva2VuTmFtZSI6ImlkX3Rva2VuIiwibm9uY2UiOiI1YTZiMGQ3ODMyYTlmYjRmODBmMTE3MGEiLCJhY3IiOiJ1cm46b3BlbmJhbmtpbmc6cHNkMjpzY2EiLCJhdWQiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJjX2hhc2giOiIxbGt1SEFuaVJDZlZNS2xEc0pxTTNBIiwib3BlbmJhbmtpbmdfaW50ZW50X2lkIjoiQTY5MDA3Nzc1LTcwZGQtNGIyMi1iZmM1LTlkNTI0YTkxZjk4MCIsInNfaGFzaCI6ImZ0OWRrQTdTWXdlb2hlZXpjOGFHeEEiLCJhenAiOiI1NGY2NDMwOS00MzNkLTQ2MTAtOTVkMi02M2QyZjUyNTM0MTIiLCJhdXRoX3RpbWUiOjE1Mzk5NDM3NzUsInJlYWxtIjoiL29wZW5iYW5raW5nIiwiZXhwIjoxNTQwMDMwMTgxLCJ0b2tlblR5cGUiOiJKV1RUb2tlbiIsImlhdCI6MTUzOTk0Mzc4MX0.8bm69KPVQIuvcTlC-p0FGcplTV1LnmtacHybV2PTb2uEgMgrL3JNA0jpT2OYO73r3zPC41mNQlMDvVOUn78osQ',
      state: '5a6b0d7832a9fb4f80f1170a',
    };
    fetch.mockResponseOnce(
      JSON.stringify(response),
      { status: 200 },
    );

    // render the component
    const wrapper = mount($route);

    // assert on the times called and arguments given to fetch
    expect(fetch.mock.calls.length).toEqual(1);
    expect(fetch.mock.calls[0][0]).toEqual('/api/redirect/query/ok');

    // assert element values
    expect(wrapper.find('.has-query').text()).toBe('true');
    expect(wrapper.find('.has-fragment').text()).toBe('false');
    expect(wrapper.find('.params').text()).toBe(JSON.stringify(response, null, 2));
    expect(wrapper.find('.is-error').text()).toBe('false');

    await flushPromises();
    expect(wrapper.find('.response').text()).toBe(JSON.stringify(response, null, 2));
  });

  it('error', async () => {
    const $route = {
      fullPath: '/conformancesuite/callback?error_description=JWT%20invalid.%20Expiration%20time%20incorrect.&state=5a6b0d7832a9fb4f80f1170a&error=invalid_request',
    };

    const response = {
      error_description: 'JWT invalid. Expiration time incorrect.',
      state: '5a6b0d7832a9fb4f80f1170a',
      error: 'invalid_request',
    };
    fetch.mockResponseOnce(
      JSON.stringify(response),
      { status: 200 },
    );

    // render the component
    const wrapper = mount($route);

    // assert on the times called and arguments given to fetch
    expect(fetch.mock.calls.length).toEqual(1);
    expect(fetch.mock.calls[0][0]).toEqual('/api/redirect/error');

    // assert element values
    expect(wrapper.find('.has-query').text()).toBe('true');
    expect(wrapper.find('.has-fragment').text()).toBe('false');
    expect(wrapper.find('.params').text()).toBe(JSON.stringify(response, null, 2));
    expect(wrapper.find('.is-error').text()).toBe('true');

    await flushPromises();
    expect(wrapper.find('.response').text()).toBe(JSON.stringify(response, null, 2));
  });

  it('is serverError', async () => {
    const $route = {
      fullPath: '/conformancesuite/callback?code=1234567890',
    };

    const response = {
      code: '1234567890',
    };
    fetch.mockResponseOnce(
      JSON.stringify(response),
      { status: 400 },
    );

    // render the component
    const wrapper = mount($route);

    // assert on the times called and arguments given to fetch
    expect(fetch.mock.calls.length).toEqual(1);
    expect(fetch.mock.calls[0][0]).toEqual('/api/redirect/query/ok');

    // assert element values
    expect(wrapper.find('.has-query').text()).toBe('true');
    expect(wrapper.find('.has-fragment').text()).toBe('false');
    expect(wrapper.find('.params').text()).toBe(JSON.stringify(response, null, 2));
    expect(wrapper.find('.is-error').text()).toBe('false');

    await flushPromises();
    expect(wrapper.find('.serverError').text()).toBe('Error processing callback - expected HTTP 200 OK');
  });
});
