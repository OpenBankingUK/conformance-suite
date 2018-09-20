use Mix.Config

config :kafka_ex,
  brokers: [
    {System.get_env("KAFKA_HOST") || "localhost", String.to_integer(System.get_env("KAFKA_PORT") || "9092")}
    # {"localhost", 9092}
    # {"localhost", 9093},
    # {"localhost", 9094}
  ],
  consumer_group: "kafka_ex",
  disable_default_worker: false,
  sync_timeout: 3000,
  max_restarts: 10,
  max_seconds: 60,
  use_ssl: false,
  ssl_options: [],
  # ssl_options: [
  #   cacertfile: System.cwd() <> "/ssl/ca-cert",
  #   certfile: System.cwd() <> "/ssl/cert.pem",
  #   keyfile: System.cwd() <> "/ssl/key.pem"
  # ],
  kafka_version: "0.9.0"
  # kafka_version: "0.9.0.1"
  # kafka_version: "1.1.0"

config :logger,
  level: :debug,
  truncate: :infinity
