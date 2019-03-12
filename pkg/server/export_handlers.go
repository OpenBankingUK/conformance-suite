package server

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
)

// ExportRequest - Request to `/api/export`.
type ExportRequest struct {
	Implementer         string `json:"implementer"`
	AuthorisedBy        string `json:"authorised_by"`
	JobTitle            string `json:"job_title"`
	HasAgreed           bool   `json:"has_agreed"`
	AddDigitalSignature bool   `json:"add_digital_signature"`
}

func (e ExportRequest) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Implementer, validation.Required),
		validation.Field(&e.AuthorisedBy, validation.Required),
		validation.Field(&e.JobTitle, validation.Required),
		validation.Field(&e.HasAgreed, validation.Required, validation.In(true)),
	)
}

// ExportResponse - Response to `/api/export`.
type ExportResponse struct {
	ExportRequest ExportRequest                `json:"export_request"`
	HasPassed     bool                         `json:"has_passed"`
	Results       []results.TestCase           `json:"results"`
	Tokens        []events.AcquiredAccessToken `json:"tokens"`
}

type exportHandlers struct {
	journey Journey
	logger  *logrus.Entry
}

func newExportHandlers(journey Journey, logger *logrus.Entry) exportHandlers {
	return exportHandlers{
		journey: journey,
		logger:  logger.WithField("handler", "exportHandlers"),
	}
}

func (h exportHandlers) postExport(c echo.Context) error {
	logger := h.logger.WithField("function", "postExport")

	exportRequest := new(ExportRequest)
	if err := c.Bind(exportRequest); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	err := exportRequest.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	logger.WithField("exportRequest", exportRequest).Info("Exporting ...")

	results := h.journey.Results().AllResults()
	tokens := h.journey.Events().AllAcquiredAccessToken()
	exportResponse := ExportResponse{
		ExportRequest: *exportRequest,
		HasPassed:     true,
		Results:       results,
		Tokens:        tokens,
	}
	logger.WithField("exportResponse", exportResponse).Info("Exported")

	return c.JSON(http.StatusOK, exportResponse)
}
