defmodule Compliance.ValidationRuns.Aggregate do
  @moduledoc """
  Module for function(s) to aggregate EndpointCalls into EndpointReports.
  """

  use Compliance.ValidationRuns.RunProcess

  alias Compliance.ValidationRuns.{EndpointCall, EndpointReport}

  @impl GenServer
  def init(validation_run_id: validation_run_id, listener_pids: listener_pids)
      when is_binary(validation_run_id) and is_list(listener_pids) do
    initial_state = %{
      validation_run_id: validation_run_id,
      listener_pids: listener_pids,
      endpoint_reports: %{}
    }

    {:ok, initial_state}
  end

  def endpoint_reports(aggregate_pid) do
    GenServer.call(aggregate_pid, {:endpoint_reports})
  end

  def add_log_item(aggregate_pid, %EndpointCall{} = endpoint_call) do
    GenServer.call(aggregate_pid, {:add_log_item, endpoint_call})
  end

  def validation_run_id(aggregate_pid) do
    GenServer.call(aggregate_pid, {:validation_run_id})
  end

  @impl GenServer
  def handle_call({:endpoint_reports}, _from, state) do
    {:reply, {:ok, state.endpoint_reports}, state}
  end

  @impl GenServer
  def handle_call({:add_log_item, endpoint_call = %EndpointCall{}}, _from, state) do
    state = update_in(state.endpoint_reports, &(&1 |> add(endpoint_call)))
    state.listener_pids |> Enum.each(&notify_listener/1)
    {:reply, {:ok, state}, state}
  end

  @impl GenServer
  def handle_call({:validation_run_id}, _from, state) do
    {:reply, {:ok, state.validation_run_id}, state}
  end

  # When listener process alive, send :report_update notification.
  defp notify_listener(pid) when is_pid(pid) do
    if Process.alive?(pid) do
      send(pid, :report_update)
    end
  end

  @doc """
  ## Examples

    iex> alias Compliance.ValidationRuns.EndpointCall
    iex> # add first endpoint call
    iex> reports = Aggregate.add(%{},
    iex>   %EndpointCall{
    iex>     failed_validation: false,
    iex>     path: "/open-banking/v1.1/accounts",
    iex>     request: %{"path" => "/open-banking/v1.1/accounts"},
    iex>     report: %{ "failedValidation" => false }
    iex>   }
    iex> )
    %{
      "/open-banking/v1.1/accounts" => %Compliance.ValidationRuns.EndpointReport{
        failed_calls: 0, failures: [],
        path: "/open-banking/v1.1/accounts",
        total_calls: 1
      }
    }
    iex> # add second endpoint call to same path
    iex> reports = Aggregate.add(reports,
    iex>   %EndpointCall{
    iex>     failed_validation: false,
    iex>     path: "/open-banking/v1.1/accounts",
    iex>     request: %{"path" => "/open-banking/v1.1/accounts"},
    iex>     report: %{ "failedValidation" => false }
    iex>   }
    iex> )
    %{
      "/open-banking/v1.1/accounts" => %Compliance.ValidationRuns.EndpointReport{
        failed_calls: 0, failures: [],
        path: "/open-banking/v1.1/accounts",
        total_calls: 2
      }
    }
    iex> # add endpoint call with new path
    iex> reports = Aggregate.add(reports,
    iex>   %EndpointCall{
    iex>     failed_validation: false,
    iex>     path: "/open-banking/v1.1/accounts/{AccountId}/balances",
    iex>     request: %{"path" => "/open-banking/v1.1/accounts/22290/balances"},
    iex>     report: %{ "failedValidation" => false }
    iex>   }
    iex> )
    %{
      "/open-banking/v1.1/accounts" => %Compliance.ValidationRuns.EndpointReport{
        failed_calls: 0, failures: [],
        path: "/open-banking/v1.1/accounts",
        total_calls: 2
      },
      "/open-banking/v1.1/accounts/{AccountId}/balances" => %Compliance.ValidationRuns.EndpointReport{
        failed_calls: 0, failures: [],
        path: "/open-banking/v1.1/accounts/{AccountId}/balances",
        total_calls: 1
      }
    }
    iex> # add endpoint call with validation failure
    iex> _reports = Aggregate.add(reports,
    iex>   %EndpointCall{
    iex>     failed_validation: true,
    iex>     path: "/open-banking/v1.1/accounts/{AccountId}/balances",
    iex>     request: %{"path" => "/open-banking/v1.1/accounts/22290/balances"},
    iex>     report: %{
    iex>       "failedValidation" => true,
    iex>       "message" => "Response validation failed: failed schema validation",
    iex>       "results" => %{
    iex>         "errors" => [
    iex>           %{
    iex>             "code" => "OBJECT_MISSING_REQUIRED_PROPERTY",
    iex>             "description" => "Amount of money of the cash balance.",
    iex>             "message" => "Missing required property: Amount",
    iex>             "path" => ["Data", "Balance", "0", "Amount"]
    iex>           }
    iex>         ],
    iex>         "warnings" => []
    iex>       }
    iex>     }
    iex>   }
    iex> )
    %{
      "/open-banking/v1.1/accounts" => %Compliance.ValidationRuns.EndpointReport{
        failed_calls: 0, failures: [],
        path: "/open-banking/v1.1/accounts",
        total_calls: 2
      },
      "/open-banking/v1.1/accounts/{AccountId}/balances" => %Compliance.ValidationRuns.EndpointReport{
        failed_calls: 1,
        failures: [
          %{
            "failedValidation" => true,
            "message" => "Response validation failed: failed schema validation",
            "results" => %{
              "errors" => [
                %{
                  "code" => "OBJECT_MISSING_REQUIRED_PROPERTY",
                  "description" => "Amount of money of the cash balance.",
                  "message" => "Missing required property: Amount",
                  "path" => ["Data", "Balance", "0", "Amount"]
                }
              ],
              "warnings" => []
            }
          }
        ],
        path: "/open-banking/v1.1/accounts/{AccountId}/balances",
        total_calls: 2
      }
    }
  """
  def add(%{} = reports, %EndpointCall{
        failed_validation: failed_validation,
        path: path,
        report: validation
      }) do
    update_in(reports[path], fn report ->
      report
      |> create_if_nil(path)
      |> update_failures(failed_validation, validation)
      |> update_total_calls()
    end)
  end

  defp create_if_nil(report, path) do
    case report do
      nil -> %EndpointReport{path: path}
      _ -> report
    end
  end

  defp update_failures(report, failed_validation, validation) do
    if failed_validation do
      report
      |> update_failed_calls()
      |> add_failure(validation)
    else
      report
    end
  end

  defp add_failure(report, validation) do
    update_in(report.failures, &[validation | &1])
  end

  defp update_failed_calls(report) do
    update_in(report.failed_calls, &(&1 + 1))
  end

  defp update_total_calls(report) do
    update_in(report.total_calls, &(&1 + 1))
  end
end
