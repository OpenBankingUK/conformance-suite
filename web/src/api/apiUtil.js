
const fetchWithTimeout = (url, timeout, options) => Promise.race([
  fetch(url, options),
  new Promise((_, reject) => setTimeout(() => reject(new Error(`Request timed out: ${url} ${JSON.stringify(options)}`)), timeout)),
]);

const FETCH_TIMEOUT = 90000; // 90 seconds

const doRequest = async (path, setShowLoading, method) => {
  if (setShowLoading) {
    setShowLoading(true);
  }
  try {
    const response = await fetchWithTimeout(path, FETCH_TIMEOUT, {
      method,
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
    });
    return response;
  } catch (e) {
    console.log(JSON.stringify(e)); //eslint-disable-line
    throw e;
  } finally {
    if (setShowLoading) {
      setShowLoading(false);
    }
  }
};

export default {
  Headers: {
    HeaderAccept: 'Accept',
    HeaderContentType: 'Content-Type',
  },
  // Async GET request to API endpoint, returns promise.
  async get(path, setShowLoading) {
    return doRequest(path, setShowLoading, 'GET');
  },
  // Async DELETE request API endpoint, returns promise.
  async delete(path, setShowLoading) {
    return doRequest(path, setShowLoading, 'DELETE');
  },
  // Async call to post API endpoint, returns promise.
  async post(path, obj, setShowLoading, headers = {
    [this.Headers.HeaderAccept]: 'application/json; charset=UTF-8',
    [this.Headers.HeaderContentType]: 'application/json; charset=UTF-8',
  }) {
    if (setShowLoading) {
      setShowLoading(true);
    }
    try {
      const response = await fetchWithTimeout(path, FETCH_TIMEOUT, {
        method: 'POST',
        headers,
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
