const debug = require('debug')('debug');
const redis = require('redis-mock');
const _ = require('lodash');

exports.store = (() => {
  const client = redis.createClient();
  const EXPIRY_DURATION = 3600; // Default to 1 hour so we don't have too many sessions stored
  const noop = (args) => {
    debug('services/ob-api-proxy/app/session/persistence.js:noop -> args=%j', args);
  };

  const set = (key, value, callback) => {
    const cbk = callback || noop;

    if (!_.isString(key)) {
      throw new Error(`services/ob-api-proxy/app/session/persistence.js:set -> key must be of type String, key=${JSON.stringify(key)}`);
    }

    if (!_.isString(value)) {
      throw new Error(`services/ob-api-proxy/app/session/persistence.js:set -> value must be of type String, value=${JSON.stringify(value)}`);
    }

    debug('services/ob-api-proxy/app/session/persistence.js:set -> key=%j, value=%j, cbk=%O', key, value, cbk);
    return client.set(key, value, 'EX', EXPIRY_DURATION, cbk);
  };

  const get = (key, callback) => {
    const cbk = callback || noop;
    debug('services/ob-api-proxy/app/session/persistence.js:get -> key=%j, cbk=%O', key, cbk);

    if (!key) {
      return cbk(null, null);
    }

    return client.get(key, cbk);
  };

  const remove = (key) => {
    debug('services/ob-api-proxy/app/session/persistence.js:remove -> key=%j', key);

    return client.del(key, noop);
  };

  const getAll = (callback) => {
    const cbk = callback || noop;
    debug('services/ob-api-proxy/app/session/persistence.js:getAll -> cbk=%O', cbk);

    return client.keys('*', cbk);
  };

  const deleteAll = async () => {
    debug('services/ob-api-proxy/app/session/persistence.js:deleteAll');

    return new Promise(resolve => client.flushall(resolve));
  };

  return {
    set,
    get,
    remove,
    getAll,
    deleteAll,
  };
})();
