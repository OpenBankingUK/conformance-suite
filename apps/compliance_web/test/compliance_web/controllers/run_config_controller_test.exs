defmodule ComplianceWeb.RunConfigControllerTest do
  use ComplianceWeb.ConnCase

  alias Compliance.Configs
  alias Compliance.Configs.RunConfig

  import Mock

  @openid_config_endpoint "https://example.com/.well-known/openid-configuration"

  describe "POST /run-configs" do
    @config RunConfig.from_openid_config(%{})

    test "returns config map", %{
      conn: conn
    } do
      with_mock Configs, from_openid_config: fn @openid_config_endpoint -> @config end do
        conn =
          conn
          |> simulate_authenticated
          |> put_req_header("content-type", "application/json")
          |> post("/run-configs", %{"openid_config_endpoint" => @openid_config_endpoint})

        body = json_response(conn, :ok)

        assert body == %{
                 "authorization_endpoint" => nil,
                 "client_id" => nil,
                 "client_secret" => nil,
                 "fapi_financial_id" => nil,
                 "issuer" => nil,
                 "redirect_uri" => nil,
                 "resource_endpoint" => nil,
                 "signing_key" => nil,
                 "signing_kid" => nil,
                 "token_endpoint" => nil,
                 "token_endpoint_auth_method" => "private_key_jwt",
                 "transport_cert" => nil,
                 "transport_key" => nil
               }
      end
    end
  end
end
