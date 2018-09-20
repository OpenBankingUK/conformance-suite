# LogConsumer

## Installation

Requires Apache Kafka be installed. For local development on a Mac:

```sh
brew install zoopkeeper
brew install kafka

brew services start kafka
brew services start zoopkeeper

brew services list
# Name      Status  ... Plist
# kafka     started ... ~/Library/LaunchAgents/homebrew.mxcl.kafka.plist
# zookeeper started ... ~/Library/LaunchAgents/homebrew.mxcl.zookeeper.plist
```

You can find a [Kafka quickstart installation guide here](https://kafka.apache.org/quickstart),
for alternate install instructions.

## Creating test topic in Kafka

To create a topic for tests:

```sh
kafka-topics --zookeeper 127.0.0.1:2181 --create --topic kafka-test-topic  --partitions 1 --replication-factor 1
```

Set the following env var:
```
VALIDATION_KAFKA_TOPIC=kafka-test-topic
```

Sometimes we needed to restart kafka after creating a new topic for Node OB proxy code to connect successfully:

```sh
brew services restart kafka
```

To list topics:

```sh
kafka-topics --zookeeper=localhost:2181 --list
```
