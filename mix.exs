defmodule Compliance.Umbrella.Mixfile do
  @moduledoc false

  use Mix.Project

  def project do
    [
      apps_path: "apps",
      start_permanent: Mix.env() == :prod,
      aliases: aliases(),
      deps: deps()
    ]
  end

  # Dependencies can be Hex packages:
  #
  #   {:mydep, "~> 0.3.0"}
  #
  # Or git/path repositories:
  #
  #   {:mydep, git: "https://github.com/elixir-lang/mydep.git", tag: "0.1.0"}
  #
  # Type "mix help deps" for more examples and options.
  #
  # Dependencies listed here are available only for this project
  # and cannot be accessed from applications inside the apps folder
  defp deps do
    [
      {:credo, "~> 0.9.2", only: [:dev, :test], runtime: false},
      {:distillery, "~> 1.5.5", runtime: false},
    ]
  end

  defp aliases do
    [
      test: &test_apps/1
    ]
  end

  defp test_apps(opts) do
    opts = Enum.join(opts, " ")
    # Don't start log_consumer app in test, as this requires kafka to be runnning
    Mix.Task.run("cmd", ["--app", "log_consumer", "mix test --no-start --color " <> opts])

    File.ls!("apps")
    |> Enum.reject(&(&1 == "log_consumer"))
    |> Enum.each(&Mix.Task.rerun("cmd", ["--app", &1, "mix test --color " <> opts]))
  end
end
