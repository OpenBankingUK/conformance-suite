defmodule Compliance.Application do
  @moduledoc """
  The Compliance Application Service.

  The compliance system business domain lives in this application.

  Exposes API to clients such as the `ComplianceWeb` application
  for use in channels, controllers, and elsewhere.
  """
  use Application

  def start(_type, _args) do
    import Supervisor.Spec, warn: false

    children = [
      supervisor(Compliance.Repo, []),
      {Registry, keys: :unique, name: Registry.ValidationRuns.AggregateSupervisor},
      Compliance.ValidationRuns.AggregateSupervisor,
      {Registry, keys: :unique, name: Registry.ValidationRuns.ValidationRunSupervisor},
      Compliance.ValidationRuns.ValidationRunSupervisor
    ]

    Supervisor.start_link(
      children,
      strategy: :one_for_one,
      name: Compliance.Supervisor
    )
  end
end
