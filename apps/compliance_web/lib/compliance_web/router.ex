defmodule ComplianceWeb.Router do
  use ComplianceWeb, :router

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_flash
    plug :protect_from_forgery
    plug :put_secure_browser_headers
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  pipeline :authorized do
    plug :fetch_session

    plug Guardian.Plug.Pipeline,
      module: ComplianceWeb.Guardian,
      error_handler: ComplianceWeb.AuthErrorHandler

    plug Guardian.Plug.VerifyHeader, realm: "Bearer"
    plug Guardian.Plug.LoadResource
  end

  scope "/", ComplianceWeb do
    # Use the default browser stack
    pipe_through :browser

    get "/", PageController, :index
  end

  scope "/", ComplianceWeb do
    # Use the default browser stack
    pipe_through :browser

    get "/account-payment-service-provider-authorisation-servers",
        AuthorisationServersController,
        :get
  end

  scope "/", ComplianceWeb do
    pipe_through :api
    pipe_through :authorized

    get "/user", UserController, :show
    # do we need that?
    delete "/auth", AuthController, :delete
    resources "/validation-runs", ValidationRunController, only: [:create, :show]
    resources "/run-configs", RunConfigController, only: [:create]
  end

  scope "/", ComplianceWeb do
    # open endpoint to create new users
    post "/auth", AuthController, :new
    # This only works in dev for the e2e tests
    get "/tokeninfo", AuthController, :tokeninfo
  end

  # Catch all routing - to accomodate browser reloading of Vue.js client app.
  # This *MUST* stay as the last route in file to allow other matches to happen first.
  scope "/", ComplianceWeb do
    pipe_through :browser

    get "/:anything", PageController, :index
    get "/:anything/:anotherthing", PageController, :index
  end
end
