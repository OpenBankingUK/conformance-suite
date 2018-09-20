defmodule Compliance.Accounts.UserValidationRun do
  @moduledoc """
  Represents a user validation run including validation_run_id.
  """

  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  schema "user_validation_runs" do
    field(:user_id, :binary_id)
    field(:validation_run_id, :string)

    timestamps()
  end

  @doc false
  def changeset(user_validation_run, attrs) do
    user_validation_run
    |> cast(attrs, [:user_id, :validation_run_id])
    |> validate_required([:user_id, :validation_run_id])
  end
end
