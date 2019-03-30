import api from './apiUtil';

const TESTCASES_URL = '/api/test-cases';
const EXECUTE_URL = '/api/run';

export default {
  /**
   * Call GET /api/test-cases
   */
  async computeTestCases(setShowLoading) {
    const response = await api.get(TESTCASES_URL, setShowLoading);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
  /**
   * Calls POST `/api/run`.
   */
  async executeTestCases(setShowLoading) {
    const response = await api.post(EXECUTE_URL, setShowLoading);
    const data = await response.text();

    // `fetch` does not throw an error even when status is not 201.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 201) {
      throw data;
    }

    return data;
  },
  /**
   * Calls DELETE `/api/run`.
   */
  async stopTestRun(setShowLoading) {
    await api.delete(EXECUTE_URL, setShowLoading);
  },
};
