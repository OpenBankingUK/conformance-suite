import constants from './constants';
import OzoneTemplate from '../../../../../pkg/discovery/templates/ob-v3.0-ozone.json';
import GenericTemplate from '../../../../../pkg/discovery/templates/ob-v3.0-generic.json';
import OzoneTemplateImg from './images/o3logo_159x159.png';
import GenericTemplateImg from './images/obie_logotype_blue_rgb-400×39.jpg';

const templates = [
  {
    model: OzoneTemplate,
    // Fetched from: 'https://o3bank.files.wordpress.com/2017/10/o3logo.png?w=159'
    image: OzoneTemplateImg,
  },
  {
    model: GenericTemplate,
    // Fetched from: 'https://openbanking.atlassian.net/wiki/download/attachments/17236165/OBIE_logotype_blue_RGB.jpg'
    image: GenericTemplateImg,
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

  testCases: [],
  testCaseResults: {},

  errors: {
    configuration: [],
    testCases: [],
    testCaseResults: [],
  },

  wizard: {
    step: constants.WIZARD.STEP_ONE,
  },
};
