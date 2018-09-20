defmodule ComplianceWeb.UserControllerTest do
  @moduledoc false

  use ComplianceWeb.ConnCase
  require Logger

  @tag capture_log: true
  test "requires user authentication on show action", %{conn: conn} do
    conn = get(conn, "/user")
    assert response(conn, 401)
    assert conn.halted
  end

  test "show action returns user when user in session", %{conn: conn} do
    conn =
      conn
      |> simulate_authenticated
      |> get("/user")

    assert user = json_response(conn, 200)
    assert user["user"]
  end
end
