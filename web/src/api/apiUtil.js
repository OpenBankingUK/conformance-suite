
const fetchWithTimeout = (url, timeout, options) => Promise.race([
  fetch(url, options),
  new Promise((_, reject) => setTimeout(() => reject(new Error(`Request timed out: ${url} ${JSON.stringify(options)}`)), timeout)),
]);

const FETCH_TIMEOUT = 30000; // 30 seconds

export default {

  // Async call to get API endpoint, returns promise.
  async get(path, setShowLoading) {
    if (setShowLoading) {
      setShowLoading(true);
    }
    try {
      const response = await fetchWithTimeout(path, FETCH_TIMEOUT, {
        method: 'GET',
        headers: {
          Accept: 'application/json; charset=UTF-8',
          'Content-Type': 'application/json; charset=UTF-8',
        },
      });
      return response;
    } catch (e) {
      throw e;
    } finally {
      if (setShowLoading) {
        setShowLoading(false);
      }
    }
  },

  // Async call to post API endpoint, returns promise.
  async post(path, obj, setShowLoading) {
    if (setShowLoading) {
      setShowLoading(true);
    }
    try {
      const response = await fetchWithTimeout(path, FETCH_TIMEOUT, {
        method: 'POST',
        headers: {
          Accept: 'application/json; charset=UTF-8',
          'Content-Type': 'application/json; charset=UTF-8',
        },
        body: obj ? JSON.stringify(obj) : null,
      });
      return response;
    } catch (e) {
      throw e;
    } finally {
      if (setShowLoading) {
        setShowLoading(false);
      }
    }
  },
};
