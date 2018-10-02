defmodule Compliance.Commands do
  @moduledoc """
  The Commands context.
  """

  require Logger

  alias Compliance.Commands.{ApiConfig, Driver}

  def request_resource(
        endpoint,
        validation_run_id,
        swagger_uris,
        permissions,
        auth_server_id,
        config = %ApiConfig{},
        revoke_consent_at_end: revoke_consent_at_end
      ) do
    authorise_function = fn session_token ->
      Driver.do_post_authorise_account_access(
        permissions,
        validation_run_id,
        swagger_uris,
        session_token,
        auth_server_id,
        config
      )
    end

    with {:ok, parsed_state} <-
           login_and_consent(authorise_function, validation_run_id, swagger_uris, config),
         {:ok, payload} <-
           Driver.do_get_resource(endpoint, validation_run_id, swagger_uris, parsed_state, config) do
      if revoke_consent_at_end do
        validation_run_id
        |> Driver.do_post_revoke_account_access_consent(swagger_uris, parsed_state, config)
        |> case do
          {:error, details} ->
            Logger.error("error revoking account access consent:
                  auth_server_id: #{auth_server_id}
                  endpoint: #{endpoint}
                  validation_run_id: #{validation_run_id}
                  details: #{inspect(details)}")

          :ok ->
            :ok
        end
      end

      {:ok, payload}
    else
      error -> error
    end
  end

  def make_payment(
        validation_run_id,
        swagger_uris,
        payment,
        auth_server_id,
        config = %ApiConfig{}
      ) do
    authorise_function = fn session_token ->
      Driver.do_post_authorise_payment(
        payment,
        validation_run_id,
        swagger_uris,
        session_token,
        auth_server_id,
        config
      )
    end

    with {:ok, parsed_state} <-
           login_and_consent(authorise_function, validation_run_id, swagger_uris, config),
         :ok <-
           Driver.do_post_complete_payment(validation_run_id, swagger_uris, parsed_state, config) do
      {:ok, validation_run_id}
    else
      error -> error
    end
  end

  defp login_and_consent(
         authorise_function,
         validation_run_id,
         swagger_uris,
         config = %ApiConfig{}
       ) do
    with {:ok, session_token} <- Driver.do_post_login(),
         {:ok, authorise_consent_url} <- authorise_function.(session_token),
         {:ok, %{"code" => auth_code, "state" => raw_state}} <-
           Driver.do_get_authorise_consent(authorise_consent_url),
         {:ok, parsed_state} <-
           Driver.do_post_consent_authorised(
             raw_state,
             auth_code,
             validation_run_id,
             swagger_uris,
             config
           ) do
      {:ok, parsed_state}
    else
      error -> error
    end
  end
end
