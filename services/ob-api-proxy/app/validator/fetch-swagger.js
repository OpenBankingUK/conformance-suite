const fs = require('fs');
const superagent = require('superagent');

const error = require('debug')('error');
const debug = require('debug')('debug');

exports.fetchSwagger = async (swaggerPath, fileName) => {
  const swagger = swaggerPath;
  if (swagger.startsWith('https')) {
    const file = `./${fileName}`;
    if (fs.existsSync(file)) {
      debug(`Swagger file exists: ${file} - so not getting: ${swagger}`);
      return file;
    }

    try {
      debug(`http get: ${swagger}`);
      fs.writeFileSync(file, '');
      const response = await superagent
        .get(swagger)
        .set({ Accept: 'application/json' });
      debug(`response: ${response.status}`);
      if (response.status === 200) {
        const swaggerJson = response.text;
        fs.appendFileSync(file, swaggerJson);
        return file;
      }
      const msg = `Swagger file ${swagger} not retrieved: ${response.status}`;
      error(msg);
      throw new Error(msg);
    } catch (e) {
      const msg = `Swagger file ${swagger} not retrieved:`;
      error(msg);
      error(e);
      throw e;
    }
  } else if (swagger.endsWith('.json') && fs.existsSync(swagger)) {
    debug(`Swagger FILE Found ${swagger}`);
    return swagger;
  } else {
    const err = `Swagger file ${swagger} in JSON format does not exist`;
    error(err);
    throw new Error(err);
  }
};
