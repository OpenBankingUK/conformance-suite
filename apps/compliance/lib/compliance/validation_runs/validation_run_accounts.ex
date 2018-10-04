defmodule Compliance.ValidationRuns.ValidationRunAccounts do
  @moduledoc """
  Functions related to creating and initiating validation runs.
  """

  alias Compliance.Permutations.Generator
  alias Compliance.SwaggerUris
  alias Compliance.Commands
  alias Compliance.Configs.RunConfig
  require Logger

  def request_account_resources(
        api_version,
        validation_run_id,
        auth_server_id,
        config = %RunConfig{}
      )
      when is_binary(api_version) do
    permutations = Generator.endpoint_permutations(api_version)

    Logger.debug(fn ->
      "Compliance.ValidationRuns.request_account_resources, permutations: #{inspect(permutations)}"
    end)

    api_version
    |> Generator.endpoint_permutations()
    |> case do
      {:ok, permutations} ->
        Logger.debug(fn ->
          "Compliance.ValidationRuns.request_account_resources, permutations: #{
            inspect(permutations)
          }"
        end)

        request_account_resource_permutations(
          permutations,
          api_version,
          validation_run_id,
          auth_server_id,
          config
        )

        :ok

      {:error, msg} ->
        {:error, msg}
    end
  end

  def request_account_resources(
        api_version,
        validation_run_id,
        auth_server_id,
        config = %{}
      )
      when is_binary(api_version) do
    params = binding()

    Logger.debug(fn ->
      "Compliance.ValidationRuns.request_account_resources, params: #{inspect(params)}"
    end)

    config = RunConfig.from_map(config)
    request_account_resources(api_version, validation_run_id, auth_server_id, config)
  end

  def request_account_resource_permutations(
        endpoint_permutations,
        api_version,
        validation_run_id,
        auth_server_id,
        config = %RunConfig{}
      )
      when is_binary(api_version) do
    api_version
    |> get_account_id(validation_run_id, auth_server_id, config)
    |> case do
      {:ok, account_id} ->
        endpoint_permutations
        |> Enum.each(fn endpoint_permutation ->
          Logger.debug(fn -> "endpoint_permutation: #{inspect(endpoint_permutation)}" end)

          request_account_resource(
            api_version,
            endpoint_permutation,
            account_id,
            validation_run_id,
            auth_server_id,
            config
          )
        end)

      {:error, msg} ->
        Logger.error(msg)
    end
  end

  defp get_account_id(
         api_version,
         validation_run_id,
         auth_server_id,
         config = %RunConfig{}
       ) do
    swagger_uris = SwaggerUris.from("accounts", api_version, ["ReadAccountsBasic"])

    "/open-banking/v#{api_version}/accounts"
    |> Commands.request_resource(
      validation_run_id,
      swagger_uris,
      "ReadAccountsBasic",
      auth_server_id,
      RunConfig.to_api_config(config, api_version),
      revoke_consent_at_end: true
    )
    |> case do
      {:ok, response} ->
        try do
          account_id =
            response
            |> Map.get("Data")
            |> Map.get("Account")
            |> List.first()
            |> Map.get("AccountId")

          {:ok, account_id}
        catch
          _ ->
            {:error,
             "unable to obtain AccountId from /accounts payload | validation_run_id: #{
               validation_run_id
             } | auth_server_id: #{auth_server_id} | response: #{inspect(response)}"}
        end

      error ->
        {:error,
         "unable to get /accounts response | validation_run_id: #{validation_run_id} | auth_server_id: #{
           auth_server_id
         } | error: #{inspect(error)}"}
    end
  end

  defp request_account_resource(
         _api_version,
         %{"endpoint" => "/accounts", "permissions" => ["ReadAccountsBasic"]},
         _account_id,
         _validation_run_id,
         _auth_server_id,
         _config = %RunConfig{}
       ) do
    # we've already requested this first resource
    nil
  end

  defp request_account_resource(
         api_version,
         %{
           "endpoint" => endpoint,
           "permissions" => permissions,
           "optional" => optional,
           "conditional" => conditional
         },
         account_id,
         validation_run_id,
         auth_server_id,
         config = %RunConfig{}
       ) do
    swagger_uris = SwaggerUris.from("accounts", api_version, permissions)
    permissions = permissions |> Enum.join(" ")
    endpoint = endpoint |> String.replace("{AccountId}", account_id)

    "/open-banking/v#{api_version}#{endpoint}"
    |> Commands.request_resource(
      validation_run_id,
      swagger_uris,
      permissions,
      auth_server_id,
      RunConfig.to_api_config(config, api_version),
      revoke_consent_at_end: true
    )
    |> case do
      {:error, result} ->
        log_error = !(conditional || optional)

        if log_error do
          Logger.error(fn ->
            "Compliance.ValidationRuns.request_account_resources failed: #{inspect(result)}"
          end)
        end

      _ ->
        nil
    end
  end
end
