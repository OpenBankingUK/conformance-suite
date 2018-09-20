defmodule Compliance.AuthServersTest do
  @moduledoc """
  Tests for auth servers context.
  """
  use ExUnit.Case, async: true

  alias Compliance.AuthServers
  alias OBApiRemote.Commands

  import Mock

  @authorisation_servers [
    %{
      "id" => "aaaj4NmBD8lQxmLh2O",
      "logoUri" => "",
      "name" => "AAA Example Bank"
    },
    %{
      "id" => "bbbX7tUB4fPIYB0k1m",
      "logoUri" => "",
      "name" => "BBB Example Bank"
    },
    %{
      "id" => "cccbN8iAsMh74sOXhk",
      "logoUri" => "",
      "name" => "CCC Example Bank"
    }
  ]

  describe "AuthServers.get_all" do
    test "returns all authorisation servers" do
      with_mocks([
        {
          Commands,
          [],
          [
            get_auth_servers: fn -> {:ok, @authorisation_servers} end
          ]
        }
      ]) do
        {:ok, auth_servers} = AuthServers.get_all()
        assert(auth_servers == @authorisation_servers)
        assert called(Commands.get_auth_servers())
      end
    end
  end
end
