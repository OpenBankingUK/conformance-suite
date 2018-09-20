defmodule Compliance.ConfigsTest do
  @moduledoc """
  Tests for Configs module.
  """
  use ExUnit.Case, async: true

  alias Compliance.Configs
  alias Compliance.Configs.RunConfig

  import Mock
  import ExUnit.CaptureLog

  @issuer "https://aspsp.example.com"
  @authorization_endpoint "https://aspsp.example.com/auth"
  @token_endpoint "https://aspsp.example.com/token"
  @openid_config %{
    issuer: @issuer,
    claims_parameter_supported: false,
    claims_supported: [],
    grant_types_supported: [
      "authorization_code",
      "client_credentials"
    ],
    response_types_supported: [
      "code",
      "code id_token"
    ],
    request_parameter_supported: true,
    request_uri_parameter_supported: false,
    require_request_uri_registration: false,
    scopes_supported: [
      "openid",
      "accounts",
      "payments"
    ],
    id_token_signing_alg_values_supported: [
      "none",
      "HS256",
      "RS256",
      "PS256"
    ],
    request_object_signing_alg_values_supported: [
      "none",
      "HS256",
      "RS256",
      "PS256"
    ],
    token_endpoint_auth_methods_supported: [
      "client_secret_basic",
      "client_secret_post",
      "client_secret_jwt",
      "private_key_jwt"
    ],
    token_endpoint_auth_signing_alg_values_supported: [
      "none",
      "HS256",
      "RS256",
      "PS256"
    ],
    userinfo_signing_alg_values_supported: [
      "none"
    ],
    claim_types_supported: [
      "normal"
    ],
    subject_types_supported: [
      "public"
    ],
    response_modes_supported: [
      "form_post",
      "fragment",
      "query"
    ],
    token_endpoint: @token_endpoint,
    authorization_endpoint: @authorization_endpoint,
    registration_endpoint: "https://aspsp.example.com/reg",
    jwks_uri: "https://aspsp.example.com/jwks"
  }

  describe "Configs.from_openid_config" do
    @openid_config_endpoint "https://example.com/.well-known/openid-configuration"

    def mocked_response_fn(mocked_response) do
      fn @openid_config_endpoint ->
        {:ok, %HTTPoison.Response{status_code: 200, body: Poison.encode!(mocked_response)}}
      end
    end

    test "calls openid_config_endpoint" do
      with_mock HTTPoison, get: mocked_response_fn(@openid_config) do
        Configs.from_openid_config(@openid_config_endpoint)
        assert called(HTTPoison.get(@openid_config_endpoint))
      end
    end

    test "returns pre-populated RunConfig struct when openid config available,
        with token_endpoint_auth_method field hardcoded to 'private_key_jwt',
        as for now we only support 'private_key_jwt'" do
      with_mock HTTPoison, get: mocked_response_fn(@openid_config) do
        assert Configs.from_openid_config(@openid_config_endpoint) ==
                 %RunConfig{
                   authorization_endpoint: @authorization_endpoint,
                   issuer: @issuer,
                   token_endpoint: @token_endpoint,
                   token_endpoint_auth_method: "private_key_jwt"
                 }
      end
    end

    test "returns pre-populated RunConfig struct when partial openid config available" do
      partial_openid_config = %{authorization_endpoint: @authorization_endpoint}

      with_mock HTTPoison, get: mocked_response_fn(partial_openid_config) do
        assert Configs.from_openid_config(@openid_config_endpoint) ==
                 %RunConfig{
                   authorization_endpoint: @authorization_endpoint,
                   issuer: nil,
                   token_endpoint: nil,
                   token_endpoint_auth_method: "private_key_jwt"
                 }
      end
    end

    test "returns empty RunConfig struct when openid config not available due to non-200 status code" do
      with_mock HTTPoison,
        get: fn @openid_config_endpoint ->
          {:ok, %HTTPoison.Response{status_code: 404}}
        end do
        assert capture_log(fn ->
                 assert Configs.from_openid_config(@openid_config_endpoint) == %RunConfig{
                          authorization_endpoint: nil,
                          issuer: nil,
                          token_endpoint: nil,
                          token_endpoint_auth_method: "private_key_jwt"
                        }
               end) =~
                 "{:warn, %{reason: {:ok, %HTTPoison.Response{body: nil, headers: [], request_url: nil, status_code: 404}}, url: \"https://example.com/.well-known/openid-configuration\"}}"
      end
    end

    test "returns empty RunConfig struct when openid config endpoint errors" do
      with_mock HTTPoison,
        get: fn @openid_config_endpoint ->
          {:error, %HTTPoison.Error{id: nil, reason: :nxdomain}}
        end do
        assert capture_log(fn ->
                 assert Configs.from_openid_config(@openid_config_endpoint) == %RunConfig{
                          authorization_endpoint: nil,
                          issuer: nil,
                          token_endpoint: nil,
                          token_endpoint_auth_method: "private_key_jwt"
                        }
               end) =~
                 "{:warn, %{reason: {:error, %HTTPoison.Error{id: nil, reason: :nxdomain}}, url: \"https://example.com/.well-known/openid-configuration\"}}"
      end
    end
  end
end
