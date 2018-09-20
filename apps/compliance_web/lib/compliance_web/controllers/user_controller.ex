defmodule ComplianceWeb.UserController do
  use ComplianceWeb, :controller
  require Logger

  def show(conn, _) do
    resource = ComplianceWeb.Guardian.Plug.current_resource(conn)

    user = %{
      first_name: resource.first_name,
      last_name: resource.last_name,
      email: resource.email
    }

    json(conn, %{user: user})
  end
end
