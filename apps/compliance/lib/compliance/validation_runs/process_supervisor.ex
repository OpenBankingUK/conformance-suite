defmodule Compliance.ValidationRuns.ProcessSupervisor do
  @moduledoc """
  Implements DynamicSupervisor and adds generic functions for managing processes
  keyed by validation_run_id.

  Usage:

    use Compliance.ValidationRuns.ProcessSupervisor,
      process_module: Compliance.ValidationRuns.Aggregate

  """

  @doc """
  The __using__ callback is implemented as a macro - i.e. inside a quote block -
  as it is used to invoke code in the original module.
  """
  defmacro __using__(opts) do
    quote do
      use DynamicSupervisor

      @doc """
      Use the local name `name: __MODULE__` to ensure there can be only
      one supervisor process for this module.
      """
      def start_link(_options) do
        DynamicSupervisor.start_link(__MODULE__, :ok, name: __MODULE__)
      end

      @doc """
      Configure restart strategy.
      """
      @impl DynamicSupervisor
      def init(:ok) do
        DynamicSupervisor.init(strategy: :one_for_one)
      end

      @doc """
      Returns a list with information about all children of the supervisor.
      """
      def children do
        DynamicSupervisor.which_children(__MODULE__)
      end

      @doc """
      Returns a map containing count values for the supervisor.
      """
      def count_children do
        DynamicSupervisor.count_children(__MODULE__)
      end

      defp start_process(arg = [{:validation_run_id, validation_run_id} | _tail])
           when is_binary(validation_run_id) do
        # Default listener_pids to empty list.
        arg = Keyword.merge(arg, listener_pids: [])
        child_spec = {unquote(opts[:process_module]), arg}
        DynamicSupervisor.start_child(__MODULE__, child_spec)
      end

      @doc """
      Stop process for given validation_run_id.
      """
      def stop_process(validation_run_id) when is_binary(validation_run_id) do
        DynamicSupervisor.terminate_child(__MODULE__, pid_for(validation_run_id))
      end

      @doc """
      Add listener pid to list of listener_pids for given validation_run_id.
      """
      def add_listener_pid(validation_run_id, listener_pid)
          when is_binary(validation_run_id) and is_pid(listener_pid) do
        GenServer.call(
          pid_for(validation_run_id),
          {:add_listener_pid, listener_pid: listener_pid}
        )
      end

      @doc """
      Return PID for process with given validation_run_id.
      Otherwise return nil.
      """
      def pid_for(validation_run_id) do
        validation_run_id
        |> unquote(opts[:process_module]).via()
        |> GenServer.whereis()
      end
    end
  end
end
