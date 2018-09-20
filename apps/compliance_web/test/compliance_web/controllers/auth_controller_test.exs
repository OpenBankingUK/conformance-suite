defmodule ComplianceWeb.AuthControllerTest do
  @moduledoc """
    Tests authentication.
  """
  use ComplianceWeb.ConnCase
  alias Compliance.Accounts
  alias ComplianceWeb.Guardian
  require Logger

  defp delete_auth(conn) do
    conn
    |> simulate_authenticated
    |> delete("/auth")
  end

  test "authorises user to access report for given validation_run_id", %{conn: conn} do
    conn = get_auth_token(conn)
    token = Guardian.Plug.current_token(conn)
    user = Guardian.Plug.current_resource(conn)
    validation_run_id = "test_validation_run_id"

    refute Guardian.authorised?(token, validation_run_id)

    Accounts.create_user_validation_run(user, validation_run_id)
    assert Guardian.authorised?(token, validation_run_id)
  end

  # test "sets jwt token on successful google authentication", %{conn: conn} do
  #   conn = get_auth_token(conn)
  #   conn =
  #     conn
  #     |> recycle
  #     |> fetch_session
  #   # I can see the session in conn, don't know why is failing
  #   IO.inspect conn
  #   assert get_session(conn, :guardian_default_token)
  # end

  test "delete action removes the access_token from session", %{conn: conn} do
    conn =
      conn
      |> delete_auth()

    assert !get_session(conn, :guardian_default_token)
  end

  test "delete action returns sign out successful json", %{conn: conn} do
    conn =
      conn
      |> delete_auth()

    assert resp = json_response(conn, 200)
    assert resp["success"] == "true"
  end
end
