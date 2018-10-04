/*
POST a username and password (form encoded)
if the username and password match - create a session cookie
 */

const { session } = require('./session');
const { credentials } = require('./credentials');
const log = require('debug')('log');

/**
 * This is obviously insanely simplified and no production code should use this method
 */
const checkCredentials = (u, p) => {
  const allow = (u && p && credentials[u] && credentials[u].p === p) || false;
  return allow;
};

exports.login = (() => {
  const authenticate = (req, res) => {
    res.setHeader('Access-Control-Allow-Origin', '*');

    const { u, p } = req.body;
    const allow = checkCredentials(u, p);
    if (allow) {
      session.newSession(u, (sid) => {
        res.setHeader('Content-Type', 'application/json');
        res.status(200).send(JSON.stringify({ sid }));
      });
    } else if (u === 'trigger-error') {
      res.status(500).send();
    } else {
      res.status(401).send();
    }
  };

  const logout = (req, res) => {
    res.setHeader('Access-Control-Allow-Origin', '*');
    const sid = req.headers['authorization'];
    log(`in logout sid is ${sid}`);
    session.destroy(sid, (sidConf) => {
      if (sidConf) {
        log(`destroying sid ${sidConf}`);
        res.setHeader('Content-Type', 'application/json');
        res.status(200).send(JSON.stringify({ sid: sidConf }));
      } else {
        res.sendStatus(204);
      }
    });
  };

  return {
    authenticate,
    logout,
  };
})();
