defmodule Compliance.ValidationRuns.AggregateTest do
  @moduledoc """
  Tests for Aggregate.
  """
  use ExUnit.Case, async: false
  alias Compliance.ValidationRuns.{Aggregate, EndpointCall}
  doctest Aggregate

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

  describe "Aggregate for given validation_run_id" do
    setup do
      aggregate =
        start_supervised!({Aggregate, [validation_run_id: @validation_run_id, listener_pids: []]})

      %{aggregate: aggregate}
    end

    test "starts with validation_run_id set", %{aggregate: aggregate} do
      assert {:ok, @validation_run_id} = Aggregate.validation_run_id(aggregate)
    end

    test "starts with empty endpoint_reports", %{aggregate: aggregate} do
      assert {:ok, %{} = reports} = Aggregate.endpoint_reports(aggregate)
      assert reports == %{}
    end

    test "add_log_item called with pid, EndpointCall updates state", %{aggregate: aggregate} do
      {:ok, endpoint_call} = EndpointCall.from_json(@valid_log_item)
      assert {:ok, %{} = result} = Aggregate.add_log_item(aggregate, endpoint_call)
      assert result == @expected_state
    end
  end
end
