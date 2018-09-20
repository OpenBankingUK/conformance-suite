# This file is responsible for configuring your application
# and its dependencies with the aid of the Mix.Config module.
#
# This configuration file is loaded before any dependency and
# is restricted to this project.
use Mix.Config

# General application configuration
config :compliance_web,
  namespace: ComplianceWeb,
  ecto_repos: [Compliance.Repo]

# Configures the endpoint
config :compliance_web, ComplianceWeb.Endpoint,
  url: [host: "localhost"],
  secret_key_base: "rOoQVTyg8RcA6RmLEhX0Fs86GYUvA3ufMCBcunLoiTCA0MNqVBRssHOJpfskMe+9",
  render_errors: [view: ComplianceWeb.ErrorView, accepts: ~w(html json)],
  pubsub: [name: ComplianceWeb.PubSub, adapter: Phoenix.PubSub.PG2]

# Configures Elixir's Logger
config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:module, :function, :request_id]

config :compliance_web, :generators, context_app: :compliance

# see config options: https://github.com/scrogson/oauth2#debug-mode
# this is the library being used behind the scenes by Ueberauth
config :oauth2, debug: true

# Configure Ueberauth authentication
config :ueberauth, Ueberauth,
  providers: [
    google: {Ueberauth.Strategy.Google, [default_scope: "emails profile plus.me"]}
  ]

# Configure Google OAuth
config :ueberauth, Ueberauth.Strategy.Google.OAuth,
  client_id: "GOOGLE_CLIENT_ID",
  client_secret: "GOOGLE_CLIENT_SECRET"

# Configure Guardian
config :compliance_web, ComplianceWeb.Guardian,
  issuer: "compliance_web",
  secret_key: "GUARDIAN_SECRET_KEY"

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{Mix.env()}.exs"
