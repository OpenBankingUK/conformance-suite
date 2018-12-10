import DiscoveryExample from '../../../../../pkg/discovery/templates/ob-v3.0-ozone.json';

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
  discoveryModel: DiscoveryExample,
  problems: null,
};
