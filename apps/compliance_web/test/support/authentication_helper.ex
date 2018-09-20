defmodule ComplianceWeb.AuthenticationHelper do
  @moduledoc """
  Test helper for mocking authenticated connections.
  """

  defmacro __using__(_) do
    quote do
      import Mock

      def get_auth_token(conn) do
        google_client_id = "GOOGLE_CLIENT_ID"
        System.put_env("GOOGLE_OAUTH_CLIENT_ID", google_client_id)
        id_token = "TOKEN"

        google_tokeninfo_url =
          "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=#{id_token}"

        body = %{
          sub: "token",
          aud: google_client_id,
          email: "test@email.com",
          given_name: "Test",
          family_name: "User"
        }

        resp = {:ok, %HTTPoison.Response{status_code: 200, body: Poison.encode!(body)}}

        with_mock HTTPoison, get: fn ^google_tokeninfo_url -> resp end do
          conn
          |> post("/auth", %{id_token: id_token})
        end
      end

      defp simulate_authenticated(conn) do
        auth = get_auth_token(conn)
        profile = Poison.decode!(auth.resp_body)["profile"]

        conn
        |> put_req_header("authorization", "Bearer #{profile["access_token"]}")
      end
    end
  end
end
