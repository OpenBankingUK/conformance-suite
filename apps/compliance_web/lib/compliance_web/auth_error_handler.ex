defmodule ComplianceWeb.AuthErrorHandler do
  @moduledoc """
  Handles authentication errors.
  """
  import Plug.Conn
  require Logger

  def auth_error(conn, {type, reason}, opts) do
    Logger.error(
      "AuthErrorHandler.auth_error -> type: #{inspect(type)}, reason: #{inspect(reason)}, opts: #{
        inspect(opts)
      }"
    )

    body = Poison.encode!(%{message: to_string(type)})

    conn
    |> send_resp(401, body)
    |> halt()
  end
end
