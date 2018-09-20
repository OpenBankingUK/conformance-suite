defmodule Compliance.ValidationRunsTest do
  @moduledoc """
  Tests for validation runs context.
  """
  use Compliance.DataCase

  alias Compliance.Accounts
  alias Compliance.ValidationRuns
  alias Compliance.ValidationRuns.{ValidationRun, ValidationRunAccounts, ValidationRunPayments}

  import Mock

  require Logger

  test "create_user_validation_run/1 returns new validation_run_id" do
    user = user_fixture()

    assert {:ok, validation_run_id} = ValidationRuns.create_user_validation_run(user)
    assert is_binary(validation_run_id)
    assert validation_run_id != ""

    assert Accounts.has_user_validation_run?(user, validation_run_id)
  end

  @validation_run_id "test-val-run-id"
  @auth_server_id "test-auth-server-id"
  @config %{}
  @payment_scenario %{
    "api_version" => "1.1",
    "account_number" => "12345678",
    "amount" => "10.00",
    "name" => "Sam Morse",
    "sort_code" => "111111",
    "type" => "payments"
  }
  @accounts_api_version "2.0"
  @accounts_scenario %{
    "type" => "accounts",
    "api_version" => @accounts_api_version
  }

  test "do_validation_run calls make_payments for payments scenario" do
    with_mock(
      ValidationRunPayments,
      make_payments: fn _payments, _validation_run_id, _auth_server_id, _config -> :ok end
    ) do
      ValidationRun.do_validation_run(%{
        validation_run_id: @validation_run_id,
        config: @config,
        scenarios: [@payment_scenario],
        auth_server_id: @auth_server_id
      })

      assert called(
               ValidationRunPayments.make_payments(
                 [@payment_scenario],
                 @validation_run_id,
                 @auth_server_id,
                 @config
               )
             )
    end
  end

  test "do_validation_run calls request_account_resources and make_payments
    for joint accounts/payments scenarios" do
    with_mocks([
      {
        ValidationRunPayments,
        [],
        [make_payments: fn _payments, _validation_run_id, _auth_server_id, _config -> :ok end]
      },
      {
        ValidationRunAccounts,
        [],
        [request_account_resources: fn _api_version, _run_id, _auth_server_id, _config -> :ok end]
      }
    ]) do
      ValidationRun.do_validation_run(%{
        validation_run_id: @validation_run_id,
        config: @config,
        scenarios: [@payment_scenario, @accounts_scenario],
        auth_server_id: @auth_server_id
      })

      assert called(
               ValidationRunPayments.make_payments(
                 [@payment_scenario],
                 @validation_run_id,
                 @auth_server_id,
                 @config
               )
             )

      assert called(
               ValidationRunAccounts.request_account_resources(
                 @accounts_api_version,
                 @validation_run_id,
                 @auth_server_id,
                 @config
               )
             )
    end
  end
end
