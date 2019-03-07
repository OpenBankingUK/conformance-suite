
import constants from './constants';

const templates = [];

export default {
  discoveryTemplates: templates,
  discoveryModel: null,
  problems: null,
  token_endpoint_auth_methods: [],

  configuration: {
    signing_private: '',
    signing_public: '',
    transport_private: '',
    transport_public: '',
    client_id: '',
    client_secret: '',
    token_endpoint: '',
    token_endpoint_auth_method: 'client_secret_basic',
    authorization_endpoint: '',
    resource_base_url: '',
    x_fapi_financial_id: '',
    issuer: '',
    redirect_url: 'https://0.0.0.0:8443/conformancesuite/callback',
    resource_account_ids: [],
    resource_statement_ids: [],
  },

  wizard: {
    step: constants.WIZARD.STEP_ONE,
  },
};
