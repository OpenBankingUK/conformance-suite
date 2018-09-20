defmodule LogConsumer.Stages.ProducerTest do
  @moduledoc false
  # @moduletag kafka

  use ExUnit.Case, async: true
  import Mock
  alias LogConsumer.Stages.Producer

  def message(offset) do
    %KafkaEx.Protocol.Fetch.Message{
      attributes: 0,
      crc: 4_264_455_069 + offset,
      key: nil,
      offset: offset,
      value: "hey#{offset}"
    }
  end

  # mock of from data
  def from do
    {:pid, :reference}
  end

  # mock of from data
  def from2 do
    {:pid, :reference2}
  end

  test "initial_state has demand: 0, message_set: [], and from: nil" do
    {:producer, state} = Producer.init(:ok)
    assert state == %{message_set: [], demand: 0, from: nil}
  end

  test "handle_call when no demand, save message_set in state, and don't reply" do
    messages = [message(1)]

    {response, to_dispatch, state} =
      Producer.handle_call({:notify, messages}, from(), %{
        message_set: [],
        demand: 0,
        from: nil
      })

    assert response == :noreply
    assert to_dispatch == []
    assert state == %{message_set: messages, demand: 0, from: from()}
  end

  test "handle_call when incoming messages greater than demand, emit some now
    but don't reply" do
    messages = [message(2), message(3)]

    {response, to_dispatch, state} =
      Producer.handle_call({:notify, messages}, from2(), %{
        message_set: [],
        demand: 1,
        from: from()
      })

    assert response == :noreply
    assert to_dispatch == [message(2).value]
    assert state == %{message_set: [message(3)], demand: 0, from: from2()}
  end

  test "handle_call when incoming messages less than or equal to demand, emit all now and reply" do
    messages = [message(2), message(3)]

    {response, status, to_dispatch, state} =
      Producer.handle_call({:notify, messages}, from2(), %{
        message_set: [],
        demand: 5,
        from: from()
      })

    assert response == :reply
    assert status == :ok
    assert to_dispatch == [message(2).value, message(3).value]
    assert state == %{message_set: [], demand: 5 - 2, from: from()}
  end

  test "handle_demand when no messages to emit, save the demand, and don't reply yet" do
    {response, to_dispatch, state} =
      Producer.handle_demand(100, %{message_set: [], demand: 0, from: nil})

    assert response == :noreply
    assert to_dispatch == []
    assert state == %{message_set: [], demand: 100, from: nil}
  end

  test "handle_demand when message_set greater than demand, emit some and don't reply yet" do
    messages = [message(2), message(3)]

    {response, to_dispatch, state} =
      Producer.handle_demand(1, %{message_set: messages, demand: 0, from: from()})

    assert response == :noreply
    assert to_dispatch == [message(2).value]
    assert state == %{message_set: [message(3)], demand: 0, from: from()}
  end

  test "handle_demand when message_set less than or equal to demand, emit what we have and reply" do
    with_mock GenStage, reply: fn {:pid, :reference}, :ok -> :ok end do
      messages = [message(2), message(3)]

      {response, to_dispatch, state} =
        Producer.handle_demand(5, %{message_set: messages, demand: 0, from: from()})

      assert response == :noreply
      assert to_dispatch == [message(2).value, message(3).value]
      assert state == %{message_set: [], demand: 5 - 2, from: from()}
    end
  end
end
