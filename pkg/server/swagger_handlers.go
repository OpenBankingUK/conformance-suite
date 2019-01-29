package server

import (
	"fmt"
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	openapi_middleware "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"

	"net/url"
)

// swaggerHandlers - maps paths (e.g., /swagger/account-transaction-v3.0/v3.0/docs) below to handlers.
func swaggerHandlers(logger *logrus.Entry) map[string]echo.HandlerFunc {
	handlers := map[string]echo.HandlerFunc{}

	specs := model.Specifications()
	for _, spec := range specs {
		path := "docs"

		// /swagger/account-transaction-v3.0/v3.0
		basePath := fmt.Sprintf("/swagger/%s/%s", spec.Identifier, spec.Version)
		basePathURL, err := url.Parse(basePath)
		if err != nil {
			logger.Errorf("swaggerHandlers -> cannot parse basePath=%+v, basePathURL=%+v, err=%+v", basePath, basePathURL, err)
			continue
		}

		// /swagger/account-transaction-v3.0/v3.0/docs
		fullPath := fmt.Sprintf("/swagger/%s/%s/%s", spec.Identifier, spec.Version, path)
		fullPathURL, err := url.Parse(fullPath)
		if err != nil {
			logger.Errorf("swaggerHandlers -> cannot parse fullPath=%+v, fullPathURL=%+v, err=%+v", fullPath, fullPathURL, err)
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
			SpecURL:  specURL,
			Path:     path,
		}, notFoundHandler)

		handlers[fullPath] = echo.WrapHandler(handler)
	}

	return handlers
}