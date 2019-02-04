package server

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"

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
	logger *logrus.Entry
}

// postFragmentOKHandler - POST /api/redirect/fragment/ok
func (h *redirectHandlers) postFragmentOKHandler(c echo.Context) error {
	fragment := new(RedirectFragment)
	if err := c.Bind(fragment); err != nil {
		return err
	}

	claim := &AuthClaim{}

	// If not providing Keyfunc (3rd param), don't check for error here
	// as it will always be error("no Keyfunc was provided")
	t, _ := jwt.ParseWithClaims(fragment.IDToken, claim, nil)

	if !t.Valid {
		logrus.Warn("token not valid")
	}

	cHash := h.calculateCHash(fragment.Code, true)
	if cHash == claim.CHash {
		return c.JSON(http.StatusOK, fragment)
	}
	return c.JSON(http.StatusBadRequest, fragment)
}

// postQueryOKHandler - POST /redirect/query/ok
func (h *redirectHandlers) postQueryOKHandler(c echo.Context) error {
	query := new(RedirectQuery)
	if err := c.Bind(query); err != nil {
		return err
	}

	claim := &AuthClaim{}

	// If not providing Keyfunc (3rd param), don't check for error here
	// as it will always be error("no Keyfunc was provided")
	t, _ := jwt.ParseWithClaims(query.IDToken, claim, nil)

	if !t.Valid {
		logrus.Warn("token not valid")
	}

	cHash := h.calculateCHash(query.Code, true)
	if cHash == claim.CHash {
		return c.JSON(http.StatusOK, query)
	}
	return c.JSON(http.StatusBadRequest, query)
}

// postErrorHandler - POST /api/redirect/error
func (h *redirectHandlers) postErrorHandler(c echo.Context) error {
	redirectError := new(RedirectError)
	if err := c.Bind(redirectError); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, redirectError)
}

// calculateCHash calculates the code hash (c_hash) value
// as described in section 3.3.2.11 (ID Token) https://openid.net/specs/openid-connect-core-1_0.html#HybridIDToken
// There is an option to trim the `==` from the end of the base64url encoded string, by setting `trim=true`
func (h *redirectHandlers) calculateCHash(code string, trim bool) string {
	digest := sha256.Sum256([]byte(code))
	left := digest[0 : len(digest)/2]
	b64 := base64.URLEncoding.EncodeToString(left)

	if trim {
		return strings.Trim(b64, "=")
	}

	return b64
}
