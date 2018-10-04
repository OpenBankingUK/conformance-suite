defmodule Compliance.Permutations.GeneratorTest do
  @moduledoc """
  Tests for Permutations.Generator.
  """
  use ExUnit.Case, async: false
  doctest Compliance.Permutations.Generator
  alias Compliance.Permutations.Generator

  def endpoint_data(overrides = %{}) do
    %{
      endpoint: "",
      required: "mandatory",
      read: "",
      readbasic: "",
      readdetail: "",
      readcredits: "",
      readdebits: "",
      readpan: ""
    }
    |> Map.merge(overrides)
  end

  describe "Generator" do
    test "given optional endpoint with single read permission,
      generates single permutation, including ReadAccountsBasic in permissions list" do
      endpoint =
        endpoint_data(%{endpoint: "/balances", required: "optional", read: "ReadBalances"})

      assert Generator.permutations(endpoint) == [
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: true,
                 permissions: ["ReadAccountsBasic", endpoint[:read]]
               }
             ]
    end

    test "given endpoint with basic/detail read permissions, generates two permutations" do
      endpoint =
        endpoint_data(%{
          endpoint: "/accounts",
          readbasic: "ReadAccountsBasic",
          readdetail: "ReadAccountsDetail"
        })

      assert Generator.permutations(endpoint) == [
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: [endpoint[:readbasic]]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: [endpoint[:readdetail]]
               }
             ]
    end

    test "given endpoint with non-ReadAccounts basic/detail read permissions,
      generates two permutations adding ReadAccounts* permission" do
      endpoint =
        endpoint_data(%{
          endpoint: "/beneficiaries",
          readbasic: "ReadBeneficiariesBasic",
          readdetail: "ReadBeneficiariesDetail"
        })

      assert Generator.permutations(endpoint) == [
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: ["ReadAccountsBasic", endpoint[:readbasic]]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: ["ReadAccountsDetail", endpoint[:readdetail]]
               }
             ]
    end

    test "given endpoint with basic/detail read permissions and read PAN permission, generates four permutations" do
      endpoint =
        endpoint_data(%{
          endpoint: "/accounts",
          readbasic: "ReadAccountsBasic",
          readdetail: "ReadAccountsDetail",
          readpan: "ReadPAN"
        })

      assert Generator.permutations(endpoint) == [
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: [endpoint[:readbasic]]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: [endpoint[:readdetail]]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: [endpoint[:readbasic], endpoint[:readpan]]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: [endpoint[:readdetail], endpoint[:readpan]]
               }
             ]
    end

    test "given endpoint with debit/credit, basic/detail read permissions,
      generates six permutations adding ReadAccounts* permission" do
      endpoint =
        endpoint_data(%{
          endpoint: "/accounts/{AccountId}/transactions",
          readbasic: "ReadTransactionsBasic",
          readdetail: "ReadTransactionsDetail",
          readcredits: "ReadTransactionsCredits",
          readdebits: "ReadTransactionsDebits"
        })

      assert Generator.permutations(endpoint) == [
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: ["ReadAccountsBasic", endpoint[:readbasic], endpoint[:readcredits]]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: [
                   "ReadAccountsDetail",
                   endpoint[:readdetail],
                   endpoint[:readcredits]
                 ]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: ["ReadAccountsBasic", endpoint[:readbasic], endpoint[:readdebits]]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: ["ReadAccountsDetail", endpoint[:readdetail], endpoint[:readdebits]]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: [
                   "ReadAccountsBasic",
                   endpoint[:readbasic],
                   endpoint[:readcredits],
                   endpoint[:readdebits]
                 ]
               },
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: false,
                 permissions: [
                   "ReadAccountsDetail",
                   endpoint[:readdetail],
                   endpoint[:readcredits],
                   endpoint[:readdebits]
                 ]
               }
             ]
    end

    test "given optional endpoint with debit/credit, basic/detail read permissions, and readpan permission, generates 12 permutations" do
      endpoint =
        endpoint_data(%{
          endpoint: "/transactions",
          required: "optional",
          readbasic: "ReadTransactionsBasic",
          readdetail: "ReadTransactionsDetail",
          readcredits: "ReadTransactionsCredits",
          readdebits: "ReadTransactionsDebits",
          readpan: "ReadPAN"
        })

      permutations = Generator.permutations(endpoint)
      assert Enum.count(permutations) == 12

      assert Enum.at(permutations, 5) ==
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: true,
                 permissions: [
                   "ReadAccountsDetail",
                   endpoint[:readdetail],
                   endpoint[:readcredits],
                   endpoint[:readdebits]
                 ]
               }

      assert List.last(permutations) ==
               %{
                 conditional: false,
                 endpoint: endpoint[:endpoint],
                 optional: true,
                 permissions: [
                   "ReadAccountsDetail",
                   endpoint[:readdetail],
                   endpoint[:readcredits],
                   endpoint[:readdebits],
                   endpoint[:readpan]
                 ]
               }
    end
  end

  test "given conditional endpoint with generic read, debit/credit,
    basic/detail read permissions, and readpan permission, generates 12 permutations" do
    endpoint =
      endpoint_data(%{
        endpoint: "/accounts/{AccountId}/statements/{StatementId}/transactions",
        required: "conditional",
        read: "ReadStatementsDetail",
        readbasic: "ReadTransactionsBasic",
        readcredits: "ReadTransactionsCredits",
        readdebits: "ReadTransactionsDebits",
        readdetail: "ReadTransactionsDetail",
        readpan: "ReadPAN"
      })

    permutations = Generator.permutations(endpoint)
    assert Enum.count(permutations) == 12

    assert Enum.at(permutations, 5) ==
             %{
               conditional: true,
               endpoint: endpoint[:endpoint],
               optional: false,
               permissions: [
                 "ReadAccountsDetail",
                 endpoint[:read],
                 endpoint[:readdetail],
                 endpoint[:readcredits],
                 endpoint[:readdebits]
               ]
             }

    assert List.last(permutations) ==
             %{
               conditional: true,
               endpoint: endpoint[:endpoint],
               optional: false,
               permissions: [
                 "ReadAccountsDetail",
                 endpoint[:read],
                 endpoint[:readdetail],
                 endpoint[:readcredits],
                 endpoint[:readdebits],
                 endpoint[:readpan]
               ]
             }
  end
end
