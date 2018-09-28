const express = require('express');
const { replayMiddleware } = require('./replay-middleware');
const { validationErrorMiddleware } = require('./validation-error-middleware');
const { swaggerMiddleware } = require('./swagger-middleware');
const { logger } = require('../utils/logger');

// validator key -> express app validator
const validators = new Map();

const validateResponseOn = () => process.env.VALIDATE_RESPONSE === 'true';

const addValidationMiddleware = async (app, swaggerUriOrFile, swaggerFile) => {
  const { metadata, validator } = await swaggerMiddleware(swaggerUriOrFile, swaggerFile);
  app.use(metadata);
  app.use(validator);
  app.use(replayMiddleware);
  app.use(validationErrorMiddleware);
};

const configureSwagger = async (swaggerUris, app) => // eslint-ignore-line
  Promise.all(swaggerUris.map(async (swaggerUri) => {
    const swaggerFile = swaggerUri.split('/').slice(-2).join('-');

    logger.log('verbose', 'configureSwagger', { swaggerUri, swaggerFile });
    return addValidationMiddleware(app, swaggerUri, swaggerFile);
  }));

const initValidatorApp = async ({ swaggerUris }) => {
  const app = express();
  app.disable('x-powered-by');
  await configureSwagger(swaggerUris, app);

  return app;
};

const makeValidatorKey = ({ swaggerUris = [], scope = '' }) => {
  const keys = swaggerUris.slice().concat([scope]);
  return keys.sort().toString();
};

const validatorApp = async (details) => {
  if (!validateResponseOn()) {
    return undefined;
  }

  const validatorKey = makeValidatorKey(details);
  const createValidator = !validators.has(validatorKey) || !validators.get(validatorKey).default;
  logger.log('verbose', 'validatorApp', {
    details, validatorKey, createValidator, validators: JSON.stringify([...validators]),
  });

  if (createValidator) {
    const app = await initValidatorApp(details);
    validators.set(validatorKey, { default: app });
  }

  const validator = validators.get(validatorKey);
  return validator.default;
};

module.exports = {
  initValidatorApp,
  validateResponseOn,
  validatorApp,
};
