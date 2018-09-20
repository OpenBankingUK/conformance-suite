defmodule Compliance.ValidationRuns.EndpointReport do
  @moduledoc """
  Struct to hold data summarising endpoint calls to a particular path.
  """
  @enforce_keys [:path]
  defstruct path: nil, total_calls: 0, failed_calls: 0, failures: []
end
