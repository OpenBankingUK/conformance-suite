export default {

  // Async call to get API endpoint, returns promise.
  get(path) {
    return fetch(path, {
      method: 'GET',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
    });
  },

  // Async call to post API endpoint, returns promise.
  post(path, obj) {
    return fetch(path, {
      method: 'POST',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: obj ? JSON.stringify(obj) : null,
    });
  },
};
