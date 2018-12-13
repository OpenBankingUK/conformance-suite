
const lowercaseFirstLetter = string => string.charAt(0).toLowerCase() + string.slice(1);

// Converts problem key to discovery model JSON path.
const parseProblem = ({ key, error }) => {
  if (key && error) {
    const parts = key
      .replace('API', 'Api')
      .replace('URL', 'Url')
      .split('.')
      .map(w => lowercaseFirstLetter(w));

    const path = parts.join('.');
    const parent = parts.slice(0, -1).join('.');

    return {
      path,
      parent,
      error,
    };
  }
  return {
    path: null,
    error,
  };
};

export default {
  getConfig: state => state.main,
  getDiscoveryModel: state => state.discoveryModel,
  problems: state => state.problems,
  discoveryProblems: state => (state.problems ? state.problems.map(p => parseProblem(p)) : null),
};
