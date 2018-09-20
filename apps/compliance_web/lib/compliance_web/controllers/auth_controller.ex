defmodule ComplianceWeb.AuthController do
  @moduledoc """
  Auth controller responsible for handling Ueberauth responses
  """

  use ComplianceWeb, :controller

  require Logger
  alias Compliance.Accounts

  def new(conn, %{"id_token" => id_token}) do
    case Accounts.create_user_from_id_token(id_token) do
      {:ok, user} ->
        auth_conn = ComplianceWeb.Guardian.Plug.sign_in(conn, user)
        jwt = ComplianceWeb.Guardian.Plug.current_token(auth_conn)

        profile = %{
          first_name: user.first_name,
          last_name: user.last_name,
          email: user.email,
          access_token: jwt
        }

        auth_conn
        |> json(%{profile: profile})

      {:error, reason} ->
        Logger.error("AuthController.new -> reason: #{inspect(reason)}")

        conn
        |> put_status(403)
        |> json(%{error: reason})
    end
  end

  def delete(conn, _params) do
    conn
    |> ComplianceWeb.Guardian.Plug.sign_out()
    |> json(%{success: "true"})
  end

  def tokeninfo(conn, _params) do
    if System.get_env("GOOGLE_OAUTH_TOKENINFO_URL") == "http://localhost:4000/tokeninfo?id_token=" do
      conn
      |> json(%{
        sub: "token",
        given_name: "Test",
        family_name: "User",
        email: "test@email.com",
        aud: System.get_env("GOOGLE_OAUTH_CLIENT_ID")
      })
    else
      conn
      |> put_status(403)
      |> json(%{})
    end
  end
end
