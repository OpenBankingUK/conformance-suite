// This is still WORK IN PROGRESS. The handlers just return either an empty
// `github.com/OpenBankingUK/conformance-suite/pkg/server/models.ImportReviewResponse` or
//  `github.com/OpenBankingUK/conformance-suite/pkg/server/models.ImportRerunResponse` and do not do the
// importing or review functionality. This will be implemented as we go along.

package server

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/OpenBankingUK/conformance-suite/pkg/server/models"
)

type authHandlers struct {
	logger       *logrus.Entry
	secretJWTKey string
}

func newAuthHandlers(logger *logrus.Entry) authHandlers {
	return authHandlers{
		secretJWTKey: "todo_secret_key_for_jwt_signing",
		logger:       logger.WithField("handler", "authHandlers"),
	}
}

// postLogin - `/api/login` POST.
func (h authHandlers) postLogin(c echo.Context) error {
	logger := h.logger.WithField("function", "postLogin")

	request := models.AuthRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	if err := request.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	// check credentials

	var jwtKey = []byte(h.secretJWTKey)

	claims := &authClaims{
		Email: request.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logger.Errorf("Error creating JWT: %v", err)
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(err))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

type authClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (c *authClaims) Valid() error {
	return nil
}
