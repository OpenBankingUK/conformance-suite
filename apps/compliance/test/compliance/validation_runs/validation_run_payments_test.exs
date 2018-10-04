defmodule Compliance.ValidationRuns.ValidationRunPaymentsTest do
  @moduledoc """
  Tests for validation runs context.
  """
  use Compliance.DataCase

  alias Compliance.ValidationRuns.ValidationRunPayments
  alias Compliance.SwaggerUris
  alias Compliance.Commands
  alias Compliance.Commands.ApiConfig

  import Mock
  import ExUnit.CaptureLog

  require Logger

  @api_version "1.1"
  @validation_run_id "validation-run-id-uuid"
  @payments [
    %{
      "api_version" => @api_version,
      "account_number" => "12345678",
      "amount" => "10.00",
      "name" => "Sam Morse",
      "sort_code" => "111111",
      "type" => "payments"
    }
  ]
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
  @api_config struct(ApiConfig, @config)

  test "make_payments/3 makes all payments" do
    with_mocks([
      {
        Commands,
        [],
        [
          make_payment: fn _validation_run_id,
                           _swagger_uris,
                           _payment,
                           _auth_server_id,
                           _config = %ApiConfig{} ->
            {:ok, @validation_run_id}
          end
        ]
      }
    ]) do
      assert :ok ==
               ValidationRunPayments.make_payments(
                 @payments,
                 @validation_run_id,
                 @auth_server_id,
                 @config
               )

      assert called(
               Commands.make_payment(
                 @validation_run_id,
                 SwaggerUris.from("payments", "1.1", "generic"),
                 Enum.at(@payments, 0),
                 @auth_server_id,
                 @api_config
               )
             )
    end
  end

  test "make_payments/3 logs payment errors" do
    with_mocks([
      {
        Commands,
        [],
        [
          make_payment: fn _validation_run_id,
                           _swagger_uri,
                           _payment,
                           _auth_server_id,
                           _config = %ApiConfig{} ->
            {:error, "AN ERROR!"}
          end
        ]
      }
    ]) do
      assert capture_log(fn ->
               ValidationRunPayments.make_payments(
                 @payments,
                 @validation_run_id,
                 @auth_server_id,
                 @config
               )
             end) =~ "Compliance.ValidationRunPayments.make_payments failed"
    end
  end
end
