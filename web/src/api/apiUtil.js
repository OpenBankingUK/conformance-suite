export default {

  // Calls post API endpoint, returns response.
  async post(path, obj) {
    return window.fetch(path, {
      method: 'POST',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: JSON.stringify(obj),
    });
  },
};
