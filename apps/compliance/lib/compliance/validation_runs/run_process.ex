defmodule Compliance.ValidationRuns.RunProcess do
  @moduledoc """
  Module for function(s) common to GenServer processes keyed by
  validation_run_id.
  """

  defmacro __using__(_opts) do
    quote do
      use GenServer

      @doc """
      Returns :via tuple used to register and access processes via Registry.
      """
      def via(validation_run_id) do
        registry_module =
          Module.concat(
            Registry,
            "#{__MODULE__}Supervisor" |> String.replace("Compliance.", "")
          )

        {:via, Registry, {registry_module, validation_run_id}}
      end

      def start_link(arg = [{:validation_run_id, validation_run_id} | _tail])
          when is_binary(validation_run_id) do
        GenServer.start_link(__MODULE__, arg, name: via(validation_run_id))
      end

      @impl GenServer
      def handle_call({:add_listener_pid, listener_pid: listener_pid}, _from, state)
          when is_pid(listener_pid) do
        state = update_in(state.listener_pids, &(&1 ++ [listener_pid]))
        {:reply, {:ok, state}, state}
      end
    end
  end
end
