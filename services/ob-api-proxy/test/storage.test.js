const {
  set, get, getAll, drop,
} = require('../app/storage.js');
const assert = require('assert');
// const error = require('debug')('error');

describe('storage set data with id', () => {
  const collection = 'testAuthServers';
  const data = {
    BaseApiDNSUri: 'http://bbb.example.com',
    CustomerFriendlyLogoUri: 'string',
    CustomerFriendlyName: 'BBB Example Bank',
  };
  const id = 'testId';

  it('then get with invalid id returns null', async () => {
    await set(collection, data, id);
    const result = await get(collection, 'bad-id');
    return assert.equal(null, result);
  });

  it('then get with id returns same data', async () => {
    await set(collection, data, id);
    const expected = Object.assign({ id }, data);
    const result = await get(collection, id);
    return assert.deepEqual(expected, result);
  });

  it('then getAll returns array containing element with same data', async () => {
    await set(collection, data, id);
    const expected = Object.assign({ id }, data);
    const result = await getAll(collection);
    return assert.deepEqual(expected, result[0]);
  });

  it('called second time overwrites data', async () => {
    const newData = Object.assign({ CustomerFriendlyName: 'New Name' }, data);

    await set(collection, data, id);
    await set(collection, newData, id);

    const expected = Object.assign({ id }, newData);
    const result = await get(collection, id);
    return assert.deepEqual(expected, result);
  });

  after(async () => {
    await drop(collection);
  });
});
