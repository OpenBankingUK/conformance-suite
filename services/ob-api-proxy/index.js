const { app } = require('./app');
const log = require('debug')('log');

const server = app.listen(process.env.PORT || 8003, '0.0.0.0', () => {
  log('listening address=%O', server.address());
});
