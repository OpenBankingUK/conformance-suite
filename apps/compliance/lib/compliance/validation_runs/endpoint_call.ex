defmodule Compliance.ValidationRuns.EndpointCall do
  @moduledoc """
  Struct to hold data from endpoint call log items.
  """
  alias __MODULE__

  @derive [Poison.Encoder]
  defstruct [
    :details,
    :failed_validation,
    :path,
    :report,
    :request,
    :response,
    :validation_run_id
  ]

  @account_id_pattern Regex.compile!("(/accounts)/([^/]+)(/.+)?$")
  @account_request_id_pattern Regex.compile!("(/account-requests)/(.+)$")
  @payment_id_pattern Regex.compile!("(/payments)/([^/]+)$")
  @payment_submission_id_pattern Regex.compile!("(/payment-submissions)/([^/]+)$")
  @statement_id_pattern Regex.compile!("(.+/statements)/([^/]+)(/.+)?$")

  def generic_path(path) do
    path
    |> String.replace(@account_id_pattern, "\\1/{AccountId}\\3")
    |> String.replace(@account_request_id_pattern, "\\1/{AccountRequestId}")
    |> String.replace(@payment_id_pattern, "\\1/{PaymentId}\\3")
    |> String.replace(@payment_submission_id_pattern, "\\1/{PaymentSubmissionId}\\3")
    |> String.replace(@statement_id_pattern, "\\1/{StatementId}\\3")
  end

  def from_json(string) do
    case Poison.decode(string, as: %EndpointCall{}) do
      {:ok, endpoint_call} ->
        with %EndpointCall{details: %{"validationRunId" => validation_run_id}} <- endpoint_call,
             %EndpointCall{request: %{"path" => path}} <- endpoint_call,
             %EndpointCall{report: %{"failedValidation" => failed_validation}} <- endpoint_call do
          endpoint_call =
            endpoint_call
            |> set_failed_validation(failed_validation)
            |> set_path(path)
            |> set_validation_run_id(validation_run_id)

          {:ok, endpoint_call}
        else
          _ ->
            {:error, {:invalid, error_messages(endpoint_call)}}
        end

      other ->
        other
    end
  end

  defp set_failed_validation(endpoint_call, failed_validation) do
    update_in(endpoint_call.failed_validation, fn _ -> failed_validation end)
  end

  defp set_path(endpoint_call, path) do
    update_in(endpoint_call.path, fn _ -> generic_path(path) end)
  end

  defp set_validation_run_id(endpoint_call, validation_run_id) do
    update_in(endpoint_call.validation_run_id, fn _ -> validation_run_id end)
  end

  defp add_error_on_nil(messages, nil, message), do: [message | messages]
  defp add_error_on_nil(messages, _, _message), do: messages

  defp error_messages(endpoint_call) do
    []
    |> add_error_on_nil(
      endpoint_call.details["validationRunId"],
      "endpoint_call details missing validationRunId"
    )
    |> add_error_on_nil(
      endpoint_call.request["path"],
      "endpoint_call request missing path"
    )
    |> add_error_on_nil(
      endpoint_call.report["failedValidation"],
      "endpoint_call report missing failedValidation boolean"
    )
  end
end
