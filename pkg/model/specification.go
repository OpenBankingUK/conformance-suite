package model

import (
	"errors"
	"net/url"
)

// Specification - Represents OB API specification.
// Fields are from the APIReference JSON-LD schema, see: https://schema.org/APIReference
type Specification struct {
	Identifier string
	Name       string
	// URL of confluence specifications file.
	URL *url.URL
	// Version of the specifications
	Version string
	// URL of OpenAPI/Swagger specifications file.
	SchemaVersion *url.URL
}

var (
	specifications = []Specification{
		{
			Identifier:    "account-transaction-v3.1.6",
			Name:          "Account and Transaction API Specification",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.6/profiles/account-and-transaction-api-profile.html"),
			Version:       "v3.1.6",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.6/dist/swagger/account-info-swagger.json"),
		},
		{
			Identifier:    "account-transaction-v3.1.5",
			Name:          "Account and Transaction API Specification",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.5/profiles/account-and-transaction-api-profile.html"),
			Version:       "v3.1.5",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.5/dist/swagger/account-info-swagger.json"),
		},
		{
			Identifier:    "account-transaction-v3.1.4",
			Name:          "Account and Transaction API Specification",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.4/profiles/account-and-transaction-api-profile.html"),
			Version:       "v3.1.4",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.4/dist/swagger/account-info-swagger.json"),
		},
		{
			Identifier:    "account-transaction-v3.1.3",
			Name:          "Account and Transaction API Specification",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.3/profiles/account-and-transaction-api-profile.html"),
			Version:       "v3.1.3",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.3/dist/account-info-swagger.json"),
		},
		{
			Identifier:    "account-transaction-v3.1.2",
			Name:          "Account and Transaction API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1077805296/Account+and+Transaction+API+Specification+-+v3.1.2"),
			Version:       "v3.1.2",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.2/dist/account-info-swagger.json"),
		},

		{
			Identifier:    "account-transaction-v3.1.1",
			Name:          "Account and Transaction API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/999622968/Account+and+Transaction+API+Specification+-+v3.1.1"),
			Version:       "v3.1.1",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.1/dist/account-info-swagger.json"),
		},
		{
			Identifier:    "payment-initiation-v3.1.6",
			Name:          "Payment Initiation API",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.6/profiles/payment-initiation-api-profile.html"),
			Version:       "v3.1.6",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.6/dist/swagger/payment-initiation-swagger.json"),
		},
		{
			Identifier:    "payment-initiation-v3.1.5",
			Name:          "Payment Initiation API",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.5/profiles/payment-initiation-api-profile.html"),
			Version:       "v3.1.5",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.5/dist/swagger/payment-initiation-swagger.json"),
		},
		{
			Identifier:    "payment-initiation-v3.1.4",
			Name:          "Payment Initiation API",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.4/profiles/payment-initiation-api-profile.html"),
			Version:       "v3.1.4",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.4/dist/swagger/payment-initiation-swagger.json"),
		},
		{
			Identifier:    "payment-initiation-v3.1.3",
			Name:          "Payment Initiation API",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.3/profiles/payment-initiation-api-profile.html"),
			Version:       "v3.1.3",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.3/dist/payment-initiation-swagger.json"),
		},
		{
			Identifier:    "payment-initiation-v3.1.2",
			Name:          "Payment Initiation API",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1077805743/Payment+Initiation+API+Specification+-+v3.1.2"),
			Version:       "v3.1.2",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.2/dist/payment-initiation-swagger.json"),
		},
		{
			Identifier:    "payment-initiation-v3.1.1",
			Name:          "Payment Initiation API",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/999426309/Payment+Initiation+API+Specification+-+v3.1.1"),
			Version:       "v3.1.1",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.1/dist/payment-initiation-swagger.json"),
		},
		{
			Identifier:    "confirmation-funds-v3.1.6",
			Name:          "Confirmation of Funds API Specification",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.6/profiles/confirmation-of-funds-api-profile.html"),
			Version:       "v3.1.6",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.6/dist/swagger/confirmation-funds-swagger.json"),
		},
		{
			Identifier:    "confirmation-funds-v3.1.5",
			Name:          "Confirmation of Funds API Specification",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.5/profiles/confirmation-of-funds-api-profile.html"),
			Version:       "v3.1.5",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.5/dist/swagger/confirmation-funds-swagger.json"),
		},
		{
			Identifier:    "confirmation-funds-v3.1.4",
			Name:          "Confirmation of Funds API Specification",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.4/profiles/confirmation-of-funds-api-profile.html"),
			Version:       "v3.1.4",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.4/dist/swagger/confirmation-funds-swagger.json"),
		},
		{
			Identifier:    "confirmation-funds-v3.1.3",
			Name:          "Confirmation of Funds API Specification",
			URL:           mustParseURL("https://openbankinguk.github.io/read-write-api-site3/v3.1.3/profiles/confirmation-of-funds-api-profile.html"),
			Version:       "v3.1.3",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.3/dist/confirmation-funds-swagger.json"),
		},
		{
			Identifier:    "confirmation-funds-v3.1.2",
			Name:          "Confirmation of Funds API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1077806537/Confirmation+of+Funds+API+Specification+-+v3.1.2"),
			Version:       "v3.1.2",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.2/dist/confirmation-funds-swagger.json"),
		},
		{
			Identifier:    "confirmation-funds-v3.1.1",
			Name:          "Confirmation of Funds API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1000015607/Confirmation+of+Funds+API+Specification+-+v3.1.1"),
			Version:       "v3.1.1",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.1/dist/confirmation-funds-swagger.json"),
		},
		{
			Identifier:    "event-notification-aspsp-v3.1.2",
			Name:          "Event Notification API Specification - ASPSP Endpoints",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1077806617/Event+Notification+API+Specification+-+v3.1.2"),
			Version:       "v3.1.2",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.2/dist/event-subscriptions-swagger.json"),
		},
		{
			Identifier:    "event-notification-aspsp-v3.1",
			Name:          "Event Notification API Specification - ASPSP Endpoints",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/1000114043/Event+Notification+API+Specification+-+v3.1.1"),
			Version:       "v3.1.1",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.1/dist/callback-urls-swagger.yaml"),
		},

		{
			Identifier:    "account-transaction-v3.1",
			Name:          "Account and Transaction API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1"),
			Version:       "v3.1.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json"),
		},
		{
			Identifier:    "payment-initiation-v3.1",
			Name:          "Payment Initiation API",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937754701/Payment+Initiation+API+Specification+-+v3.1"),
			Version:       "v3.1.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/payment-initiation-swagger.json"),
		},
		{
			Identifier:    "confirmation-funds-v3.1",
			Name:          "Confirmation of Funds API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937951380/Confirmation+of+Funds+API+Specification+-+v3.1"),
			Version:       "v3.1.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/confirmation-funds-swagger.json"),
		},
		{
			Identifier:    "event-notification-aspsp-v3.1",
			Name:          "Event Notification API Specification - ASPSP Endpoints",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937951397/Event+Notification+API+Specification+-+v3.1"),
			Version:       "v3.1.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/callback-urls-swagger.yaml"),
		},
		{
			Identifier:    "event-notification-tpp-v3.1",
			Name:          "Event Notification API Specification - TPP Endpoints",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937951397/Event+Notification+API+Specification+-+v3.1"),
			Version:       "v3.1.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/event-notifications-swagger.json"),
		},
		{
			Identifier:    "account-transaction-v3.0",
			Name:          "Account and Transaction API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0"),
			Version:       "v3.0.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"),
		},
		{
			Identifier:    "payment-initiation-v3.0",
			Name:          "Payment Initiation API",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/645367011/Payment+Initiation+API+Specification+-+v3.0"),
			Version:       "v3.0.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/payment-initiation-swagger.json"),
		},
		{
			Identifier:    "confirmation-funds-v3.0",
			Name:          "Confirmation of Funds API Specification",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/645203467/Confirmation+of+Funds+API+Specification+-+v3.0"),
			Version:       "v3.0.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/confirmation-funds-swagger.json"),
		},
		{
			Identifier:    "event-notification-aspsp-v3.0",
			Name:          "Event Notification API Specification - ASPSP Endpoints",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/645367055/Event+Notification+API+Specification+-+v3.0"),
			Version:       "v3.0.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/callback-urls-swagger.yaml"),
		},
		{
			Identifier:    "event-notification-tpp-v3.0",
			Name:          "Event Notification API Specification - TPP Endpoints",
			URL:           mustParseURL("https://openbanking.atlassian.net/wiki/spaces/DZ/pages/645367055/Event+Notification+API+Specification+-+v3.0"),
			Version:       "v3.0.0",
			SchemaVersion: mustParseURL("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/event-notifications-swagger.yaml"),
		},
	}
)

// Specifications - get a clone of the `specifications` array.
func Specifications() []Specification {
	clone := make([]Specification, len(specifications))
	copy(clone, specifications)
	return clone
}

// SpecificationFromSchemaVersion - returns specification struct
// for given schema version URL, or nil when there is no match.
func SpecificationFromSchemaVersion(schemaVersion string) (Specification, error) {
	var spec Specification
	for _, specification := range specifications {
		if specification.SchemaVersion.String() == schemaVersion {
			return specification, nil
		}
	}
	return spec, errors.New("no specifications found for schema version: " + schemaVersion)
}

func mustParseURL(rawurl string) *url.URL {
	parsedUrl, err := url.Parse(rawurl)
	if err != nil {
		panic(rawurl)
	}
	return parsedUrl
}
