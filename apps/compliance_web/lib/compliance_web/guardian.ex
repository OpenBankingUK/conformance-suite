defmodule ComplianceWeb.Guardian do
  @moduledoc """
  Implementation of the Guardian library for app.
  """
  use Guardian, otp_app: :compliance_web
  alias Compliance.Accounts

  require Logger

  @doc """
  In our case resource is a %User{id: 'id'}.
  ## Compliances

      iex> subject_for_token(%User{}, _)
      {:ok, value}

      iex> subject_for_token(_, _)
      {:error, :reason_for_error}
  """
  def subject_for_token(resource, _claims) do
    {:ok, to_string(resource.id)}
  end

  @doc """
  ## Compliances

      iex> resource_from_claims(_)
      {:ok, %{id: _}}

      iex> resource_from_claims(_)
      {:error, :reason_for_error}
  """
  def resource_from_claims(claims) do
    id = claims["sub"]
    resource = Accounts.get_user!(id)
    {:ok, resource}
  end

  @doc """
  Return true when user identified by token is authorised to access report for
  given validation_run_id.
  """
  def authorised?(token, validation_run_id) do
    with {:ok, claims} <- decode_and_verify(token),
         {:ok, user} <- resource_from_claims(claims) do
      Accounts.has_user_validation_run?(user, validation_run_id)
    else
      error ->
        Logger.warn(inspect(error))
        false
    end
  end
end
