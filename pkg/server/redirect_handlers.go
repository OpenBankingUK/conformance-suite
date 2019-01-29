package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type RedirectFragment struct {
	Code    string `json:"code" form:"code" query:"code"`
	Scope   string `json:"scope" form:"scope" query:"scope"`
	IDToken string `json:"id_token" form:"id_token" query:"id_token"`
	State   string `json:"state" form:"state" query:"state"`
}

type RedirectQuery struct {
	Code  string `json:"code" form:"code" query:"code"`
	State string `json:"state" form:"state" query:"state"`
}

type RedirectError struct {
	ErrorDescription string `json:"error_description" form:"error_description" query:"error_description"`
	Error            string `json:"error" form:"error" query:"error"`
	State            string `json:"state" form:"state" query:"state"`
}

type redirectHandlers struct {
	logger *logrus.Entry
}

// postFragmentOKHandler - POST /api/redirect/fragment/ok
func (h *redirectHandlers) postFragmentOKHandler(c echo.Context) error {
	fragment := new(RedirectFragment)
	if err := c.Bind(fragment); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fragment)
}

// postQueryOKHandler - POST /redirect/query/ok
func (h *redirectHandlers) postQueryOKHandler(c echo.Context) error {
	query := new(RedirectQuery)
	if err := c.Bind(query); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, query)
}

// postErrorHandler - POST /api/redirect/error
func (h *redirectHandlers) postErrorHandler(c echo.Context) error {
	redirectError := new(RedirectError)
	if err := c.Bind(redirectError); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, redirectError)
}
