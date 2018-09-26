const example = {
  config: {
    authorization_endpoint: 'http://localhost:8001/aaaj4NmBD8lQxmLh2O/authorize',
    client_id: 'spoofClientId',
    client_secret: 'spoofClientSecret',
    fapi_financial_id: 'aaax5nTR33811Qy',
    issuer: 'http://aspsp.example.com',
    redirect_uri: 'http://localhost:8080/tpp/authorized',
    resource_endpoint: 'http://localhost:8001',
    signing_key: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----',
    signing_kid: 'XXXXXX-XXXXxxxXxXXXxxx_xxxx',
    token_endpoint: 'http://localhost:8001/aaaj4NmBD8lQxmLh2O/token',
    token_endpoint_auth_method: 'client_secret_basic',
    transport_cert: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----',
    transport_key: '-----BEGIN PRIVATE KEY----------END PRIVATE KEY-----',
  },
  payload: [
    {
      api_version: '1.1',
      name: 'Sam Morse',
      sort_code: '111111',
      account_number: '12345678',
      amount: '10.00',
      type: 'payments',
    },
    {
      api_version: '1.1',
      name: 'Michael Burnham',
      sort_code: '222222',
      account_number: '87654321',
      amount: '200.00',
      type: 'payments',
    },
    {
      api_version: '2.0',
      type: 'accounts',
    },
  ],
};

export default {
  main: example.config,
  payload: example.payload,
};
