package server

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/pkg/errors"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
)

type configHandlers struct {
	logger  *logrus.Entry
	journey Journey
}

type GlobalConfiguration struct {
	SigningPrivate        string `json:"signing_private" validate:"not_empty"`
	SigningPublic         string `json:"signing_public" validate:"not_empty"`
	TransportPrivate      string `json:"transport_private" validate:"not_empty"`
	TransportPublic       string `json:"transport_public" validate:"not_empty"`
	ClientID              string `json:"client_id" validate:"not_empty"`
	ClientSecret          string `json:"client_secret" validate:"not_empty"`
	TokenEndpoint         string `json:"token_endpoint" validate:"valid_url"`
	AuthorizationEndpoint string `json:"authorization_endpoint" validate:"valid_url"`
	ResourceBaseURL       string `json:"resource_base_url" validate:"valid_url"`
	XFAPIFinancialID      string `json:"x_fapi_financial_id" validate:"not_empty"`
	Issuer                string `json:"issuer" validate:"valid_url"`
	RedirectURL           string `json:"redirect_url" validate:"valid_url"`
}

// POST /api/config/global
func (h *configHandlers) configGlobalPostHandler(c echo.Context) error {
	config := new(GlobalConfiguration)
	if err := c.Bind(config); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "error with Bind")))
	}

	ok, message := validateConfig(config)
	if !ok {
		return c.JSON(http.StatusBadRequest, NewErrorMessageResponse(message))
	}

	certificateSigning, err := authentication.NewCertificate(config.SigningPublic, config.SigningPrivate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "error with signing certificate")))
	}

	certificateTransport, err := authentication.NewCertificate(config.TransportPublic, config.TransportPrivate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "error with transport certificate")))
	}

	jConfig := JourneyConfig{
		certificateSigning:    certificateSigning,
		certificateTransport:  certificateTransport,
		clientID:              config.ClientID,
		clientSecret:          config.ClientSecret,
		tokenEndpoint:         config.TokenEndpoint,
		authorizationEndpoint: config.AuthorizationEndpoint,
		resourceBaseURL:       config.ResourceBaseURL,
		xXFAPIFinancialID:     config.XFAPIFinancialID,
		issuer:                config.Issuer,
		redirectURL:           config.RedirectURL,
	}
	err = h.journey.SetConfig(jConfig)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, config)
}

func validateConfig(config *GlobalConfiguration) (bool, string) {
	rules := parseRules(config)
	for _, rule := range rules {
		ok, message := rule.validateFunc(rule.property, rule.value)
		if !ok {
			return false, message
		}
	}
	return true, ""
}

type validationRule struct {
	property     string
	value        string
	validateFunc validateFunc
}

type validateFunc func(key, value string) (bool, string)

func notEmpty(key, value string) (bool, string) {
	if value == "" {
		return false, fmt.Sprintf("%s is empty", key)
	}
	return true, ""
}

func validURL(key, value string) (bool, string) {
	if _, err := url.Parse(value); err != nil {
		return false, fmt.Sprintf("invalid %s url: %s", key, err.Error())
	}
	return true, ""
}

func and(left, right validateFunc) validateFunc {
	return func(key, value string) (bool, string) {
		ok, msg := left(key, value)
		if !ok {
			return false, msg
		}
		ok, msg = right(key, value)
		if !ok {
			return false, msg
		}
		return true, ""
	}
}

var rulesFunc = map[string]validateFunc{
	"not_empty": notEmpty,
	"valid_url": and(notEmpty, validURL),
}

func parseRules(config *GlobalConfiguration) []validationRule {
	var rules []validationRule
	val := reflect.ValueOf(config).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		if tag.Get("validate") == "" {
			// no validate tag
			continue
		}

		validate, ok := rulesFunc[tag.Get("validate")]
		if !ok {
			// no rule func found
			continue
		}

		rules = append(rules, validationRule{
			property:     tag.Get("json"),
			value:        valueField.Interface().(string),
			validateFunc: validate,
		})
	}
	return rules
}
