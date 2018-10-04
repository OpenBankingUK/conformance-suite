defmodule Compliance.ValidationRuns.AggregateSupervisorTest do
  @moduledoc """
  Tests for AggregateSupervisor.
  """
  use ExUnit.Case, async: false
  alias Compliance.ValidationRuns.AggregateSupervisor

  @validation_run_id "validation-run-id-123"

  @valid_log_item Poison.encode!(%{
                    "details" => %{"validationRunId" => @validation_run_id},
                    "request" => %{"path" => "/open-banking/v1.1/accounts/22290/balances"},
                    "report" => %{"failedValidation" => false},
                    "response" => "d"
                  })

  @expected_state %{
    validation_run_id: @validation_run_id,
    endpoint_reports: %{
      "/open-banking/v1.1/accounts/{AccountId}/balances" =>
        %Compliance.ValidationRuns.EndpointReport{
          failed_calls: 0,
          failures: [],
          path: "/open-banking/v1.1/accounts/{AccountId}/balances",
          total_calls: 1
        }
    },
    listener_pids: []
  }

  def stop_process(), do: stop_process(AggregateSupervisor.count_children())

  def stop_process(%{active: 0}), do: nil

  def stop_process(%{active: _}) do
    if AggregateSupervisor.pid_for(@validation_run_id) do
      AggregateSupervisor.stop_process(@validation_run_id)
      stop_process()
    end
  end

  describe "AggregateSupervisor" do
    setup do: on_exit(&stop_process/0)

    test "start and stop Aggregate for given validation-run-id" do
      started = AggregateSupervisor.process_for(@validation_run_id)
      stopped = AggregateSupervisor.stop_process(@validation_run_id)
      assert {:ok, pid} = started
      assert is_pid(pid)
      assert :ok == stopped
    end

    test "add_listener_pid called with listener pid updates listener_pids list" do
      # call process_for() to start process:
      AggregateSupervisor.process_for(@validation_run_id)
      listener_pid = self()

      assert {:ok, %{} = result} =
               AggregateSupervisor.add_listener_pid(@validation_run_id, listener_pid)

      assert result == %{
               validation_run_id: @validation_run_id,
               listener_pids: [listener_pid],
               endpoint_reports: %{}
             }
    end

    test "add_log_item called with valid log item notifies listener of :report_update" do
      # call process_for() to start process:
      AggregateSupervisor.process_for(@validation_run_id)
      AggregateSupervisor.add_listener_pid(@validation_run_id, self())
      AggregateSupervisor.add_log_item(@valid_log_item)
      assert_received(:report_update)
    end

    test "add_log_item called with JSON string updates state via Aggregate process for validation_run_id" do
      assert %{active: 0} = AggregateSupervisor.count_children()

      {:ok, %{} = result} = AggregateSupervisor.add_log_item(@valid_log_item)
      assert result == @expected_state
      assert %{active: 1} = AggregateSupervisor.count_children()
    end

    test "add_log_item called with invalid JSON string returns error" do
      invalid_json = @valid_log_item |> String.replace("path", "bad_field")

      assert {:error, msg} = AggregateSupervisor.add_log_item(invalid_json)
      assert msg == {:invalid, ["endpoint_call request missing path"]}
    end

    test "process_for returns existing aggregate process when it exists for given validation_run_id" do
      {:ok, existing_pid} = AggregateSupervisor.process_for(@validation_run_id)

      {:ok, pid} = AggregateSupervisor.process_for(@validation_run_id)
      assert pid == existing_pid
    end

    test "process_for starts new aggregate process when no process started for given validation_run_id" do
      {:ok, pid} = AggregateSupervisor.process_for(@validation_run_id)
      assert is_pid(pid)
    end

    test "children() returns list of child process information" do
      assert [] = AggregateSupervisor.children()
      {:ok, new_pid} = AggregateSupervisor.process_for(@validation_run_id)

      assert [
               {:undefined, pid, :worker, [Compliance.ValidationRuns.Aggregate]}
             ] = AggregateSupervisor.children()

      assert new_pid == pid
    end

    test "count_children() returns map of child counts" do
      assert %{active: 0, specs: 0, supervisors: 0, workers: 0} =
               AggregateSupervisor.count_children()

      AggregateSupervisor.process_for(@validation_run_id)

      assert %{active: 1, specs: 1, supervisors: 0, workers: 1} =
               AggregateSupervisor.count_children()
    end
  end
end
