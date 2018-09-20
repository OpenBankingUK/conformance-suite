defmodule Compliance.ValidationRuns do
  @moduledoc """
  Functions related to creating and initiating validation runs.
  """

  alias Compliance.Accounts
  alias Compliance.Accounts.User

  alias Compliance.ValidationRuns.ValidationRunSupervisor

  require Logger

  @doc """
  Creates a user validation run and returns new validation_run_id.

  ## Compliances

      iex> create_user_validation_run(%User{})
      {:ok, validation_run_id}

      iex> create_user_validation_run(%User{})
      {:error, %Ecto.Changeset{}}

  """
  def create_user_validation_run(user = %User{}) do
    validation_run_id = UUID.uuid1()

    case Accounts.create_user_validation_run(user, validation_run_id) do
      {:ok, %{validation_run_id: validation_run_id}} ->
        {:ok, validation_run_id}

      {:error, details} ->
        {:error, details}
    end
  end

  defdelegate start_validation_run(
                validation_run_id,
                config,
                scenarios,
                auth_server_id
              ),
              to: ValidationRunSupervisor,
              as: :initiate_validation_run
end
