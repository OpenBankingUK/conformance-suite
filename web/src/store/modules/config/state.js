import constants from './constants';

const templates = [];

export default {
  discoveryTemplates: templates,
  discoveryModel: null,
  problems: null,

  configuration: {
    signing_private: '',
    signing_public: '',
    transport_private: '',
    transport_public: '',
  },

  wizard: {
    step: constants.WIZARD.STEP_ONE,
  },
};
