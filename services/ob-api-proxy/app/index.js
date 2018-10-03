if (!process.env.DEBUG) process.env.DEBUG = 'error,log';

const express = require('express');
const morgan = require('morgan');
const cors = require('cors');
const bodyParser = require('body-parser');
const { requireAuthorization } = require('./session');
const { login } = require('./session');
const { resourceRequestHandler } = require('./request-data/ob-proxy.js');
const { accountRequestAuthoriseConsent, accountRequestRevokeConsent } = require('./setup-account-request');
const { paymentAuthoriseConsent } = require('./setup-payment');
const { paymentSubmission } = require('./submit-payment');
const { authorisationCodeGrantedHandler } = require('./authorise');

const app = express();

// don't log requests when testing
if (process.env.NODE_ENV !== 'test') {
  // // Log twice once for the request and once for the response.
  // // Immediate means log as soon as the request arrives
  // app.use(morgan('common', {
  //   immediate: true,
  // }));
  app.use(morgan('combined', {
    immediate: false,
  }));
}

const requireAuthorisationServerId = async (req, res, next) => {
  const authServerId = req.headers['x-authorization-server-id'];
  if (!authServerId) {
    return res.status(400).send('request missing x-authorization-server-id header');
  }
  return next();
};

app.options('*', cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

app.use('/login', login.authenticate);
app.use('/logout', login.logout);

app.all('/account-request-authorise-consent', requireAuthorization, requireAuthorisationServerId);
app.post('/account-request-authorise-consent', accountRequestAuthoriseConsent);

app.all('/account-request-revoke-consent', requireAuthorization, requireAuthorisationServerId);
app.post('/account-request-revoke-consent', accountRequestRevokeConsent);

app.all('/payment-authorise-consent', requireAuthorization, requireAuthorisationServerId);
app.post('/payment-authorise-consent', paymentAuthoriseConsent);

app.all('/payment-submissions', requireAuthorization, requireAuthorisationServerId);
app.post('/payment-submissions', paymentSubmission);

app.all('/tpp/authorized', requireAuthorization);
app.post('/tpp/authorized', authorisationCodeGrantedHandler);

app.all(
  '/open-banking/*',
  requireAuthorization,
  requireAuthorisationServerId,
);
app.use('/open-banking', resourceRequestHandler);

exports.app = app;
