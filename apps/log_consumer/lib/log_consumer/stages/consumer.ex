defmodule LogConsumer.Stages.Consumer do
  @moduledoc """
  Consumer
  """

  use GenStage

  require Logger

  alias Compliance.ValidationRuns.AggregateSupervisor

  @doc "Starts the consumer."
  def start_link() do
    GenStage.start_link(__MODULE__, :ok, name: __MODULE__)
  end

  def init(:ok) do
    {:consumer, nil, subscribe_to: [{LogConsumer.Stages.Producer, max_demand: 10}]}
  end

  def handle_events(log_items, _from, state) do
    log_items
    |> Stream.map(&log_item/1)
    |> Stream.map(&add_log_item/1)
    |> Enum.each(&log_report/1)

    {:noreply, [], state}
  end

  defp add_log_item(json_string) when is_binary(json_string) do
    case AggregateSupervisor.add_log_item(json_string) do
      {:ok, report} ->
        report

      error ->
        log_error(error, json_string)
        nil
    end
  end

  defp log_item(json_string) when is_binary(json_string) do
    Logger.debug(fn -> "item: #{inspect(json_string)}" end)
    json_string
  end

  defp log_report(report) do
    Logger.debug(fn -> "report: #{inspect(report)}" end)
    report
  end

  defp log_error(error, json_string) do
    Logger.error(inspect(error) <> " --- " <> inspect(json_string))
  end
end
