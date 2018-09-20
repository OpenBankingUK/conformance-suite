defmodule ComplianceWeb.RunConfigController do
  use ComplianceWeb, :controller

  require Logger

  alias Compliance.Configs

  def create(conn, %{
        "openid_config_endpoint" => openid_config_endpoint
      }) do
    run_config = Configs.from_openid_config(openid_config_endpoint)
    map = Map.from_struct(run_config)
    json(conn, map)
  end
end
