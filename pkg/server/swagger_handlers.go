package server

import (
	"fmt"
	"net/http"
	"net/url"

	openapi_middleware "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

// swaggerHandlers - maps paths (e.g., /swagger/account-transaction-v3.0/v3.0/docs) below to handlers.
func swaggerHandlers(logger *logrus.Entry) map[string]echo.HandlerFunc {
	const path = "docs"
	const redocURL = "/static/redoc/bundles/redoc.standalone.js" // we copy ReDoc in `web/vue.config.js`

	ctxLogger := logger.WithField("module", "swaggerHandlers")

	handlers := map[string]echo.HandlerFunc{}
	specs := model.Specifications()
	for _, spec := range specs {
		// /swagger/account-transaction-v3.0/v3.0
		basePath := fmt.Sprintf("/swagger/%s/%s", spec.Identifier, spec.Version)
		basePathURL, err := url.Parse(basePath)
		if err != nil {
			ctxLogger.WithFields(logrus.Fields{
				"basePath":    basePath,
				"basePathURL": basePathURL,
			}).WithError(err)
			continue
		}

		// /swagger/account-transaction-v3.0/v3.0/docs
		fullPath := fmt.Sprintf("/swagger/%s/%s/%s", spec.Identifier, spec.Version, path)
		fullPathURL, err := url.Parse(fullPath)
		if err != nil {
			ctxLogger.WithFields(logrus.Fields{
				"fullPath":    fullPath,
				"fullPathURL": fullPathURL,
			}).WithError(err)
			continue
		}

		// https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json
		specURL := spec.SchemaVersion.String()

		// See https://github.com/go-swagger/go-swagger/blob/master/cmd/swagger/commands/serve.go
		// for how `swagger serve 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json'`
		// works internally.

		// Somehow wrap https://github.com/labstack/echo/blob/master/echo.go#L276
		var notFoundHandler http.Handler = nil // http.NotFoundHandler()
		handler := openapi_middleware.Redoc(openapi_middleware.RedocOpts{
			BasePath: basePath,
			Path:     path,
			SpecURL:  specURL,
			RedocURL: redocURL,
		}, notFoundHandler)

		handlers[fullPath] = echo.WrapHandler(handler)
	}
	return handlers
}
