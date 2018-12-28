import discovery from './discovery';

export default {
  // eslint-disable-next-line camelcase
  async validateConfiguration(configuration) {
    const input = '/api/config/global';
    const init = {
      method: 'POST',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: JSON.stringify(configuration),
    };
    const response = await fetch(input, init);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    //
    // Not too sure about the last check for I've commented it out for now.
    if (response.status !== 201 /* || !response.ok */) {
      throw data;
    }

    return data;
  },
  ...discovery,
};
