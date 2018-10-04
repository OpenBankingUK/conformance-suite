defmodule Compliance.CommandsTest do
  @moduledoc """
  Tests for Commands.
  """
  use ExUnit.Case, async: false

  alias Compliance.Commands
  alias Compliance.Commands.{ApiConfig, Driver}

  import Mock

  @auth_server_id "aaaj4NmBD8lQxmLh2O"
  @session_token "58b47d20-591c-11e8-950f-b72f26bb29de"
  @validation_run_id "validation-run-id-uuid"

  @authorise_consent_url "http://authorise-consent-url.com/user/1"
  @interaction_id "testInteractionId"
  @authorisation_code "spoofAuthorisationCode"

  @mock_swagger_uris "swaggeruris"

  def parsed_state(scope) do
    %{
      "authorisationServerId" => @auth_server_id,
      "interactionId" => @interaction_id,
      "sessionId" => @session_token,
      "scope" => scope
    }
  end

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

  @payment1 %{
    "account_number" => "12345678",
    "amount" => "10.00",
    "name" => "Sam Morse",
    "sort_code" => "111111",
    "type" => "payments"
  }

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

  defp make_state(parsed) do
    parsed
    |> Map.put("scope", "openid " <> parsed["scope"])
    |> Poison.encode!()
    |> Base.encode64()
  end

  describe "Commands.read_resource" do
    @scope "accounts"
    @permissions "ReadAccountsBasic"
    @endpoint "/accounts"
    @response %{"Data" => %{}}

    def driver_mocks(state) do
      [
        do_post_login: fn -> {:ok, @session_token} end,
        do_post_authorise_account_access: fn _permissions,
                                             _validation_run_id,
                                             _swagger_uris,
                                             _session_token,
                                             _auth_server_id,
                                             _config = %ApiConfig{} ->
          {:ok, @authorise_consent_url}
        end,
        do_get_authorise_consent: fn _authorise_consent_url ->
          {:ok, %{"code" => @authorisation_code, "state" => state}}
        end,
        do_post_consent_authorised: fn _raw_state,
                                       _auth_code,
                                       _run_id,
                                       _swagger_uris,
                                       _config = %ApiConfig{} ->
          {:ok, parsed_state(@scope)}
        end,
        do_get_resource: fn _endpoint,
                            _validation_run_id,
                            _swagger_uris,
                            _parsed_state,
                            _config = %ApiConfig{} ->
          {:ok, @response}
        end,
        do_post_revoke_account_access_consent: fn _validation_run_id,
                                                  _swagger_uris,
                                                  _parsed_state,
                                                  _api_config ->
          :ok
        end
      ]
    end

    def request_resource(revoke_consent_at_end: revoke_consent_at_end) do
      Commands.request_resource(
        @endpoint,
        @validation_run_id,
        @mock_swagger_uris,
        @permissions,
        @auth_server_id,
        @api_config,
        revoke_consent_at_end: revoke_consent_at_end
      )
    end

    test "requests resource using all required driver calls" do
      state = make_state(parsed_state(@scope))

      with_mock(Driver, driver_mocks(state)) do
        assert {:ok, @response} == request_resource(revoke_consent_at_end: false)

        assert called(Driver.do_post_login())

        assert called(
                 Driver.do_post_authorise_account_access(
                   @permissions,
                   @validation_run_id,
                   @mock_swagger_uris,
                   @session_token,
                   @auth_server_id,
                   @api_config
                 )
               )

        assert called(Driver.do_get_authorise_consent(@authorise_consent_url))

        state_map = parsed_state(@scope)

        assert called(
                 Driver.do_get_resource(
                   @endpoint,
                   @validation_run_id,
                   @mock_swagger_uris,
                   state_map,
                   @api_config
                 )
               )

        refute called(
                 Driver.do_post_revoke_account_access_consent(
                   @validation_run_id,
                   @mock_swagger_uris,
                   state_map,
                   @api_config
                 )
               )
      end
    end

    test "revokes consent when revoke_consent_at_end true" do
      state = make_state(parsed_state(@scope))

      with_mock(Driver, driver_mocks(state)) do
        assert {:ok, @response} == request_resource(revoke_consent_at_end: true)

        assert called(
                 Driver.do_post_revoke_account_access_consent(
                   @validation_run_id,
                   @mock_swagger_uris,
                   parsed_state(@scope),
                   @api_config
                 )
               )
      end
    end
  end

  describe "Commands.make_payment" do
    @scope "payments"

    test "makes a payment using all required driver calls" do
      state = make_state(parsed_state(@scope))

      with_mock(
        Driver,
        do_post_login: fn -> {:ok, @session_token} end,
        do_post_authorise_payment: fn _payment,
                                      _validation_run_id,
                                      _swagger_uris,
                                      _session_token,
                                      _auth_server_id,
                                      _config = %ApiConfig{} ->
          {:ok, @authorise_consent_url}
        end,
        do_get_authorise_consent: fn _authorise_consent_url ->
          {:ok, %{"code" => @authorisation_code, "state" => state}}
        end,
        do_post_consent_authorised: fn _raw_state,
                                       _auth_code,
                                       _run_id,
                                       _swagger_uris,
                                       _config = %ApiConfig{} ->
          {:ok, parsed_state(@scope)}
        end,
        do_post_complete_payment: fn _validation_run_id,
                                     _swagger_uris,
                                     _parsed_state,
                                     _config = %ApiConfig{} ->
          :ok
        end
      ) do
        Commands.make_payment(
          @validation_run_id,
          @mock_swagger_uris,
          @payment1,
          @auth_server_id,
          @api_config
        )

        assert called(Driver.do_post_login())

        assert called(
                 Driver.do_post_authorise_payment(
                   @payment1,
                   @validation_run_id,
                   @mock_swagger_uris,
                   @session_token,
                   @auth_server_id,
                   @api_config
                 )
               )

        assert called(Driver.do_get_authorise_consent(@authorise_consent_url))

        assert called(
                 Driver.do_post_consent_authorised(
                   state,
                   @authorisation_code,
                   @validation_run_id,
                   @mock_swagger_uris,
                   @api_config
                 )
               )

        assert called(
                 Driver.do_post_complete_payment(
                   @validation_run_id,
                   @mock_swagger_uris,
                   parsed_state(@scope),
                   @api_config
                 )
               )
      end
    end
  end
end
