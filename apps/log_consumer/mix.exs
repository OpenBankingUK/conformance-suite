defmodule LogConsumer.MixProject do
  @moduledoc false

  use Mix.Project

  def project do
    [
      app: :log_consumer,
      version: "0.1.0",
      build_path: "../../_build",
      config_path: "../../config/config.exs",
      deps_path: "../../deps",
      lockfile: "../../mix.lock",
      elixir: "~> 1.7",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  def application do
    [
      extra_applications: [:logger],
      mod: {LogConsumer.Application, []}
    ]
  end

  defp deps do
    [
      {:compliance, in_umbrella: true},
      {:gen_stage, "~> 0.13.1"},
      {:kafka_ex, "~> 0.8.2"},
      {:mix_test_watch, "~> 0.9", only: :dev},
      {:mock, "~> 0.3.2", only: :test}
    ]
  end
end
