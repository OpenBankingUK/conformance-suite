package schemaprops

import (
	"errors"
	"regexp"

	"github.com/sirupsen/logrus"
)

var (
	acctPay  = append(accountsRegex, paymentsRegex...)
	allregex = append(acctPay, cbpiiRegex...)
)

func pathToSwagger(path string) (string, error) {
	for _, regPath := range allregex {
		matched, err := regexp.MatchString(regPath.Regex, path)
		if err != nil {
			return "", errors.New("path mapping error")
		}
		if matched {
			return regPath.Mapping, nil
		}
	}
	logrus.Tracef("Unknown swagger path for %s", path)
	return "", errors.New("Unknown swaggerPath for " + path)
}

func mapPathsToSwagger(endpoints []string) []string {
	lookupMap := make(map[string]string)
	for _, ep := range endpoints {
		for _, regPath := range allregex {
			matched, err := regexp.MatchString(regPath.Regex, ep)
			if err != nil {
				continue
			}
			if matched {
				lookupMap[regPath.Mapping] = ep
			}
		}
	}

	paths := sortPathStrings(lookupMap)
	return paths
}

//Note: these regexs are similar to those found in manifest/script.go
// however these require a different start of the expression (.*) rather than (^)
var accountsRegex = []PathRegex{
	{
		Regex:   ".*/accounts$",
		Name:    "Get Accounts",
		Mapping: "/accounts",
	},
	{
		Regex:   ".*/account-access-consents$",
		Name:    "Get Accounts",
		Mapping: "/account-access-consents",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "$",
		Name:    "Get Accounts Resource",
		Mapping: "/accounts/{AccountId}",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/balances$",
		Name:    "Get Balances Resource",
		Mapping: "/accounts/{AccountId}/balances",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/beneficiaries$",
		Name:    "Get Beneficiaries Resource",
		Mapping: "/accounts/{AccountId}/beneficiaries",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/direct-debits$",
		Name:    "Get Direct Debits Resource",
		Mapping: "/accounts/{AccountId}/direct-debits",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/offers$",
		Name:    "Get Offers Resource",
		Mapping: "/accounts/{AccountId}/offers",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/party$",
		Name:    "Get Party Resource",
		Mapping: "/accounts/{AccountId}/party",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/product$",
		Name:    "Get Product Resource",
		Mapping: "/accounts/{AccountId}/products",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/scheduled-payments$",
		Name:    "Get Scheduled Payment resource",
		Mapping: "/accounts/{AccountId}/scheduled-payments",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/standing-orders$",
		Name:    "Get Standing Orders resource",
		Mapping: "/accounts/{AccountId}/standing-orders",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/statements$",
		Name:    "Get Statements Resource",
		Mapping: "/accounts/{AccountId}/statements",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/statements/" + subPathx + "/file$",
		Name:    "Get statement files resource",
		Mapping: "/accounts/{AccountId}/statements/{StatementId}/file",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/statements/" + subPathx + "/transactions$",
		Name:    "Get statement transactions resource",
		Mapping: "/accounts/{AccountId}/statements/{StatementId}/transactions",
	},
	{
		Regex:   ".*/accounts/" + subPathx + "/transactions$",
		Name:    "Get transactions resource",
		Mapping: "/accounts/{AccountId}/transactions",
	},
	{
		Regex:   ".*/balances$",
		Name:    "Get Balances",
		Mapping: "/balances",
	},
	{
		Regex:   ".*/beneficiaries$",
		Name:    "Get Beneficiaries",
		Mapping: "/beneficiaries",
	},
	{
		Regex:   ".*/direct-debits$",
		Name:    "Get directory debits",
		Mapping: "/direct-debits",
	},
	{
		Regex:   ".*/offers$",
		Name:    "Get Offers",
		Mapping: "/offers",
	},
	{
		Regex:   ".*/party$",
		Name:    "Get party",
		Mapping: "/party",
	},
	{
		Regex:   ".*/products$",
		Name:    "Get Products",
		Mapping: "/products",
	},

	{
		Regex:   ".*/scheduled-payments$",
		Name:    "Get Payments",
		Mapping: "/scheduled-payments",
	},
	{
		Regex: ".*/standing-orders$",
		Name:  "/standing-orders",
	},
	{
		Regex:   ".*/statements$",
		Name:    "Get Statements",
		Mapping: "/statements",
	},
	{
		Regex:   ".*/transactions$",
		Name:    "Get Transactions",
		Mapping: "/transactions",
	},
}

var paymentsRegex = []PathRegex{
	{
		Regex:   ".*/domestic-payment-consents$",
		Method:  "POST",
		Name:    "Create a domestic payment consent",
		Mapping: "/domestic-payment-consents",
	},
	{
		Regex:   ".*/domestic-payment-consents/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get domestic payment consent by by consent ID",
		Mapping: "/domestic-payment-consents/{ConsentId}",
	},
	{
		Regex:   ".*/domestic-payment-consents/" + subPathx + "/funds-confirmation$",
		Method:  "GET",
		Name:    "Get domestic payment consents funds confirmation, by consentID",
		Mapping: "/domestic-payment-consents/{ConsentId}/funds-confirmation",
	},
	{
		Regex:   ".*/domestic-payments$",
		Method:  "POST",
		Name:    "Create a domestic payment",
		Mapping: "/domestic-payments",
	},
	{
		Regex:   ".*/domestic-payments/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get domestic payment by domesticPaymentID",
		Mapping: "/domestic-payments/{DomesticPaymentId}",
	},
	{
		Regex:   ".*/domestic-scheduled-payment-consents$",
		Method:  "POST",
		Name:    "Create a domestic scheduled payment consent",
		Mapping: "/domestic-scheduled-payment-consents",
	},
	{
		Regex:   ".*/domestic-scheduled-payment-consents/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get domestic scheduled payment consent by consentID",
		Mapping: "/domestic-scheduled-payment-consents/{ConsentId}",
	},
	{
		Regex:   ".*/domestic-scheduled-payments$",
		Method:  "POST",
		Name:    "Create a domestic scheduled payment",
		Mapping: "/domestic-scheduled-payments",
	},
	{
		Regex:   ".*/domestic-scheduled-payment/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get domestic scheduled payments by consentID",
		Mapping: "/domestic-scheduled-payment/{ConsentId}",
	},
	{
		Regex:   ".*/domestic-standing-order-consents$",
		Method:  "POST",
		Name:    "Create a domestic standing order consent",
		Mapping: "/domestic-standing-order-consents",
	},
	{
		Regex:   ".*/domestic-standing-order-consents/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get domestic standing order consent by consentID",
		Mapping: "/domestic-standing-order-consents/(ConsentId}",
	},
	{
		Regex:   ".*/domestic-standing-orders$",
		Method:  "POST",
		Name:    "Create a domestic standing order",
		Mapping: "/domestic-standing-order",
	},
	{
		Regex:   ".*/domestic-standing-orders/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get domestic standing order by domesticStandingOrderID",
		Mapping: "/domestic-standing-orders/{DomesticStandingOrderId}",
	},
	{
		Regex:   ".*/international-payment-consents$",
		Method:  "POST",
		Name:    "Create an international payment consent",
		Mapping: "/international-payment-consents",
	},
	{
		Regex:   ".*/international-payment-consents/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get international payment consent by consentID",
		Mapping: "/international-payment-consents/{ConsentId}",
	},
	{
		Regex:   ".*/international-payment-consents/" + subPathx + "/funds-confirmation$",
		Method:  "GET",
		Name:    "Get international payment consent funds confirmation by consentID",
		Mapping: "/international-payment-consents/{ConsentId}/funds-confirmation",
	},
	{
		Regex:   ".*/international-payments$",
		Method:  "POST",
		Name:    "Create an international payment",
		Mapping: "/international-payments",
	},
	{
		Regex:   ".*/international-payments/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get international payment by internationalPaymentID",
		Mapping: "/international-payments/{InternationalPaymentId}",
	},
	{
		Regex:   ".*/international-scheduled-payment-consents$",
		Method:  "POST",
		Name:    "Create an international scheduled payment consent",
		Mapping: "/international-scheduled-payment-consent",
	},
	{
		Regex:   ".*/international-scheduled-payment-consents/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get international scheduled payment consents by consentID",
		Mapping: "/international-scheduled-payment-consent/{ConsentId}",
	},
	{
		Regex:   ".*/international-scheduled-payments/" + subPathx + "/funds-confirmation$",
		Method:  "GET",
		Name:    "Get international scheduled payment funds confirmation by consentID",
		Mapping: "/international-scheduled-payments/{ConsentId}/funds-confirmation",
	},
	{
		Regex:   ".*/international-scheduled-payments$",
		Method:  "POST",
		Name:    "Create an international scheduled payment",
		Mapping: "/international-scheduled-payments",
	},
	{
		Regex:   ".*/international-scheduled-payments/" + subPathx + "$",
		Method:  "GET",
		Name:    "Create an international scheduled payment by internationalScheduledPaymentID",
		Mapping: "/international-scheduled-payments/{InternationalSchedulatedPaymentId}",
	},
	{
		Regex:   ".*/international-standing-order-consents$",
		Method:  "POST",
		Name:    "Create international standing order consent",
		Mapping: "/international-standing-order-consents",
	},
	{
		Regex:   ".*/international-standing-order-consents/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get international standing order consent by consentID",
		Mapping: "/international-standing-order-consents/{ConsentId}",
	},
	{
		Regex:   ".*/international-standing-orders$",
		Method:  "POST",
		Name:    "Create international standing order",
		Mapping: "/international-standing-orders",
	},
	{
		Regex:   ".*/international-standing-orders/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get an international standing order by internationalStandingOrderID",
		Mapping: "/international-standing-orders/{InternationalStandingOrderPaymentId}",
	},
	{
		Regex:   ".*/file-payment-consents$",
		Method:  "POST",
		Name:    "Create a file payment consent",
		Mapping: "/file-payment-consents",
	},
	{
		Regex:   ".*/file-payment-consents/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get a file payment consent by consentID",
		Mapping: "/file-payment-consents/{ConsentId}",
	},
	{
		Regex:   ".*/file-payment-consents/" + subPathx + "/file$",
		Method:  "POST",
		Name:    "Create a file payment consent file by consentID",
		Mapping: "/file-payment-consents/{ConsentId}/file",
	},
	{
		Regex:   ".*/file-payment-consents/" + subPathx + "/file$",
		Method:  "GET",
		Name:    "Get a file payment consents file by consentID",
		Mapping: "/file-payment-consents/{ConsentId}/file",
	},
	{
		Regex:   ".*/file-payments$",
		Method:  "POST",
		Name:    "Create a file payment",
		Mapping: "/file-payments",
	},
	{
		Regex:   ".*/file-payments/" + subPathx + "$",
		Method:  "GET",
		Name:    "Get a file payment by filePaymentID",
		Mapping: "/file-payments/{FilePaymentId}",
	},
	{
		Regex:   ".*/file-payments/" + subPathx + "/report-file$",
		Method:  "GET",
		Name:    "Get a file payment report file by filePaymentID",
		Mapping: "/file-payments/{FilePaymentId}/report-files",
	},
}

var cbpiiRegex = []PathRegex{
	{
		Regex:   ".*/funds-confirmation-consents$",
		Method:  "POST",
		Name:    "Create Funds Confirmation Consent",
		Mapping: "/funds-confirmation-consents",
	},
	{
		Regex:   ".*/funds-confirmation-consents/" + subPathx + "$",
		Method:  "GET",
		Name:    "Retrieve Funds Confirmation Consent",
		Mapping: "/funds-confirmation-consents/{ConsentId}",
	},
	{
		Regex:   ".*/funds-confirmation-consents/" + subPathx + "$",
		Method:  "DELETE",
		Name:    "Delete Funds Confirmation Consent",
		Mapping: "/funds-confirmation-consents/{ConsentId}",
	},
	{
		Regex:   ".*/funds-confirmations$",
		Method:  "POST",
		Name:    "Create Funds Confirmation",
		Mapping: "",
	},
}
