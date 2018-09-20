defmodule ComplianceWeb.ValidationRunControllerTest do
  @moduledoc false
  use ComplianceWeb.ConnCase
  import Mock

  alias Compliance.ValidationRuns

  @config %{
    "client_id" => "testClientId",
    "client_secret" => "",
    "resource_endpoint" => "http://example.com"
  }
  @validation_run_id "validation-run-id-uuid"

  @payments_api_version "1.1"

  @payment_scenario %{
    "api_version" => @payments_api_version,
    "account_number" => "12345678",
    "amount" => "10.00",
    "name" => "Sam Morse",
    "sort_code" => "111111",
    "type" => "payments"
  }

  @payments_params %{
    "config" => @config,
    "payload" => [@payment_scenario]
  }

  @accounts_api_version "2.0"

  @accounts_scenario %{
    "type" => "accounts",
    "api_version" => @accounts_api_version
  }

  @accounts_params %{
    "config" => @config,
    "payload" => [@accounts_scenario]
  }

  @accounts_plus_payments_params %{
    "config" => @config,
    "payload" => [@accounts_scenario, @payment_scenario]
  }

  describe "not authenticated" do
    @tag capture_log: true
    test "create/2 responds with 401", %{conn: conn} do
      conn = post(conn, "/validation-runs")
      assert response(conn, 401)
      assert conn.halted
    end

    @tag capture_log: true
    test "show/2 responds with 401", %{conn: conn} do
      conn = get(conn, "/validation-runs/2130040")
      assert response(conn, 401)
      assert conn.halted
    end
  end

  describe "authenticated but hardcoded" do
    test "create/2 initiates request account resources validation run for accounts payload", %{
      conn: conn
    } do
      with_mock(
        ValidationRuns,
        create_user_validation_run: fn _user -> {:ok, @validation_run_id} end,
        start_validation_run: fn _run_id, _config, _scenarios, _auth_server_id -> :ok end
      ) do
        conn
        |> simulate_authenticated
        |> post("/validation-runs", @accounts_params)

        assert called(
                 ValidationRuns.start_validation_run(
                   @validation_run_id,
                   @config,
                   [@accounts_scenario],
                   nil
                 )
               )
      end
    end

    test "create/2 initiates make payments validation run for payments payload", %{conn: conn} do
      with_mock(
        ValidationRuns,
        create_user_validation_run: fn _user -> {:ok, @validation_run_id} end,
        start_validation_run: fn _run_id, _config, _scenarios, _auth_server_id -> :ok end
      ) do
        conn
        |> simulate_authenticated
        |> post("/validation-runs", @payments_params)

        assert called(
                 ValidationRuns.start_validation_run(
                   @validation_run_id,
                   @config,
                   [@payment_scenario],
                   nil
                 )
               )
      end
    end

    test "create/2 initiates request account resources validation run and initiates make payments validation run for joint accounts/payments payload",
         %{
           conn: conn
         } do
      with_mock(
        ValidationRuns,
        create_user_validation_run: fn _user -> {:ok, @validation_run_id} end,
        start_validation_run: fn _run_id, _config, _scenarios, _auth_server_id -> :ok end
      ) do
        conn
        |> simulate_authenticated
        |> post("/validation-runs", @accounts_plus_payments_params)

        assert called(
                 ValidationRuns.start_validation_run(
                   @validation_run_id,
                   @config,
                   [@accounts_scenario, @payment_scenario],
                   nil
                 )
               )
      end
    end

    test "create/2 responds with http status accepted and json payload", %{conn: conn} do
      with_mock(
        ValidationRuns,
        create_user_validation_run: fn _user -> {:ok, @validation_run_id} end,
        start_validation_run: fn _run_id, _config, _scenarios, _auth_server_id -> :ok end
      ) do
        response =
          conn
          |> simulate_authenticated
          |> post("/validation-runs", @payments_params)

        expected = %{
          "data" => %{
            "href" => "/validation-runs/#{@validation_run_id}",
            "id" => @validation_run_id
          }
        }

        body = json_response(response, :accepted)
        assert body == expected
      end
    end

    test "show/2 responds with http status 200 and json payload", %{conn: conn} do
      response =
        conn
        |> simulate_authenticated
        |> get("/validation-runs/#{@validation_run_id}")

      expected = %{
        "data" => %{
          "href" => "/validation-runs/#{@validation_run_id}",
          "id" => @validation_run_id,
          "status" => "PROCESSING",
          "summary" => %{
            "payload" => [
              %{
                "name" => "Sam Morse",
                "sort_code" => "111111",
                "account_number" => "12345678",
                "amount" => "10.00",
                "type" => "payments"
              },
              %{
                "name" => "Michael Burnham",
                "sort_code" => "222222",
                "account_number" => "87654321",
                "amount" => "200.00",
                "type" => "payments"
              }
            ]
          }
        }
      }

      body = json_response(response, :ok)
      assert body == expected
    end
  end
end
