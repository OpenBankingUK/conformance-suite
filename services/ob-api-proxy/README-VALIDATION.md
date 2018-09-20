# Validation against Swagger specifications

## Validate ASPSP API responses

You can configure validation of ASPSP API responses against swagger specifications.
This is optional.

To turn on validation, set the following env vars:

```sh
VALIDATE_RESPONSE=true
```

Set `x-swagger-uris` header on request, containing list of space separated
swagger file URIs that you wish to validate against.

```
'x-swagger-uris': "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/ee715e094a59b37aeec46aef278f528f5d89eb03/dist/v1.1/account-info-swagger.json"
```

# Logging to Kafka

## Publish validation results to Kafka

You can configure the swagger validation results to be published to
[Kafka](https://kafka.apache.org). You may want to do this if you have
written a consumer for processing the validation results. This is optional.

To configure logging to Kafka, set the following env vars:

```sh
VALIDATION_KAFKA_TOPIC=kafka-test-topic
VALIDATION_KAFKA_BROKER=127.0.0.1:9092
```

## Install local Kakfa install

On Mac OSX you can install Kakfa using homebrew:

```sh
brew install zookeeper kafka
```

To check services running:

```sh
brew services list
```

You can find a [Kafka quickstart installation guide here](https://kafka.apache.org/quickstart),
for alternate install instructions.

We expect to add steps for setting up a docker container for running Kafka soon.

## Creating test topic in Kafka

To create a topic for tests:

```sh
kafka-topics --zookeeper 127.0.0.1:2181 --create --topic kafka-test-topic  --partitions 1 --replication-factor 1
```

Sometimes we needed to restart kafka after creating a new topic for Node code to connect successfully:

```sh
brew services restart kafka
```

To list topics:

```sh
kafka-topics --zookeeper=localhost:2181 --list
```
