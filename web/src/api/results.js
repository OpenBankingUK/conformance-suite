import api from './apiUtil';

const REPORT_URL = '/api/report';

export default {
  /**
   * Call GET /api/report
   */
  async computeTestCaseResults() {
    const response = await api.get(REPORT_URL);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
};
