defmodule ComplianceWeb.AuthorisationServersController do
  use ComplianceWeb, :controller

  require Logger
  alias Compliance.AuthServers

  def get(conn, _params) do
    case AuthServers.get_all() do
      {:ok, results} ->
        conn
        |> put_status(:ok)
        |> json(%{results: results})

      {:error, reason} ->
        conn
        |> put_status(500)
        |> json(Poison.encode!(%{"error" => reason}))
    end
  end
end
