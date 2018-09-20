const winston = require('winston');
const _ = require('lodash');

// Config:
// https://github.com/winstonjs/winston/blob/master/examples/custom-timestamp.js
// https://github.com/winstonjs/winston/issues/1135
const level = _.get(process, 'env.LOG_LEVEL', 'info');

exports.logger = winston.createLogger({
  level,
  transports: [
    new winston.transports.Console({
      level,
      format: winston.format.combine(
        winston.format.colorize({ all: true }),
        winston.format.align(),
        winston.format.simple(),
      ),
    }),
  ],
});
