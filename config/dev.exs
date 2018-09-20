use Mix.Config

config :logger, :console,
  format: "$metadata[$level] $message\n",
  metadata: [:module, :function]

# Set a higher stacktrace during development. Avoid configuring such
# in production as building large stacktraces may be expensive.
config :phoenix, :stacktrace_depth, 20
