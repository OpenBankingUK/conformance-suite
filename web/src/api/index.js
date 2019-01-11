import discovery from './discovery';

const INPUT_PREFIX = '/api';

export default {
  /**
   * Call POST /api/config/global
   * @param {*} configuration Object containing signing_private,
   * signing_public, transport_private and transport_public fields.
   */
  async validateConfiguration(configuration) {
    const input = `${INPUT_PREFIX}/config/global`;
    const init = {
      method: 'POST',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: JSON.stringify(configuration),
    };
    const response = await fetch(input, init);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 201) {
      throw data;
    }

    return data;
  },
  /**
   * Call GET /api/test-cases
   */
  async computeTestCases() {
    const input = `${INPUT_PREFIX}/test-cases`;
    const init = {
      method: 'GET',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
    };
    const response = await fetch(input, init);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
  /**
   * Call GET /api/report
   */
  async computeTestCaseResults() {
    const input = `${INPUT_PREFIX}/report`;
    const init = {
      method: 'GET',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
    };
    const response = await fetch(input, init);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
  ...discovery,
};
