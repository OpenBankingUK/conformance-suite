package server

import (
	"fmt"
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type PostDiscoveryModelResponse struct {
	TokenEndpoints         map[string]string `json:"token_endpoints"`
	AuthorizationEndpoints map[string]string `json:"authorization_endpoints"`
	Issuers                map[string]string `json:"issuers"`
}

type validationFailuresResponse struct {
	Error discovery.ValidationFailures `json:"error"`
}

type discoveryHandlers struct {
	webJourney Journey
	logger     *logrus.Entry
}

func newDiscoveryHandlers(webJourney Journey, logger *logrus.Entry) discoveryHandlers {
	return discoveryHandlers{webJourney, logger.WithField("handler", "discoveryHandler")}
}

func (d discoveryHandlers) setDiscoveryModelHandler(c echo.Context) error {
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
		TokenEndpoints:         map[string]string{},
		AuthorizationEndpoints: map[string]string{},
		Issuers:                map[string]string{},
	}
	for discoveryItemIndex, discoveryItem := range discoveryModel.DiscoveryModel.DiscoveryItems {
		key := fmt.Sprintf("schema_version=%s", discoveryItem.APISpecification.SchemaVersion)

		url := discoveryItem.OpenidConfigurationURI
		d.logger.Info(fmt.Sprintf("GET OpenID config: %s", url))
		config, e := authentication.OpenIdConfig(url)
		if e != nil {
			d.logger.WithError(e).Warn(fmt.Sprintf("Failed to GET %s", url))
			failures = append(failures, newOpenidConfigurationURIFailure(discoveryItemIndex, e))
		} else {
			response.TokenEndpoints[key] = config.TokenEndpoint
			response.AuthorizationEndpoints[key] = config.AuthorizationEndpoint
			response.Issuers[key] = config.Issuer
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
