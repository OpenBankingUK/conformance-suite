import * as _ from 'lodash';
import constants from './constants';

// Converts problem key to discovery model JSON path.
const parseProblem = ({ key, error }) => {
  if (key && error) {
    const parts = key
      .replace('API', 'Api')
      .replace('URL', 'Url')
      .split('.')
      .map(w => _.lowerFirst(w));

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
  discoveryModel: state => state.discoveryModel,
  discoveryModelString: state => JSON.stringify(state.discoveryModel, null, 2),
  discoveryTemplates: state => state.discoveryTemplates,
  tokenAcquisition: state => (state.discoveryModel ? state.discoveryModel.discoveryModel.tokenAcquisition : null),
  problems: state => state.problems,
  discoveryProblems: state => (state.problems ? state.problems.map(p => parseProblem(p)) : null),
  configuration: state => state.configuration,
  configurationString: state => JSON.stringify(state.configuration, null, 2),
  /**
   * Computes what the user can navigate to based on the current step they are on.
   */
  navigation: (state) => {
    const { step } = state.wizard;
    const navigation = {
      '/wizard/continue-or-start': step > 0,
      '/wizard/discovery-config': step > constants.WIZARD.STEP_ONE,
      '/wizard/configuration': step > constants.WIZARD.STEP_TWO,
      '/wizard/overview-run': step > constants.WIZARD.STEP_THREE,
      '/wizard/export': step > constants.WIZARD.STEP_FOUR,
    };
    return navigation;
  },
};
