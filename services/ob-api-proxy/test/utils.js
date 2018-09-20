const assert = require('assert');

const checkErrorThrown = async (testFn, status, message) => {
  try {
    await testFn();
    assert.ok(false);
  } catch (err) {
    if (err.code && err.code === 'ERR_ASSERTION') throw err;
    assert.equal(err.message, message);
    assert.equal(err.status, status);
  }
};

exports.checkErrorThrown = checkErrorThrown;
