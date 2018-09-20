use Mix.Config

if System.get_env("DEBUG_ENVS") == true || System.get_env("DEBUG_ENVS") == "true" do
  IO.inspect(
    System.get_env("OB_API_PROXY_URL"),
    label: "env[apps/ob_api_remote/config/prod.exs] => OB_API_PROXY_URL"
  )
end

config :ob_api_remote, proxy_url: System.get_env("OB_API_PROXY_URL")
