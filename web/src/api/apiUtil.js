export default {

  // Async call to get API endpoint, returns promise.
  async get(path, setShowLoading) {
    if (setShowLoading) {
      setShowLoading(true);
    }
    try {
      const response = await fetch(path, {
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
      const response = await fetch(path, {
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
