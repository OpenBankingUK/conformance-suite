import api from './apiUtil';
import jsonLocation from './jsonLocation';

export default {

  // Calls validate endpoint, returns {success, problemsArray}.
  async validateDiscoveryConfig(discoveryModel) {
    const response = await api.post('/api/discovery-model/validate', discoveryModel);
    const { status } = response;

    if (status !== 200 && status !== 400) {
      throw new Error('Expected 200 OK or 400 BadRequest Status.');
    }

    const validationFailed = status === 400;
    if (validationFailed) {
      const json = await response.json();
      if (json.error) {
        const problems = json.error;
        return { success: false, problems };
      }
    }
    return { success: true, problems: [] };
  },

  annotationsAndMarkers(discoveryProblems, discoveryModelString) {
    if (discoveryProblems === null) {
      return {
        annotations: [],
      };
    }
    const paths = jsonLocation.parse(discoveryModelString);
    const locatableProblems = discoveryProblems.filter(p =>
      p.path &&
      (paths[p.path] || paths[p.parent]));

    if (locatableProblems.length === 0) {
      return {
        annotations: [],
      };
    }
    const annotations = locatableProblems.map((problem) => {
      const { path, parent, error } = problem;
      const { start } = paths[path] || paths[parent];
      const { column, line } = start;
      const row = line - 1;
      return {
        row, column, type: 'error', text: error,
      };
    });
    return {
      annotations,
    };
  },
};
