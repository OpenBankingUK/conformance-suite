defmodule ComplianceWeb.AuthorisationServersControllerTest do
  use ComplianceWeb.ConnCase
  use ExUnit.Case, async: false

  import Mock

  alias Compliance.AuthServers

  @authorisation_servers [
    %{
      "id" => "aaaj4NmBD8lQxmLh2O",
      "logoUri" => "",
      "name" => "AAA Example Bank"
    },
    %{
      "id" => "bbbX7tUB4fPIYB0k1m",
      "logoUri" => "",
      "name" => "BBB Example Bank"
    },
    %{
      "id" => "cccbN8iAsMh74sOXhk",
      "logoUri" => "",
      "name" => "CCC Example Bank"
    }
  ]

  describe "/account-payment-service-provider-authorisation-servers" do
    test "returns list of APSPS", %{
      conn: conn
    } do
      with_mock AuthServers, get_all: fn -> {:ok, @authorisation_servers} end do
        conn =
          conn
          |> put_req_header("content-type", "application/json")
          |> get("/account-payment-service-provider-authorisation-servers")

        body = json_response(conn, :ok)
        assert body == %{"results" => @authorisation_servers}
      end
    end

    test "handles errors", %{
      conn: conn
    } do
      with_mock AuthServers, get_all: fn -> {:error, "ERROR"} end do
        conn =
          conn
          |> put_req_header("content-type", "application/json")
          |> get("/account-payment-service-provider-authorisation-servers")

        assert json_response(conn, 500) =~ "ERROR"
      end
    end
  end
end
