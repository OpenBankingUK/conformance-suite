defmodule OBApiRemote.Commands.ApiConfig do
  @moduledoc """
  Represents a config required to use OB APIs.

  ## Example

    iex> %ApiConfig{
    iex>   api_version: "1.1",
    iex>   authorization_endpoint: "http://aspsp.example.com/auth"
    iex>   client_id: "clientId",
    iex>   client_secret: "clientSecret",
    iex>   fapi_financial_id: "xyz",
    iex>   issuer: "http://aspsp.example.com",
    iex>   redirect_uri: "http://tpp.example.com/redirect"
    iex>   resource_endpoint: "http://aspsp.example.com" # without "/openbanking/v*.*" prefix
    iex>   signing_key: "-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n",
    iex>   signing_kid: "XXXXXX-XXXXxxxXxXXXxxx_xxxx",
    iex>   token_endpoint: "http://aspsp.example.com/token",
    iex>   token_endpoint_auth_method: "private_key_jwt",
    iex>   transport_cert: "-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n",
    iex>   transport_key: "-----BEGIN PRIVATE KEY-----\nexample\nexample\nexample\n-----END PRIVATE KEY-----\n"
    iex> }
  """
  alias __MODULE__

  defstruct api_version: "",
            authorization_endpoint: "",
            client_id: "",
            client_secret: nil,
            fapi_financial_id: "",
            issuer: "",
            redirect_uri: "",
            resource_endpoint: "",
            signing_key: "",
            signing_kid: "",
            token_endpoint: "",
            token_endpoint_auth_method: "",
            transport_cert: "",
            transport_key: ""

  def base64_encode_json(config = %ApiConfig{}) do
    config
    |> Map.from_struct()
    |> base64_encode_json()
  end

  def base64_encode_json(config = %{}) do
    config
    |> Poison.encode!()
    |> Base.encode64()
  end
end
