package server

import (
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/pkg/errors"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
)

var (
	ErrEmptyClientID              = errors.New("client_id is empty")
	ErrEmptyClientSecret          = errors.New("client_secret is empty")
	ErrEmptyTokenEndpoint         = errors.New("token_endpoint is empty")
	ErrEmptyAuthorizationEndpoint = errors.New("authorization_endpoint is empty")
	ErrEmptyXFAPIFinancialID      = errors.New("x_fapi_financial_id is empty")
	ErrEmptyRedirectURL           = errors.New("redirect_url is empty")
)

type configHandlers struct {
	logger  *logrus.Entry
	journey Journey
}

type GlobalConfiguration struct {
	SigningPrivate        string `json:"signing_private"`
	SigningPublic         string `json:"signing_public"`
	TransportPrivate      string `json:"transport_private"`
	TransportPublic       string `json:"transport_public"`
	ClientID              string `json:"client_id"`
	ClientSecret          string `json:"client_secret"`
	TokenEndpoint         string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	XFAPIFinancialID      string `json:"x_fapi_financial_id"`
	RedirectURL           string `json:"redirect_url"`
}

// POST /api/config/global
func (h *configHandlers) configGlobalPostHandler(c echo.Context) error {
	config := new(GlobalConfiguration)
	if err := c.Bind(config); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "error with Bind")))
	}

	certificateSigning, err := authentication.NewCertificate(config.SigningPublic, config.SigningPrivate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "error with signing certificate")))
	}

	certificateTransport, err := authentication.NewCertificate(config.TransportPublic, config.TransportPrivate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "error with transport certificate")))
	}

	if config.ClientID == "" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(ErrEmptyClientID))
	}
	if config.ClientSecret == "" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(ErrEmptyClientSecret))
	}
	if config.TokenEndpoint == "" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(ErrEmptyTokenEndpoint))
	}
	if config.XFAPIFinancialID == "" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(ErrEmptyXFAPIFinancialID))
	}
	if config.RedirectURL == "" {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(ErrEmptyRedirectURL))
	}

	if _, err := url.Parse(config.TokenEndpoint); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "token_endpoint")))
	}
	if _, err := url.Parse(config.AuthorizationEndpoint); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "authorization_endpoint")))
	}
	if _, err := url.Parse(config.RedirectURL); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "redirect_url")))
	}

	h.journey.SetConfig(certificateSigning, certificateTransport, config.ClientID, config.ClientSecret, config.TokenEndpoint, config.AuthorizationEndpoint, config.XFAPIFinancialID, config.RedirectURL)

	return c.JSON(http.StatusCreated, config)
}
