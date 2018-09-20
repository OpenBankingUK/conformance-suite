defmodule Compliance.ValidationRuns.ValidationRunTest do
  @moduledoc """
  Tests for ValidationRun.
  """
  use ExUnit.Case, async: true
  alias Compliance.ValidationRuns.{ValidationRun, ValidationRunSupervisor}

  @validation_run_id "validation-run-id-123"
  @scenarios []
  @config %{}

  describe "ValidationRun for given validation_run_id" do
    setup do
      pid =
        start_supervised!(
          {ValidationRun,
           [
             validation_run_id: @validation_run_id,
             config: @config,
             scenarios: @scenarios,
             auth_server_id: nil,
             listener_pids: []
           ]}
        )

      %{validation_run: pid}
    end

    test "process_for called with validation_run_id returns process pid", %{
      validation_run: pid
    } do
      assert {:ok, ^pid} = ValidationRunSupervisor.process_for(@validation_run_id)
    end

    test "process_for called with invalid validation_run_id returns error tuple" do
      assert {:error, msg} = ValidationRunSupervisor.process_for("invalid-run-id")
      assert is_binary(msg)
    end

    test "add_listener_pid called with listener pid updates listener_pids list" do
      listener_pid = self()

      assert {:ok, %{} = result} =
               ValidationRunSupervisor.add_listener_pid(@validation_run_id, listener_pid)

      assert result == %{
               validation_run_id: @validation_run_id,
               listener_pids: [listener_pid],
               auth_server_id: nil,
               config: @config,
               scenarios: @scenarios
             }
    end
  end
end
