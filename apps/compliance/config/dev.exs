use Mix.Config

config :compliance, proxy_url: System.get_env("OB_API_PROXY_URL") || "http://localhost:8003"

# Configure dev database for container and local use.
config :compliance, Compliance.Repo,
  adapter: Mongo.Ecto,
  database: System.get_env("DATA_DB_NAME") || "compliance_dev",
  # DATA_DB_HOST is a Nanobox auto-generated environment variable
  hostname: System.get_env("DATA_DB_HOST") || "localhost"
