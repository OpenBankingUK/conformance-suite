import constants from './constants';
import OzoneTemplate from '../../../../../pkg/discovery/templates/ob-v3.0-ozone.json';
import GenericTemplate from '../../../../../pkg/discovery/templates/ob-v3.0-generic.json';

const templates = [
  {
    model: OzoneTemplate,
  },
  {
    model: GenericTemplate,
  },
];

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
