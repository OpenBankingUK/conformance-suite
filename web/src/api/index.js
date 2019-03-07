import config from './config';
import discovery from './discovery';
import results from './results';
import testcases from './testcases';
import apiUtil from './apiUtil';

const EXPORT_URL = '/api/export';

export default {
  ...config,
  ...discovery,
  ...results,
  ...testcases,
  /**
   * Call GET /api/export
   */
  async exportResults(payload) {
    const response = await apiUtil.post(EXPORT_URL, payload);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
};
