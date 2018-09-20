const express = require('express');
const { replayMiddleware } = require('./replay-middleware');
const { validationErrorMiddleware } = require('./validation-error-middleware');
const { swaggerMiddleware } = require('./swagger-middleware');
const { KafkaStream } = require('./kafka-stream');
const { logger } = require('../utils/logger');

// validator key -> express app validator
const validators = new Map();
let _kafkaStream; // eslint-disable-line

const validateResponseOn = () => process.env.VALIDATE_RESPONSE === 'true';
const logTopic = () => process.env.VALIDATION_KAFKA_TOPIC;
const connectionString = () => process.env.VALIDATION_KAFKA_BROKER;

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

const kakfaConfigured = () =>
  !!(logTopic() && connectionString());

const initKafkaStream = async () => {
  const kafkaStream = new KafkaStream({
    kafkaOpts: {
      connectionString: connectionString(),
    },
    topic: logTopic(),
  });
  await kafkaStream.init();

  return kafkaStream;
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

const kafkaStream = async () => {
  if (!(validateResponseOn() && kakfaConfigured())) {
    return undefined;
  }
  if (!_kafkaStream) {
    _kafkaStream = await initKafkaStream();
  }
  return _kafkaStream;
};

module.exports = {
  initValidatorApp,
  validateResponseOn,
  validatorApp,
  kafkaStream,
  kakfaConfigured,
};
