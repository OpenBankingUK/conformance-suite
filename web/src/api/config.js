const CONFIG_URL = '/api/config/global';

export default {
  /**
   * Call POST /api/config/global
   * @param {*} configuration Object containing signing_private,
   * signing_public, transport_private and transport_public fields.
   */
  async validateConfiguration(configuration) {
    const init = {
      method: 'POST',
      headers: {
        Accept: 'application/json; charset=UTF-8',
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: JSON.stringify(configuration),
    };
    const response = await fetch(CONFIG_URL, init);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 201) {
      throw data;
    }

    return data;
  },
};
