defmodule Compliance.AuthServers do
  @moduledoc false
  alias OBApiRemote.Commands

  @doc """
  Gets all available Authorisation Servers available in OB Directory.
  """
  def get_all() do
    Commands.get_auth_servers()
  end
end
