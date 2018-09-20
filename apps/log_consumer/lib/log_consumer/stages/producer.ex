defmodule LogConsumer.Stages.Producer do
  @moduledoc """
  Producer
  """

  use GenStage

  require Logger

  @doc "Starts the producer."
  def start_link() do
    GenStage.start_link(__MODULE__, :ok, name: __MODULE__)
  end

  def initial_state do
    %{message_set: [], demand: 0, from: nil}
  end

  def init(:ok) do
    {:producer, initial_state()}
  end

  @doc """
  Called with message_set from KafkaConsumer.
  """
  def notify(message_set, timeout \\ 5000) do
    GenStage.call(__MODULE__, {:notify, message_set}, timeout)
  end

  @doc """
  When we have no demand, save message_set for later, and don't reply.
  """
  def handle_call({:notify, message_set}, from, %{demand: 0} = state) do
    new_state = %{state | message_set: message_set, from: from}
    Logger.debug(fn -> "from1: #{inspect(from)}" end)
    {:noreply, [], new_state}
  end

  @doc """
  When we have incoming messages greater than demand, emit some now
  but don't reply yet.
  """
  def handle_call({:notify, message_set}, from, %{demand: demand} = state)
      when length(message_set) > demand do
    {to_dispatch, remaining} = Enum.split(message_set, demand)
    to_dispatch = to_dispatch |> Enum.map(&"#{&1.value}")

    new_demand = demand - length(to_dispatch)
    Logger.debug(fn -> "from2: #{inspect(from)}" end)
    new_state = %{state | message_set: remaining, demand: new_demand, from: from}

    {:noreply, to_dispatch, new_state}
  end

  @doc """
  When we have incoming messages less than or equal to demand, emit all now and reply.
  """
  def handle_call({:notify, message_set}, _from, %{demand: demand} = state) do
    to_dispatch = message_set |> Enum.map(&"#{&1.value}")

    new_demand = demand - length(to_dispatch)
    # Logger.debug(fn -> "from3: #{inspect(from)}" end)
    new_state = %{state | demand: new_demand}

    {:reply, :ok, to_dispatch, new_state}
  end

  @doc """
  When we can't emit anything, save the demand don't reply yet
  """
  def handle_demand(demand, %{message_set: []} = state) when demand > 0 do
    {:noreply, [], %{state | demand: demand}}
  end

  @doc """
  When we can immediately satisfy demand, don't reply yet
  """
  def handle_demand(demand, %{message_set: message_set} = state)
      when demand > 0 and length(message_set) > demand do
    {to_dispatch, remaining} = Enum.split(message_set, demand)
    to_dispatch = Enum.map(to_dispatch, &"#{&1.value}")
    {:noreply, to_dispatch, %{state | message_set: remaining, demand: 0}}
  end

  @doc """
  When we can't immediately satisfy demand, emit what we have and reply
  """
  def handle_demand(demand, %{message_set: message_set} = state) when demand > 0 do
    to_dispatch = Enum.map(message_set, &"#{&1.value}")
    new_demand = demand - length(to_dispatch)
    new_state = %{state | message_set: [], demand: new_demand}
    GenStage.reply(state.from, :ok)
    {:noreply, to_dispatch, new_state}
  end
end
