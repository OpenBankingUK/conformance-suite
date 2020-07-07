package server

import (
	"bytes"
	"net/http"

	"github.com/pkg/errors"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/report"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
)

// MIME types
const (
	MIMEApplicationZIP = "application/zip"
)

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

	request := models.ExportRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	logger.WithField("request", request).Info("Exporting ...")

	results := h.journey.Results().AllResultsGrouped()
	responseFields := h.journey.Results().ResponseFieldsJSON()
	tokens := h.journey.Events().AllAcquiredAccessToken()
	discovery, err := h.journey.DiscoveryModel()

	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(errors.Wrap(err, "exporting report-get journey discovery model")))
	}

	exportResults := models.ExportResults{
		ExportRequest:    request,
		HasPassed:        false,
		Results:          results,
		Tokens:           tokens,
		DiscoveryModel:   discovery,
		TLSVersionResult: h.journey.TLSVersionResult(),
		ResponseFields:   responseFields,
		JWSStatus:        model.JWSStatus(),
	}

	r, err := report.NewReport(exportResults, request.Environment)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	buff := bytes.NewBuffer([]byte{})
	exporter := report.NewZipExporter(r, buff)
	if err := exporter.Export(); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	// TODO(mbana): Might help to return these, if not remove in the future.
	// name := "report.zip"
	// dispositionType := "attachment"
	// c.Response().Header().Set(HeaderContentDisposition, fmt.Sprintf("%s; filename=%q", dispositionType, name))
	// c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="report.zip"`)
	return c.Blob(http.StatusOK, MIMEApplicationZIP, buff.Bytes())
}
