defmodule Compliance.Configs do
  @moduledoc """
  The Configs context.
  """

  require Logger

  alias Compliance.Configs.RunConfig

  def aspsp_host, do: System.get_env("ASPSP_AUTH_HOST_IP") || "reference-mock-server"

  @doc """
  Returns pre-populated RunConfig struct based on values in
  openid_config_endpoint resource,
  or empty RunConfig struct when openid_config_endpoint resource
  is not available.
  """
  def from_openid_config(openid_config_endpoint)
      when is_binary(openid_config_endpoint) do
    openid_config_endpoint
    |> HTTPoison.get()
    |> case do
      {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
        case Poison.decode(body) do
          {:ok, response} ->
            response

          {:error, reason} ->
            Logger.warn(inspect({:warn, %{reason: reason, url: openid_config_endpoint}}))
            %{}
        end

      other ->
        Logger.warn(inspect({:warn, %{reason: other, url: openid_config_endpoint}}))
        %{}
    end
    |> RunConfig.from_openid_config()
  end
end
