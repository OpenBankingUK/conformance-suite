use Mix.Config

# We don't run a server during test. If one is required,
# you can enable the server option below.
config :compliance_web, ComplianceWeb.Endpoint,
  http: [port: 4001],
  server: false

# see config options: https://github.com/scrogson/oauth2#debug-mode
# this is the library being used behind the scenes by Ueberauth
config :oauth2, debug: true
