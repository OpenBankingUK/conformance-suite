const assert = require('assert');
const Kafka = require('no-kafka');
const test = require('debug')('test');
const util = require('util');

const { KafkaStream, RoundRobinPartitioner } = require('../../app/validator/kafka-stream');

const logTopic = 'kafka-test-topic';
const connectionString = '127.0.0.1:9092';

/**
  Create topic for test with:
  kafka-topics --zookeeper 127.0.0.1:2181 --create --topic kafka-test-topic
               --partitions 1 --replication-factor 1
*/

describe('KafkaStream', () => {
  it('fails to construct without kafkaOpts', () => {
    try {
      new KafkaStream({ topic: logTopic }); // eslint-disable-line
      assert.fail('expected construction to fail');
    } catch (err) {
      assert(err);
      assert.equal(err.message, 'Kafka options must be provided');
    }
  });

  it('fails to construct without topic', () => {
    try {
      new KafkaStream({ // eslint-disable-line
        kafkaOpts: {
          connectionString,
        },
      });
      assert.fail('expected construction to fail');
    } catch (err) {
      assert(err);
      assert.equal(err.message, 'KafkaStream must have a topic');
    }
  });

  it('gets partition in round-robin when there is more than 1 partition', () => {
    const partitions = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
    let curPartition = 0;
    const partitioner = new RoundRobinPartitioner();
    for (let i = 0; i < 100; i++) { // eslint-disable-line
      assert.equal(partitioner.partition('', partitions), curPartition);
      curPartition++; // eslint-disable-line
      if (curPartition >= partitions.length) {
        curPartition = 0;
      }
    }
  });

  // const partitions = [0, 1, 2, 3];
  const partitions = [0];

  const setupConsumer = () => new Kafka.SimpleConsumer({
    kafkaOpts: {
      connectionString,
    },
  });

  it('gets producer error properly; requires Kafka be running', function (done) { // eslint-disable-line
    this.timeout(5000);
    const kafkaConsumer = setupConsumer();
    kafkaConsumer.init().then(() => {
      kafkaConsumer.subscribe(logTopic, partitions, (messageSet, topic, partition) => { // eslint-disable-line
        assert.fail('did not expect to receive messageSet');
      });
    }).then(() => {
      const kafkaStream = new KafkaStream({
        topic: 'bad-kafka-test-topic',
        kafkaOpts: {
          connectionString,
        },
      });
      kafkaStream.init().then((kStream) => {
        kStream.write({ msg: 'test error' }).then(
          () => assert.fail('expected error to be thrown'),
          (err) => {
            assert(err);
            test(JSON.stringify(err));
            assert.equal(err.message, 'This request is for a topic or partition that does not exist on this broker.');
            test('passed');
            kafkaConsumer.unsubscribe(logTopic, partitions)
              .then(() => kafkaConsumer.end())
              .then(() => done());
          },
        );
      });
    });
  });

  it('writes to stream properly; requires Kafka be running', (done) => {
    const kafkaConsumer = setupConsumer();

    kafkaConsumer.init().then(() => {
      kafkaConsumer.subscribe(logTopic, partitions, (messageSet, topic, partition) => {
        messageSet.forEach((m) => {
          kafkaConsumer.commitOffset({ topic, partition, offset: m.offset });
          const logMsg = JSON.parse(m.message.value);
          test(util.inspect(logMsg));
          assert.equal(logMsg.msg, 'test success');
          test('passed');
        });
        kafkaConsumer.unsubscribe(logTopic, partitions)
          .then(() => done());
      });
    }).then(() => {
      const kafkaStream = new KafkaStream({
        topic: logTopic,
        kafkaOpts: { connectionString },
      });

      kafkaStream.init().then((kStream) => {
        try {
          kStream.write({ msg: 'test success' });
        } catch (err) {
          assert.fail(err.message);
        }
      });
    });
  });
});
