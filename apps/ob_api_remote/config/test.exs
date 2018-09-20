use Mix.Config

config :ob_api_remote, proxy_url: System.get_env("OB_API_PROXY_URL") || "http://localhost:8003"
