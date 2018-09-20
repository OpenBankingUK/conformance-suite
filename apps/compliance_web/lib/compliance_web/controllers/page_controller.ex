defmodule ComplianceWeb.PageController do
  use ComplianceWeb, :controller

  require Logger

  def index(conn, _params) do
    render(conn, "index.html", google_client_id: System.get_env("GOOGLE_OAUTH_CLIENT_ID"))
  end
end
