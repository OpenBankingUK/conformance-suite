defmodule OBApiRemote.Commands.Proxied do
  @moduledoc """
  Contains configuration for HTTP calls to OB API proxy.
  """
  require Logger

  def base_url do
    Application.get_env(:ob_api_remote, :proxy_url)
  end

  def urls do
    %{
      login: "#{base_url()}/login",
      logout: "#{base_url()}/logout",
      get_auth_servers: "#{base_url()}/account-payment-service-provider-authorisation-servers",
      authorise_account_access: "#{base_url()}/account-request-authorise-consent",
      authorise_payment: "#{base_url()}/payment-authorise-consent",
      account_request_revoke_consent: "#{base_url()}/account-request-revoke-consent",
      consent_authorised: "#{base_url()}/tpp/authorized",
      complete_payment: "#{base_url()}/payment-submissions",
      accounts: "#{base_url()}/accounts"
    }
  end

  def template(name) do
    case name do
      :login ->
        fn %{username: username, password: password} ->
          %{
            url: urls().login,
            method: :post,
            body: Poison.encode!(%{u: username, p: password})
          }
        end

      :logout ->
        %{
          url: urls().logout,
          method: :get,
          body: ""
        }

      :get_auth_servers ->
        %{
          url: urls().get_auth_servers,
          method: :get,
          body: ""
        }

      :authorise_account_access ->
        %{
          url: urls().authorise_account_access,
          method: :post,
          body: ""
        }

      :authorise_payment ->
        fn payment ->
          %{
            url: urls().authorise_payment,
            method: :post,
            body:
              Poison.encode!(%{
                authorisationServerId: payment["auth_server_id"],
                InstructedAmount: %{
                  Amount: payment["amount"],
                  Currency: "GBP"
                },
                CreditorAccount: %{
                  SchemeName: "SortCodeAccountNumber",
                  Identification: payment["sort_code"] <> payment["account_number"],
                  Name: payment["name"]
                }
              })
          }
        end

      :consent_authorised ->
        fn params ->
          %{
            url: urls().consent_authorised,
            method: :post,
            body:
              Poison.encode!(%{
                accountRequestId: params.account_request_id,
                authorisationServerId: params.auth_server_id,
                authorisationCode: params.authorisation_code,
                scope: params.scope
              })
          }
        end

      :account_request_revoke_consent ->
        %{
          url: urls().account_request_revoke_consent,
          method: :post,
          body: ""
        }

      :complete_payment ->
        %{
          url: urls().complete_payment,
          method: :post,
          body: ""
        }

      endpoint ->
        %{
          url: "#{base_url()}#{endpoint}",
          method: :get,
          body: ""
        }
    end
  end

  def get(name), do: template(name)

  def get(name, params), do: template(name).(params)
end
