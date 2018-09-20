defmodule Compliance.Mixfile do
  @moduledoc false
  use Mix.Project

  def project do
    [
      app: :compliance,
      version: "0.0.1",
      build_path: "../../_build",
      config_path: "../../config/config.exs",
      deps_path: "../../deps",
      lockfile: "../../mix.lock",
      elixir: "~> 1.7",
      elixirc_paths: elixirc_paths(Mix.env()),
      start_permanent: Mix.env() == :prod,
      aliases: aliases(),
      deps: deps()
    ]
  end

  # Configuration for the OTP application.
  #
  # Type `mix help compile.app` for more information.
  def application do
    [
      mod: {Compliance.Application, []},
      extra_applications: [:logger, :runtime_tools]
    ]
  end

  # Specifies which paths to compile per environment.
  defp elixirc_paths(:test), do: ["lib", "test/support"]
  defp elixirc_paths(_), do: ["lib"]

  # Specifies your project dependencies.
  #
  # Type `mix help deps` for examples and options.
  defp deps do
    [
      {:data_morph, "~> 0.0.8"},
      {:ecto, "~> 2.1.0"},
      {:httpoison, "~> 1.0"},
      {:mix_test_watch, "~> 0.9", only: :dev},
      {:mock, "~> 0.3.2", only: :test},
      {:mongodb_ecto, "~> 0.2"},
      {:poison, "~> 3.1.0"},
      {:uuid, "~> 1.1"},
      {:ob_api_remote, in_umbrella: true}
    ]
  end

  # Aliases are shortcuts or tasks specific to the current project.
  # For example, to create, migrate and run the seeds file at once:
  #
  #     $ mix ecto.setup
  #
  # See the documentation for `Mix` for more info on aliases.
  defp aliases do
    [
      "ecto.setup": ["ecto.create", "ecto.migrate", "run priv/repo/seeds.exs"],
      "ecto.reset": ["ecto.drop", "ecto.setup"],
      "compliance.permutations": ["run priv/endpoints/generate.exs"],
      test: ["ecto.create --quiet", "ecto.migrate", "test"]
    ]
  end
end
