defmodule Compliance.ValidationRuns.ValidationRunPayments do
  @moduledoc """
  Functions related to making payments in validation run.
  """

  alias Compliance.SwaggerUris
  alias Compliance.Commands
  alias Compliance.Configs.RunConfig
  require Logger

  def make_payments(
        payments,
        validation_run_id,
        auth_server_id,
        config = %{}
      ) do
    # params = binding()
    # Logger.debug(fn -> "Compliance.ValidationRunPayments.make_payments, #{inspect(params)}" end)
    Logger.debug(fn -> "Compliance.ValidationRunPayments.make_payments" end)

    run_config = RunConfig.from_map(config)

    payments
    |> Enum.each(&make_payment(&1, validation_run_id, run_config, auth_server_id))
  end

  defp make_payment(payment, validation_run_id, config, auth_server_id) do
    api_version = payment["api_version"]
    swagger_uris = SwaggerUris.from("payments", api_version, "generic")

    validation_run_id
    |> Commands.make_payment(
      swagger_uris,
      payment,
      auth_server_id,
      RunConfig.to_api_config(config, api_version)
    )
    |> case do
      {:error, result} ->
        Logger.error("Compliance.ValidationRunPayments.make_payments failed: #{inspect(result)}")

      _ ->
        nil
    end
  end
end
