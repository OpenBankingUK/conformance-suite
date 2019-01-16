const TESTCASES_URL = '/api/test-cases';
const EXECUTE_URL = '/api//run/start';

export default {
  /**
   * Call GET /api/test-cases
   */
  async computeTestCases() {
    const init = {
      method: 'GET',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
    };
    const response = await fetch(TESTCASES_URL, init);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
  /**
   * Calls POST `/api/run/start`.
   */
  async executeTestCases() {
    const init = {
      method: 'POST',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
    };
    const response = await fetch(EXECUTE_URL, init);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 201.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 201) {
      throw data;
    }

    return data;
  },
};
