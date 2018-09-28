use Mix.Config

if System.get_env("DEBUG_ENVS") == true || System.get_env("DEBUG_ENVS") == "true" do
  IO.inspect(
    System.get_env("DATA_DB_NAME"),
    label: "env[apps/compliance/config/prod.exs] => DATA_DB_NAME"
  )

  IO.inspect(
    System.get_env("DATA_DB_HOST"),
    label: "env[apps/compliance/config/prod.exs] => DATA_DB_HOST"
  )

  IO.inspect(
    System.get_env("OB_API_PROXY_URL"),
    label: "env[apps/compliance/config/prod.exs] => OB_API_PROXY_URL"
  )
end

config :compliance, proxy_url: System.get_env("OB_API_PROXY_URL")

# Configure prod database for remote container.
config :compliance, Compliance.Repo,
  adapter: Mongo.Ecto,
  database: System.get_env("DATA_DB_NAME") || "compilance_prod",
  # DATA_DB_HOST is a Nanobox auto-generated environment variable
  hostname: System.get_env("DATA_DB_HOST") || "localhost"
