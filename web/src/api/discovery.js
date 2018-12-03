import api from './apiUtil';

export default {

  // Calls validate endpoint, returns {success, problemsArray}.
  async validateDiscoveryConfig(discoveryModel) {
    const response = await api.post('/api/discovery-model/validate', discoveryModel);

    if (response.status !== 200 && response.status !== 400) {
      throw new Error('Expected 200 OK or 400 BadRequest Status.');
    }

    const validationFailed = response.status === 400;
    if (validationFailed) {
      const json = await response.json();
      if (json.error) {
        const problems = json.error;
        return { success: false, problems };
      }
    }
    return { success: true, problems: [] };
  },
};
