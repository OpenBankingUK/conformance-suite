defmodule OBApiRemote.Commands.Driver do
  @moduledoc """
  The Commands context.
  """

  alias OBApiRemote.Commands.{ApiConfig, Proxied}
  require Logger

  def do_post_login() do
    Logger.debug(fn -> "do_post_login" end)
    handler = fn {:ok, decoded} -> {:ok, decoded["sid"]} end

    :login
    |> execute(handler, params: %{username: "alice", password: "wonderland"})
  end

  def do_get_logout(session_token) do
    Logger.debug(fn -> "do_get_logout" end)
    handler = fn {:ok, decoded} -> {:ok, decoded["sid"]} end

    :logout
    |> execute(handler, headers: [{"authorization", session_token}])
  end

  def do_get_authorisation_servers(session_token) do
    Logger.debug(fn -> "do_get_authorisation_servers" end)

    handler = fn {:ok, auth_servers} ->
      results = auth_servers |> Enum.map(&(&1 |> Map.delete("accountsConsentGranted")))
      {:ok, results}
    end

    :get_auth_servers
    |> execute(handler, headers: [{"authorization", session_token}])
  end

  @doc """
  Pass permissions as a space separated list: "ReadAccountsBasic ReadTransactionsDebits"
  Pass swagger_uris as a space separated list: "http://example/com/1 http://example.com/2"
  """
  def do_post_authorise_account_access(
        permissions,
        validation_run_id,
        swagger_uris,
        session_token,
        auth_server_id,
        config = %ApiConfig{}
      ) do
    Logger.debug(fn ->
      "do_post_authorise_account_access, #{
        inspect({permissions, validation_run_id, session_token, auth_server_id})
      }"
    end)

    handler = fn {:ok, decoded_response} ->
      uri = decoded_response["uri"]

      validation_result = decoded_response["validation_result"]
      json_string = Poison.encode!(validation_result)

      case Compliance.ValidationRuns.AggregateSupervisor.add_log_item(json_string) do
        {:ok, report} ->
          # Logger.info("handler -> report=#{inspect(report)}")
          report

        error ->
          Logger.info("handler -> error=#{inspect(error)}")
          nil
      end

      {:ok, uri}
    end

    :authorise_account_access
    |> execute(
      handler,
      headers:
        authorise_headers(validation_run_id, swagger_uris, session_token, auth_server_id, config) ++
          [{"x-permissions", permissions}]
    )
  end

  def do_post_authorise_payment(
        payment,
        validation_run_id,
        swagger_uris,
        session_token,
        auth_server_id,
        config = %ApiConfig{}
      ) do
    Logger.debug(fn -> "do_post_authorise_payment" end)

    handler = fn {:ok, decoded_response} ->
      uri = decoded_response["uri"]

      validation_result = decoded_response["validation_result"]
      json_string = Poison.encode!(validation_result)

      case Compliance.ValidationRuns.AggregateSupervisor.add_log_item(json_string) do
        {:ok, report} ->
          # Logger.info("handler -> report=#{inspect(report)}")
          report

        error ->
          Logger.info("handler -> error=#{inspect(error)}")
          nil
      end

      {:ok, uri}
    end

    :authorise_payment
    |> execute(
      handler,
      headers:
        authorise_headers(validation_run_id, swagger_uris, session_token, auth_server_id, config),
      params: payment |> Map.put("auth_server_id", auth_server_id)
    )
  end

  def do_get_authorise_consent(authorise_consent_url) do
    Logger.debug(fn -> "do_get_authorise_consent" end)

    HTTPoison.request(:get, authorise_consent_url)
    |> case do
      {:ok, %HTTPoison.Response{status_code: 302, headers: headers}} ->
        parsed_query =
          headers
          |> List.keyfind("Location", 0)
          |> Tuple.to_list()
          |> List.last()
          |> URI.parse()
          |> Map.get(:query)
          |> URI.decode_query()

        {:ok, parsed_query}

      {:ok, %HTTPoison.Response{status_code: 400, body: body}} ->
        reason =
          case Poison.decode(body) do
            {:ok, response} ->
              response

            {:error, _} ->
              body
          end

        {:error, %{reason: reason, url: authorise_consent_url}}

      {:error, %HTTPoison.Error{id: _id, reason: reason}} ->
        {:error, %{reason: reason, url: authorise_consent_url}}
    end
  end

  def do_post_consent_authorised(
        state,
        authorisation_code,
        validation_run_id,
        swagger_uris,
        config = %ApiConfig{}
      ) do
    Logger.debug(fn -> "do_post_consent_authorised" end)
    parsed_state = parse_state(state)
    handler = fn {:ok, _} -> {:ok, parsed_state} end

    :consent_authorised
    |> execute(
      handler,
      headers: [
        {"authorization", parsed_state["sessionId"]},
        {"x-authorization-server-id", parsed_state["authorisationServerId"]},
        {"x-swagger-uris", swagger_uris},
        {"x-validation-run-id", validation_run_id},
        {"x-config", ApiConfig.base64_encode_json(config)}
      ],
      params: %{
        account_request_id: parsed_state["accountRequestId"],
        auth_server_id: parsed_state["authorisationServerId"],
        authorisation_code: authorisation_code,
        scope: parsed_state["scope"]
      }
    )
  end

  def do_get_resource(
        endpoint,
        validation_run_id,
        swagger_uris,
        parsed_state,
        config = %ApiConfig{}
      ) do
    Logger.debug(fn -> "do_get_resource: #{endpoint}" end)
    handler = fn {:ok, decoded_response} ->
      validation_result = decoded_response["validation_result"]
      json_string = Poison.encode!(validation_result)

      case Compliance.ValidationRuns.AggregateSupervisor.add_log_item(json_string) do
        {:ok, report} ->
          # Logger.info("handler -> report=#{inspect(report)}")
          report

        error ->
          Logger.info("handler -> error=#{inspect(error)}")
          nil
      end

      {:ok, decoded_response}
    end

    endpoint
    |> execute(
      handler,
      headers: validation_run_headers(validation_run_id, swagger_uris, parsed_state, config)
    )
  end

  def do_post_complete_payment(
        validation_run_id,
        swagger_uris,
        parsed_state,
        config = %ApiConfig{}
      ) do
    Logger.debug(fn -> "do_post_complete_payment" end)
    handler = fn {:ok, response} ->
      decoded_response = Poison.decode!(response)

      validation_result = decoded_response["validation_result"]
      json_string = Poison.encode!(validation_result)

      case Compliance.ValidationRuns.AggregateSupervisor.add_log_item(json_string) do
        {:ok, report} ->
          # Logger.info("handler -> report=#{inspect(report)}")
          report

        error ->
          Logger.info("handler -> error=#{inspect(error)}")
          nil
      end

      :ok
    end

    :complete_payment
    |> execute(
      handler,
      headers: validation_run_headers(validation_run_id, swagger_uris, parsed_state, config)
    )
  end

  def do_post_revoke_account_access_consent(
        validation_run_id,
        swagger_uris,
        %{
          "sessionId" => session_id,
          "authorisationServerId" => auth_server_id
        },
        config = %ApiConfig{}
      ) do
    do_post_revoke_account_access_consent(
      validation_run_id,
      swagger_uris,
      session_id,
      auth_server_id,
      config
    )
  end

  defp do_post_revoke_account_access_consent(
         validation_run_id,
         swagger_uris,
         session_token,
         auth_server_id,
         config = %ApiConfig{}
       ) do
    Logger.debug(fn -> "do_post_revoke_account_access_consent" end)
    handler = fn {:ok, _} -> :ok end

    :account_request_revoke_consent
    |> execute(
      handler,
      headers:
        authorise_headers(validation_run_id, swagger_uris, session_token, auth_server_id, config)
    )
  end

  defp execute(cmd, handler, headers: headers, params: params) do
    cmd
    |> Proxied.get(params)
    |> execute_cmd(handler, headers)
  end

  defp execute(cmd, handler, headers: headers) do
    cmd
    |> Proxied.get()
    |> execute_cmd(handler, headers)
  end

  defp execute(cmd, handler, params: params) do
    cmd
    |> Proxied.get(params)
    |> execute_cmd(handler)
  end

  defp execute_cmd(cmd, handler, optionals \\ []) do
    request_headers = headers(optionals)

    response = HTTPoison.request(cmd.method, cmd.url, cmd.body, request_headers)

    case response do
      {:ok, %HTTPoison.Response{status_code: 200, body: body}} ->
        Logger.info("status_code=200, cmd.method=#{cmd.method}, cmd.url=#{cmd.url}")
        # Logger.info("cmd.body=#{cmd.body}")
        # Logger.info("body=#{body}")
        # Logger.info("response=#{inspect(response)}")
        handle_json(body, handler)

      {:ok, %HTTPoison.Response{status_code: 201, body: body}} ->
        Logger.info("status_code=201, cmd.method=#{cmd.method}, cmd.url=#{cmd.url}")
        # Logger.info("cmd.body=#{cmd.body}")
        # Logger.info("body=#{body}")
        # Logger.info("response=#{inspect(response)}")
        handler.({:ok, body})

      {:ok, %HTTPoison.Response{status_code: 204, body: body}} ->
        Logger.info("status_code=204, cmd.method=#{cmd.method}, cmd.url=#{cmd.url}")
        # Logger.info("cmd.body=#{cmd.body}")
        # Logger.info("body=#{body}")
        # Logger.info("response=#{inspect(response)}")
        handler.({:ok, body})

      {:ok, %HTTPoison.Response{status_code: 400, body: body}} ->
        Logger.info("status_code=400, cmd.method=#{cmd.method}, cmd.url=#{cmd.url}")
        {:error, %{reason: body, cmd: cmd, status_code: 400, headers: request_headers}}

      {:ok, %HTTPoison.Response{status_code: 404, body: body}} ->
        Logger.info("status_code=404, cmd.method=#{cmd.method}, cmd.url=#{cmd.url}")
        {:error, %{reason: body, cmd: cmd, status_code: 404, headers: request_headers}}

      {:ok, %HTTPoison.Response{status_code: 500, body: body}} ->
        Logger.info("status_code=500, cmd.method=#{cmd.method}, cmd.url=#{cmd.url}")
        {:error, %{reason: body, cmd: cmd, status_code: 500, headers: request_headers}}

      {:error, %HTTPoison.Error{id: _id, reason: reason}} ->
        {:error, %{reason: reason, cmd: cmd, headers: request_headers}}
    end
  end

  defp authorise_headers(validation_run_id, swagger_uris, session_token, auth_server_id, config) do
    [
      {"authorization", session_token},
      {"x-authorization-server-id", auth_server_id},
      {"x-swagger-uris", swagger_uris},
      {"x-validation-run-id", validation_run_id},
      {"x-config", ApiConfig.base64_encode_json(config)}
    ]
  end

  defp validation_run_headers(validation_run_id, swagger_uris, parsed_state, config) do
    [
      {"authorization", parsed_state["sessionId"]},
      {"x-authorization-server-id", parsed_state["authorisationServerId"]},
      {"x-fapi-interaction-id", parsed_state["interactionId"]},
      {"x-swagger-uris", swagger_uris},
      {"x-validation-run-id", validation_run_id},
      {"x-config", ApiConfig.base64_encode_json(config)}
    ]
  end

  defp headers(optionals) do
    optionals ++
      [
        {"content-type", "application/json; charset=utf-8"},
        {"accept", "application/json; charset=utf-8"}
      ]
  end

  defp handle_json(json, handler) do
    case Poison.decode(json) do
      {:ok, value} ->
        handler.({:ok, value})

      other ->
        other
    end
  end

  defp parse_state(state) do
    state
    |> Base.decode64!()
    |> Poison.decode!()
    |> Map.update!("scope", &String.trim(String.replace(&1, "openid", "")))
  end
end
