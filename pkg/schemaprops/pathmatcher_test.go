package schemaprops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchCallPathsToSwaggerPaths(t *testing.T) {
	swaggerPaths := mapPathsToSwagger(sourcepaths)
	assert.Equal(t, swaggerPaths, resultPaths)
}

var resultPaths = []string{"/account-access-consents", "/accounts", "/accounts/{AccountId}", "/accounts/{AccountId}/balances",
	"/accounts/{AccountId}/beneficiaries", "/accounts/{AccountId}/direct-debits",
	"/accounts/{AccountId}/offers", "/accounts/{AccountId}/party", "/accounts/{AccountId}/products",
	"/accounts/{AccountId}/scheduled-payments", "/accounts/{AccountId}/standing-orders",
	"/accounts/{AccountId}/transactions", "/balances", "/beneficiaries", "/direct-debits",
	"/domestic-payment-consents", "/domestic-payment-consents/{ConsentId}",
	"/domestic-payment-consents/{ConsentId}/funds-confirmation", "/domestic-payments",
	"/domestic-payments/{DomesticPaymentId}", "/domestic-scheduled-payment-consents",
	"/domestic-scheduled-payment-consents/{ConsentId}", "/domestic-standing-order",
	"/domestic-standing-order-consents", "/domestic-standing-order-consents/(ConsentId}",
	"/domestic-standing-orders/{DomesticStandingOrderId}", "/funds-confirmation-consents",
	"/funds-confirmation-consents/{ConsentId}", "/international-payment-consents",
	"/international-payment-consents/{ConsentId}", "/international-payments",
	"/international-payments/{InternationalPaymentId}", "/international-scheduled-payment-consent",
	"/international-scheduled-payment-consent/{ConsentId}", "/international-scheduled-payments",
	"/international-scheduled-payments/{InternationalSchedulatedPaymentId}", "/offers",
	"/party", "/products", "/scheduled-payments", "/statements", "/transactions"}

var sourcepaths = []string{
	"/open-banking/v3.1/cbpii/funds-confirmation-consents/42",
	"/open-banking/v3.1/cbpii/funds-confirmation-consents/fcc-c581c9c2-d3d1-4771-8402-524f3a9ed29d",
	"/open-banking/v3.1/aisp/accounts",
	"/open-banking/v3.1/aisp/accounts",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/balances",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/balances",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/balances/foobar",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/beneficiaries",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/beneficiaries",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/beneficiaries/foobar",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/direct-debits",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/direct-debits/foobar",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/foobar",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/offers",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/offers/foobar",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/party",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/party/foobar",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/product",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/product/foobar",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/scheduled-payments",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/scheduled-payments/foobar",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/standing-orders",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/standing-orders/foobar",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/transactions",
	"/open-banking/v3.1/aisp/accounts/700004000000000000000005/transactions",
	"/open-banking/v3.1/aisp/accounts/foobar/balances",
	"/open-banking/v3.1/aisp/accounts/foobar/beneficiaries",
	"/open-banking/v3.1/aisp/accounts/foobar/direct-debits",
	"/open-banking/v3.1/aisp/accounts/foobar/offers",
	"/open-banking/v3.1/aisp/accounts/foobar/party",
	"/open-banking/v3.1/aisp/accounts/foobar/product",
	"/open-banking/v3.1/aisp/accounts/foobar/scheduled-payments",
	"/open-banking/v3.1/aisp/accounts/foobar/standing-orders",
	"/open-banking/v3.1/aisp/accounts/foobar/transactions",
	"/open-banking/v3.1/aisp/balances",
	"/open-banking/v3.1/aisp/balances",
	"/open-banking/v3.1/aisp/beneficiaries",
	"/open-banking/v3.1/aisp/beneficiaries",
	"/open-banking/v3.1/aisp/direct-debits",
	"/open-banking/v3.1/aisp/direct-debits",
	"/open-banking/v3.1/aisp/foobar",
	"/open-banking/v3.1/aisp/offers",
	"/open-banking/v3.1/aisp/offers",
	"/open-banking/v3.1/aisp/party",
	"/open-banking/v3.1/aisp/party",
	"/open-banking/v3.1/aisp/products",
	"/open-banking/v3.1/aisp/products",
	"/open-banking/v3.1/aisp/scheduled-payments",
	"/open-banking/v3.1/aisp/scheduled-payments",
	"/open-banking/v3.1/aisp/standing-orders",
	"/open-banking/v3.1/aisp/standing-orders",
	"/open-banking/v3.1/aisp/statements",
	"/open-banking/v3.1/aisp/transactions",
	"/open-banking/v3.1/aisp/transactions",
	"/open-banking/v3.1/cbpii/funds-confirmation-consents/fcc-c581c9c2-d3d1-4771-8402-524f3a9ed29d",
	"/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-bbc1af3e-ca2e-4102-8f66-69b06d168816",
	"/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-bbc1af3e-ca2e-4102-8f66-69b06d168816/funds-confirmation",
	"/open-banking/v3.1/pisp/domestic-payments/pv3-6505bb74-f271-4d6f-a52e-240629c357c2",
	"/open-banking/v3.1/pisp/domestic-scheduled-payment-consents/sdp-2-17785de3-e0ad-4105-b5e8-9e5f522e2f91",
	"/open-banking/v3.1/pisp/domestic-scheduled-payment-consents/sdp-2-caf40218-92d0-444c-a00b-a86efd44e830",
	"/open-banking/v3.1/pisp/domestic-standing-order-consents/sdp-3-16c5e60a-ac52-48b7-9c0b-3031626deee0",
	"/open-banking/v3.1/pisp/domestic-standing-orders/pv3-862c4dde-8ad4-417e-b28e-b9fc6db4d7d3",
	"/open-banking/v3.1/pisp/international-payment-consents/sdp-4-05c9fc22-37aa-458b-ad6e-ab1889765d13",
	"/open-banking/v3.1/pisp/international-payments/pv3-60fbf2f7-440d-45a0-bbb9-24756be4c744",
	"/open-banking/v3.1/pisp/international-scheduled-payment-consents/sdp-5-00beaf93-bac1-43af-9c18-3398e758ec6f",
	"/open-banking/v3.1/pisp/international-scheduled-payments/pv3-09aa13bc-5e9a-46c4-bb72-83c28f4515ed",
	"/open-banking/https://as1.obie.uk.ozoneapi.io/token",
	"/open-banking/v3.1/aisp/account-access-consents",
	"/open-banking/v3.1/cbpii/funds-confirmation-consents",
	"/open-banking/v3.1/cbpii/funds-confirmation-consents",
	"/open-banking/v3.1/cbpii/funds-confirmations",
	"/open-banking/v3.1/pisp/domestic-payment-consents",
	"/open-banking/v3.1/pisp/domestic-payments",
	"/open-banking/v3.1/pisp/domestic-scheduled-payment-consents",
	"/open-banking/v3.1/pisp/domestic-standing-order-consents",
	"/open-banking/v3.1/pisp/domestic-standing-orders",
	"/open-banking/v3.1/pisp/domestic-standing-orders",
	"/open-banking/v3.1/pisp/international-payment-consents",
	"/open-banking/v3.1/pisp/international-payments",
	"/open-banking/v3.1/pisp/international-scheduled-payment-consents",
	"/open-banking/v3.1/pisp/international-scheduled-payments",
}
