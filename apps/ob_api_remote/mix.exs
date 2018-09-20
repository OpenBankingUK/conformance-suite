defmodule OBApiRemote.MixProject do
  @moduledoc false
  use Mix.Project

  def project do
    [
      app: :ob_api_remote,
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

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger],
      mod: {OBApiRemote.Application, []}
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:poison, "~> 3.1.0"},
      {:httpoison, "~> 1.0"},
      {:mix_test_watch, "~> 0.9", only: :dev},
      {:mock, "~> 0.3.2", only: :test}
    ]
  end
end
