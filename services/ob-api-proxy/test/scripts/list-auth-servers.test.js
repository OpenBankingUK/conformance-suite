const assert = require('assert'); // eslint-disable-line
const proxyquire = require('proxyquire'); // eslint-disable-line
const sinon = require('sinon'); //eslint-disable-line

const authorisationServersData = [
  {
    id: 'testId',
    obDirectoryConfig: {
      id: 'testId',
      OBOrganisationId: 'testOrdId',
      CustomerFriendlyName: 'testName',
      OrganisationCommonName: 'testOrg',
      AuthorityId: 'FCA',
      MemberState: 'GB',
      RegistrationId: '123',
    },
    clientCredentials: [{ ex: 'ample' }],
    openIdConfig: { ex: 'ample' },
    registeredConfigs: [{ ex: 'ample' }],
  },
  {
    id: 'testId2',
    obDirectoryConfig: {
      id: 'testId2',
      OBOrganisationId: 'testOrdId',
      CustomerFriendlyName: 'testName2',
      OrganisationCommonName: 'testOrg2',
      AuthorityId: 'FCA',
      MemberState: 'GB',
      RegistrationId: '456',
    },
  },
];
let authServerRows;

const authServerRowsFn = (authServerList) => {
  const allAuthorisationServersStub = sinon.stub().returns(authServerList);
  return proxyquire('../../scripts/list-auth-servers', // eslint-disable-line
    { '../app/authorisation-servers': { allAuthorisationServers: allAuthorisationServersStub } },
  ).authServerRows;
};

describe('authServerRows', () => {
  describe('when no auth servers present', () => {
    beforeEach(() => {
      authServerRows = authServerRowsFn([]);
    });

    it('returns tsv headers', async () => {
      assert.deepEqual(
        await authServerRows(),
        ['id\tCustomerFriendlyName\tOrganisationCommonName\tAuthority\tOBOrganisationId\tclientCredentialsPresent\topenIdConfigPresent\tregisteredConfigsPresent'],
      );
    });
  });

  describe('when auth servers present', () => {
    beforeEach(() => {
      authServerRows = authServerRowsFn(authorisationServersData);
    });

    it('returns tsv of auth servers', async () => {
      const rows = await authServerRows();
      assert.deepEqual(
        rows[0],
        'id\tCustomerFriendlyName\tOrganisationCommonName\tAuthority\tOBOrganisationId\tclientCredentialsPresent\topenIdConfigPresent\tregisteredConfigsPresent',
      );
      assert.deepEqual(
        rows[1],
        'testId\ttestName\ttestOrg\tGB:FCA:123\ttestOrdId\ttrue\ttrue\ttrue',
      );
      assert.deepEqual(
        rows[2],
        'testId2\ttestName2\ttestOrg2\tGB:FCA:456\ttestOrdId\tfalse\tfalse\tfalse',
      );
    });
  });
});
