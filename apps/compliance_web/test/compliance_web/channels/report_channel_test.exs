defmodule ComplianceWeb.ReportChannelTest do
  use ComplianceWeb.ChannelCase
  use ExUnit.Case, async: false

  import Mock

  alias ComplianceWeb.Guardian
  alias ComplianceWeb.ReportChannel

  @validation_run_id "58b47d20-591c-11e8-950f-b72f26bb29de"
  @token "test-token"
  @current_user %{
    user_id: "test-user-id",
    access_token: @token
  }
  @socket_assigns %{
    current_user: @current_user
  }
  @socket_id "user_id"

  describe "when the socket is invalid" do
    test "cannot join without socket.assigns.current_user" do
      sock =
        @socket_id
        |> socket(%{})
        |> subscribe_and_join(ReportChannel, "report:#{@validation_run_id}")

      assert sock == {:error, %{reason: "unauthorized: :current_user missing"}}
    end

    test "cannot join without validation_run_id" do
      sock =
        @socket_id
        |> socket(%{})
        |> subscribe_and_join(ReportChannel, "report:")

      assert sock == {:error, %{reason: "invalid: validation_run_id cannot be empty"}}
    end

    test "cannot join without validation_run_id and socket.assigns.current_user" do
      sock =
        @socket_id
        |> socket(@socket_assigns)
        |> subscribe_and_join(ReportChannel, "report:")

      assert sock == {:error, %{reason: "invalid: validation_run_id cannot be empty"}}
    end

    test "user can view validation runs they have started" do
      with_mock Guardian,
        authorised?: fn @token, @validation_run_id -> true end do
        sock =
          @socket_id
          |> socket(@socket_assigns)
          |> subscribe_and_join(ReportChannel, "report:#{@validation_run_id}")

        refute sock == {:error}

        {:ok, _, sock} = sock
        leave(sock)
        close(sock)
      end
    end

    test "user cannot view validation runs they have not started" do
      with_mock Guardian,
        authorised?: fn @token, @validation_run_id -> false end do
        sock =
          @socket_id
          |> socket(@socket_assigns)
          |> subscribe_and_join(ReportChannel, "report:#{@validation_run_id}")

        assert sock ==
                 {:error,
                  %{
                    reason:
                      "unauthorized: you are not authorized to view the report #{
                        @validation_run_id
                      }"
                  }}
      end
    end
  end

  describe "when the socket is valid" do
    @aggregate_pid "<aggregate_pid>"
    @reports %{
      "/open-banking/v1.1/payment-submissions" => %Compliance.ValidationRuns.EndpointReport{
        failed_calls: 0,
        failures: [],
        path: "/open-banking/v1.1/payment-submissions",
        total_calls: 2
      },
      "/open-banking/v1.1/payments" => %Compliance.ValidationRuns.EndpointReport{
        failed_calls: 0,
        failures: [],
        path: "/open-banking/v1.1/payments",
        total_calls: 2
      }
    }
    @pushed_message %{
      payload: @reports
    }

    defp mocked_modules do
      [
        {
          Guardian,
          [],
          [
            authorised?: fn _token, _validation_run_id -> true end
          ]
        },
        {
          Compliance.ValidationRuns.AggregateSupervisor,
          [],
          [
            add_listener_pid: fn validation_run_id, listener_pid ->
              {:ok,
               %{
                 validation_run_id: validation_run_id,
                 listener_pids: [listener_pid],
                 endpoint_reports: %{}
               }}
            end,
            process_for: fn _ -> {:ok, @aggregate_pid} end
          ]
        },
        {
          Compliance.ValidationRuns.Aggregate,
          [],
          [
            endpoint_reports: fn _ -> {:ok, @reports} end
          ]
        }
      ]
    end

    defp setup_socket do
      {:ok, _, socket} =
        @socket_id
        |> socket(@socket_assigns)
        |> subscribe_and_join(ReportChannel, "report:#{@validation_run_id}")

      socket
    end

    test "`started` message is only sent on join not again" do
      with_mocks(mocked_modules()) do
        setup_socket()

        assert_push "started", @pushed_message
        refute_push "started", @pushed_message

        # leave here for now: this simulates leaving and closing the socket
        # leave(socket)
        # close(socket)
      end
    end

    test "when :report_update true then `updated` message containing updated report is sent once within report_update_ms" do
      with_mocks(mocked_modules()) do
        # set report_update_ms to 50 ms to speed up test run
        System.put_env("report_update_ms", "50")
        socket = setup_socket()

        send(socket.channel_pid, :report_update)
        assert_push "updated", @pushed_message, 60

        send(socket.channel_pid, :report_update)
        assert_push "updated", @pushed_message, 60

        refute_push "updated", @pushed_message, 60
        System.delete_env("report_update_ms")
      end
    end

    test "when :report_update not true then `updated` message is not sent" do
      with_mocks(mocked_modules()) do
        # set report_update_ms to 50 ms to speed up test run
        System.put_env("report_update_ms", "50")
        setup_socket()
        refute_push "updated", @pushed_message, 60
        System.delete_env("report_update_ms")
      end
    end

    test "ping replies with status ok" do
      with_mocks(mocked_modules()) do
        socket = setup_socket()
        ref = push(socket, "ping", %{"hello" => "there"})
        assert_reply ref, :ok, %{"hello" => "there"}
      end
    end

    test "shout broadcasts to report:lobby" do
      with_mocks(mocked_modules()) do
        socket = setup_socket()
        push(socket, "shout", %{"hello" => "all"})
        assert_broadcast "shout", %{"hello" => "all"}
      end
    end

    test "broadcasts are pushed to the client" do
      with_mocks(mocked_modules()) do
        socket = setup_socket()
        broadcast_from!(socket, "broadcast", %{"some" => "data"})
        assert_push "broadcast", %{"some" => "data"}
      end
    end
  end
end
