defmodule Compliance.ValidationRuns.ValidationRunSupervisor do
  @moduledoc """
  DynamicSupervisor to manage validation run report aggregation processes.
  """
  alias Compliance.ValidationRuns.ValidationRun

  require Logger

  use Compliance.ValidationRuns.ProcessSupervisor,
    process_module: ValidationRun

  @doc """
  Return existing ValidationRun process for given validation_run_id, or
  when none exists return error.
  """
  def process_for(validation_run_id) when is_binary(validation_run_id) do
    case pid_for(validation_run_id) do
      nil ->
        {:error, "ValidationRun process not found for validation_run_id: #{validation_run_id}"}

      pid ->
        {:ok, pid}
    end
  end

  @doc """
  Asynchronously initiate validation run for given
  validation_run_id, config, scenarios, auth_server_id.
  Starts validation run process.
  Cast start_asynchronous_validation_run on process.
  """
  def initiate_validation_run(validation_run_id, config, scenarios, auth_server_id)
      when is_binary(validation_run_id) and is_list(scenarios) do
    case start_process(
           validation_run_id: validation_run_id,
           config: config,
           scenarios: scenarios,
           auth_server_id: auth_server_id
         ) do
      {:ok, pid} ->
        GenServer.cast(pid, {:start_asynchronous_validation_run})
        :ok

      other ->
        Logger.error("ValidationRunSupervisor: initiate_validation_run: #{inspect(other)}")
        other
    end
  end
end
