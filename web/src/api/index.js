import config from './config';
import consentCallback from './consentCallback';
import discovery from './discovery';
import testcases from './testcases';
import apiUtil from './apiUtil';

const REPORT_URL = '/api/report';
const EXPORT_URL = '/api/export';
const IMPORT_REVIEW = '/api/import/review';
const IMPORT_RERUN = '/api/import/rerun';

export default {
  ...apiUtil,
  ...config,
  ...consentCallback,
  ...discovery,
  ...testcases,
  /**
   * Call GET `/api/report`.
   */
  async computeTestCaseResults(setShowLoading) {
    const response = await apiUtil.get(REPORT_URL, setShowLoading);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
  /**
   * Call GET `/api/export`.
   * @param {*} payload See `ExportRequest` in `pkg/server/models/export.go`.
   */
  async exportResults(payload) {
    const headers = {
      [apiUtil.Headers.HeaderAccept]: 'application/zip',
      [apiUtil.Headers.HeaderContentType]: 'application/json; charset=UTF-8',
    };
    const response = await apiUtil.post(EXPORT_URL, payload, null, headers);
    const data = await response.blob();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
  /**
   * Call POST `/api/import/review`.
   * @param {*} payload See `ImportRequest` in `pkg/server/models/import.go`.
   */
  async importReview(payload) {
    const response = await apiUtil.post(IMPORT_REVIEW, payload);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
  /**
   * Call POST `/api/import/rerun`.
   * @param {*} payload See `ImportRequest` in `pkg/server/models/import.go`.
   */
  async importRerun(payload) {
    const response = await apiUtil.post(IMPORT_RERUN, payload);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 200) {
      throw data;
    }

    return data;
  },
};
