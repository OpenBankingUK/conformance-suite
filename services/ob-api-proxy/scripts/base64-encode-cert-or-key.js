const fs = require('fs');
const debug = require('debug')('debug');
const log = require('debug')('log');
const error = require('debug')('error');

const encoder = async (certFile) => {
  if (!certFile) {
    throw new Error('Please include a path to a CERT or KEY file,\n<<e.g. npm run base64-cert full/path/to/file>>');
  }

  log('Running encoder');
  const cert = fs.readFileSync(certFile);
  const result = Buffer.from(cert).toString('base64');
  debug(`Encoded CERT: ${result}`);
  return result;
};

encoder(process.argv.slice(2)[0]).then((encodedCert) => {
  if (process.env.NODE_ENV !== 'test') {
    console.log('\nBASE64 ENCODING COMPLETE (Please copy the text below to the required ENV):\n'); // eslint-disable-line
    console.log(encodedCert); // eslint-disable-line
    process.exit();
  }
}).catch((err) => {
  if (process.env.NODE_ENV !== 'test') {
    console.log(err.message); // eslint-disable-line
    error(err);
    process.exit();
  }
});

exports.base64EncodeCertOrKey = encoder;
