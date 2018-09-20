defmodule Compliance.SwaggerUris do
  @moduledoc """
  Configuration of Swagger URIs for various API versions and permission levels.
  """

  @uris %{
    "accounts" => %{
      "1.1" => %{
        "generic" =>
          "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger.json",
        "basic" =>
          "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger-basic.json",
        "detail" =>
          "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger-detail.json"
      },
      "2.0" => %{
        "generic" =>
          "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v2.0.0/account-info-swagger.json",
        "basic" =>
          "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v2.0.0/account-info-swagger-basic.json",
        "detail" =>
          "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v2.0.0/account-info-swagger-detail.json"
      },
      "3.0.0" => %{
        "generic" =>
          "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json",
        "basic" =>
          "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-361-permission-specific-accounts-v3-swagger/dist/v3.0.0/account-info-swagger-basic.json",
        "detail" =>
          "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-361-permission-specific-accounts-v3-swagger/dist/v3.0.0/account-info-swagger-detail.json"
      }
    },
    "payments" => %{
      "1.1" => %{
        "generic" =>
          "https://raw.githubusercontent.com/OpenBankingUK/payment-initiation-api-spec/master/dist/v1.1/payment-initiation-swagger.json"
      },
      "3.0.0" => %{
        "generic" =>
          "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/payment-initiation-swagger.json"
      }
    }
  }

  @doc """
  Get a specific swagger URI for given type, api_version, level.

  ## Example

    iex> Compliance.SwaggerUris.from("accounts", "1.1", "generic")
    "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger.json"
  """
  def from(type, api_version, level)
      when is_binary(type) and is_binary(api_version) and is_binary(level) do
    @uris[type][api_version][level]
  end

  @separator " "

  @doc """
  Get a space separated string of swagger URIS for give type, api_version,
  and permissions list.

  ## Examples

    iex> Compliance.SwaggerUris.from("accounts", "1.1", ["ReadBalances"])
    "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger.json"

    iex> Compliance.SwaggerUris.from("accounts", "1.1", ["ReadAccountsBasic"])
    "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger.json https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger-basic.json"

    iex> Compliance.SwaggerUris.from("accounts", "1.1", ["ReadAccountsDetail"])
    "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger.json https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger-detail.json"

    iex> Compliance.SwaggerUris.from("accounts", "1.1",
    iex>   [
    iex>     "ReadStatementsDetail",
    iex>     "ReadTransactionsBasic",
    iex>     "ReadTransactionsCredits"
    iex>   ]
    iex> )
    "https://raw.githubusercontent.com/OpenBankingUK/account-info-api-spec/refapp-295-permission-specific-swagger-files/dist/v1.1.1/account-info-swagger.json"
  """
  def from(type, api_version, permissions) when is_list(permissions) do
    case [has_basic(permissions), has_detail(permissions)] do
      [true, false] ->
        [
          from(type, api_version, "generic"),
          from(type, api_version, "basic")
        ]
        |> Enum.join(@separator)

      [false, true] ->
        [
          from(type, api_version, "generic"),
          from(type, api_version, "detail")
        ]
        |> Enum.join(@separator)

      _ ->
        from(type, api_version, "generic")
    end
  end

  defp has_basic(permissions) do
    permissions |> Enum.any?(&(&1 |> String.ends_with?("Basic")))
  end

  defp has_detail(permissions) do
    permissions |> Enum.any?(&(&1 |> String.ends_with?("Detail")))
  end
end
