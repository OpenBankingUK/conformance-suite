defmodule ComplianceWeb.PageControllerTest do
  use ComplianceWeb.ConnCase

  test "injects up GOOGLE_OAUTH_CLIENT_ID on index page", %{conn: conn} do
    google_client_id = "Some google id"
    System.put_env("GOOGLE_OAUTH_CLIENT_ID", google_client_id)
    conn = get(conn, "/")
    assert response(conn, 200)

    assert String.contains?(
             conn.resp_body,
             "<meta name=\"google-signin-client_id\" content=\"#{google_client_id}\">"
           )
  end
end
