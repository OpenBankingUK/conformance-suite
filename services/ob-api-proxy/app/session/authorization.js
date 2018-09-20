const { session } = require('./session');

const validSession = (candidate, callback) => {
  session.getData(candidate, (err, data) => callback(!!data));
};

exports.requireAuthorization = (req, res, next) => {
  const sid = req.headers.authorization;
  if (sid) {
    validSession(sid, (valid) => {
      if (!valid) {
        res.setHeader('Access-Control-Allow-Origin', '*');
        res.status(401).send();
      } else {
        next();
      }
    });
  } else {
    res.setHeader('Access-Control-Allow-Origin', '*');
    res.status(401).send();
  }
};
