package server

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"

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

	// If ID Token has not been set in the query, then there is no need to validate
	// (Nothing to validate)
	if fragment.IDToken == "" {
		if fragment.Code != "" {
			return c.JSON(http.StatusOK, nil)
		}
		return c.JSON(http.StatusBadRequest, errors.New("code not set"))
	}

	claim := &AuthClaim{}

	t, err := jwt.ParseWithClaims(fragment.IDToken, claim, nil)
	if err != nil {
		// If not providing Keyfunc (3rd param), don't check for error here
		// as it will always be error("no Keyfunc was provided")
		h.logger.Debug("Keyfunc not provided")
	}

	cHash, err := calculateCHash(t.Header["alg"].(string), fragment.Code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	if cHash == claim.CHash {
		return c.JSON(http.StatusOK, nil)
	}

	resp := NewErrorResponse(errors.New("c_hash invalid"))
	return c.JSON(http.StatusBadRequest, resp)
}

// postQueryOKHandler - POST /redirect/query/ok
func (h redirectHandlers) postQueryOKHandler(c echo.Context) error {
	query := new(RedirectQuery)
	if err := c.Bind(query); err != nil {
		return err
	}

	// If ID Token has not been set in the query, then there is no need to validate
	// (Nothing to validate)
	if query.IDToken == "" {
		if query.Code != "" {
			err := h.handleCodeExchange(query)
			if err != nil {
				resp := NewErrorResponse(errors.Wrap(err, "unable to handle redirect"))
				return c.JSON(http.StatusBadRequest, resp)
			}
			return c.JSON(http.StatusOK, nil)
		}
		return c.JSON(http.StatusBadRequest, errors.New("code not set"))
	}

	claim := &AuthClaim{}

	t, err := jwt.ParseWithClaims(query.IDToken, claim, nil)
	if err != nil {
		// If not providing Keyfunc (3rd param), don't check for error here
		// as it will always be error("no Keyfunc was provided")
		h.logger.Debug("Keyfunc not provided")
	}

	cHash, err := calculateCHash(t.Header["alg"].(string), query.Code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	if cHash == claim.CHash {
		err = h.handleCodeExchange(query)
		if err != nil {
			resp := NewErrorResponse(errors.Wrap(err, "unable to handle redirect"))
			return c.JSON(http.StatusBadRequest, resp)
		}
		return c.JSON(http.StatusOK, nil)
	}

	resp := NewErrorResponse(errors.New("c_hash invalid"))
	return c.JSON(http.StatusBadRequest, resp)
}

func (h redirectHandlers) handleCodeExchange(query *RedirectQuery) error {
	logrus.StandardLogger().Warnf("received redirect Query %#v", query)
	return h.journey.CollectToken(query.Code, query.State, query.Scope)
}

// postErrorHandler - POST /api/redirect/error
func (h redirectHandlers) postErrorHandler(c echo.Context) error {
	redirectError := new(RedirectError)
	if err := c.Bind(redirectError); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, redirectError)
}

// calculateCHash calculates the code hash (c_hash) value
// as described in section 3.3.2.11 (ID Token) https://openid.net/specs/openid-connect-core-1_0.html#HybridIDToken
// List of valid algorithms https://openid.net/specs/openid-financial-api-part-2.html#jws-algorithm-considerations
// At the time of writing, the list shows "PS256", "ES256"
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/83919096/Open+Banking+Security+Profile+-+Implementer+s+Draft+v1.1.2#OpenBankingSecurityProfile-Implementer'sDraftv1.1.2-Step2:FormtheJOSEHeader
func calculateCHash(alg string, code string) (string, error) {
	var digest []byte

	switch alg {
	case "ES256", "PS256":
		d := sha256.Sum256([]byte(code))
		//left most 256 bits.. 256/8 = 32bytes
		// no need to validate length as sha256.Sum256 returns fixed length
		digest = d[0:32]
	default:
		return "", fmt.Errorf("%s algorithm not supported", alg)
	}

	left := digest[0 : len(digest)/2]
	return base64.RawURLEncoding.EncodeToString(left), nil
}
