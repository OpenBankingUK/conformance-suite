// This is still WORK IN PROGRESS. The handlers just return either an empty
// `github.com/OpenBankingUK/conformance-suite/pkg/server/models.ImportReviewResponse` or
//  `github.com/OpenBankingUK/conformance-suite/pkg/server/models.ImportRerunResponse` and do not do the
// importing or review functionality. This will be implemented as we go along.

package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/OpenBankingUK/conformance-suite/pkg/discovery"
	"github.com/OpenBankingUK/conformance-suite/pkg/server/models"
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

	discovery, err := h.doImport(request, logger);
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	response := models.ImportRerunResponse{
		Discovery: discovery,
	}
	logger.Info("Imported")

	return c.JSON(http.StatusOK, response)
}

// nolint:unparam
func (h importHandlers) doImport(request models.ImportRequest, logger *logrus.Entry) (discovery.ModelDiscovery, error) {
	var discoveryModel discovery.ModelDiscovery
	logger.WithField("len(request.Report)", len(request.Report)).Info("Importing ...")

	// Create a reader for the zip file
	zipReader, err := zip.NewReader(bytes.NewReader([]byte(request.Report)), int64(len(request.Report)))
	if err != nil {
		return discoveryModel, fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Search for discovery.json file
	var discoveryFile *zip.File
	for _, file := range zipReader.File {
		if file.Name == "discovery.json" {
			discoveryFile = file
			break
		}
	}

	if discoveryFile == nil {
		return discoveryModel, fmt.Errorf("discovery.json not found in the zip file")
	}

	// Open the discovery.json file
	rc, err := discoveryFile.Open()
	if err != nil {
		return discoveryModel, fmt.Errorf("failed to open discovery.json: %w", err)
	}
	defer rc.Close()

	// Read the content of discovery.json
	content, err := io.ReadAll(rc)
	if err != nil {
		return discoveryModel, fmt.Errorf("failed to read discovery.json: %w", err)
	}

	// Parse the JSON content
	err = json.Unmarshal(content, &discoveryModel)
	if err != nil {
		return discoveryModel, fmt.Errorf("failed to parse discovery.json: %w", err)
	}

	// Log the parsed data
	logger.WithField("discovery", discoveryModel).Info("Discovery data parsed successfully")

	return discoveryModel, nil
}
