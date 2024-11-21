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
	journey      Journey
}

func newAuthHandlers(logger *logrus.Entry, journey Journey) authHandlers {
	return authHandlers{
		secretJWTKey: "todo_secret_key_for_jwt_signing",
		journey:      journey,
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

	user, err := h.journey.GetUserByEmail(request.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(err))
	}

	var jwtKey = []byte(h.secretJWTKey)

	claims := &authClaims{
		Email:  user.Email,
		UserID: user.ID,
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
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (c *authClaims) Valid() error {
	return nil
}

func (h authHandlers) JWTFromCookie() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("ob_jwt")
			if err != nil {
				return next(c)
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return []byte(h.secretJWTKey), nil
			})

			if err != nil || !token.Valid {
				return next(c)
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if email, ok := claims["email"].(string); ok {
					c.Set("user_id", email)
				}
			}

			return next(c)
		}
	}
}
