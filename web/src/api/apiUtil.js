export default {

  // Async call to post API endpoint, returns promise.
  post(path, obj) {
    return fetch(path, {
      method: 'POST',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: JSON.stringify(obj),
    });
  },
};
