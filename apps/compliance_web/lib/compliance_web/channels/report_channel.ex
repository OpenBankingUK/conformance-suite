defmodule ComplianceWeb.ReportChannel do
  @moduledoc """
  The ReportChannel channel allows listening to updates to individual validation runs.
  """

  require Logger
  use ComplianceWeb, :channel
  alias ComplianceWeb.Guardian
  alias Compliance.ValidationRuns.{Aggregate, AggregateSupervisor}

  def join("report:", _payload, _socket) do
    {:error, %{reason: "invalid: validation_run_id cannot be empty"}}
  end

  def join(
        "report:" <> validation_run_id,
        payload,
        %{:assigns => %{:current_user => current_user}} = socket
      ) do
    log_join(validation_run_id, payload, socket)

    if authorised?(current_user, validation_run_id) do
      case AggregateSupervisor.process_for(validation_run_id) do
        {:ok, aggregate_pid} ->
          Logger.info(["aggregate_pid: ", inspect(aggregate_pid)])

          AggregateSupervisor.add_listener_pid(validation_run_id, self())
          send(self(), {:after_join, aggregate_pid})
          {:ok, socket}

        error ->
          {:error,
           %{
             reason:
               "invalid: cannot find aggregate for #{validation_run_id}, error: #{inspect(error)}"
           }}
      end
    else
      {:error,
       %{reason: "unauthorized: you are not authorized to view the report #{validation_run_id}"}}
    end
  end

  def join("report:" <> _validation_run_id, _payload, _socket) do
    {:error, %{reason: "unauthorized: :current_user missing"}}
  end

  defp log_join(validation_run_id, payload, socket) do
    Logger.info([
      "validation_run_id: ",
      inspect(validation_run_id),
      "payload: ",
      inspect(payload),
      "socket.assigns: ",
      inspect(socket.assigns)
    ])
  end

  def terminate(reason, _socket) do
    Logger.info("reason: #{inspect(reason)}")
    :ok
  end

  # Channels can be used in a request/response fashion
  # by sending replies to requests from the client
  def handle_in("ping", payload, socket) do
    {:reply, {:ok, payload}, socket}
  end

  # It is also common to receive messages from the client and
  # broadcast to everyone in the current topic (report:lobby).
  def handle_in("shout", payload, socket) do
    broadcast(socket, "shout", payload)

    {:noreply, socket}
  end

  @doc """
  Receive notification of report_update and
  set :report_update to true in socket.assigns map.
  """
  def handle_info(:report_update, socket) do
    socket = assign(socket, :report_update, true)
    {:noreply, socket}
  end

  @doc """
  After join:
  - push endpoint reports message
  - schedule next update check
  - set :report_update false in socket.assigns map.
  """
  def handle_info({:after_join, aggregate_pid}, socket) do
    push_endpoint_reports(:after_join, aggregate_pid, socket, "started")
    schedule_update(aggregate_pid)
    {:noreply, assign(socket, :report_update, false)}
  end

  @doc """
  After started:
  - push endpoint reports message when :report_update true
  - schedule next update check
  - set :report_update false in socket.assigns map.
  """
  def handle_info({:after_started, aggregate_pid}, socket) do
    if socket.assigns.report_update do
      push_endpoint_reports(:after_started, aggregate_pid, socket, "updated")
    end

    schedule_update(aggregate_pid)
    {:noreply, assign(socket, :report_update, false)}
  end

  defp push_endpoint_reports(key, aggregate_pid, socket, type) do
    message =
      case Aggregate.endpoint_reports(aggregate_pid) do
        {:ok, %{} = reports} ->
          %{
            payload: reports
          }

        error ->
          %{
            error: error
          }
      end

    log_message(key, aggregate_pid, message)
    push(socket, type, message)
  end

  defp log_message(name, aggregate_pid, message) do
    Logger.info([
      inspect(name),
      " aggregate_pid: ",
      inspect(aggregate_pid),
      " message: ",
      inspect(message)
    ])
  end

  defp schedule_update(aggregate_pid) do
    pid = self()

    report_update_ms =
      case System.get_env("report_update_ms") do
        nil ->
          1000

        val ->
          case Integer.parse(val) do
            {ms, ""} -> ms
            _ -> 1000
          end
      end

    Process.send_after(
      pid,
      {:after_started, aggregate_pid},
      report_update_ms
    )
  end

  defp authorised?(%{:access_token => token}, validation_run_id)
       when is_binary(validation_run_id) do
    Guardian.authorised?(token, validation_run_id)
  end

  defp authorised?(_, _), do: false
end
