defmodule Compliance.ValidationRuns.AggregateSupervisor do
  @moduledoc """
  DynamicSupervisor to manage validation run report aggregation processes.
  """
  alias Compliance.ValidationRuns.{Aggregate, EndpointCall}

  use Compliance.ValidationRuns.ProcessSupervisor,
    process_module: Aggregate

  @doc """
  Return existing Aggregate process for given validation_run_id, or
  when none exists start Aggregate process for given validation_run_id
  and add it to supervision.
  """
  def process_for(validation_run_id) when is_binary(validation_run_id) do
    case pid_for(validation_run_id) do
      nil -> start_process(validation_run_id: validation_run_id)
      pid -> {:ok, pid}
    end
  end

  @doc """
  Add log item JSON string to Aggregate process for given
  validation_run_id found in the log item.

  Starts Aggregate process if not already started.
  """
  def add_log_item(json_string) when is_binary(json_string) do
    with {:ok, endpoint_call} <- EndpointCall.from_json(json_string),
         {:ok, aggregate_pid} <- process_for(endpoint_call.validation_run_id) do
      Aggregate.add_log_item(aggregate_pid, endpoint_call)
    else
      {:error, msg} ->
        {:error, msg}

      other ->
        {:error, other}
    end
  end
end
