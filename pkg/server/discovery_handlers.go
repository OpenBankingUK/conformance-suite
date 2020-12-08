package server

import (
	"fmt"
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/sets"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
)

const (
	defaultTxnFrom = "2016-01-01T10:40:00+02:00"
	defaultTxnTo   = "2025-12-31T10:40:00+02:00"
)

type PostDiscoveryModelResponse struct {
	TokenEndpoints                                map[string]string   `json:"token_endpoints"`
	TokenEndpointAuthMethods                      map[string][]string `json:"token_endpoint_auth_methods"`
	DefaultTokenEndpointAuthMethod                map[string]string   `json:"default_token_endpoint_auth_method"`
	RequestObjectSigningAlgValuesSupported        map[string][]string `json:"request_object_signing_alg_values_supported"`
	DefaultRequestObjectSigningAlgValuesSupported map[string]string   `json:"default_request_object_signing_alg_values_supported"`
	AuthorizationEndpoints                        map[string]string   `json:"authorization_endpoints"`
	Issuers                                       map[string]string   `json:"issuers"`
	DefaultTxnFromDateTime                        string              `json:"default_transaction_from_date"`
	DefaultTxnToDateTime                          string              `json:"default_transaction_to_date"`
	ResponseTypesSupported                        []string            `json:"response_types_supported"`
	AcrValuesSupported                            []string            `json:"acr_values_supported"`
}

type validationFailuresResponse struct {
	Error discovery.ValidationFailures `json:"error"`
}

type discoveryHandlers struct {
	webJourney Journey
	logger     *logrus.Entry
}

func newDiscoveryHandlers(webJourney Journey, logger *logrus.Entry) discoveryHandlers {
	return discoveryHandlers{webJourney, logger.WithField("handler", "discoveryHandlers")}
}

func (d discoveryHandlers) setDiscoveryModelHandler(c echo.Context) error {
	ctxLogger := d.logger.WithFields(logrus.Fields{
		"function": "setDiscoveryModelHandler",
	})

	discoveryModel := &discovery.Model{}
	if err := c.Bind(discoveryModel); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	failures, err := d.webJourney.SetDiscoveryModel(discoveryModel)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	if !failures.Empty() {
		return c.JSON(http.StatusBadRequest, validationFailuresResponse{failures})
	}

	failures = discovery.ValidationFailures{}
	response := PostDiscoveryModelResponse{
		TokenEndpoints:                                map[string]string{},
		TokenEndpointAuthMethods:                      map[string][]string{},
		DefaultTokenEndpointAuthMethod:                map[string]string{},
		RequestObjectSigningAlgValuesSupported:        map[string][]string{},
		DefaultRequestObjectSigningAlgValuesSupported: map[string]string{},
		AuthorizationEndpoints:                        map[string]string{},
		Issuers:                                       map[string]string{},
		ResponseTypesSupported:                        []string{},
		AcrValuesSupported:                            []string{},
	}

	configGetter := authentication.NewOpenIdConfigGetter()
	for discoveryItemIndex, discoveryItem := range discoveryModel.DiscoveryModel.DiscoveryItems {
		key := fmt.Sprintf("schema_version=%s", discoveryItem.APISpecification.SchemaVersion)

		url := discoveryItem.OpenidConfigurationURI
		ctxLogger.WithFields(logrus.Fields{
			"url": url,
		}).Info("GET /.well-known/openid-configuration")
		config, e := configGetter.Get(url)
		if e != nil {
			ctxLogger.WithFields(logrus.Fields{
				"url": url,
				"err": e,
			}).Error("Error on /.well-known/openid-configuration")
			failures = append(failures, newOpenidConfigurationURIFailure(discoveryItemIndex, e))
		} else {
			var SupportedRequestSignAlgValues = []string{"PS256", "RS256", "NONE"}
			requestObjectSigningAlgValuesSupported := sets.InsensitiveIntersection(config.RequestObjectSigningAlgValuesSupported, SupportedRequestSignAlgValues)
			if len(requestObjectSigningAlgValuesSupported) == 0 {
				return errors.New("no supported request object signing alg found")
			}

			response.TokenEndpoints[key] = config.TokenEndpoint
			response.AuthorizationEndpoints[key] = config.AuthorizationEndpoint
			response.Issuers[key] = config.Issuer
			response.TokenEndpointAuthMethods[key] = authentication.SuiteSupportedAuthMethodsMostSecureFirst()
			response.DefaultTokenEndpointAuthMethod[key] = authentication.DefaultAuthMethod(config.TokenEndpointAuthMethodsSupported, d.logger)
			response.RequestObjectSigningAlgValuesSupported[key] = requestObjectSigningAlgValuesSupported
			response.DefaultRequestObjectSigningAlgValuesSupported[key] = config.RequestObjectSigningAlgValuesSupported[0]
			response.DefaultTxnFromDateTime = defaultTxnFrom
			response.DefaultTxnToDateTime = defaultTxnTo
			response.ResponseTypesSupported = config.ResponseTypesSupported
			response.AcrValuesSupported = config.AcrValuesSupported
		}
	}

	if !failures.Empty() {
		return c.JSON(http.StatusBadRequest, validationFailuresResponse{failures})
	}
	return c.JSON(http.StatusCreated, response)
}

func newOpenidConfigurationURIFailure(discoveryItemIndex int, err error) discovery.ValidationFailure {
	return discovery.ValidationFailure{
		Key:   fmt.Sprintf("DiscoveryModel.DiscoveryItems[%d].OpenidConfigurationURI", discoveryItemIndex),
		Error: err.Error(),
	}
}
