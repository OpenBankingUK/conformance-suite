defmodule ComplianceWeb.ValidationRunController do
  use ComplianceWeb, :controller

  require Logger

  alias Compliance.ValidationRuns

  def create(conn, %{
        "payload" => scenarios,
        "config" => config
      })
      when is_list(scenarios) and is_map(config) do
    conn |> start_run(scenarios, config, nil)
  end

  defp start_run(conn, scenarios, config, auth_server_id) do
    with user = ComplianceWeb.Guardian.Plug.current_resource(conn),
         {:ok, validation_run_id} <- ValidationRuns.create_user_validation_run(user) do
      # Asynchronously initiate run
      ValidationRuns.start_validation_run(
        validation_run_id,
        config,
        scenarios,
        auth_server_id
      )

      conn
      |> put_status(:accepted)
      |> json(%{
        data: %{
          href: "/validation-runs/#{validation_run_id}",
          id: validation_run_id
        }
      })
    end
  end

  def show(conn, %{"id" => id} = params) do
    Logger.info("SHOW received params: #{inspect(params)}")

    hardcoded = %{
      href: "/validation-runs/#{id}",
      id: id,
      status: "PROCESSING",
      summary: %{
        payload: [
          %{
            name: "Sam Morse",
            sort_code: "111111",
            account_number: "12345678",
            amount: "10.00",
            type: "payments"
          },
          %{
            name: "Michael Burnham",
            sort_code: "222222",
            account_number: "87654321",
            amount: "200.00",
            type: "payments"
          }
        ]
      }
    }

    conn
    |> put_status(:ok)
    |> json(%{data: hardcoded})
  end
end
