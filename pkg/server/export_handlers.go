package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
)

// ExportRequest - Request to `/api/export`.
type ExportRequest struct {
	Implementer         string `json:"implementer" validate:"not_empty"`
	AuthorisedBy        string `json:"authorised_by" validate:"not_empty"`
	JobTitle            string `json:"job_title" validate:"not_empty"`
	HasAgreed           bool   `json:"has_agreed" validate:"not_empty"`
	AddDigitalSignature bool   `json:"add_digital_signature" validate:"not_empty"`
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
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "error with Bind")))
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
