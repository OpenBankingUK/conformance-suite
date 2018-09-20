const { app } = require('./app');
const log = require('debug')('log');

if (process.env.TPP_REF_SERVER_PORT) {
  app.listen(process.env.TPP_REF_SERVER_PORT, '0.0.0.0');
  log(` App listening on port ${process.env.TPP_REF_SERVER_PORT}`);
} else {
  const port = process.env.PORT || 8003;
  app.listen(port);
  log(` App listening on port ${port}`);
}
