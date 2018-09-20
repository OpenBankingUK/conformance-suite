const assert = require('assert');

const { scopeAndUrl } = require('../../app/request-data/ob-proxy');

describe('scopeAndUrl', () => {
  const assertScope = (path, expectedScope) => {
    const { scope } = scopeAndUrl(path, 'mock-host');
    assert.equal(scope, expectedScope);
  };

  it('returns correct host for API path', async () => {
    const { proxiedUrl } = scopeAndUrl('/v1.1/accounts', 'mock-host');
    assert.equal(proxiedUrl, 'mock-host/open-banking/v1.1/accounts');
  });

  it('returns correct scope for accounts API paths', async () => {
    assertScope('/v1.1/accounts', 'accounts');
    assertScope('/v1.1/accounts/123', 'accounts');
    assertScope('/v1.1/balances', 'accounts');
    assertScope('/v1.1/standing-orders', 'accounts');
  });

  it('returns correct scope for payments API paths', async () => {
    assertScope('/v1.1/payments', 'payments');
    assertScope('/v1.1/payment-submissions', 'payments');
  });
});
