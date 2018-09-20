defmodule Compliance.Permutations.Generator do
  @moduledoc """
  Generate permutations for given endpoint definition.
  """
  require DataMorph
  require Poison
  require Logger

  @doc """
  Examples

  iex> {:ok, list} = Compliance.Permutations.Generator.endpoint_permutations("1.1")
  iex> list |> List.first
  %{
    "conditional" => false,
    "endpoint" => "/accounts",
    "optional" => false,
    "permissions" => ["ReadAccountsBasic"]
  }

  iex> Compliance.Permutations.Generator.endpoint_permutations("bad-version")
  {:error, "endpoint_permutations file not found: #{Application.app_dir(:compliance)}/priv/endpoints/permutations-bad-version.json"}
  """
  def endpoint_permutations(api_version) when is_binary(api_version) do
    app_dir = Application.app_dir(:compliance)
    file_exists = File.exists?("#{app_dir}/priv/endpoints")

    file =
      if file_exists do
        "#{app_dir}/priv/endpoints/permutations-#{api_version}.json"
      else
        "apps/compliance/priv/endpoints/permutations-#{api_version}.json"
      end

    Logger.debug(fn ->
      "Compliance.Permutations.Generator, api_version: #{inspect(api_version)}, file_exists: #{
        inspect(file_exists)
      }, file: #{inspect(file)}, app_dir: #{inspect(app_dir)}"
    end)

    with {:ok, json} <- File.read(file),
         {:ok, list} <- Poison.decode(json) do
      Logger.debug(fn -> "Compliance.Permutations.Generator, json: #{inspect(json)}" end)
      Logger.debug(fn -> "Compliance.Permutations.Generator, list: #{inspect(list)}" end)
      {:ok, list}
    else
      {:error, :enoent} ->
        msg = "endpoint_permutations file not found: #{file}"
        Logger.debug("Compliance.Permutations.Generator, msg: #{inspect(msg)}")
        {:error, msg}

      {:error, details} ->
        msg = "endpoint_permutations json not valid: #{file} - #{inspect(details)}"
        Logger.debug("Compliance.Permutations.Generator, msg: #{inspect(msg)}")
        {:error, msg}
    end
  end

  def process_file(filename) do
    list = permutations(filename)

    outfile =
      filename
      |> String.replace("endpoints-", "permutations-")
      |> String.replace("csv", "json")

    outfile
    |> File.write(Poison.encode!(list, pretty: true))

    {:ok, outfile, list |> Enum.count()}
  end

  def permutations(filename) when is_binary(filename) do
    filename
    |> File.stream!()
    |> DataMorph.maps_from_csv()
    |> Enum.map(&permutations/1)
    |> List.flatten()
  end

  def permutations(%{
        endpoint: endpoint,
        required: required,
        read: "",
        readbasic: "ReadAccountsBasic",
        readdetail: "ReadAccountsDetail",
        readcredits: "",
        readdebits: "",
        readpan: readpan
      }) do
    list = [
      permutation(endpoint, required, ["ReadAccountsBasic"]),
      permutation(endpoint, required, ["ReadAccountsDetail"])
    ]

    case readpan do
      "" -> list
      "ReadPAN" -> permutations_toggle_permission(list, readpan)
    end
  end

  def permutations(%{
        endpoint: endpoint,
        required: required,
        read: "",
        readbasic: readbasic,
        readdetail: readdetail,
        readcredits: "",
        readdebits: "",
        readpan: readpan
      }) do
    list = [
      permutation(endpoint, required, ["ReadAccountsBasic", readbasic]),
      permutation(endpoint, required, ["ReadAccountsDetail", readdetail])
    ]

    case readpan do
      "" -> list
      "ReadPAN" -> permutations_toggle_permission(list, readpan)
    end
  end

  def permutations(%{
        endpoint: endpoint,
        required: required,
        read: "",
        readbasic: readbasic,
        readdetail: readdetail,
        readcredits: readcredits,
        readdebits: readdebits,
        readpan: readpan
      }) do
    list = [
      permutation(endpoint, required, ["ReadAccountsBasic", readbasic, readcredits]),
      permutation(endpoint, required, ["ReadAccountsDetail", readdetail, readcredits]),
      permutation(endpoint, required, ["ReadAccountsBasic", readbasic, readdebits]),
      permutation(endpoint, required, ["ReadAccountsDetail", readdetail, readdebits]),
      permutation(endpoint, required, ["ReadAccountsBasic", readbasic, readcredits, readdebits]),
      permutation(endpoint, required, ["ReadAccountsDetail", readdetail, readcredits, readdebits])
    ]

    case readpan do
      "" -> list
      "ReadPAN" -> permutations_toggle_permission(list, readpan)
    end
  end

  def permutations(%{
        endpoint: endpoint,
        required: required,
        read: read,
        readbasic: "",
        readdetail: ""
      }) do
    [
      permutation(endpoint, required, ["ReadAccountsBasic", read])
    ]
  end

  def permutations(%{
        endpoint: endpoint,
        required: required,
        read: read,
        readbasic: readbasic,
        readdetail: readdetail,
        readcredits: readcredits,
        readdebits: readdebits,
        readpan: readpan
      }) do
    list = [
      permutation(endpoint, required, ["ReadAccountsBasic", read, readbasic, readcredits]),
      permutation(endpoint, required, ["ReadAccountsDetail", read, readdetail, readcredits]),
      permutation(endpoint, required, ["ReadAccountsBasic", read, readbasic, readdebits]),
      permutation(endpoint, required, ["ReadAccountsDetail", read, readdetail, readdebits]),
      permutation(endpoint, required, [
        "ReadAccountsBasic",
        read,
        readbasic,
        readcredits,
        readdebits
      ]),
      permutation(endpoint, required, [
        "ReadAccountsDetail",
        read,
        readdetail,
        readcredits,
        readdebits
      ])
    ]

    case readpan do
      "" -> list
      "ReadPAN" -> permutations_toggle_permission(list, readpan)
    end
  end

  defp permutation(endpoint, required, permissions) do
    %{
      conditional: required == "conditional",
      endpoint: endpoint,
      optional: required == "optional",
      permissions: permissions
    }
  end

  defp permutations_toggle_permission(permutations, toggle_permission) do
    permutations_with_readpan =
      permutations
      |> Enum.map(&(&1 |> Map.merge(%{permissions: &1.permissions ++ [toggle_permission]})))

    permutations ++ permutations_with_readpan
  end
end
