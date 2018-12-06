
const lowercaseFirstLetter = string => string.charAt(0).toLowerCase() + string.slice(1);

const parseProblem = (problem) => {
  if (problem.indexOf('Key') !== -1 && problem.indexOf('Error') !== -1) {
    const parts = problem.split('Error');

    const path = parts[0]
      .replace(/^Key: ?'?/, '')
      .replace(/'? ?$/, '')
      .replace(/^Model\./, '')
      .replace('API', 'Api')
      .split('.')
      .map(w => lowercaseFirstLetter(w))
      .join('.');

    const error = parts[1].replace(/^:/, '');

    return {
      path,
      error,
    };
  }
  return {
    path: null,
    error: problem,
  };
};

export default {
  getConfig: state => state.main,
  getDiscoveryModel: state => state.discoveryModel,
  problems: state => state.problems,
  discoveryProblems: state => state.problems.map(p => parseProblem(p)),
};
