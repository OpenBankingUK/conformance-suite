defmodule Compliance.ValidationRuns.ValidationRun do
  @moduledoc """
  Process module for asynchronously initiating validation run with
  supplied scenarios.
  """
  use Compliance.ValidationRuns.RunProcess

  alias Compliance.ValidationRuns.{
    ValidationRunAccounts,
    ValidationRunPayments
  }

  @impl GenServer
  def init(
        [
          validation_run_id: validation_run_id,
          config: config,
          scenarios: scenarios,
          auth_server_id: _auth_server_id,
          listener_pids: listener_pids
        ] = arg
      )
      when is_binary(validation_run_id) and is_map(config) and is_list(scenarios) and
             is_list(listener_pids) do
    initial_state =
      arg
      |> Enum.into(%{})

    {:ok, initial_state}
  end

  @impl GenServer
  def handle_cast({:start_asynchronous_validation_run}, state) do
    do_validation_run(state)
    {:noreply, state}
  end

  @doc """
  Do a validation run with given validation_run_id, config, scenarios.
  """
  def do_validation_run(%{
        validation_run_id: validation_run_id,
        config: config,
        scenarios: scenarios,
        auth_server_id: auth_server_id
      })
      when is_list(scenarios) do
    payment_scenarios = scenarios |> Enum.filter(&(&1["type"] == "payments"))
    account_scenarios = scenarios |> Enum.filter(&(&1["type"] == "accounts"))

    if Enum.count(payment_scenarios) > 0 do
      ValidationRunPayments.make_payments(
        payment_scenarios,
        validation_run_id,
        auth_server_id,
        config
      )
    end

    account_scenarios
    |> Enum.each(
      &ValidationRunAccounts.request_account_resources(
        &1["api_version"],
        validation_run_id,
        auth_server_id,
        config
      )
    )
  end
end
