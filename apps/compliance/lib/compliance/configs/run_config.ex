defmodule Compliance.Configs.RunConfig do
  @moduledoc """
  Represents configuration for a validation run.

  Has same keys as Compliance.Commands.ApiConfig, apart from api_version.

  The Acccounts API and Payments API can be at different version numbers.
  So the api_version is not set a the global RunConfig level.
  """
  alias Compliance.Commands.ApiConfig
  use Ecto.Schema
  import Ecto.Changeset

  @config_map %ApiConfig{} |> Map.from_struct()
  @config_keys @config_map |> Map.keys() |> List.delete(:api_version)

  @primary_key false
  embedded_schema do
    @config_keys
    |> Enum.each(&field(&1, :string))
  end

  @doc false
  def changeset(user_validation_run, attrs) do
    user_validation_run
    |> cast(attrs, @config_keys)
    |> validate_required(@config_keys)
  end

  defp string_keys do
    %__MODULE__{}
    |> Map.keys()
    |> Enum.map(&Atom.to_string/1)
  end

  @doc """
  Returns RunConfig prepopulated with values from an openid config map.
  """
  def from_openid_config(data) when is_map(data) do
    data
    |> Map.take(string_keys())
    |> Map.put("token_endpoint_auth_method", "private_key_jwt")
    |> from_map()
  end

  @doc false
  def from_map(data) when is_map(data) do
    %__MODULE__{}
    |> cast(data, @config_keys)
    |> apply_changes
  end

  @doc """
  Creates ApiConfig from supplied RunConfig and api_version e.g. "1.1".
  """
  def to_api_config(config = %__MODULE__{}, api_version) when is_binary(api_version) do
    map =
      config
      |> Map.from_struct()
      |> Map.put(:api_version, api_version)

    struct(ApiConfig, map)
  end
end
