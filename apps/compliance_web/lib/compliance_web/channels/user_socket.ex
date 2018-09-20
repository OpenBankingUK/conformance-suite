defmodule ComplianceWeb.UserSocket do
  use Phoenix.Socket
  require Logger

  ## Channels
  channel "report:*", ComplianceWeb.ReportChannel

  ## Transports
  transport(:websocket, Phoenix.Transports.WebSocket)
  # http://elixirbridge.org/01_Installfest/10-deploy-a-phoenix-app.html
  # transport :websocket, Phoenix.Transports.WebSocket, check_origin: false
  # transport :websocket, Phoenix.Transports.WebSocket, check_origin: false, timeout: :infinity
  # transport :longpoll, Phoenix.Transports.LongPoll

  # Socket params are passed from the client and can
  # be used to verify and authenticate a user. After
  # verification, you can put default assigns into
  # the socket that will be set for all channels, ie
  #
  #     {:ok, assign(socket, :user_id, verified_user_id)}
  #
  # To deny connection, return `:error`.
  #
  # See `Phoenix.Token` documentation for examples in
  # performing token verification on connect.
  def connect(%{"token" => token}, socket) do
    Logger.info(
      "CONNECT(%{\"TOKEN\"}) token: #{inspect(token)}, socket.assigns: #{inspect(socket.assigns)}"
    )

    # %{id: user_id, token: token} = Compliance.Accounts.get_user(token: token)

    current_user = %{
      access_token: token,
    }

    {:ok, assign(socket, :current_user, current_user)}
  end

  def connect(_params, _socket) do
    :error
  end

  # Socket id's are topics that allow you to identify all sockets for a
  # given user:
  #
  #   def id(socket), do: "user_socket:#{socket.assigns.user_id}"
  #
  # Would allow you to broadcast a "disconnect" event and terminate
  # all active sockets and channels for a given user:
  #
  #   ComplianceWeb.Endpoint.broadcast(
  #     "user_socket:#{user.id}", "disconnect", %{}
  #   )
  #
  # Returning `nil` makes this socket anonymous.
  def id(socket) do
    "user_socket:#{socket.assigns.current_user.access_token}"
  end
end
