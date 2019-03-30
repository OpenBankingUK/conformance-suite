import * as _ from 'lodash';
import api from './apiUtil';
import jsonLocation from './jsonLocation';

const BLANK_ANNOTATION_MARKER = {
  annotations: [],
  markers: [],
};

const calculateAnnotationsAndMarkers = (locatableProblems, paths) => {
  const annotations = [];
  const markers = [];
  locatableProblems.forEach((problem) => {
    const { path, parent, error } = problem;
    const { start, end } = paths[path] || paths[parent];
    const { column, line } = start;
    const row = line - 1;
    annotations.push({
      row,
      column,
      type: 'error',
      text: error,
    });
    markers.push({
      startRow: row,
      startCol: column,
      endRow: end.line - 1,
      endCol: end.column,
      className: 'ace_error-marker',
      type: 'background',
    });
  });
  return {
    annotations,
    markers,
  };
};

const formatProblems = problems => problems.map(({ key, error }) => {
  const reformattedKey = key
    .replace('API', 'Api')
    .replace('URL', 'Url')
    .split('.')
    .map(w => _.lowerFirst(w))
    .join('.');
  const reformatted = error.replace(key, reformattedKey);
  return {
    key,
    error: reformatted,
  };
});

export default {
  // Calls validate endpoint, returns {success, problemsArray}.
  async validateDiscoveryConfig(discoveryModel, setShowLoading) {
    const response = await api.post('/api/discovery-model', discoveryModel, setShowLoading);
    const { status } = response;

    if (status !== 201 && status !== 400) {
      throw new Error('Expected 201 OK or 400 BadRequest Status.');
    }

    const validationFailed = status === 400;
    if (validationFailed) {
      const json = await response.json();
      if (json.error) {
        const problems = formatProblems(json.error);
        return { success: false, problems };
      }
    }

    return { success: true, problems: [], response: await response.json() };
  },

  annotationsAndMarkers(discoveryProblems, discoveryModelString) {
    if (discoveryProblems === null) {
      return BLANK_ANNOTATION_MARKER;
    }
    const paths = jsonLocation.parse(discoveryModelString);
    const locatableProblems = discoveryProblems.filter(p => p.path
      && (paths[p.path] || paths[p.parent]));

    if (locatableProblems.length === 0) {
      return BLANK_ANNOTATION_MARKER;
    }
    return calculateAnnotationsAndMarkers(locatableProblems, paths);
  },
};
