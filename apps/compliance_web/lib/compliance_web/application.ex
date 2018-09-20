defmodule ComplianceWeb.Application do
  @moduledoc """
  Application configuration.
  """
  use Application

  def start(_type, _args) do
    import Supervisor.Spec

    # Define workers and child supervisors to be supervised
    children = [
      # Start the endpoint when the application starts
      supervisor(ComplianceWeb.Endpoint, [])
      # Start your own worker by calling: ComplianceWeb.Worker.start_link(arg1, arg2, arg3)
      # worker(ComplianceWeb.Worker, [arg1, arg2, arg3]),
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: ComplianceWeb.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  def config_change(changed, _new, removed) do
    ComplianceWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
