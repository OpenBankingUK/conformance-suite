defmodule LogConsumer.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application
  # alias KafkaEx.ConsumerGroup.PartitionAssignment
  alias LogConsumer.KafkaConsumer

  def start(_type, _args) do
    import Supervisor.Spec

    consumer_group_opts = [
      # commit relatively often to make demonstration easy
      commit_interval: 1_000,
      # same with a relatively quick heartbeat rate
      heartbeat_interval: 1_000,
      # name for process registration
      name: LogConsumerGroup,
      # name for the Manager process (for convenience)
      gen_server_opts: [name: LogConsumerGroup.Manager],
      # override the partition assignment callback (optional, see below)
      # partition_assignment_callback: &assign_partitions/2,
      # how long before Kafka considers a consumer gone
      # must be >= group.min.session.timeout.ms from broker config
      session_timeout: 6_000
    ]

    # List all child processes to be supervised
    # children = [
    #   # Starts a worker by calling: LogConsumer.Worker.start_link(arg)
    #   {KafkaEx.ConsumerGroup,
    #    [
    #      LogConsumer.KafkaConsumer,
    #      "log_consumer_group",
    #      ["a-kafka-topic"],
    #      consumer_group_opts
    #    ]}
    # ]

    children = [
      worker(LogConsumer.Stages.Producer, []),
      worker(LogConsumer.Stages.Consumer, []),
      supervisor(KafkaEx.ConsumerGroup, [
        KafkaConsumer,
        "kafka_ex",
        [System.get_env("VALIDATION_KAFKA_TOPIC")],
        consumer_group_opts
      ])
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: LogConsumer.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
