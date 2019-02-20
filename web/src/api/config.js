import api from './apiUtil';

const CONFIG_URL = '/api/config/global';

export default {
  /**
   * Call POST /api/config/global
   * @param {*} configuration Object containing signing_private,
   * signing_public, transport_private and transport_public fields.
   * @param setShowLoading function to handle setShowLoading(true/false) calls
   */
  async validateConfiguration(configuration, setShowLoading) {
    const response = await api.post(CONFIG_URL, configuration, setShowLoading);
    const data = await response.json();

    // `fetch` does not throw an error even when status is not 200.
    // See: https://github.com/whatwg/fetch/issues/18
    if (response.status !== 201) {
      throw data;
    }

    return data;
  },
};
