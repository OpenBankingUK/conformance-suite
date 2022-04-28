package model

// conditionalityStaticData - Get /accounts example json response from ozone
func conditionalityStaticData() []byte {

	// @NEW-SPEC-VERSION - Update this json structure using existing convention
	return []byte(
		`{
		  "account-transaction-v3.1": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.1": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.2": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.3": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.4": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.5": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.6": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.7": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.8": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.9": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		  "account-transaction-v3.1.10": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/account-access-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/account-access-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/balances"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/balances"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/transactions"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/beneficiaries"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/beneficiaries"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/direct-debits"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/direct-debits"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/standing-orders"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/standing-orders"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/product"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/products"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/offers"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/offers"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/party"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/party"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/scheduled-payments"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/scheduled-payments"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/file"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/accounts/{AccountId}/statements/{StatementId}/transactions"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/statements"
			}
		  ],
		
		  "payment-initiation-v3.1": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.1": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.2": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.3": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.4": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.5": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.6": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.7": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.8": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.9": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		  "payment-initiation-v3.1.10": [
			{
			  "endpoint": "/domestic-payment-consents",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments",
			  "method": "POST",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-payments/{DomesticPaymentId}",
			  "method": "GET",
			  "condition": "mandatory"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-scheduled-payments/{DomesticScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/domestic-standing-orders/{DomesticStandingOrderId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-payments/{InternationalPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payment-consents/{ConsentId}/funds-confirmation",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-scheduled-payments/{InternationalScheduledPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-order-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/international-standing-orders/{InternationalStandingOrderPaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payment-consents/{ConsentId}/file",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments",
			  "method": "POST",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}",
			  "method": "GET",
			  "condition": "conditional"
			},
			{
			  "endpoint": "/file-payments/{FilePaymentId}/report-file",
			  "method": "GET",
			  "condition": "conditional"
			}
		  ],
		
		  "confirmation-funds-v3.1": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.1": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.2": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.3": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.4": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.5": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.6": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.7": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.8": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.9": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		  "confirmation-funds-v3.1.10": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmation-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/funds-confirmation-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/funds-confirmations"
			}
		  ],
		
		  "event-notification-aspsp-v3.1": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/callback-urls"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/callback-urls"
			},
			{
			  "condition": "mandatory",
			  "method": "PUT",
			  "endpoint": "/callback-urls/{CallbackUrlId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/callback-urls/{CallbackUrlId}"
			}
		  ],
		  "event-notification-aspsp-v3.1.1": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/callback-urls"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/callback-urls"
			},
			{
			  "condition": "mandatory",
			  "method": "PUT",
			  "endpoint": "/callback-urls/{CallbackUrlId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/callback-urls/{CallbackUrlId}"
			}
		  ],
		  "event-notification-aspsp-v3.1.2": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/callback-urls"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/callback-urls"
			},
			{
			  "condition": "mandatory",
			  "method": "PUT",
			  "endpoint": "/callback-urls/{CallbackUrlId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/callback-urls/{CallbackUrlId}"
			}
		  ],
		
		  "variable-recurring-payments-v3.1.8": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/domestic-vrp-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/domestic-vrp-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/domestic-vrp-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/domestic-vrp-consents/{ConsentId}/funds-confirmation"
			},
			{
			  "condition": "conditional",
			  "method": "POST",
			  "endpoint": "/domestic-vrps"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/domestic-vrps/{DomesticVRPId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/domestic-vrps/{DomesticVRPId}/payment-details"
			}
		  ],
		  "variable-recurring-payments-v3.1.9": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/domestic-vrp-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/domestic-vrp-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/domestic-vrp-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/domestic-vrp-consents/{ConsentId}/funds-confirmation"
			},
			{
			  "condition": "conditional",
			  "method": "POST",
			  "endpoint": "/domestic-vrps"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/domestic-vrps/{DomesticVRPId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/domestic-vrps/{DomesticVRPId}/payment-details"
			}
		  ],
		  "variable-recurring-payments-v3.1.10": [
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/domestic-vrp-consents"
			},
			{
			  "condition": "mandatory",
			  "method": "GET",
			  "endpoint": "/domestic-vrp-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "DELETE",
			  "endpoint": "/domestic-vrp-consents/{ConsentId}"
			},
			{
			  "condition": "mandatory",
			  "method": "POST",
			  "endpoint": "/domestic-vrp-consents/{ConsentId}/funds-confirmation"
			},
			{
			  "condition": "conditional",
			  "method": "POST",
			  "endpoint": "/domestic-vrps"
			},
			{
			  "condition": "conditional",
			  "method": "GET",
			  "endpoint": "/domestic-vrps/{DomesticVRPId}"
			},
			{
			  "condition": "optional",
			  "method": "GET",
			  "endpoint": "/domestic-vrps/{DomesticVRPId}/payment-details"
			}
		  ]
		}
    `)
}
