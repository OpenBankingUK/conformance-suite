use Mix.Config

config :compliance, ecto_repos: [Compliance.Repo]

import_config "#{Mix.env()}.exs"
