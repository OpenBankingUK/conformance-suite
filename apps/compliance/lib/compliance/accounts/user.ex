defmodule Compliance.Accounts.User do
  @moduledoc """
  User details and token from third party authentication provider.
  """
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  schema "users" do
    field(:email, :string)
    field(:first_name, :string)
    field(:last_name, :string)
    field(:provider, :string)
    field(:token, :string)

    timestamps()
  end

  @doc false
  def changeset(user, %{provider: provider} = attrs) when is_atom(provider) do
    changeset(user, attrs |> Map.put(:provider, Atom.to_string(provider)))
  end

  @doc false
  def changeset(user, attrs) do
    user
    |> cast(attrs, [:first_name, :last_name, :email, :provider, :token])
    |> validate_required([:first_name, :last_name, :email, :provider, :token])
  end
end
