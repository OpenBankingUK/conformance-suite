const monk = require('monk');
const error = require('debug')('error');
const log = require('debug')('log');

let mongodbUri = 'localhost:27017/compliance_dev';
if (process.env.MONGODB_URI) {
  mongodbUri = process.env.MONGODB_URI.replace('mongodb://', '');
}
log(`MONGODB_URI: ${mongodbUri}`);
const db = monk(mongodbUri);

/**
 * Get document with given `id` field from collection.
 * @param {string} collection - name of collection.
 * @param {string} id - `id` field to query.
 * @return {object} document with given `id` field, or null if none found.
 */
const get = async (collection, id, fields = []) => {
  try {
    const store = await db.get(collection);
    return await store.findOne({ id }, ['-_id'].concat(fields));
  } catch (e) {
    error(`error in storage get: ${e.stack}`);
    throw e;
  }
};

/**
 * Get all documents from a collection.
 * @param {string} collection - name of collection.
 * @return {array} array of documents in collection, or empty array if none found.
 */
const getAll = async (collection) => {
  try {
    const store = await db.get(collection);
    return await store.find({}, ['-_id']);
  } catch (e) {
    error(`error in storage getAll: ${e.stack}`);
    throw e;
  }
};

/**
 * Get document with given `id` field from collection.
 * @param {string} collection - name of collection.
 * @param {object} object - document to store.
 * @param {string} id - supply `id` field to set on document.
 */
const set = async (collection, object, id) => {
  try {
    const item = Object.assign({ id }, object);
    const store = await db.get(collection);
    await store.createIndex('id'); // ensure id index exists
    if (await get(collection, id)) {
      return await store.findOneAndUpdate({ id }, item);
    }
    return await store.insert(item);
  } catch (e) {
    error(`error in storage set: ${e.stack}`);
    throw e;
  }
};

/**
 * Drop collection from store.
 * @param {string} collection - name of collection.
 */
const drop = async (collection) => {
  try {
    const store = await db.get(collection);
    return await store.drop();
  } catch (e) {
    error(e);
    throw e;
  }
};

/**
 * Remove document with given id from collection
 * @param collection
 * @param id
 * @returns {Promise.<*>}
 */
const remove = async (collection, id) => {
  try {
    const store = await db.get(collection);
    return await store.remove({ id });
  } catch (e) {
    error(e);
    throw e;
  }
};

const close = async () => {
  try {
    await db.close();
  } catch (e) {
    error(e);
  }
};

module.exports = {
  get,
  getAll,
  set,
  drop,
  remove,
  close,
};
