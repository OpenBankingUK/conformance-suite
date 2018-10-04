defmodule Compliance.ValidationRuns.ValidationRunAccountsTest do
  @moduledoc """
  Tests for validation runs context.
  """
  use Compliance.DataCase

  alias Compliance.ValidationRuns.ValidationRunAccounts
  alias Compliance.Permutations.Generator
  alias Compliance.SwaggerUris
  alias Compliance.Configs.RunConfig
  alias Compliance.Commands
  alias Compliance.Commands.ApiConfig

  import Mock
  import ExUnit.CaptureLog

  require Logger

  @api_version "1.1"
  @validation_run_id "validation-run-id-uuid"
  @auth_server_id "xxxj4NmBD8lQxmLh2O"

  @config %{
    api_version: @api_version,
    authorization_endpoint: "http://example.com/auth",
    client_id: "testClientId",
    client_secret: nil,
    fapi_financial_id: "testFapiFinancialId",
    issuer: "http://aspsp.example.com",
    redirect_uri: "http://tpp.example.com/redirect",
    resource_endpoint: "http://example.com",
    signing_key:
      "-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n",
    signing_kid: "XXXXXX-XXXXxxxXxXXXxxx_xxxx",
    token_endpoint: "http://example.com/token",
    token_endpoint_auth_method: "private_key_jwt",
    transport_cert:
      "-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n",
    transport_key:
      "-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n"
  }
  @run_config struct(RunConfig, @config |> Map.delete(:api_version))
  @api_config struct(ApiConfig, @config)

  @account_id "22290"
  @accounts_basic_response %{
    "Data" => %{
      "Account" => [
        %{
          "AccountId" => @account_id,
          "Currency" => "GBP"
        }
      ]
    },
    "Links" => %{"Self" => "/accounts"},
    "Meta" => %{"TotalPages" => 1},
    "failedValidation" => false
  }
  @basic_accounts_swaggers [
                             SwaggerUris.from("accounts", "1.1", "generic"),
                             SwaggerUris.from("accounts", "1.1", "basic")
                           ]
                           |> Enum.join(" ")
  @detail_accounts_swaggers [
                              SwaggerUris.from("accounts", "1.1", "generic"),
                              SwaggerUris.from("accounts", "1.1", "detail")
                            ]
                            |> Enum.join(" ")

  @accounts_detail_response "mock-accounts-detail-response"
  @account_basic_response "mock-account-basic-response"

  @basic_permission "ReadAccountsBasic"
  @detail_permission "ReadAccountsDetail"

  describe "request_account_resources/4" do
    defp endpoint_permutations do
      {:ok,
       [
         %{
           "endpoint" => "/accounts",
           "permissions" => [@basic_permission],
           "optional" => false,
           "conditional" => false
         },
         %{
           "endpoint" => "/accounts",
           "permissions" => [@detail_permission],
           "optional" => false,
           "conditional" => false
         },
         %{
           "endpoint" => "/accounts/{AccountId}",
           "permissions" => [@basic_permission],
           "optional" => false,
           "conditional" => false
         }
       ]}
    end

    defp mocks do
      [
        {
          Generator,
          [],
          [
            endpoint_permutations: fn _version -> endpoint_permutations() end
          ]
        },
        {
          Commands,
          [],
          [
            request_resource: fn endpoint,
                                 _run_id,
                                 _swagger_uris,
                                 permissions,
                                 _auth_server_id,
                                 _config = %ApiConfig{},
                                 revoke_consent_at_end: _flag ->
              case endpoint do
                "/open-banking/v1.1/accounts" ->
                  case permissions do
                    @basic_permission -> {:ok, @accounts_basic_response}
                    @detail_permission -> {:ok, @accounts_detail_response}
                  end

                "/open-banking/v1.1/accounts/#{@account_id}" ->
                  {:ok, @account_basic_response}
              end
            end
          ]
        }
      ]
    end

    test "uses api_version to obtain endpoint permutations" do
      with_mocks(mocks()) do
        assert :ok ==
                 ValidationRunAccounts.request_account_resources(
                   @api_version,
                   @validation_run_id,
                   @auth_server_id,
                   @config
                 )

        assert called(Generator.endpoint_permutations(@api_version))
      end
    end

    test "requests /accounts to obtain AccountId" do
      with_mocks(mocks()) do
        assert :ok ==
                 ValidationRunAccounts.request_account_resources(
                   @api_version,
                   @validation_run_id,
                   @auth_server_id,
                   @config
                 )

        assert called(
                 Commands.request_resource(
                   "/open-banking/v1.1/accounts",
                   @validation_run_id,
                   @basic_accounts_swaggers,
                   @basic_permission,
                   @auth_server_id,
                   @api_config,
                   revoke_consent_at_end: true
                 )
               )
      end
    end

    test "requests second endpoint from supplied permutations" do
      with_mocks(mocks()) do
        {:ok, permutations} = endpoint_permutations()

        assert :ok ==
                 ValidationRunAccounts.request_account_resource_permutations(
                   permutations |> Enum.take(2),
                   @api_version,
                   @validation_run_id,
                   @auth_server_id,
                   @run_config
                 )

        assert called(
                 Commands.request_resource(
                   "/open-banking/v1.1/accounts",
                   @validation_run_id,
                   @detail_accounts_swaggers,
                   @detail_permission,
                   @auth_server_id,
                   @api_config,
                   revoke_consent_at_end: true
                 )
               )
      end
    end

    test "requests endpoint from supplied permutations substituting AccountId into the url" do
      with_mocks(mocks()) do
        assert :ok ==
                 ValidationRunAccounts.request_account_resources(
                   @api_version,
                   @validation_run_id,
                   @auth_server_id,
                   @config
                 )

        assert called(
                 Commands.request_resource(
                   "/open-banking/v1.1/accounts/#{@account_id}",
                   @validation_run_id,
                   @basic_accounts_swaggers,
                   @basic_permission,
                   @auth_server_id,
                   @api_config,
                   revoke_consent_at_end: true
                 )
               )
      end
    end
  end

  describe "request_account_resources/3 error logging" do
    defp endpoint_permutations_error_logging do
      {:ok,
       [
         %{
           "endpoint" => "/accounts",
           "permissions" => [@basic_permission],
           "optional" => false,
           "conditional" => false
         },
         %{
           "endpoint" => "/accounts",
           "permissions" => [@detail_permission],
           "optional" => true,
           "conditional" => false
         },
         %{
           "endpoint" => "/accounts/{AccountId}",
           "permissions" => [@basic_permission],
           "optional" => false,
           "conditional" => true
         },
         %{
           "endpoint" => "/accounts/{AccountId}",
           "permissions" => [@detail_permission],
           "optional" => true,
           "conditional" => true
         }
       ]}
    end

    defp mocks_error_logging do
      [
        {
          Generator,
          [],
          [
            endpoint_permutations: fn _version -> endpoint_permutations_error_logging() end
          ]
        },
        {
          Commands,
          [],
          [
            request_resource: fn endpoint,
                                 _run_id,
                                 _swagger_uris,
                                 permissions,
                                 _auth_server_id,
                                 _config = %ApiConfig{},
                                 revoke_consent_at_end: _flag ->
              case endpoint do
                "/open-banking/v1.1/accounts" ->
                  case permissions do
                    @basic_permission ->
                      {:ok, @accounts_basic_response}

                    @detail_permission ->
                      {:error, "Error /open-banking/v1.1/accounts: ReadAccountsDetail"}
                  end

                "/open-banking/v1.1/accounts/#{@account_id}" ->
                  case permissions do
                    @basic_permission ->
                      {:error,
                       "Error /open-banking/v1.1/accounts/#{@account_id}: ReadAccountsBasic"}

                    @detail_permission ->
                      {:error,
                       "Error /open-banking/v1.1/accounts/#{@account_id}: ReadAccountsDetail"}
                  end
              end
            end
          ]
        }
      ]
    end

    test "conditional endpoint does not log error" do
      with_mocks(mocks_error_logging()) do
        assert capture_log(fn ->
                 ValidationRunAccounts.request_account_resources(
                   @api_version,
                   @validation_run_id,
                   @auth_server_id,
                   @config
                 )
               end) =~ ""
      end
    end
  end
end
