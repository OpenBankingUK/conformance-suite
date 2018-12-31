import OzoneTemplate from '../../../../../pkg/discovery/templates/ob-v3.0-ozone.json';
import GenericTemplate from '../../../../../pkg/discovery/templates/ob-v3.0-generic.json';

const templates = [
  {
    model: OzoneTemplate,
    image: 'https://o3bank.files.wordpress.com/2017/10/o3logo.png?w=159',
  },
  {
    model: GenericTemplate,
    image: 'https://openbanking.atlassian.net/wiki/download/attachments/17236165/OBIE_logotype_blue_RGB.jpg',
  },
];
const defaultTemplate = templates.find(t => t.model.discoveryModel.name === 'ob-v3.0-ozone');

const example = {
  config: {
    accountAccessToken: 'access-token',
    certificateSigning: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----',
    certificateTransport: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----',
    clientScopes: 'AuthoritiesReadAccess ASPSPReadAccess TPPReadAccess',
    keyId: 'key-id',
    privateKeySigning: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----',
    privateKeyTransport: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----',
    softwareStatementId: 'software-statement-id',
    targetHost: 'https://resourceserver.example.com/',
  },
};

export default {
  main: example.config,
  discoveryTemplates: templates,
  discoveryModel: JSON.parse(JSON.stringify(defaultTemplate.model)), // JSON parse to make copy of template model
  problems: null,
  configuration: {
    signing_private: '',
    signing_public: '',
    transport_private: '',
    transport_public: '',
  },
  errors: {
    configuration: [],
  },
};
