defmodule Compliance.ValidationRuns.EndpointCallTest do
  @moduledoc """
  Tests for EndPointCall.
  """
  use ExUnit.Case, async: true
  alias Compliance.ValidationRuns.EndpointCall

  describe "EndpointCall" do
    @valid_json %{
      "details" => %{"validationRunId" => "test-run-id"},
      "request" => %{"path" => "/open-banking/v1.1/accounts/22290/balances"},
      "report" => %{"failedValidation" => false},
      "response" => "d"
    }

    @missing_path_json update_in(@valid_json["request"], fn _ -> %{} end)

    @missing_failed_validation_json update_in(@valid_json["report"], fn _ -> %{} end)

    @missing_validation_run_id_json update_in(@valid_json["details"], fn _ -> %{} end)

    test "from_json/1 with valid JSON creates a struct with generic path property and validation_run_id" do
      assert {:ok, %EndpointCall{} = endpoint_call} =
               EndpointCall.from_json(Poison.encode!(@valid_json))

      assert endpoint_call.details == @valid_json["details"]
      assert endpoint_call.request == @valid_json["request"]
      assert endpoint_call.report == @valid_json["report"]
      assert endpoint_call.response == @valid_json["response"]
      assert endpoint_call.path == "/open-banking/v1.1/accounts/{AccountId}/balances"
      assert endpoint_call.validation_run_id == @valid_json["details"]["validationRunId"]
      assert endpoint_call.failed_validation == @valid_json["report"]["failedValidation"]
    end

    test "from_json/1 with invalid JSON returns error" do
      assert {:error, msg, _} = EndpointCall.from_json("{")
      assert msg == :invalid
    end

    test "from_json/1 with JSON request missing a path" do
      assert {:error, {:invalid, msgs}} =
               EndpointCall.from_json(Poison.encode!(@missing_path_json))

      assert msgs == ["endpoint_call request missing path"]
    end

    test "from_json/1 with JSON report missing failedValidation boolean" do
      assert {:error, {:invalid, msgs}} =
               EndpointCall.from_json(Poison.encode!(@missing_failed_validation_json))

      assert msgs == ["endpoint_call report missing failedValidation boolean"]
    end

    test "from_json/1 with JSON report missing validationRunId" do
      assert {:error, {:invalid, msgs}} =
               EndpointCall.from_json(Poison.encode!(@missing_validation_run_id_json))

      assert msgs == ["endpoint_call details missing validationRunId"]
    end
  end

  describe "EndpointCall.generic_path" do
    test "substitutes AccountRequestId correctly" do
      assert EndpointCall.generic_path("/account-requests/abc123") ==
               "/account-requests/{AccountRequestId}"
    end

    test "substitutes AccountId correctly" do
      assert EndpointCall.generic_path("/accounts/abc123") == "/accounts/{AccountId}"

      assert EndpointCall.generic_path("/accounts/456xyz/transactions") ==
               "/accounts/{AccountId}/transactions"
    end

    test "substitutes StatementId correctly" do
      assert EndpointCall.generic_path("/accounts/abc123/statements/xyz456") ==
               "/accounts/{AccountId}/statements/{StatementId}"

      assert EndpointCall.generic_path("/accounts/xyz456/statements/abc123/transactions") ==
               "/accounts/{AccountId}/statements/{StatementId}/transactions"
    end

    test "substitutes PaymentId correctly" do
      assert EndpointCall.generic_path("/payments/xyz456") == "/payments/{PaymentId}"
    end

    test "substitutes PaymentSubmissionId correctly" do
      assert EndpointCall.generic_path("/payment-submissions/abc123") ==
               "/payment-submissions/{PaymentSubmissionId}"
    end
  end
end
