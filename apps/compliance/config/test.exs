use Mix.Config

# Configure test database for CI and local use.
config :compliance, Compliance.Repo,
  adapter: Mongo.Ecto,
  database: System.get_env("DATA_DB_NAME") || "compliance_test",
  # DATA_DB_HOST is a Nanobox auto-generated environment variable
  hostname: System.get_env("DATA_DB_HOST") || "localhost"
