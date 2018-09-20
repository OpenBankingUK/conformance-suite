defmodule OBApiRemote.Commands.DriverTest do
  @moduledoc """
  Tests for Commands.
  """
  use ExUnit.Case, async: true
  alias OBApiRemote.Commands.{ApiConfig, Driver, Proxied}

  import Mock

  @urls Proxied.urls()

  @session_token "58b47d20-591c-11e8-950f-b72f26bb29de"
  @aspsp_auth_server_id "aaaj4NmBD8lQxmLh2O"
  @interaction_id "testInteractionId"
  @authorisation_code "spoofAuthorisationCode"
  @account_request_id "mock-account-request-id"

  @validation_run_id "validation-run-id-xxx"

  @generic_account_swagger_uri "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger.json"
  @basic_account_swagger_uri "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger-basic.json"
  @account_swagger_uris "#{@generic_account_swagger_uri} #{@basic_account_swagger_uri}"

  @payment_swagger_uri "https://raw.githubusercontent.com/OpenBankingUK/payment-initiation-api-spec/master/dist/v1.1/payment-initiation-swagger.json"

  @payment %{
    "name" => "Sam Morse",
    "sort_code" => "111111",
    "account_number" => "12345678",
    "amount" => "10.00"
  }

  def parsed_state(scope) do
    %{
      "accountRequestId" => @account_request_id,
      "authorisationServerId" => @aspsp_auth_server_id,
      "interactionId" => @interaction_id,
      "sessionId" => @session_token,
      "scope" => scope
    }
  end

  @type_headers [
    {"content-type", "application/json; charset=utf-8"},
    {"accept", "application/json; charset=utf-8"}
  ]

  @config %{
    api_version: "2.0",
    authorization_endpoint: "http://example.com/auth",
    client_id: "testClientId",
    client_secret: nil,
    fapi_financial_id: "testFapiId",
    issuer: "http://aspsp.example.com",
    redirect_uri: "http://tpp.example.com/redirect",
    resource_endpoint: "http://example.com",
    signing_key:
      "-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n",
    signing_kid: "XXXXXX-XXXXxxxXxXXXxxx_xxxx",
    token_endpoint: "http: //example.com/token",
    token_endpoint_auth_method: "private_key_jwt",
    transport_cert:
      "-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n",
    transport_key:
      "-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n"
  }
  @api_config struct(ApiConfig, @config)

  @config_header [
    {"x-config", ApiConfig.base64_encode_json(@config)}
  ]

  @session_headers [
                     {"authorization", "#{@session_token}"}
                   ] ++ @type_headers

  defp make_state(parsed) do
    parsed
    |> Map.put("scope", "openid " <> parsed["scope"])
    |> Poison.encode!()
    |> Base.encode64()
  end

  defp make_authorise_consent_url(aspsp_auth_server_id, raw_state, scope) do
    %{
      state: raw_state,
      client_id: "spoofClientId",
      response_type: "code",
      request:
        "eyJhbGciOiJub25lIn0.eyJpc3MiOiJzcG9vZkNsaWVudElkIiwicmVzcG9uc2VfdHlwZSI6ImNvZGUiLCJjbGllbnRfaWQiOiJzcG9vZkNsaWVudElkIiwicmVkaXJlY3RfdXJpIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3RwcC9hdXRob3JpemVkIiwic2NvcGUiOiJvcGVuaWQgcGF5bWVudHMiLCJzdGF0ZSI6ImV5SmhkWFJvYjNKcGMyRjBhVzl1VTJWeWRtVnlTV1FpT2lKaFlXRnFORTV0UWtRNGJGRjRiVXhvTWs4aUxDSnBiblJsY21GamRHbHZia2xrSWpvaVlUTTFZMlF6TkdRdFl6QTNZaTAwTVdaaExXSmpaR1F0WWpjNVlUUTVOR0U0TlRFNElpd2ljMlZ6YzJsdmJrbGtJam9pWkRsbFpUSm1OekF0WldNNU5DMHhNV1UzTFdFeFpqWXRNR1l3T1RjNVlqYzNaVEppSWl3aWMyTnZjR1VpT2lKdmNHVnVhV1FnY0dGNWJXVnVkSE1pZlE9PSIsIm5vbmNlIjoiZHVtbXktbm9uY2UiLCJtYXhfYWdlIjo4NjQwMCwiY2xhaW1zIjp7InVzZXJpbmZvIjp7Im9wZW5iYW5raW5nX2ludGVudF9pZCI6eyJ2YWx1ZSI6IjVmMGNiYjAxLTQzOTctNDhmZi04MDE3LTQ3OTA4YmU0NWNlYiIsImVzc2VudGlhbCI6dHJ1ZX19LCJpZF90b2tlbiI6eyJvcGVuYmFua2luZ19pbnRlbnRfaWQiOnsidmFsdWUiOiI1ZjBjYmIwMS00Mzk3LTQ4ZmYtODAxNy00NzkwOGJlNDVjZWIiLCJlc3NlbnRpYWwiOnRydWV9LCJhY3IiOnsiZXNzZW50aWFsIjp0cnVlfX19fQ.",
      scope: scope
    }
    |> URI.encode_query()
    |> (&("http://localhost:8080/tpp/authorized&" <> &1)).()
    |> URI.encode_www_form()
    |> (&("http://localhost:8001/#{aspsp_auth_server_id}/authorize?redirect_uri=" <> &1)).()
  end

  defp assert_error({_, %{reason: reason, cmd: desc}}, target_url) do
    assert(reason == "ERROR")
    assert(desc.url =~ target_url)
  end

  defp assert_error({_, %{reason: reason, url: url}}, target_url) do
    assert_error({nil, %{reason: reason, cmd: %{url: url}}}, target_url)
  end

  describe "Driver.do_post_login" do
    @mock_response "{\"sid\":\"#{@session_token}\"}"
    @request_body %{u: "alice", p: "wonderland"}

    @request_headers @type_headers

    test "returns session token on success" do
      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:ok, %HTTPoison.Response{status_code: 200, body: @mock_response}}
        end do
        assert({:ok, @session_token} == Driver.do_post_login())

        assert called(
                 HTTPoison.request(
                   :post,
                   @urls.login,
                   Poison.encode!(@request_body),
                   @request_headers
                 )
               )
      end
    end

    test "returns error reason on failure" do
      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:error, %HTTPoison.Error{id: "an-id", reason: "ERROR"}}
        end do
        assert_error(
          Driver.do_post_login(),
          @urls.login
        )

        assert called(
                 HTTPoison.request(
                   :post,
                   @urls.login,
                   Poison.encode!(@request_body),
                   @request_headers
                 )
               )
      end
    end
  end

  describe "Driver.do_get_logout" do
    @mock_response "{\"sid\":\"#{@session_token}\"}"
    @request_headers @session_headers

    test "returns session token on success" do
      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:ok, %HTTPoison.Response{status_code: 200, body: @mock_response}}
        end do
        assert({:ok, @session_token} == Driver.do_get_logout(@session_token))
        assert called(HTTPoison.request(:get, @urls.logout, "", @request_headers))
      end
    end

    test "returns error reason on failure" do
      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:error, %HTTPoison.Error{id: "an-id", reason: "ERROR"}}
        end do
        assert_error(
          Driver.do_get_logout(@session_token),
          @urls.logout
        )

        assert called(HTTPoison.request(:get, @urls.logout, "", @request_headers))
      end
    end
  end

  describe "Driver.do_get_authorisation_servers" do
    @mock_response "[{\"name\":\"AAA Example Bank\",\"logoUri\":\"\",\"id\":\"aaaj4NmBD8lQxmLh2O\",\"accountsConsentGranted\":false}]"
    @auth_servers [
      %{
        "id" => @aspsp_auth_server_id,
        "logoUri" => "",
        "name" => "AAA Example Bank"
      }
    ]
    @request_headers @session_headers

    test "returns authorization servers on success" do
      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:ok, %HTTPoison.Response{status_code: 200, body: @mock_response}}
        end do
        assert({:ok, @auth_servers} == Driver.do_get_authorisation_servers(@session_token))
        assert called(HTTPoison.request(:get, @urls.get_auth_servers, "", @request_headers))
      end
    end

    test "returns error reason on failure" do
      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:error, %HTTPoison.Error{id: "an-id", reason: "ERROR"}}
        end do
        assert_error(
          Driver.do_get_authorisation_servers(@session_token),
          @urls.get_auth_servers
        )

        assert called(HTTPoison.request(:get, @urls.get_auth_servers, "", @request_headers))
      end
    end
  end

  describe "Driver.do_post_authorise_account_access" do
    @scope "accounts"
    @permissions "ReadAccountsBasic ReadTransactionsDebits"
    @request_body ""
    @request_headers [
                       {"authorization", @session_token},
                       {"x-authorization-server-id", @aspsp_auth_server_id},
                       {"x-swagger-uris", @account_swagger_uris},
                       {"x-validation-run-id", @validation_run_id}
                     ] ++ @config_header ++ [{"x-permissions", @permissions}] ++ @type_headers

    test "returns authorise consent url on success" do
      state = make_state(parsed_state(@scope))
      authorise_consent_url = make_authorise_consent_url(@aspsp_auth_server_id, state, @scope)

      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:ok,
           %HTTPoison.Response{
             status_code: 200,
             body: Poison.encode!(%{uri: authorise_consent_url})
           }}
        end do
        assert(
          {:ok, authorise_consent_url} ==
            Driver.do_post_authorise_account_access(
              @permissions,
              @validation_run_id,
              @account_swagger_uris,
              @session_token,
              @aspsp_auth_server_id,
              @api_config
            )
        )

        assert called(
                 HTTPoison.request(
                   :post,
                   @urls.authorise_account_access,
                   @request_body,
                   @request_headers
                 )
               )
      end
    end
  end

  describe "Driver.do_post_revoke_account_access_consent" do
    @request_headers [
                       {"authorization", "#{@session_token}"},
                       {"x-authorization-server-id", "#{@aspsp_auth_server_id}"},
                       {"x-swagger-uris", @account_swagger_uris},
                       {"x-validation-run-id", "#{@validation_run_id}"}
                     ] ++ @config_header ++ @type_headers

    test "returns :ok on success" do
      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:ok,
           %HTTPoison.Response{
             status_code: 201,
             body: ""
           }}
        end do
        assert(
          :ok ==
            Driver.do_post_revoke_account_access_consent(
              @validation_run_id,
              @account_swagger_uris,
              %{
                "sessionId" => @session_token,
                "authorisationServerId" => @aspsp_auth_server_id
              },
              @api_config
            )
        )

        assert called(
                 HTTPoison.request(
                   :post,
                   @urls.account_request_revoke_consent,
                   "",
                   @request_headers
                 )
               )
      end
    end
  end

  describe "Driver.do_post_authorise_payment" do
    @scope "payments"
    @request_headers [
                       {"authorization", "#{@session_token}"},
                       {"x-authorization-server-id", "#{@aspsp_auth_server_id}"},
                       {"x-swagger-uris", @payment_swagger_uri},
                       {"x-validation-run-id", "#{@validation_run_id}"}
                     ] ++ @config_header ++ @type_headers

    @request_body %{
      authorisationServerId: @aspsp_auth_server_id,
      InstructedAmount: %{
        Amount: @payment["amount"],
        Currency: "GBP"
      },
      CreditorAccount: %{
        SchemeName: "SortCodeAccountNumber",
        Identification: "#{@payment["sort_code"]}#{@payment["account_number"]}",
        Name: @payment["name"]
      }
    }

    test "returns authorise consent url on success" do
      state = make_state(parsed_state(@scope))
      authorise_consent_url = make_authorise_consent_url(@aspsp_auth_server_id, state, @scope)

      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:ok,
           %HTTPoison.Response{
             status_code: 200,
             body: Poison.encode!(%{uri: authorise_consent_url})
           }}
        end do
        assert(
          {:ok, authorise_consent_url} ==
            Driver.do_post_authorise_payment(
              @payment,
              @validation_run_id,
              @payment_swagger_uri,
              @session_token,
              @aspsp_auth_server_id,
              @api_config
            )
        )

        assert called(
                 HTTPoison.request(
                   :post,
                   @urls.authorise_payment,
                   Poison.encode!(@request_body),
                   @request_headers
                 )
               )
      end
    end

    test "returns error reason on failure" do
      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:error, %HTTPoison.Error{id: "an-id", reason: "ERROR"}}
        end do
        assert_error(
          Driver.do_post_authorise_payment(
            @payment,
            @validation_run_id,
            @payment_swagger_uri,
            @session_token,
            @aspsp_auth_server_id,
            @api_config
          ),
          @urls.authorise_payment
        )

        assert called(
                 HTTPoison.request(
                   :post,
                   @urls.authorise_payment,
                   Poison.encode!(@request_body),
                   @request_headers
                 )
               )
      end
    end
  end

  describe "Driver.do_get_authorise_consent" do
    @scope "payments"
    @mock_response ""

    test "returns raw state and authorisation code on success" do
      state = make_state(parsed_state(@scope))
      authorise_consent_url = make_authorise_consent_url(@aspsp_auth_server_id, state, @scope)

      with_mock HTTPoison,
        request: fn _method, _url ->
          {:ok,
           %HTTPoison.Response{
             body:
               "Found. Redirecting to http://localhost:8080/tpp/authorized?code=#{
                 @authorisation_code
               }&state=#{state}",
             headers: [
               {"Location",
                "http://localhost:8080/tpp/authorized?code=#{@authorisation_code}&state=#{state}"},
               {"Content-Type", "text/plain; charset=utf-8"}
             ],
             request_url: authorise_consent_url,
             status_code: 302
           }}
        end do
        assert(
          {:ok, %{"code" => @authorisation_code, "state" => state}} ==
            Driver.do_get_authorise_consent(authorise_consent_url)
        )

        assert called(HTTPoison.request(:get, authorise_consent_url))
      end
    end

    test "returns error reason on failure" do
      state = make_state(parsed_state(@scope))
      authorise_consent_url = make_authorise_consent_url(@aspsp_auth_server_id, state, @scope)

      with_mock HTTPoison,
        request: fn _method, _url ->
          {:error, %HTTPoison.Error{id: "an-id", reason: "ERROR"}}
        end do
        assert_error(
          Driver.do_get_authorise_consent(authorise_consent_url),
          authorise_consent_url
        )

        assert called(HTTPoison.request(:get, authorise_consent_url))
      end
    end
  end

  describe "Driver.do_post_consent_authorised" do
    @scope "payments"
    @mock_response ""
    @request_headers [
                       {"authorization", "#{@session_token}"},
                       {"x-authorization-server-id", "#{@aspsp_auth_server_id}"},
                       {"x-swagger-uris", @account_swagger_uris},
                       {"x-validation-run-id", "#{@validation_run_id}"}
                     ] ++ @config_header ++ @type_headers
    @request_body %{
      accountRequestId: @account_request_id,
      authorisationServerId: @aspsp_auth_server_id,
      authorisationCode: @authorisation_code,
      scope: @scope
    }

    test "ensures consent authorisation is acknowledged" do
      result_state = parsed_state(@scope)
      state = make_state(result_state)

      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:ok, %HTTPoison.Response{status_code: 204, body: @mock_response}}
        end do
        assert(
          {:ok, result_state} ==
            Driver.do_post_consent_authorised(
              state,
              @authorisation_code,
              @validation_run_id,
              @account_swagger_uris,
              @api_config
            )
        )

        assert called(
                 HTTPoison.request(
                   :post,
                   @urls.consent_authorised,
                   Poison.encode!(@request_body),
                   @request_headers
                 )
               )
      end
    end
  end

  describe "Driver.do_get_resource" do
    @scope "accounts"
    @mock_response "{ \"Data\": { \"Account\": [ { \"AccountId\": \"22290\" } ] } }"
    @payload %{"Data" => %{"Account" => [%{"AccountId" => "22290"}]}}
    @request_headers [
                       {"authorization", "#{@session_token}"},
                       {"x-authorization-server-id", "#{@aspsp_auth_server_id}"},
                       {"x-fapi-interaction-id", "#{@interaction_id}"},
                       {"x-swagger-uris", @account_swagger_uris},
                       {"x-validation-run-id", "#{@validation_run_id}"}
                     ] ++ @config_header ++ @type_headers
    @request_body ""

    test "requests an account resource after consent" do
      endpoint = "/accounts"

      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:ok, %HTTPoison.Response{status_code: 200, body: @mock_response}}
        end do
        assert {:ok, @payload} ==
                 Driver.do_get_resource(
                   endpoint,
                   @validation_run_id,
                   @account_swagger_uris,
                   parsed_state(@scope),
                   @api_config
                 )

        assert called(
                 HTTPoison.request(
                   :get,
                   @urls.accounts,
                   @request_body,
                   @request_headers
                 )
               )
      end
    end
  end

  describe "Driver.do_post_complete_payment" do
    @scope "payments"
    @mock_response ""
    @request_headers [
                       {"authorization", "#{@session_token}"},
                       {"x-authorization-server-id", "#{@aspsp_auth_server_id}"},
                       {"x-fapi-interaction-id", "#{@interaction_id}"},
                       {"x-swagger-uris", @payment_swagger_uri},
                       {"x-validation-run-id", "#{@validation_run_id}"}
                     ] ++ @config_header ++ @type_headers
    @request_body ""

    test "completes an already setup payment" do
      with_mock HTTPoison,
        request: fn _method, _url, _body, _headers ->
          {:ok, %HTTPoison.Response{status_code: 201, body: @mock_response}}
        end do
        assert(
          :ok ==
            Driver.do_post_complete_payment(
              @validation_run_id,
              @payment_swagger_uri,
              parsed_state(@scope),
              @api_config
            )
        )

        assert called(
                 HTTPoison.request(
                   :post,
                   @urls.complete_payment,
                   @request_body,
                   @request_headers
                 )
               )
      end
    end
  end
end
