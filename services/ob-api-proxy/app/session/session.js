const uuidv1 = require('uuid/v1'); // Timestamp based UUID
const { store } = require('./persistence.js');
const util = require('util');
const log = require('debug')('log');

const session = (() => {
  const setData = (sid, username, callback) => {
    store.set(sid, JSON.stringify({ sid, username }), () => {
      callback(sid);
    });
  };
  const getData = (sid, cb) => store.get(sid, cb);
  const getDataAsync = util.promisify(getData);
  const setAccessToken = accessToken => store.set('ob_directory_access_token', JSON.stringify(accessToken));
  const getAccessToken = cb => store.get('ob_directory_access_token', cb);

  const getUsername = async (sessionId) => {
    const sessionData = JSON.parse(await session.getDataAsync(sessionId));
    return sessionData.username;
  };

  const destroy = (candidate, cb) => {
    const sessHandler = (err, data) => {
      const sid = data && JSON.parse(data).sid;
      log(`in sessHandler sid is ${sid}, candidate:[${candidate}]`);
      if (sid !== candidate) {
        return cb(null);
      }
      store.remove(candidate); // Async but we kinda don't care :-/
      return cb(sid);
    };
    store.get(candidate, sessHandler);
  };

  const newSession = (username, callback) => {
    const mySid = uuidv1();
    setData(mySid, username, callback);
  };

  const deleteAll = async () => {
    await store.deleteAll();
  };

  return {
    setData,
    getData,
    getDataAsync,
    setAccessToken,
    getAccessToken,
    getUsername,
    destroy,
    newSession,
    deleteAll,
  };
})();

module.exports = {
  session,
  getUsername: session.getUsername,
};
