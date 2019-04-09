package server

import (
	"fmt"
	"gopkg.in/resty.v1"
	"net/http"
	"net/url"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
)

// Needs to be a interface{} slice, see the official test for an example
// https://github.com/go-ozzo/ozzo-validation/blob/master/in_test.go
type ResponseType = interface{}

var (
	// responseTypesSupported REQUIRED. JSON array containing a list of the OAuth 2.0 response_type values that this OP supports. Dynamic OpenID Providers MUST support the code, id_token, and the token id_token Response Type values
	responseTypesSupported = [3]ResponseType{
		"code",
		"code id_token",
		"id_token",
	}
)

type configHandlers struct {
	logger  *logrus.Entry
	journey Journey
}

type GlobalConfiguration struct {
	SigningPrivate          string            `json:"signing_private" validate:"not_empty"`
	SigningPublic           string            `json:"signing_public" validate:"not_empty"`
	TransportPrivate        string            `json:"transport_private" validate:"not_empty"`
	TransportPublic         string            `json:"transport_public" validate:"not_empty"`
	ClientID                string            `json:"client_id" validate:"not_empty"`
	ClientSecret            string            `json:"client_secret" validate:"not_empty"`
	TokenEndpoint           string            `json:"token_endpoint" validate:"valid_url"`
	ResponseType            string            `json:"response_type" validate:"not_empty"`
	TokenEndpointAuthMethod string            `json:"token_endpoint_auth_method" validate:"not_empty"`
	AuthorizationEndpoint   string            `json:"authorization_endpoint" validate:"valid_url"`
	ResourceBaseURL         string            `json:"resource_base_url" validate:"valid_url"`
	XFAPIFinancialID        string            `json:"x_fapi_financial_id" validate:"not_empty"`
	Issuer                  string            `json:"issuer" validate:"valid_url"`
	RedirectURL             string            `json:"redirect_url" validate:"valid_url"`
	ResourceIDs             model.ResourceIDs `json:"resource_ids" validate:"not_empty"`
	CreditorAccount         models.Payment    `json:"creditor_account"`
	TransactionFromDate     string            `json:"transaction_from_date" validate:"not_empty"`
	TransactionToDate       string            `json:"transaction_to_date" validate:"not_empty"`
}

// Validate - used by https://github.com/go-ozzo/ozzo-validation to validate struct.
func (c GlobalConfiguration) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.CreditorAccount, validation.Required),
		validation.Field(&c.ResponseType, validation.Required, validation.In(responseTypesSupported[:]...)),
	)
}

func newConfigHandlers(journey Journey, logger *logrus.Entry) configHandlers {
	return configHandlers{
		journey: journey,
		logger:  logger.WithField("module", "configHandlers"),
	}
}

// POST /api/config/global
func (h configHandlers) configGlobalPostHandler(c echo.Context) error {
	config := new(GlobalConfiguration)
	if err := c.Bind(config); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "error with Bind")))
	}

	if err := config.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	journeyConfig, err := MakeJourneyConfig(config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	// Use the transport keys for MATLS as some endpoints require this
	resty.SetCertificates(journeyConfig.certificateTransport.TLSCert())

	err = h.journey.SetConfig(journeyConfig)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, config)
}

func MakeJourneyConfig(config *GlobalConfiguration) (JourneyConfig, error) {
	ok, message := validateConfig(config)
	if !ok {
		return JourneyConfig{}, errors.New(message)
	}

	certificateSigning, err := authentication.NewCertificate(config.SigningPublic, config.SigningPrivate)
	if err != nil {
		return JourneyConfig{}, errors.Wrap(err, "error with signing certificate")
	}

	certificateTransport, err := authentication.NewCertificate(config.TransportPublic, config.TransportPrivate)
	if err != nil {
		return JourneyConfig{}, errors.Wrap(err, "error with transport certificate")
	}

	return JourneyConfig{
		certificateSigning:      certificateSigning,
		certificateTransport:    certificateTransport,
		clientID:                config.ClientID,
		clientSecret:            config.ClientSecret,
		tokenEndpoint:           config.TokenEndpoint,
		tokenEndpointAuthMethod: config.TokenEndpointAuthMethod,
		authorizationEndpoint:   config.AuthorizationEndpoint,
		resourceBaseURL:         config.ResourceBaseURL,
		xXFAPIFinancialID:       config.XFAPIFinancialID,
		issuer:                  config.Issuer,
		redirectURL:             config.RedirectURL,
		resourceIDs:             config.ResourceIDs,
		creditorAccount:         config.CreditorAccount,
		transactionFromDate:     config.TransactionFromDate,
		transactionToDate:       config.TransactionToDate,
	}, nil
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
	value        interface{}
	validateFunc validateFunc
}

type validateFunc func(key, value interface{}) (bool, string)

func notEmpty(key, value interface{}) (bool, string) {
	switch v := value.(type) {
	case string:
		if v == "" {
			return false, fmt.Sprintf("%s is empty", key)
		}
		return true, ""
	case model.ResourceIDs:

		emAccts := nilOrEmpty(v.AccountIDs)
		emStmts := nilOrEmpty(v.StatementIDs)

		if emAccts && emStmts {
			return false, fmt.Sprintf("%s is empty", key)
		}

		if emAccts {
			return false, fmt.Sprintf("%s.AccountIDs is empty", key)
		}
		// Some nested validation here, not great but need to think about validation for nested values
		for i, v := range v.AccountIDs {
			if v.AccountID == "" {
				return false, fmt.Sprintf("%s.AccountIDs contains an empty value at index %d", key, i)
			}
		}

		if emStmts {
			return false, fmt.Sprintf("%s.StatementIDs is empty", key)
		}
		// Some nested validation here, not great but need to think about validation for nested values
		for i, v := range v.StatementIDs {
			if v.StatementID == "" {
				return false, fmt.Sprintf("%s.StatementIDs contains an empty value at index %d", key, i)
			}
		}

		return true, ""
	}

	return false, fmt.Sprintf("%s type not found", key)
}

func validURL(key, value interface{}) (bool, string) {
	if _, err := url.Parse(value.(string)); err != nil {
		return false, fmt.Sprintf("invalid %s url: %s", key, err.Error())
	}
	return true, ""
}

func and(left, right validateFunc) validateFunc {
	return func(key, value interface{}) (bool, string) {
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
			value:        valueField.Interface(),
			validateFunc: validate,
		})
	}
	return rules
}

func nilOrEmpty(v interface{}) bool {
	return v == nil || reflect.ValueOf(v).Len() == 0
}
