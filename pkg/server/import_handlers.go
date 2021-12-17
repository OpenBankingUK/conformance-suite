// This is still WORK IN PROGRESS. The handlers just return either an empty
// `bitbucket.org/openbankingteam/conformance-suite/pkg/server/models.ImportReviewResponse` or
//  `bitbucket.org/openbankingteam/conformance-suite/pkg/server/models.ImportRerunResponse` and do not do the
// importing or review functionality. This will be implemented as we go along.

package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
)

type importHandlers struct {
	journey Journey
	logger  *logrus.Entry
}

func newImportHandlers(journey Journey, logger *logrus.Entry) importHandlers {
	return importHandlers{
		journey: journey,
		logger:  logger.WithField("handler", "importHandlers"),
	}
}

// postImportReview - `/api/import/review` POST.
func (h importHandlers) postImportReview(c echo.Context) error {
	logger := h.logger.WithField("function", "postImportReview")

	request := models.ImportRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	if err := h.doImport(request, logger); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	response := models.ImportReviewResponse{}
	logger.Info("Imported")

	return c.JSON(http.StatusOK, response)
}

// postImportRerun - `/api/import/rerun` POST.
func (h importHandlers) postImportRerun(c echo.Context) error {
	logger := h.logger.WithField("function", "postImportRerun")

	request := models.ImportRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	if err := h.doImport(request, logger); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	response := models.ImportRerunResponse{}
	logger.Info("Imported")

	return c.JSON(http.StatusOK, response)
}

// nolint:unparam
func (h importHandlers) doImport(request models.ImportRequest, logger *logrus.Entry) error {
	// logger.WithField("request", request).Info("Importing ...")
	logger.WithField("len(request.Report)", len(request.Report)).Info("Importing ...")
	// TODO(mbana): Do something with `report.zip`
	return nil
}
