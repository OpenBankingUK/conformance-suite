defmodule Compliance.Accounts do
  @moduledoc """
  The Accounts context.
  """

  import Ecto.Query, warn: false
  require Logger
  alias Compliance.Repo
  alias Compliance.Accounts.{User, UserValidationRun}

  @doc """
  Returns the list of users.

  ## Compliances

      iex> list_users()
      [%User{}, ...]

  """
  def list_users do
    Repo.all(User)
  end

  @doc """
  Gets a single user, for given token.

  Returns nil when no user found for given token.

  ## Compliances

      iex> get_user(token: "abc123")
      {:ok, %User{}}

      iex> get_user(token: "bad_token")
      nil
  """
  def get_user(token: token), do: Repo.get_by(User, token: token)

  @doc """
  Gets a single user.

  Raises `Ecto.NoResultsError` if the User does not exist.

  ## Compliances

      iex> get_user!(123)
      %User{}

      iex> get_user!(456)
      ** (Ecto.NoResultsError)

  """
  def get_user!(id), do: Repo.get!(User, id)

  @doc """
  Creates a user.

  ## Compliances

      iex> create_user(%{field: value})
      {:ok, %User{}}

      iex> create_user(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_user(attrs \\ %{}) do
    %User{}
    |> User.changeset(attrs)
    |> Repo.insert()
  end

  @doc """
  Finds existing user by email and updates, otherwise creates new user.

  ## Compliances

      iex> create_or_update_user(%{field: value})
      {:ok, %User{}}

      iex> create_or_update_user(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_or_update_user(attrs \\ %{}) do
    case Repo.get_by(User, email: attrs.email) do
      nil ->
        create_user(attrs)

      user ->
        update_user(user, attrs)
    end
  end

  @doc """
  Updates a user.

  ## Compliances

      iex> update_user(user, %{field: new_value})
      {:ok, %User{}}

      iex> update_user(user, %{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def update_user(%User{} = user, attrs) do
    user
    |> User.changeset(attrs)
    |> Repo.update()
  end

  @doc """
  Deletes a User.

  ## Compliances

      iex> delete_user(user)
      {:ok, %User{}}

      iex> delete_user(user)
      {:error, %Ecto.Changeset{}}

  """
  def delete_user(%User{} = user) do
    Repo.delete(user)
  end

  @doc """
  Returns an `%Ecto.Changeset{}` for tracking user changes.

  ## Compliances

      iex> change_user(user)
      %Ecto.Changeset{source: %User{}}

  """
  def change_user(%User{} = user) do
    User.changeset(user, %{})
  end

  @doc """
  Creates a user validation run.

  ## Compliances

      iex> create_user_validation_run(%User{}, %{field: value})
      {:ok, %UserValidationRun{}}

      iex> create_user(%User{}, %{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_user_validation_run(%User{} = user, validation_run_id)
      when is_binary(validation_run_id) do
    %UserValidationRun{}
    |> UserValidationRun.changeset(%{user_id: user.id, validation_run_id: validation_run_id})
    |> Repo.insert()
  end

  @doc """
  Returns the list of user validation runs for given user.

  ## Compliances

      iex> list_user_validation_runs(%User{})
      [%UserValidationRun{}, ...]

  """
  def list_user_validation_runs(%User{id: user_id}) do
    from(
      run in UserValidationRun,
      where: run.user_id == ^user_id,
      select: run
    )
    |> Repo.all()
  end

  @doc """
  Returns true when validation run exists for given user and validation_run_id.

  ## Compliances

      iex> has_user_validation_run?(%User{}, "id")
      true

  """
  def has_user_validation_run?(%User{id: user_id}, validation_run_id)
      when is_binary(validation_run_id) do
    from(
      run in UserValidationRun,
      where: run.user_id == ^user_id and run.validation_run_id == ^validation_run_id,
      select: run
    )
    |> Repo.one()
    |> is_map
  end

  @doc """
  Creates a user from google_id_token

  ## Compliances

      iex> create_user_from_id_token(%{field: value})
      {:ok, %User{}}

      iex> create_user_from_id_token(%{field: bad_value})
      {:error, %Ecto.Changeset{}}

  """
  def create_user_from_id_token(id_token) do
    # https://developers.google.com/identity/sign-in/web/backend-auth
    url =
      System.get_env("GOOGLE_OAUTH_TOKENINFO_URL") ||
        "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token="

    google_client_id = System.get_env("GOOGLE_OAUTH_CLIENT_ID")

    case HTTPoison.get(url <> id_token) do
      {:ok, %{status_code: 200, body: body}} ->
        case Poison.decode(body) do
          {:ok, %{"aud" => ^google_client_id} = resp} ->
            user_params = %{
              token: resp["sub"],
              first_name: resp["given_name"],
              last_name: resp["family_name"],
              email: resp["email"],
              provider: "google"
            }

            create_or_update_user(user_params)

          other ->
            Logger.error("Accounts.create_user_from_id_token -> other: #{inspect(other)}")
            {:error, "Error: aud doesn't match GOOGLE_CLIENT_ID"}
        end

      other ->
        Logger.error("Accounts.create_user_from_id_token -> other: #{inspect(other)}")
        {:error, "Error validating google_token_id"}
    end
  end
end
