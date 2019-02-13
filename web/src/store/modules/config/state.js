
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
    client_id: '',
    client_secret: '',
    token_endpoint: '',
    x_fapi_financial_id: '',
    redirect_url: 'https://0.0.0.0:8443/conformancesuite/callback',
  },

  wizard: {
    step: constants.WIZARD.STEP_ONE,
  },
};
