const util = require('util');
const { fetchSwagger } = require('./fetch-swagger');
const { initializeMiddleware } = require('swagger-tools');

initializeMiddleware[util.promisify.custom] = (swaggerObj) => { // eslint-disable-line
  return new Promise((resolve) => {
    initializeMiddleware(swaggerObj, resolve);
  });
};

const initMiddleware = util.promisify(initializeMiddleware);

exports.swaggerMiddleware = async (swaggerUriOrFile, swaggerFile) => {
  const file = await fetchSwagger(swaggerUriOrFile, swaggerFile);
  const swaggerTools = await initMiddleware(require('../../'+file)); // eslint-disable-line
  return {
    metadata: swaggerTools.swaggerMetadata(),
    validator: swaggerTools.swaggerValidator({ validateResponse: true }),
  };
};
