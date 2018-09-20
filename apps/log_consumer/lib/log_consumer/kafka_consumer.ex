defmodule LogConsumer.KafkaConsumer do
  @moduledoc """
  Consumes test log messages from Kafka.

  A single consumer process consumes from a
  single partition of a Kafka topic.

  KafkaEx.GenConsumer default implementation
  takes care of determining a starting offset,
  fetching messages from a Kafka broker, and
  committing offsets for consumed messages.
  """

  use KafkaEx.GenConsumer

  # alias KafkaEx.Protocol.Fetch.Message

  require Logger

  # note - messages are delivered in batches
  def handle_message_set(message_set, state) do
    Logger.debug(fn -> "received message set: #{inspect(message_set)}" end)
    LogConsumer.Stages.Producer.notify(message_set)

    {:async_commit, state}
  end
end
