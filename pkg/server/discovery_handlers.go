package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/labstack/echo"
	"net/http"
)

var validDiscoveryModel *discovery.Model

type discoveryHandlers struct {
	validator discovery.Validator
}

func newDiscoveryHandlers(checker model.ConditionalityChecker) discoveryHandlers {
	return discoveryHandlers{
		discovery.NewFuncValidator(checker),
	}
}

func (d discoveryHandlers) discoveryModelValidateHandler(c echo.Context) error {
	discoveryModel := &discovery.Model{}
	if err := c.Bind(discoveryModel); err != nil {
		return badRequestErrorResponse(c, err.Error())
	}

	failures, err := d.validator.Validate(discoveryModel)
	if err != nil {
		return badRequestErrorResponse(c, err.Error())
	}

	if !failures.Empty() {
		return badRequestErrorResponse(c, failures)
	}
	return c.JSONPretty(http.StatusOK, discoveryModel, "  ")
}

func (d discoveryHandlers) persistDiscoveryModelHandler(c echo.Context) error {
	discoveryModel := &discovery.Model{}
	if err := c.Bind(discoveryModel); err != nil {
		return badRequestErrorResponse(c, err.Error())
	}

	failures, err := d.validator.Validate(discoveryModel)
	if err != nil {
		return badRequestErrorResponse(c, err.Error())
	}

	if !failures.Empty() {
		return badRequestErrorResponse(c, failures)
	}

	validDiscoveryModel = discoveryModel

	return c.JSONPretty(http.StatusCreated, discoveryModel, "  ")
}

func badRequestErrorResponse(c echo.Context, err interface{}) error {
	return c.JSONPretty(http.StatusBadRequest, &ErrorResponse{Error: err}, "  ")
}
