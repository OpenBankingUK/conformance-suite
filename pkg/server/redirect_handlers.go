package server

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RedirectFragment struct {
	Code    string `json:"code" form:"code" query:"code"`
	Scope   string `json:"scope" form:"scope" query:"scope"`
	IDToken string `json:"id_token" form:"id_token" query:"id_token"`
	State   string `json:"state" form:"state" query:"state"`
}

type RedirectQuery struct {
	Code    string `json:"code" form:"code" query:"code"`
	Scope   string `json:"scope" form:"scope" query:"scope"`
	IDToken string `json:"id_token" form:"id_token" query:"id_token"`
	State   string `json:"state" form:"state" query:"state"`
}

type RedirectError struct {
	ErrorDescription string `json:"error_description" form:"error_description" query:"error_description"`
	Error            string `json:"error" form:"error" query:"error"`
	State            string `json:"state" form:"state" query:"state"`
}

// AuthClaim represents an in coming JWT from third part ASPSP as part of authentication/consent
// process during `Hybrid Flow Authentication`
// https://openid.net/specs/openid-connect-core-1_0.html#HybridFlowAuth
type AuthClaim struct {
	jwt.StandardClaims
	AuditTrackingID     string `json:"auditTrackingId"`
	TokenName           string `json:"tokenName"`
	Nonce               string `json:"nonce"`
	Acr                 string `json:"acr"`
	CHash               string `json:"c_hash"`
	OpenBankingIntentID string `json:"openbanking_intent_id"`
	SHash               string `json:"s_hash"`
	Azp                 string `json:"azp"`
	AuthTime            int    `json:"auth_time"`
	Realm               string `json:"realm"`
	TokenType           string `json:"tokenType"`
}

type redirectHandlers struct {
	journey Journey
	logger  *logrus.Entry
}

func newRedirectHandlers(journey Journey, logger *logrus.Entry) redirectHandlers {
	return redirectHandlers{
		journey: journey,
		logger:  logger.WithField("module", "redirectHandlers"),
	}
}

// postFragmentOKHandler - POST /api/redirect/fragment/ok
func (h redirectHandlers) postFragmentOKHandler(c echo.Context) error {
	fragment := new(RedirectFragment)
	if err := c.Bind(fragment); err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"function":         "postFragmentOKHandler",
		"RedirectFragment": fragment,
	}).Warn("Received fragment in redirect")

	// If ID Token has not been set in the fragment, then there is no need to validate
	// (Nothing to validate)
	if fragment.IDToken == "" {
		if fragment.Code != "" {
			err := h.handleCodeExchange(fragment.Code, fragment.State, fragment.Scope)
			if err != nil {
				resp := NewErrorResponse(errors.Wrap(err, "unable to handle redirect"))
				return c.JSON(http.StatusBadRequest, resp)
			}
			return c.JSON(http.StatusOK, nil)
		}
		return c.JSON(http.StatusBadRequest, errors.New("code not set"))
	}

	err := h.handleCodeExchange(fragment.Code, fragment.State, fragment.Scope)
	if err != nil {
		resp := NewErrorResponse(errors.Wrap(err, "unable to handle redirect"))
		return c.JSON(http.StatusBadRequest, resp)
	}
	return c.JSON(http.StatusOK, nil)

	// TODO(mbana): Turned off validation for now.
	// claim := &AuthClaim{}

	// t, err := jwt.ParseWithClaims(fragment.IDToken, claim, nil)
	// if err != nil {
	// 	// If not providing Keyfunc (3rd param), don't check for error here
	// 	// as it will always be error("no Keyfunc was provided")
	// 	h.logger.Debug("Keyfunc not provided")
	// }

	// cHash, err := calculateCHash(t.Header["alg"].(string), fragment.Code)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	// }
	// if cHash == claim.CHash {
	// 	err := h.handleCodeExchange(fragment.Code, fragment.State, fragment.Scope)
	// 	if err != nil {
	// 		resp := NewErrorResponse(errors.Wrap(err, "unable to handle redirect"))
	// 		return c.JSON(http.StatusBadRequest, resp)
	// 	}
	// 	return c.JSON(http.StatusOK, nil)
	// }

	// resp := NewErrorResponse(errors.New("c_hash invalid"))
	// return c.JSON(http.StatusBadRequest, resp)
}

// postQueryOKHandler - POST /redirect/query/ok
func (h redirectHandlers) postQueryOKHandler(c echo.Context) error {
	query := new(RedirectQuery)
	if err := c.Bind(query); err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"function":      "postQueryOKHandler",
		"RedirectQuery": query,
	}).Warn("Received query in redirect")

	// If ID Token has not been set in the query, then there is no need to validate
	// (Nothing to validate)
	if query.IDToken == "" {
		if query.Code != "" {
			err := h.handleCodeExchange(query.Code, query.State, query.Scope)
			if err != nil {
				resp := NewErrorResponse(errors.Wrap(err, "unable to handle redirect"))
				return c.JSON(http.StatusBadRequest, resp)
			}
			return c.JSON(http.StatusOK, nil)
		}
		return c.JSON(http.StatusBadRequest, errors.New("code not set"))
	}

	err := h.handleCodeExchange(query.Code, query.State, query.Scope)
	if err != nil {
		resp := NewErrorResponse(errors.Wrap(err, "unable to handle redirect"))
		return c.JSON(http.StatusBadRequest, resp)
	}
	return c.JSON(http.StatusOK, nil)

	// TODO(mbana): Turned off validation for now.
	// claim := &AuthClaim{}

	// t, err := jwt.ParseWithClaims(query.IDToken, claim, nil)
	// if err != nil {
	// 	// If not providing Keyfunc (3rd param), don't check for error here
	// 	// as it will always be error("no Keyfunc was provided")
	// 	h.logger.Debug("Keyfunc not provided")
	// }

	// cHash, err := calculateCHash(t.Header["alg"].(string), query.Code)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	// }
	// if cHash == claim.CHash {
	// 	err := h.handleCodeExchange(query.Code, query.State, query.Scope)
	// 	if err != nil {
	// 		resp := NewErrorResponse(errors.Wrap(err, "unable to handle redirect"))
	// 		return c.JSON(http.StatusBadRequest, resp)
	// 	}
	// 	return c.JSON(http.StatusOK, nil)
	// }

	// resp := NewErrorResponse(errors.New("c_hash invalid"))
	// return c.JSON(http.StatusBadRequest, resp)
}

func (h redirectHandlers) handleCodeExchange(code string, state string, scope string) error {
	h.logger.WithFields(logrus.Fields{
		"function": "handleCodeExchange",
		"code":     code,
		"state":    state,
		"scope":    scope,
	}).Info("h.journey.CollectToken ...")
	return h.journey.CollectToken(code, state, scope)
}

// postErrorHandler - POST /api/redirect/error
func (h redirectHandlers) postErrorHandler(c echo.Context) error {
	redirectError := new(RedirectError)
	if err := c.Bind(redirectError); err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"function":      "postErrorHandler",
		"RedirectError": redirectError,
	}).Warn("Received error in redirect")

	return c.JSON(http.StatusOK, redirectError)
}
