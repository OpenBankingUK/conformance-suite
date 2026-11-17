package server

import "github.com/labstack/echo/v4"

// jsonContentTypeUTF8 is the Content-Type Echo used for JSON responses before v4.12.
const jsonContentTypeUTF8 = "application/json; charset=UTF-8"

// legacyJSONSerializer keeps the "; charset=UTF-8" suffix on JSON responses that
// newer Echo versions dropped, so API responses stay identical after the upgrade.
type legacyJSONSerializer struct {
	echo.DefaultJSONSerializer
}

// Serialize restores the legacy Content-Type, then delegates to Echo's default serializer.
func (s legacyJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	if h := c.Response().Header(); h.Get(echo.HeaderContentType) == echo.MIMEApplicationJSON {
		h.Set(echo.HeaderContentType, jsonContentTypeUTF8)
	}
	return s.DefaultJSONSerializer.Serialize(c, i, indent)
}
