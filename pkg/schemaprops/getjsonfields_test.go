package schemaprops

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCollectReturnedJSONFields(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	c := GetPropertyCollector()
	c.SetCollectorAPIDetails("myapi", "v3.1.0")
	c.CollectProperties("GET", "/accounts", string(tdata1), 200)
	c.OutputJSON()
}

func TestTransactionsJSONFields(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	c := GetPropertyCollector()
	c.SetCollectorAPIDetails("myapi", "v3.1.0")
	c.CollectProperties("GET", "https://myserver/open-banking/3.1/aisp/accounts/1234567853/transactions", string(atransaction), 200)
	c.SetCollectorAPIDetails("yourapi", "v3.1.1")
	c.CollectProperties("GET", "https://myserver/open-banking/3.1/aisp/accounts", string(accounts), 200)

	fmt.Println(c.OutputJSON())
}

func TestAccountsJSONFields(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	c := GetPropertyCollector()
	c.SetCollectorAPIDetails("myapi", "v3.1.0")
	c.CollectProperties("GET", "/open-banking/3.1/aisp/accounts", string(accounts), 200)
	result := c.OutputJSON()
	fmt.Println(result)
}

func TestAddEmptyAPI(t *testing.T) {
	c := GetPropertyCollector()
	c.CollectProperties("GET", "https://myserver/open-banking/3.1/aisp/accounts/1234567853/transactions", string(atransaction), 200)
	apitype, err := FindApi("https://myserver/open-banking/3.1/aisp/accounts/1234567853/transactions")
	assert.Equal(t, "accounts", apitype)
	assert.Nil(t, err)
	apitype, err = FindApi("https://myserver/open-banking/3.1/pisp/domestic-payment-consents")
	assert.Equal(t, "payments", apitype)
	assert.Nil(t, err)
	apitype, err = FindApi("https://myserver/open-banking/3.1/cbpii/funds-confirmation-consents")
	assert.Equal(t, "cbpii", apitype)
	assert.Nil(t, err)
	apitype, err = FindApi("https://myserver/open-banking/3.1/cbpii/funds-confirmations")
	assert.Equal(t, "cbpii", apitype)
	assert.Nil(t, err)
	apitype, err = FindApi("this does not exist")
	assert.Equal(t, "", apitype)
	assert.NotNil(t, err)
}

func TestAddUnnamedApiThenMerge(t *testing.T) {
	c := GetPropertyCollector()
	c.SetCollectorAPIDetails(ConsentGathering, "1")

	c.CollectProperties("GET", "https://myserver/open-banking/3.1/aisp/accounts", string(accounts), 200)

	c.SetCollectorAPIDetails("Accounts and Trasactions", "v3.1.0")
	c.CollectProperties("GET", "https://myserver/open-banking/3.1/aisp/accounts/1234567853/transactions", string(atransaction), 200)

	result := c.OutputJSON()
	assert.Contains(t, result, "Data.Transaction.Balance.CreditDebitIndicator")

	fmt.Println(result)

}

func TestMergeUnnamedApiThenMerge(t *testing.T) {

	c := GetPropertyCollector()
	c.SetCollectorAPIDetails("ConsentGathering", "")
	c.CollectProperties("GET", "https://myserver/open-banking/3.1/aisp/accounts", string(accountsmerge), 200)

	c.SetCollectorAPIDetails("Accounts and Trasactions", "v3.1.0")
	c.CollectProperties("GET", "https://myserver/open-banking/3.1/aisp/accounts", string(accounts), 200)
	c.CollectProperties("GET", "https://myserver/open-banking/3.1/aisp/accounts/1234567853/transactions", string(atransaction), 200)

	result := c.OutputJSON()
	assert.Contains(t, result, "MyFatMergeField")
	fmt.Println(result)

}

var (
	acctPay1  = append(accountsRegex, paymentsRegex...)
	allregex1 = append(acctPay, cbpiiRegex...)
)

func toSwagger(path string) (string, error) {
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

func pathsToSwagger(endpoints []string) []string {
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

func (c *Collector) dumpProperties() {
	logrus.SetLevel(logrus.TraceLevel)

	if logrus.GetLevel() == logrus.TraceLevel {
		logrus.Debug("Dump Properties===============")
		endpoints := sortEndpoints(c.Apis[c.currentApi].endpoints)
		for _, k := range endpoints {
			logrus.Debugf("%s", k)
			v := c.Apis[c.currentApi].endpoints[k]
			sortedv := sortPaths(v)
			for _, x := range sortedv {
				logrus.Debugf("%s", x)
			}
		}
		logrus.Debug("End Dump Properties===============")
	}
}

var (
	accounts = []byte(`{
		"Data": {
		   "Account": [
			  {
				 "AccountId": "700004000000000000000005",
				 "Currency": "GBP",
				 "Nickname": "xxxx0005",
				 "AccountType": "Business",
				 "AccountSubType": "CurrentAccount",
				 "Account": [
					{
					   "SchemeName": "UK.OBIE.SortCodeAccountNumber",
					   "Identification": "70000170000005",
					   "Name": "Octon Inc"
					}
				 ]
			  },
			  {
				 "AccountId": "700004000000000000000001",
				 "Currency": "GBP",
				 "Nickname": "xxxx0001",
				 "AccountType": "Personal",
				 "AccountSubType": "CurrentAccount",
				 "Account": [
					{
					   "SchemeName": "UK.OBIE.SortCodeAccountNumber",
					   "Identification": "70000170000001",
					   "SecondaryIdentification": "Roll No. 001"
					}
				 ]
			  },
			  {
				 "AccountId": "700004000000000000000006",
				 "Currency": "GBP",
				 "Nickname": "xxxx0002",
				 "AccountType": "Business",
				 "AccountSubType": "Other",
				 "Account": [
					{
					   "SchemeName": "UK.OBIE.SortCodeAccountNumber",
					   "Identification": "70000170000006",
					   "SecondaryIdentification": "Roll No. 002"
					}
				 ]
			  },
			  {
				 "AccountId": "700004000000000000000003",
				 "Currency": "GBP",
				 "Nickname": "xxxx0003",
				 "AccountType": "Business",
				 "AccountSubType": "CurrentAccount",
				 "Account": [
					{
					   "SchemeName": "UK.OBIE.IBAN",
					   "Identification": "GB29OBI170000170000001",
					   "Name": "Mario Carpentry"
					}
				 ],
				 "Servicer": {
					"SchemeName": "UK.OBIE.BICFI",
					"Identification": "GB29OBI1XXX"
				 }
			  },
			  {
				 "AccountId": "700004000000000000000002",
				 "Currency": "GBP",
				 "Nickname": "xxxx0006",
				 "AccountType": "Personal",
				 "AccountSubType": "CurrentAccount",
				 "Account": [
					{
					   "SchemeName": "UK.OBIE.SortCodeAccountNumber",
					   "Identification": "70000170000002",
					   "SecondaryIdentification": "Roll No. 002"
					}
				 ]
			  }
		   ]
		},
		"Links": {
		   "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/accounts"
		},
		"Meta": {
		   "TotalPages": 1
		}
	 }`)

	accountsmerge = []byte(`{
			"Data": {
			   "Account": [
				  {
					 "AccountId": "700004000000000000000005",
					 "MyFatMergeField":"hellothere",
					 "Currency": "GBP",
					 "Nickname": "xxxx0005",
					 "AccountType": "Business",
					 "AccountSubType": "CurrentAccount",
					 "Account": [
						{
						   "SchemeName": "UK.OBIE.SortCodeAccountNumber",
						   "Identification": "70000170000005",
						   "Name": "Octon Inc"
						}
					 ]
				  },
				  {
					 "AccountId": "700004000000000000000001",
					 "Currency": "GBP",
					 "Nickname": "xxxx0001",
					 "AccountType": "Personal",
					 "AccountSubType": "CurrentAccount",
					 "Account": [
						{
						   "SchemeName": "UK.OBIE.SortCodeAccountNumber",
						   "Identification": "70000170000001",
						   "SecondaryIdentification": "Roll No. 001"
						}
					 ]
				  },
				  {
					 "AccountId": "700004000000000000000006",
					 "Currency": "GBP",
					 "Nickname": "xxxx0002",
					 "AccountType": "Business",
					 "AccountSubType": "Other",
					 "Account": [
						{
						   "SchemeName": "UK.OBIE.SortCodeAccountNumber",
						   "Identification": "70000170000006",
						   "SecondaryIdentification": "Roll No. 002"
						}
					 ]
				  },
				  {
					 "AccountId": "700004000000000000000003",
					 "Currency": "GBP",
					 "Nickname": "xxxx0003",
					 "AccountType": "Business",
					 "AccountSubType": "CurrentAccount",
					 "Account": [
						{
						   "SchemeName": "UK.OBIE.IBAN",
						   "Identification": "GB29OBI170000170000001",
						   "Name": "Mario Carpentry"
						}
					 ],
					 "Servicer": {
						"SchemeName": "UK.OBIE.BICFI",
						"Identification": "GB29OBI1XXX"
					 }
				  },
				  {
					 "AccountId": "700004000000000000000002",
					 "Currency": "GBP",
					 "Nickname": "xxxx0006",
					 "AccountType": "Personal",
					 "AccountSubType": "CurrentAccount",
					 "Account": [
						{
						   "SchemeName": "UK.OBIE.SortCodeAccountNumber",
						   "Identification": "70000170000002",
						   "SecondaryIdentification": "Roll No. 002"
						}
					 ]
				  }
			   ]
			},
			"Links": {
			   "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/accounts"
			},
			"Meta": {
			   "TotalPages": 1
			}
		 }`)

	atransaction = []byte(`{
		"Data": {
		   "Transaction": [
			  {
				 "AccountId": "700004000000000000000005",
				 "TransactionId": "ffcb7ac9-b178-4093-8a10-b7281c39cf15",
				 "TransactionReference": "60c09267-21a0-4fc2-8a08-551a839bb77",
				 "BookingDateTime": "2019-09-16T01:59:39.629Z",
				 "ValueDateTime": "2019-09-16T01:59:39.629Z",
				 "ProprietaryBankTransactionCode": {
					"Code": "PMT"
				 },
				 "TransactionInformation": "Payment Id: pmt-cfd4aa73-8aa8-41c6-9342-66d3bf37bd9b",
				 "Amount": {
					"Amount": "0.10",
					"Currency": "GBP"
				 },
				 "CreditDebitIndicator": "Debit",
				 "Status": "Booked",
				 "Balance": {
					"Amount": {
					   "Amount": "0.10",
					   "Currency": "GBP"
					},
					"CreditDebitIndicator": "Debit",
					"Type": "ClosingAvailable"
				 }
			  }
		   ]
		},
		"Links": {
		   "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/accounts/700004000000000000000005/transactions"
		},
		"Meta": {
		   "TotalPages": 1,
		   "FirstAvailableDateTime": "2016-01-01T08:40:00.000Z",
		   "LastAvailableDateTime": "2025-12-31T08:40:00.000Z"
		}
	 }`)

	tdata1 = []byte(`
	{
	   "Data": {
		  "Product": [
			 {
				"ProductName": "Modelator por Business",
				"ProductId": "700005000000000000000003",
				"AccountId": "700004000000000000000005",
				"SecondaryProductId": "700005000000000000000003",
				"ProductType": "BusinessCurrentAccount",
				"MarketingStateId": "Active"
			 },
			 {
				"ProductName": "Personal Current Account Advantage",
				"ProductId": "700005000000000000000001",
				"AccountId": "700004000000000000000001",
				"SecondaryProductId": "700005000000000000000001",
				"ProductType": "PersonalCurrentAccount",
				"MarketingStateId": "Active",
				"PCA": {
				   "CreditInterest": {
					  "TierBandSet": {
						 "TierBandMethod": "Tiered",
						 "CalculationMethod": "Compound",
						 "Notes": [
							"From the 11th June 2017, we will be changing how we pat credit interest on our Bank of Scotland current accounts with Vantage.",
							"We are replacing the current tiered rates with a single interest rate of 2% AER(1.98% gross) (variable) on credit balances between £1 and £5,000. Depending on the balance of your account this may be an increase or decrease to the rate your currently receive."
						 ],
						 "TierBand": [
							{
							   "Identification": "1",
							   "TierValueMinimum": "1.00",
							   "TierValueMaximum": "999.99",
							   "CalculationFrequency": "Monthly",
							   "ApplicationFrequency": "Monthly",
							   "DepositInterestAppliedCoverage": "Tiered",
							   "FixedVariableInterestRateType": "Variable",
							   "AER": "1.50",
							   "BankInterestRateType": "Gross"
							},
							{
							   "Identification": "2",
							   "TierValueMinimum": "1000.00",
							   "TierValueMaximum": "2999.99",
							   "CalculationFrequency": "Monthly",
							   "ApplicationFrequency": "Monthly",
							   "DepositInterestAppliedCoverage": "Tiered",
							   "FixedVariableInterestRateType": "Variable",
							   "AER": "2.00",
							   "BankInterestRateType": "Gross"
							},
							{
							   "TierValueMinimum": "3000.00",
							   "TierValueMaximum": "5000.00",
							   "CalculationFrequency": "Monthly",
							   "ApplicationFrequency": "Monthly",
							   "DepositInterestAppliedCoverage": "Tiered",
							   "FixedVariableInterestRateType": "Variable",
							   "AER": "3.00",
							   "BankInterestRateType": "Gross"
							}
						 ]
					  }
				   }
				}
			 },
			 {
				"ProductName": "Business credit card",
				"ProductId": "700005000000000000000004",
				"AccountId": "700004000000000000000006",
				"SecondaryProductId": "700005000000000000000004",
				"ProductType": "Other",
				"MarketingStateId": "Active",
				"OtherProductDetails": {
				   "Name": "BasicSavingAccount",
				   "Description": "HSBC fee free saving basic saving account",
				   "OtherProductDetails": {
					  "OtherFeesCharges": {
						 "FeeChargeDetail": [
							{
							   "FeeCategory": "Servicing",
							   "FeeType": "ServiceCAccountFeeMonthly",
							   "FeeAmount": "12.500",
							   "ApplicationFrequency": "Monthly",
							   "CalculationFrequency": "Daily",
							   "Notes": [
								  "Our tariff includes:\n* depositing and sending cheques\n* cash deposits up to the limit your tariff allows\n* withdrawals\n* Direct Debits, standing orders, bill payments\n* Bas credits\n* debit card payments"
							   ]
							}
						 ]
					  }
				   }
				}
			 },
			 {
				"ProductName": "Modelator por Business",
				"ProductId": "700005000000000000000003",
				"AccountId": "700004000000000000000003",
				"SecondaryProductId": "700005000000000000000003",
				"ProductType": "BusinessCurrentAccount",
				"MarketingStateId": "Active"
			 },
			 {
				"ProductName": "Personal Current Account Advantage",
				"ProductId": "700005000000000000000001",
				"AccountId": "700004000000000000000002",
				"SecondaryProductId": "700005000000000000000001",
				"ProductType": "PersonalCurrentAccount",
				"MarketingStateId": "Active",
				"PCA": {
				   "CreditInterest": {
					  "TierBandSet": {
						 "TierBandMethod": "Tiered",
						 "CalculationMethod": "Compound",
						 "Notes": [
							"From the 11th June 2017, we will be changing how we pat credit interest on our Bank of Scotland current accounts with Vantage.",
							"We are replacing the current tiered rates with a single interest rate of 2% AER(1.98% gross) (variable) on credit balances between £1 and £5,000. Depending on the balance of your account this may be an increase or decrease to the rate your currently receive."
						 ],
						 "TierBand": [
							{
							   "Identification": "1",
							   "TierValueMinimum": "1.00",
							   "TierValueMaximum": "999.99",
							   "CalculationFrequency": "Monthly",
							   "ApplicationFrequency": "Monthly",
							   "DepositInterestAppliedCoverage": "Tiered",
							   "FixedVariableInterestRateType": "Variable",
							   "AER": "1.50",
							   "BankInterestRateType": "Gross"
							},
							{
							   "Identification": "2",
							   "TierValueMinimum": "1000.00",
							   "TierValueMaximum": "2999.99",
							   "CalculationFrequency": "Monthly",
							   "ApplicationFrequency": "Monthly",
							   "DepositInterestAppliedCoverage": "Tiered",
							   "FixedVariableInterestRateType": "Variable",
							   "AER": "2.00",
							   "BankInterestRateType": "Gross"
							},
							{
							   "TierValueMinimum": "3000.00",
							   "TierValueMaximum": "5000.00",
							   "CalculationFrequency": "Monthly",
							   "ApplicationFrequency": "Monthly",
							   "DepositInterestAppliedCoverage": "Tiered",
							   "FixedVariableInterestRateType": "Variable",
							   "AER": "3.00",
							   "BankInterestRateType": "Gross"
							}
						 ]
					  }
				   }
				}
			 }
		  ]
	   },
	   "Links": {
		  "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/products"
	   },
	   "Meta": {
		  "TotalPages": 1
	   }
	}`)

	tdata = []byte(`{
		"Data": {
		   "Transaction": [
			  {
				 "AccountId": "700004000000000000000003",
				 "BookingDateTime": "2017-01-02T00:00:00.000Z",
				 "ValueDateTime": "2017-01-02T00:00:00.000Z",
				 "Amount": {
					"Amount": "6.20",
					"Currency": "GBP"
				 },
				 "CreditDebitIndicator": "Debit",
				 "Status": "Booked",
				 "BankTransactionCode": {
					"Code": "CounterTransactions",
					"SubCode": "BranchDeposit"
				 }
			  },
			  {
				 "AccountId": "700004000000000000000002",
				 "BookingDateTime": "2017-01-02T00:00:00.000Z",
				 "ValueDateTime": "2017-01-02T00:00:00.000Z",
				 "Amount": {
					"Amount": "6.20",
					"Currency": "GBP"
				 },
				 "CreditDebitIndicator": "Debit",
				 "Status": "Booked",
				 "BankTransactionCode": {
					"Code": "CounterTransactions",
					"SubCode": "BranchDeposit"
				 }
			  },
			  {
				 "AccountId": "700004000000000000000001",
				 "BookingDateTime": "2017-01-02T00:00:00.000Z",
				 "ValueDateTime": "2017-01-02T00:00:00.000Z",
				 "Amount": {
					"Amount": "19.88",
					"Currency": "GBP"
				 },
				 "CreditDebitIndicator": "Debit",
				 "Status": "Booked",
				 "BankTransactionCode": {
					"Code": "IssuedCheques",
					"SubCode": "BankCheque"
				 }
			  },
			  {
				 "AccountId": "700004000000000000000002",
				 "TransactionId": "e11bc082-f680-4101-8509-70a0254e811a",
				 "TransactionReference": "5",
				 "BookingDateTime": "2019-09-24T22:43:48.901Z",
				 "ValueDateTime": "2019-09-24T22:43:48.901Z",
				 "ProprietaryBankTransactionCode": {
					"Code": "PMT"
				 },
				 "Amount": {
					"Amount": "3.5",
					"Currency": "GBP"
				 },
				 "CreditDebitIndicator": "Debit",
				 "Status": "Booked"
			  }
		   ]
		},
		"Links": {
		   "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/transactions"
		},
		"Meta": {
		   "TotalPages": 1,
		   "FirstAvailableDateTime": "2016-01-01T08:40:00.000Z",
		   "LastAvailableDateTime": "2025-12-31T08:40:00.000Z"
		}
	 }`)

	transdata = []byte(`{
	"Data": {
	   "Transaction": [
		  {
			 "AccountId": "700004000000000000000003",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.20",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.20",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "19.88",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.20",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-03T00:00:00.000Z",
			 "ValueDateTime": "2017-01-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "56.66",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-03T00:00:00.000Z",
			 "ValueDateTime": "2017-01-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "9.85",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-27T00:00:00.000Z",
			 "ValueDateTime": "2017-01-27T00:00:00.000Z",
			 "Amount": {
				"Amount": "5.50",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-06T00:00:00.000Z",
			 "ValueDateTime": "2017-01-06T00:00:00.000Z",
			 "Amount": {
				"Amount": "64.50",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-07T00:00:00.000Z",
			 "ValueDateTime": "2017-01-07T00:00:00.000Z",
			 "Amount": {
				"Amount": "9.75",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-07T00:00:00.000Z",
			 "ValueDateTime": "2017-01-07T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.60",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-08T00:00:00.000Z",
			 "ValueDateTime": "2017-01-08T00:00:00.000Z",
			 "Amount": {
				"Amount": "13.20",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-08T00:00:00.000Z",
			 "ValueDateTime": "2017-01-08T00:00:00.000Z",
			 "Amount": {
				"Amount": "42.95",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-10T00:00:00.000Z",
			 "ValueDateTime": "2017-01-10T00:00:00.000Z",
			 "Amount": {
				"Amount": "48.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-13T00:00:00.000Z",
			 "ValueDateTime": "2017-01-13T00:00:00.000Z",
			 "Amount": {
				"Amount": "13.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-13T00:00:00.000Z",
			 "ValueDateTime": "2017-01-13T00:00:00.000Z",
			 "Amount": {
				"Amount": "93.04",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-14T00:00:00.000Z",
			 "ValueDateTime": "2017-01-14T00:00:00.000Z",
			 "Amount": {
				"Amount": "8.95",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-16T00:00:00.000Z",
			 "ValueDateTime": "2017-01-16T00:00:00.000Z",
			 "Amount": {
				"Amount": "24.91",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-20T00:00:00.000Z",
			 "ValueDateTime": "2017-01-20T00:00:00.000Z",
			 "Amount": {
				"Amount": "2512.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-21T00:00:00.000Z",
			 "ValueDateTime": "2017-01-21T00:00:00.000Z",
			 "Amount": {
				"Amount": "8.10",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-22T00:00:00.000Z",
			 "ValueDateTime": "2017-01-22T00:00:00.000Z",
			 "Amount": {
				"Amount": "45.49",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-23T00:00:00.000Z",
			 "ValueDateTime": "2017-01-23T00:00:00.000Z",
			 "Amount": {
				"Amount": "25.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-23T00:00:00.000Z",
			 "ValueDateTime": "2017-01-23T00:00:00.000Z",
			 "Amount": {
				"Amount": "7.04",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-28T00:00:00.000Z",
			 "ValueDateTime": "2017-01-28T00:00:00.000Z",
			 "Amount": {
				"Amount": "11200.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-29T00:00:00.000Z",
			 "ValueDateTime": "2017-01-29T00:00:00.000Z",
			 "Amount": {
				"Amount": "170.99",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-30T00:00:00.000Z",
			 "ValueDateTime": "2017-01-30T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-03T00:00:00.000Z",
			 "ValueDateTime": "2017-02-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "34.60",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-03T00:00:00.000Z",
			 "ValueDateTime": "2017-02-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "6792.48",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-03T00:00:00.000Z",
			 "ValueDateTime": "2017-02-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "67.60",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-05T00:00:00.000Z",
			 "ValueDateTime": "2017-02-05T00:00:00.000Z",
			 "Amount": {
				"Amount": "59.98",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-07T00:00:00.000Z",
			 "ValueDateTime": "2017-02-07T00:00:00.000Z",
			 "Amount": {
				"Amount": "46.99",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-07T00:00:00.000Z",
			 "ValueDateTime": "2017-02-07T00:00:00.000Z",
			 "Amount": {
				"Amount": "29.07",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-10T00:00:00.000Z",
			 "ValueDateTime": "2017-02-10T00:00:00.000Z",
			 "Amount": {
				"Amount": "35.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-10T00:00:00.000Z",
			 "ValueDateTime": "2017-02-10T00:00:00.000Z",
			 "Amount": {
				"Amount": "20.72",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-10T00:00:00.000Z",
			 "ValueDateTime": "2017-02-10T00:00:00.000Z",
			 "Amount": {
				"Amount": "69.06",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-10T00:00:00.000Z",
			 "ValueDateTime": "2017-02-10T00:00:00.000Z",
			 "Amount": {
				"Amount": "9.19",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-10T00:00:00.000Z",
			 "ValueDateTime": "2017-02-10T00:00:00.000Z",
			 "Amount": {
				"Amount": "23011.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-12T00:00:00.000Z",
			 "ValueDateTime": "2017-02-12T00:00:00.000Z",
			 "Amount": {
				"Amount": "8.36",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-14T00:00:00.000Z",
			 "ValueDateTime": "2017-02-14T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.80",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-17T00:00:00.000Z",
			 "ValueDateTime": "2017-02-17T00:00:00.000Z",
			 "Amount": {
				"Amount": "23.59",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-17T00:00:00.000Z",
			 "ValueDateTime": "2017-02-17T00:00:00.000Z",
			 "Amount": {
				"Amount": "7.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-17T00:00:00.000Z",
			 "ValueDateTime": "2017-02-17T00:00:00.000Z",
			 "Amount": {
				"Amount": "2512.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-17T00:00:00.000Z",
			 "ValueDateTime": "2017-02-17T00:00:00.000Z",
			 "Amount": {
				"Amount": "2512.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-24T00:00:00.000Z",
			 "ValueDateTime": "2017-02-24T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.60",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-24T00:00:00.000Z",
			 "ValueDateTime": "2017-02-24T00:00:00.000Z",
			 "Amount": {
				"Amount": "5.50",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-28T00:00:00.000Z",
			 "ValueDateTime": "2017-02-28T00:00:00.000Z",
			 "Amount": {
				"Amount": "53.54",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-28T00:00:00.000Z",
			 "ValueDateTime": "2017-02-28T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-28T00:00:00.000Z",
			 "ValueDateTime": "2017-02-28T00:00:00.000Z",
			 "Amount": {
				"Amount": "1.50",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-03T00:00:00.000Z",
			 "ValueDateTime": "2017-03-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "4.77",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-03T00:00:00.000Z",
			 "ValueDateTime": "2017-03-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "173.74",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-03T00:00:00.000Z",
			 "ValueDateTime": "2017-03-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "7.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-06T00:00:00.000Z",
			 "ValueDateTime": "2017-03-06T00:00:00.000Z",
			 "Amount": {
				"Amount": "143.22",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-06T00:00:00.000Z",
			 "ValueDateTime": "2017-03-06T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-11T00:00:00.000Z",
			 "ValueDateTime": "2017-03-11T00:00:00.000Z",
			 "Amount": {
				"Amount": "71.99",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-13T00:00:00.000Z",
			 "ValueDateTime": "2017-03-13T00:00:00.000Z",
			 "Amount": {
				"Amount": "9.95",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-13T00:00:00.000Z",
			 "ValueDateTime": "2017-03-13T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-18T00:00:00.000Z",
			 "ValueDateTime": "2017-03-18T00:00:00.000Z",
			 "Amount": {
				"Amount": "25.09",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-18T00:00:00.000Z",
			 "ValueDateTime": "2017-03-18T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.51",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-20T00:00:00.000Z",
			 "ValueDateTime": "2017-03-20T00:00:00.000Z",
			 "Amount": {
				"Amount": "2512.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-20T00:00:00.000Z",
			 "ValueDateTime": "2017-03-20T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-24T00:00:00.000Z",
			 "ValueDateTime": "2017-03-24T00:00:00.000Z",
			 "Amount": {
				"Amount": "247.67",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-25T00:00:00.000Z",
			 "ValueDateTime": "2017-03-25T00:00:00.000Z",
			 "Amount": {
				"Amount": "25.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-26T00:00:00.000Z",
			 "ValueDateTime": "2017-03-26T00:00:00.000Z",
			 "Amount": {
				"Amount": "133.74",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-26T23:00:00.000Z",
			 "ValueDateTime": "2017-03-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "5.50",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-27T23:00:00.000Z",
			 "ValueDateTime": "2017-03-27T23:00:00.000Z",
			 "Amount": {
				"Amount": "750.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-01T23:00:00.000Z",
			 "ValueDateTime": "2017-04-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "5.32",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-01T23:00:00.000Z",
			 "ValueDateTime": "2017-04-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "29.33",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-01T23:00:00.000Z",
			 "ValueDateTime": "2017-04-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "51.54",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-02T23:00:00.000Z",
			 "ValueDateTime": "2017-04-02T23:00:00.000Z",
			 "Amount": {
				"Amount": "19.36",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-02T23:00:00.000Z",
			 "ValueDateTime": "2017-04-02T23:00:00.000Z",
			 "Amount": {
				"Amount": "13.06",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-02T23:00:00.000Z",
			 "ValueDateTime": "2017-04-02T23:00:00.000Z",
			 "Amount": {
				"Amount": "29.31",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-05T23:00:00.000Z",
			 "ValueDateTime": "2017-04-05T23:00:00.000Z",
			 "Amount": {
				"Amount": "1.38",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-06T23:00:00.000Z",
			 "ValueDateTime": "2017-04-06T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.19",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-07T23:00:00.000Z",
			 "ValueDateTime": "2017-04-07T23:00:00.000Z",
			 "Amount": {
				"Amount": "47.26",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-07T23:00:00.000Z",
			 "ValueDateTime": "2017-04-07T23:00:00.000Z",
			 "Amount": {
				"Amount": "58.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-08T23:00:00.000Z",
			 "ValueDateTime": "2017-04-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "13.19",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-09T23:00:00.000Z",
			 "ValueDateTime": "2017-04-09T23:00:00.000Z",
			 "Amount": {
				"Amount": "12.18",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-09T23:00:00.000Z",
			 "ValueDateTime": "2017-04-09T23:00:00.000Z",
			 "Amount": {
				"Amount": "1.78",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-12T23:00:00.000Z",
			 "ValueDateTime": "2017-04-12T23:00:00.000Z",
			 "Amount": {
				"Amount": "6.99",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-12T23:00:00.000Z",
			 "ValueDateTime": "2017-04-12T23:00:00.000Z",
			 "Amount": {
				"Amount": "76.87",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-13T23:00:00.000Z",
			 "ValueDateTime": "2017-04-13T23:00:00.000Z",
			 "Amount": {
				"Amount": "9.40",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-15T23:00:00.000Z",
			 "ValueDateTime": "2017-04-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "23.93",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-16T23:00:00.000Z",
			 "ValueDateTime": "2017-04-16T23:00:00.000Z",
			 "Amount": {
				"Amount": "9.73",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-19T23:00:00.000Z",
			 "ValueDateTime": "2017-04-19T23:00:00.000Z",
			 "Amount": {
				"Amount": "3683.38",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-20T23:00:00.000Z",
			 "ValueDateTime": "2017-04-20T23:00:00.000Z",
			 "Amount": {
				"Amount": "5.14",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-22T23:00:00.000Z",
			 "ValueDateTime": "2017-04-22T23:00:00.000Z",
			 "Amount": {
				"Amount": "8.43",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-22T23:00:00.000Z",
			 "ValueDateTime": "2017-04-22T23:00:00.000Z",
			 "Amount": {
				"Amount": "54.18",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-26T23:00:00.000Z",
			 "ValueDateTime": "2017-04-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.64",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-27T23:00:00.000Z",
			 "ValueDateTime": "2017-04-27T23:00:00.000Z",
			 "Amount": {
				"Amount": "23.12",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-27T23:00:00.000Z",
			 "ValueDateTime": "2017-04-27T23:00:00.000Z",
			 "Amount": {
				"Amount": "69.03",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-28T23:00:00.000Z",
			 "ValueDateTime": "2017-04-28T23:00:00.000Z",
			 "Amount": {
				"Amount": "230.79",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-29T23:00:00.000Z",
			 "ValueDateTime": "2017-04-29T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.61",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-03T23:00:00.000Z",
			 "ValueDateTime": "2017-05-03T23:00:00.000Z",
			 "Amount": {
				"Amount": "5186.68",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-03T23:00:00.000Z",
			 "ValueDateTime": "2017-05-03T23:00:00.000Z",
			 "Amount": {
				"Amount": "2.14",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-05T23:00:00.000Z",
			 "ValueDateTime": "2017-05-05T23:00:00.000Z",
			 "Amount": {
				"Amount": "4.71",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-06T23:00:00.000Z",
			 "ValueDateTime": "2017-05-06T23:00:00.000Z",
			 "Amount": {
				"Amount": "48.04",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-07T23:00:00.000Z",
			 "ValueDateTime": "2017-05-07T23:00:00.000Z",
			 "Amount": {
				"Amount": "39.49",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-07T23:00:00.000Z",
			 "ValueDateTime": "2017-05-07T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.17",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-10T23:00:00.000Z",
			 "ValueDateTime": "2017-05-10T23:00:00.000Z",
			 "Amount": {
				"Amount": "38.15",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-10T23:00:00.000Z",
			 "ValueDateTime": "2017-05-10T23:00:00.000Z",
			 "Amount": {
				"Amount": "25.19",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-10T23:00:00.000Z",
			 "ValueDateTime": "2017-05-10T23:00:00.000Z",
			 "Amount": {
				"Amount": "4.64",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-12T23:00:00.000Z",
			 "ValueDateTime": "2017-05-12T23:00:00.000Z",
			 "Amount": {
				"Amount": "4.85",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-12T23:00:00.000Z",
			 "ValueDateTime": "2017-05-12T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.80",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-13T23:00:00.000Z",
			 "ValueDateTime": "2017-05-13T23:00:00.000Z",
			 "Amount": {
				"Amount": "50.94",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-14T23:00:00.000Z",
			 "ValueDateTime": "2017-05-14T23:00:00.000Z",
			 "Amount": {
				"Amount": "9.11",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-17T23:00:00.000Z",
			 "ValueDateTime": "2017-05-17T23:00:00.000Z",
			 "Amount": {
				"Amount": "15.28",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-17T23:00:00.000Z",
			 "ValueDateTime": "2017-05-17T23:00:00.000Z",
			 "Amount": {
				"Amount": "1067.42",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-20T23:00:00.000Z",
			 "ValueDateTime": "2017-05-20T23:00:00.000Z",
			 "Amount": {
				"Amount": "41.83",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-24T23:00:00.000Z",
			 "ValueDateTime": "2017-05-24T23:00:00.000Z",
			 "Amount": {
				"Amount": "9.25",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-24T23:00:00.000Z",
			 "ValueDateTime": "2017-05-24T23:00:00.000Z",
			 "Amount": {
				"Amount": "4.83",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-28T23:00:00.000Z",
			 "ValueDateTime": "2017-05-28T23:00:00.000Z",
			 "Amount": {
				"Amount": "67.07",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-28T23:00:00.000Z",
			 "ValueDateTime": "2017-05-28T23:00:00.000Z",
			 "Amount": {
				"Amount": "28.09",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-28T23:00:00.000Z",
			 "ValueDateTime": "2017-05-28T23:00:00.000Z",
			 "Amount": {
				"Amount": "5035.98",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-31T23:00:00.000Z",
			 "ValueDateTime": "2017-05-31T23:00:00.000Z",
			 "Amount": {
				"Amount": "36.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-31T23:00:00.000Z",
			 "ValueDateTime": "2017-05-31T23:00:00.000Z",
			 "Amount": {
				"Amount": "6.35",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-31T23:00:00.000Z",
			 "ValueDateTime": "2017-05-31T23:00:00.000Z",
			 "Amount": {
				"Amount": "9.29",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-03T23:00:00.000Z",
			 "ValueDateTime": "2017-06-03T23:00:00.000Z",
			 "Amount": {
				"Amount": "133.65",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-04T23:00:00.000Z",
			 "ValueDateTime": "2017-06-04T23:00:00.000Z",
			 "Amount": {
				"Amount": "29.01",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-08T23:00:00.000Z",
			 "ValueDateTime": "2017-06-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "57.60",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-10T23:00:00.000Z",
			 "ValueDateTime": "2017-06-10T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.11",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-10T23:00:00.000Z",
			 "ValueDateTime": "2017-06-10T23:00:00.000Z",
			 "Amount": {
				"Amount": "57.09",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-15T23:00:00.000Z",
			 "ValueDateTime": "2017-06-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "22.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-16T23:00:00.000Z",
			 "ValueDateTime": "2017-06-16T23:00:00.000Z",
			 "Amount": {
				"Amount": "43.01",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-17T23:00:00.000Z",
			 "ValueDateTime": "2017-06-17T23:00:00.000Z",
			 "Amount": {
				"Amount": "1783.50",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-17T23:00:00.000Z",
			 "ValueDateTime": "2017-06-17T23:00:00.000Z",
			 "Amount": {
				"Amount": "52.19",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-21T23:00:00.000Z",
			 "ValueDateTime": "2017-06-21T23:00:00.000Z",
			 "Amount": {
				"Amount": "148.25",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-22T23:00:00.000Z",
			 "ValueDateTime": "2017-06-22T23:00:00.000Z",
			 "Amount": {
				"Amount": "20.42",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-24T23:00:00.000Z",
			 "ValueDateTime": "2017-06-24T23:00:00.000Z",
			 "Amount": {
				"Amount": "20.93",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-24T23:00:00.000Z",
			 "ValueDateTime": "2017-06-24T23:00:00.000Z",
			 "Amount": {
				"Amount": "7.78",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-25T23:00:00.000Z",
			 "ValueDateTime": "2017-06-25T23:00:00.000Z",
			 "Amount": {
				"Amount": "410.55",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-30T23:00:00.000Z",
			 "ValueDateTime": "2017-06-30T23:00:00.000Z",
			 "Amount": {
				"Amount": "7.85",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-30T23:00:00.000Z",
			 "ValueDateTime": "2017-06-30T23:00:00.000Z",
			 "Amount": {
				"Amount": "43.11",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-30T23:00:00.000Z",
			 "ValueDateTime": "2017-06-30T23:00:00.000Z",
			 "Amount": {
				"Amount": "52.06",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-01T23:00:00.000Z",
			 "ValueDateTime": "2017-07-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "20.30",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-01T23:00:00.000Z",
			 "ValueDateTime": "2017-07-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.54",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-01T23:00:00.000Z",
			 "ValueDateTime": "2017-07-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "38.49",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-04T23:00:00.000Z",
			 "ValueDateTime": "2017-07-04T23:00:00.000Z",
			 "Amount": {
				"Amount": "57.15",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-05T23:00:00.000Z",
			 "ValueDateTime": "2017-07-05T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.84",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-06T23:00:00.000Z",
			 "ValueDateTime": "2017-07-06T23:00:00.000Z",
			 "Amount": {
				"Amount": "20.27",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-06T23:00:00.000Z",
			 "ValueDateTime": "2017-07-06T23:00:00.000Z",
			 "Amount": {
				"Amount": "23.26",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-06T23:00:00.000Z",
			 "ValueDateTime": "2017-07-06T23:00:00.000Z",
			 "Amount": {
				"Amount": "31.44",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-07T23:00:00.000Z",
			 "ValueDateTime": "2017-07-07T23:00:00.000Z",
			 "Amount": {
				"Amount": "2.65",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-08T23:00:00.000Z",
			 "ValueDateTime": "2017-07-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "7.15",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-11T23:00:00.000Z",
			 "ValueDateTime": "2017-07-11T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.79",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-11T23:00:00.000Z",
			 "ValueDateTime": "2017-07-11T23:00:00.000Z",
			 "Amount": {
				"Amount": "2.40",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-12T23:00:00.000Z",
			 "ValueDateTime": "2017-07-12T23:00:00.000Z",
			 "Amount": {
				"Amount": "9.38",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-14T23:00:00.000Z",
			 "ValueDateTime": "2017-07-14T23:00:00.000Z",
			 "Amount": {
				"Amount": "26.62",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-15T23:00:00.000Z",
			 "ValueDateTime": "2017-07-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "13.50",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-18T23:00:00.000Z",
			 "ValueDateTime": "2017-07-18T23:00:00.000Z",
			 "Amount": {
				"Amount": "2916.04",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-19T23:00:00.000Z",
			 "ValueDateTime": "2017-07-19T23:00:00.000Z",
			 "Amount": {
				"Amount": "2.56",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-21T23:00:00.000Z",
			 "ValueDateTime": "2017-07-21T23:00:00.000Z",
			 "Amount": {
				"Amount": "2.46",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-21T23:00:00.000Z",
			 "ValueDateTime": "2017-07-21T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.06",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-21T23:00:00.000Z",
			 "ValueDateTime": "2017-07-21T23:00:00.000Z",
			 "Amount": {
				"Amount": "20.66",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-25T23:00:00.000Z",
			 "ValueDateTime": "2017-07-25T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.74",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-26T23:00:00.000Z",
			 "ValueDateTime": "2017-07-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "27.24",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-27T23:00:00.000Z",
			 "ValueDateTime": "2017-07-27T23:00:00.000Z",
			 "Amount": {
				"Amount": "118.05",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-28T23:00:00.000Z",
			 "ValueDateTime": "2017-07-28T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.29",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-01T23:00:00.000Z",
			 "ValueDateTime": "2017-08-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "5.23",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-01T23:00:00.000Z",
			 "ValueDateTime": "2017-08-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "7713.05",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-01T23:00:00.000Z",
			 "ValueDateTime": "2017-08-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "2.91",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-03T23:00:00.000Z",
			 "ValueDateTime": "2017-08-03T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.19",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-04T23:00:00.000Z",
			 "ValueDateTime": "2017-08-04T23:00:00.000Z",
			 "Amount": {
				"Amount": "47.65",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-05T23:00:00.000Z",
			 "ValueDateTime": "2017-08-05T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.12",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-08T23:00:00.000Z",
			 "ValueDateTime": "2017-08-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "50.53",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-08T23:00:00.000Z",
			 "ValueDateTime": "2017-08-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "19.41",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-08T23:00:00.000Z",
			 "ValueDateTime": "2017-08-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "31.49",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-08T23:00:00.000Z",
			 "ValueDateTime": "2017-08-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.42",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-08T23:00:00.000Z",
			 "ValueDateTime": "2017-08-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "1670.84",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-10T23:00:00.000Z",
			 "ValueDateTime": "2017-08-10T23:00:00.000Z",
			 "Amount": {
				"Amount": "6.21",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-11T23:00:00.000Z",
			 "ValueDateTime": "2017-08-11T23:00:00.000Z",
			 "Amount": {
				"Amount": "37.21",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-15T23:00:00.000Z",
			 "ValueDateTime": "2017-08-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "13.61",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-15T23:00:00.000Z",
			 "ValueDateTime": "2017-08-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "8.44",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-15T23:00:00.000Z",
			 "ValueDateTime": "2017-08-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "3560.64",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-18T23:00:00.000Z",
			 "ValueDateTime": "2017-08-18T23:00:00.000Z",
			 "Amount": {
				"Amount": "28.42",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-22T23:00:00.000Z",
			 "ValueDateTime": "2017-08-22T23:00:00.000Z",
			 "Amount": {
				"Amount": "16.66",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-26T23:00:00.000Z",
			 "ValueDateTime": "2017-08-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "80.98",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-26T23:00:00.000Z",
			 "ValueDateTime": "2017-08-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "11.56",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-26T23:00:00.000Z",
			 "ValueDateTime": "2017-08-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.48",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-26T23:00:00.000Z",
			 "ValueDateTime": "2017-08-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "621.05",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-29T23:00:00.000Z",
			 "ValueDateTime": "2017-08-29T23:00:00.000Z",
			 "Amount": {
				"Amount": "75.88",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-29T23:00:00.000Z",
			 "ValueDateTime": "2017-08-29T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.39",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-30T23:00:00.000Z",
			 "ValueDateTime": "2017-08-30T23:00:00.000Z",
			 "Amount": {
				"Amount": "2.59",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-01T23:00:00.000Z",
			 "ValueDateTime": "2017-09-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "4.69",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-02T23:00:00.000Z",
			 "ValueDateTime": "2017-09-02T23:00:00.000Z",
			 "Amount": {
				"Amount": "39.05",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-06T23:00:00.000Z",
			 "ValueDateTime": "2017-09-06T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.12",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-08T23:00:00.000Z",
			 "ValueDateTime": "2017-09-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "42.71",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-11T23:00:00.000Z",
			 "ValueDateTime": "2017-09-11T23:00:00.000Z",
			 "Amount": {
				"Amount": "3835.25",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-13T23:00:00.000Z",
			 "ValueDateTime": "2017-09-13T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.57",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-14T23:00:00.000Z",
			 "ValueDateTime": "2017-09-14T23:00:00.000Z",
			 "Amount": {
				"Amount": "17.45",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-15T23:00:00.000Z",
			 "ValueDateTime": "2017-09-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "2940.64",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-19T23:00:00.000Z",
			 "ValueDateTime": "2017-09-19T23:00:00.000Z",
			 "Amount": {
				"Amount": "181.76",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-19T23:00:00.000Z",
			 "ValueDateTime": "2017-09-19T23:00:00.000Z",
			 "Amount": {
				"Amount": "8.27",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-21T23:00:00.000Z",
			 "ValueDateTime": "2017-09-21T23:00:00.000Z",
			 "Amount": {
				"Amount": "18.89",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-22T23:00:00.000Z",
			 "ValueDateTime": "2017-09-22T23:00:00.000Z",
			 "Amount": {
				"Amount": "24.88",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-22T23:00:00.000Z",
			 "ValueDateTime": "2017-09-22T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.20",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-23T23:00:00.000Z",
			 "ValueDateTime": "2017-09-23T23:00:00.000Z",
			 "Amount": {
				"Amount": "4.07",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-23T23:00:00.000Z",
			 "ValueDateTime": "2017-09-23T23:00:00.000Z",
			 "Amount": {
				"Amount": "612.49",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "723.10",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-02T00:00:00.000Z",
			 "ValueDateTime": "2017-01-02T00:00:00.000Z",
			 "Amount": {
				"Amount": "452.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-03T00:00:00.000Z",
			 "ValueDateTime": "2017-01-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "23.76",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-06T00:00:00.000Z",
			 "ValueDateTime": "2017-01-06T00:00:00.000Z",
			 "Amount": {
				"Amount": "1.77",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-08T00:00:00.000Z",
			 "ValueDateTime": "2017-01-08T00:00:00.000Z",
			 "Amount": {
				"Amount": "139.37",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-09T00:00:00.000Z",
			 "ValueDateTime": "2017-01-09T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-10T00:00:00.000Z",
			 "ValueDateTime": "2017-01-10T00:00:00.000Z",
			 "Amount": {
				"Amount": "1.34",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-13T00:00:00.000Z",
			 "ValueDateTime": "2017-01-13T00:00:00.000Z",
			 "Amount": {
				"Amount": "2.55",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-17T00:00:00.000Z",
			 "ValueDateTime": "2017-01-17T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-20T00:00:00.000Z",
			 "ValueDateTime": "2017-01-20T00:00:00.000Z",
			 "Amount": {
				"Amount": "2512.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-23T00:00:00.000Z",
			 "ValueDateTime": "2017-01-23T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-01-28T00:00:00.000Z",
			 "ValueDateTime": "2017-01-28T00:00:00.000Z",
			 "Amount": {
				"Amount": "170.99",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-06T00:00:00.000Z",
			 "ValueDateTime": "2017-02-06T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-12T00:00:00.000Z",
			 "ValueDateTime": "2017-02-12T00:00:00.000Z",
			 "Amount": {
				"Amount": "6.51",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-13T00:00:00.000Z",
			 "ValueDateTime": "2017-02-13T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-20T00:00:00.000Z",
			 "ValueDateTime": "2017-02-20T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-24T00:00:00.000Z",
			 "ValueDateTime": "2017-02-24T00:00:00.000Z",
			 "Amount": {
				"Amount": "25.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-28T00:00:00.000Z",
			 "ValueDateTime": "2017-02-28T00:00:00.000Z",
			 "Amount": {
				"Amount": "38.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-02-28T00:00:00.000Z",
			 "ValueDateTime": "2017-02-28T00:00:00.000Z",
			 "Amount": {
				"Amount": "8933.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-03T00:00:00.000Z",
			 "ValueDateTime": "2017-03-03T00:00:00.000Z",
			 "Amount": {
				"Amount": "30.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-04T00:00:00.000Z",
			 "ValueDateTime": "2017-03-04T00:00:00.000Z",
			 "Amount": {
				"Amount": "7.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-07T00:00:00.000Z",
			 "ValueDateTime": "2017-03-07T00:00:00.000Z",
			 "Amount": {
				"Amount": "32.40",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-11T00:00:00.000Z",
			 "ValueDateTime": "2017-03-11T00:00:00.000Z",
			 "Amount": {
				"Amount": "1.97",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-16T00:00:00.000Z",
			 "ValueDateTime": "2017-03-16T00:00:00.000Z",
			 "Amount": {
				"Amount": "19005.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-19T00:00:00.000Z",
			 "ValueDateTime": "2017-03-19T00:00:00.000Z",
			 "Amount": {
				"Amount": "35.95",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-20T00:00:00.000Z",
			 "ValueDateTime": "2017-03-20T00:00:00.000Z",
			 "Amount": {
				"Amount": "2512.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-24T00:00:00.000Z",
			 "ValueDateTime": "2017-03-24T00:00:00.000Z",
			 "Amount": {
				"Amount": "70.27",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-26T23:00:00.000Z",
			 "ValueDateTime": "2017-03-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "38.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-03-27T23:00:00.000Z",
			 "ValueDateTime": "2017-03-27T23:00:00.000Z",
			 "Amount": {
				"Amount": "5.50",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-01T23:00:00.000Z",
			 "ValueDateTime": "2017-04-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "1041.84",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-01T23:00:00.000Z",
			 "ValueDateTime": "2017-04-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "634.57",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-05T23:00:00.000Z",
			 "ValueDateTime": "2017-04-05T23:00:00.000Z",
			 "Amount": {
				"Amount": "56.83",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-06T23:00:00.000Z",
			 "ValueDateTime": "2017-04-06T23:00:00.000Z",
			 "Amount": {
				"Amount": "9.70",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-07T23:00:00.000Z",
			 "ValueDateTime": "2017-04-07T23:00:00.000Z",
			 "Amount": {
				"Amount": "17.38",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-12T23:00:00.000Z",
			 "ValueDateTime": "2017-04-12T23:00:00.000Z",
			 "Amount": {
				"Amount": "2.54",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-19T23:00:00.000Z",
			 "ValueDateTime": "2017-04-19T23:00:00.000Z",
			 "Amount": {
				"Amount": "2699.80",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-21T23:00:00.000Z",
			 "ValueDateTime": "2017-04-21T23:00:00.000Z",
			 "Amount": {
				"Amount": "38.29",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-04-22T23:00:00.000Z",
			 "ValueDateTime": "2017-04-22T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.14",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-03T23:00:00.000Z",
			 "ValueDateTime": "2017-05-03T23:00:00.000Z",
			 "Amount": {
				"Amount": "39.14",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-10T23:00:00.000Z",
			 "ValueDateTime": "2017-05-10T23:00:00.000Z",
			 "Amount": {
				"Amount": "59.77",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-17T23:00:00.000Z",
			 "ValueDateTime": "2017-05-17T23:00:00.000Z",
			 "Amount": {
				"Amount": "6.86",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-17T23:00:00.000Z",
			 "ValueDateTime": "2017-05-17T23:00:00.000Z",
			 "Amount": {
				"Amount": "2666.03",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-24T23:00:00.000Z",
			 "ValueDateTime": "2017-05-24T23:00:00.000Z",
			 "Amount": {
				"Amount": "20.50",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-28T23:00:00.000Z",
			 "ValueDateTime": "2017-05-28T23:00:00.000Z",
			 "Amount": {
				"Amount": "12.84",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-28T23:00:00.000Z",
			 "ValueDateTime": "2017-05-28T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.39",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-05-31T23:00:00.000Z",
			 "ValueDateTime": "2017-05-31T23:00:00.000Z",
			 "Amount": {
				"Amount": "134.40",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-01T23:00:00.000Z",
			 "ValueDateTime": "2017-06-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "6.93",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-03T23:00:00.000Z",
			 "ValueDateTime": "2017-06-03T23:00:00.000Z",
			 "Amount": {
				"Amount": "24.33",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-08T23:00:00.000Z",
			 "ValueDateTime": "2017-06-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.11",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-13T23:00:00.000Z",
			 "ValueDateTime": "2017-06-13T23:00:00.000Z",
			 "Amount": {
				"Amount": "1484.02",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-15T23:00:00.000Z",
			 "ValueDateTime": "2017-06-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "1.92",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-17T23:00:00.000Z",
			 "ValueDateTime": "2017-06-17T23:00:00.000Z",
			 "Amount": {
				"Amount": "3702.84",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-21T23:00:00.000Z",
			 "ValueDateTime": "2017-06-21T23:00:00.000Z",
			 "Amount": {
				"Amount": "27.28",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-23T23:00:00.000Z",
			 "ValueDateTime": "2017-06-23T23:00:00.000Z",
			 "Amount": {
				"Amount": "49.13",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-25T23:00:00.000Z",
			 "ValueDateTime": "2017-06-25T23:00:00.000Z",
			 "Amount": {
				"Amount": "8.06",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-30T23:00:00.000Z",
			 "ValueDateTime": "2017-06-30T23:00:00.000Z",
			 "Amount": {
				"Amount": "895.26",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-06-30T23:00:00.000Z",
			 "ValueDateTime": "2017-06-30T23:00:00.000Z",
			 "Amount": {
				"Amount": "10.37",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-04T23:00:00.000Z",
			 "ValueDateTime": "2017-07-04T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.24",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-05T23:00:00.000Z",
			 "ValueDateTime": "2017-07-05T23:00:00.000Z",
			 "Amount": {
				"Amount": "7.27",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-08T23:00:00.000Z",
			 "ValueDateTime": "2017-07-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "1.21",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-11T23:00:00.000Z",
			 "ValueDateTime": "2017-07-11T23:00:00.000Z",
			 "Amount": {
				"Amount": "115.05",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-18T23:00:00.000Z",
			 "ValueDateTime": "2017-07-18T23:00:00.000Z",
			 "Amount": {
				"Amount": "4138.21",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-20T23:00:00.000Z",
			 "ValueDateTime": "2017-07-20T23:00:00.000Z",
			 "Amount": {
				"Amount": "32.93",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-07-26T23:00:00.000Z",
			 "ValueDateTime": "2017-07-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "11.57",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "ReceivedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-05T23:00:00.000Z",
			 "ValueDateTime": "2017-08-05T23:00:00.000Z",
			 "Amount": {
				"Amount": "38.11",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-10T23:00:00.000Z",
			 "ValueDateTime": "2017-08-10T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.48",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-12T23:00:00.000Z",
			 "ValueDateTime": "2017-08-12T23:00:00.000Z",
			 "Amount": {
				"Amount": "4.05",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-15T23:00:00.000Z",
			 "ValueDateTime": "2017-08-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "121.92",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-22T23:00:00.000Z",
			 "ValueDateTime": "2017-08-22T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.74",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-22T23:00:00.000Z",
			 "ValueDateTime": "2017-08-22T23:00:00.000Z",
			 "Amount": {
				"Amount": "4.12",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-26T23:00:00.000Z",
			 "ValueDateTime": "2017-08-26T23:00:00.000Z",
			 "Amount": {
				"Amount": "31.82",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-29T23:00:00.000Z",
			 "ValueDateTime": "2017-08-29T23:00:00.000Z",
			 "Amount": {
				"Amount": "51.76",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-08-29T23:00:00.000Z",
			 "ValueDateTime": "2017-08-29T23:00:00.000Z",
			 "Amount": {
				"Amount": "6.49",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedCheques",
				"SubCode": "BankCheque"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-01T23:00:00.000Z",
			 "ValueDateTime": "2017-09-01T23:00:00.000Z",
			 "Amount": {
				"Amount": "3.99",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-06T23:00:00.000Z",
			 "ValueDateTime": "2017-09-06T23:00:00.000Z",
			 "Amount": {
				"Amount": "63.73",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-08T23:00:00.000Z",
			 "ValueDateTime": "2017-09-08T23:00:00.000Z",
			 "Amount": {
				"Amount": "0.12",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CounterTransactions",
				"SubCode": "BranchDeposit"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-13T23:00:00.000Z",
			 "ValueDateTime": "2017-09-13T23:00:00.000Z",
			 "Amount": {
				"Amount": "27.56",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-15T23:00:00.000Z",
			 "ValueDateTime": "2017-09-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "1212.44",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "IssuedDirectDebits",
				"SubCode": "DirectDebitPayment"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-15T23:00:00.000Z",
			 "ValueDateTime": "2017-09-15T23:00:00.000Z",
			 "Amount": {
				"Amount": "13.34",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "BookingDateTime": "2017-09-20T23:00:00.000Z",
			 "ValueDateTime": "2017-09-20T23:00:00.000Z",
			 "Amount": {
				"Amount": "17.58",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked",
			 "BankTransactionCode": {
				"Code": "CustomerCardTransactions",
				"SubCode": "CashWithdrawal"
			 }
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "12639ced-68e3-4785-a99b-6d2e582c83f5",
			 "TransactionReference": "6",
			 "BookingDateTime": "2019-09-02T12:52:41.495Z",
			 "ValueDateTime": "2019-09-02T12:52:41.495Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "15.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "cc4887e3-d162-4b88-859e-654de49a8555",
			 "TransactionReference": "4",
			 "BookingDateTime": "2019-09-04T10:07:11.961Z",
			 "ValueDateTime": "2019-09-04T10:07:11.961Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "1.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "1105bb5c-1985-4280-9e8a-9a86b14f369d",
			 "TransactionReference": "9",
			 "BookingDateTime": "2019-09-04T10:07:11.961Z",
			 "ValueDateTime": "2019-09-04T10:07:11.961Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "1.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "c4399b0e-22fc-46f7-a78e-e32ce456aba3",
			 "TransactionReference": "5",
			 "BookingDateTime": "2019-09-05T10:34:29.745Z",
			 "ValueDateTime": "2019-09-05T10:34:29.745Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "1.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "80810793-9467-4fd0-8ff3-f1d9b6c11345",
			 "TransactionReference": "9",
			 "BookingDateTime": "2019-09-04T10:07:12.684Z",
			 "ValueDateTime": "2019-09-04T10:07:12.684Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "3aab2f35-5887-46d6-b540-0fda41d2b89a",
			 "TransactionReference": "4",
			 "BookingDateTime": "2019-09-04T10:07:12.893Z",
			 "ValueDateTime": "2019-09-04T10:07:12.893Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "cc7cbde7-1e22-49c9-8087-0e8e28fd8268",
			 "TransactionReference": "7",
			 "BookingDateTime": "2019-09-05T10:26:04.149Z",
			 "ValueDateTime": "2019-09-05T10:26:04.149Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "1.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "67e8faf7-8555-4d0f-8c11-ea63fa98fd05",
			 "TransactionReference": "5",
			 "BookingDateTime": "2019-09-05T10:26:04.149Z",
			 "ValueDateTime": "2019-09-05T10:26:04.149Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "1.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "a86ac417-6daa-4fc8-b9ab-a1cbec8a483b",
			 "TransactionReference": "9",
			 "BookingDateTime": "2019-09-05T10:26:04.709Z",
			 "ValueDateTime": "2019-09-05T10:26:04.709Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "189186ad-8394-4513-b7f6-545bb5e99637",
			 "TransactionReference": "6",
			 "BookingDateTime": "2019-09-05T10:26:04.854Z",
			 "ValueDateTime": "2019-09-05T10:26:04.854Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "312bf1a3-f54c-46f0-99b6-8342ad40c371",
			 "TransactionReference": "1",
			 "BookingDateTime": "2019-09-05T10:34:29.745Z",
			 "ValueDateTime": "2019-09-05T10:34:29.745Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "1.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "483e6b66-14a7-4d41-b8d8-bead86163c79",
			 "TransactionReference": "a",
			 "BookingDateTime": "2019-09-05T10:34:30.422Z",
			 "ValueDateTime": "2019-09-05T10:34:30.422Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "885b457d-13a2-4e76-bdef-d16bfc8c757a",
			 "TransactionReference": "9",
			 "BookingDateTime": "2019-09-05T10:34:30.559Z",
			 "ValueDateTime": "2019-09-05T10:34:30.559Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "e07c005c-3145-41ae-b559-47d82d1f3df1",
			 "TransactionReference": "e",
			 "BookingDateTime": "2019-09-09T22:02:42.013Z",
			 "ValueDateTime": "2019-09-09T22:02:42.013Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "15.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "3c22be35-7d05-467a-a859-c1f9ff843d79",
			 "TransactionReference": "f",
			 "BookingDateTime": "2019-09-10T17:09:03.289Z",
			 "ValueDateTime": "2019-09-10T17:09:03.289Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "5.9",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "e16c6d8e-997c-4d46-bffb-a2f1b25d7a82",
			 "TransactionReference": "2",
			 "BookingDateTime": "2019-09-10T17:19:35.874Z",
			 "ValueDateTime": "2019-09-10T17:19:35.874Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "3.7",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "bb36ceee-2f1b-48f6-aec6-3bfd14e80872",
			 "TransactionReference": "2",
			 "BookingDateTime": "2019-09-11T08:15:10.337Z",
			 "ValueDateTime": "2019-09-11T08:15:10.337Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "3.7",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "be2da1d2-6b73-4142-b664-e25881aeb88c",
			 "TransactionReference": "8",
			 "BookingDateTime": "2019-09-12T06:14:29.174Z",
			 "ValueDateTime": "2019-09-12T06:14:29.174Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "1.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "cb2e0259-7482-46e0-9735-627cfe1a7ed6",
			 "TransactionReference": "a",
			 "BookingDateTime": "2019-09-17T00:20:10.944Z",
			 "ValueDateTime": "2019-09-17T00:20:10.944Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "165.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "3a330796-1345-4f78-aa57-efa9d6345225",
			 "TransactionReference": "5",
			 "BookingDateTime": "2019-09-12T06:14:29.174Z",
			 "ValueDateTime": "2019-09-12T06:14:29.174Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "1.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Credit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "b3c9d6fd-614a-48cd-b485-1f146b013e6f",
			 "TransactionReference": "3",
			 "BookingDateTime": "2019-09-12T06:14:29.933Z",
			 "ValueDateTime": "2019-09-12T06:14:29.933Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "91877ef2-08be-4e0c-b9f7-805d031dab9f",
			 "TransactionReference": "3",
			 "BookingDateTime": "2019-09-16T19:03:34.742Z",
			 "ValueDateTime": "2019-09-16T19:03:34.742Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "3.7",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "db7ab89d-d6f3-49a1-821c-4434b498f271",
			 "TransactionReference": "a",
			 "BookingDateTime": "2019-09-17T13:26:54.035Z",
			 "ValueDateTime": "2019-09-17T13:26:54.035Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "2.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "662b2819-d06e-44a4-ae09-c7db9f424d1a",
			 "TransactionReference": "b",
			 "BookingDateTime": "2019-09-17T21:40:32.112Z",
			 "ValueDateTime": "2019-09-17T21:40:32.112Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "2.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000006",
			 "TransactionId": "2859bd17-bb44-4332-b43f-2323d5446460",
			 "TransactionReference": "b",
			 "BookingDateTime": "2019-09-19T23:17:29.143Z",
			 "ValueDateTime": "2019-09-19T23:17:29.143Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "34.97",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "f79b496c-34bd-4bf8-91a6-2e610d322cc0",
			 "TransactionReference": "c",
			 "BookingDateTime": "2019-09-12T06:14:30.179Z",
			 "ValueDateTime": "2019-09-12T06:14:30.179Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.90",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "76fe16a3-76fc-48ad-a40b-122732847a6e",
			 "TransactionReference": "f",
			 "BookingDateTime": "2019-09-13T15:21:35.757Z",
			 "ValueDateTime": "2019-09-13T15:21:35.757Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "3.7",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "9587f903-4f8a-4fb8-8764-9ff7517bc3c9",
			 "TransactionReference": "f",
			 "BookingDateTime": "2019-09-16T16:06:29.477Z",
			 "ValueDateTime": "2019-09-16T16:06:29.477Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "3.75",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000005",
			 "TransactionId": "ffcb7ac9-b178-4093-8a10-b7281c39cf15",
			 "TransactionReference": "60c09267-21a0-4fc2-8a08-551a839bb77",
			 "BookingDateTime": "2019-09-16T01:59:39.629Z",
			 "ValueDateTime": "2019-09-16T01:59:39.629Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.10",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "c26354b4-9b57-451e-b4a5-e5ac4b1c080d",
			 "TransactionReference": "f",
			 "BookingDateTime": "2019-09-16T04:19:58.082Z",
			 "ValueDateTime": "2019-09-16T04:19:58.082Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.10",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "8e984096-fad0-4a86-82d5-1c009bddf392",
			 "TransactionReference": "3",
			 "BookingDateTime": "2019-09-16T16:08:15.389Z",
			 "ValueDateTime": "2019-09-16T16:08:15.389Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "5.95",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "c7288e94-7af4-40de-86a2-4140dfa9a7f1",
			 "TransactionReference": "5",
			 "BookingDateTime": "2019-09-16T23:06:14.644Z",
			 "ValueDateTime": "2019-09-16T23:06:14.644Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "0.99",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "3c0d54dd-ecc3-40ad-b59a-bf747509caa6",
			 "TransactionReference": "3",
			 "BookingDateTime": "2019-09-17T06:58:44.268Z",
			 "ValueDateTime": "2019-09-17T06:58:44.268Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "3.5",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "b6862a09-b575-4253-96fb-fbfc2afc5c6b",
			 "TransactionReference": "b",
			 "BookingDateTime": "2019-09-17T07:39:41.686Z",
			 "ValueDateTime": "2019-09-17T07:39:41.686Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "4.2",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "56e34e6f-b5ee-4be9-af36-cbc022c4e95a",
			 "TransactionReference": "0",
			 "BookingDateTime": "2019-09-17T13:19:24.649Z",
			 "ValueDateTime": "2019-09-17T13:19:24.649Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "2.00",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "25c85400-2b2e-4464-9f62-8eb7e6c6f241",
			 "TransactionReference": "8",
			 "BookingDateTime": "2019-09-18T08:48:32.061Z",
			 "ValueDateTime": "2019-09-18T08:48:32.061Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "6.2",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000001",
			 "TransactionId": "9664c500-004d-40d8-b6b3-ae914b89f7dc",
			 "TransactionReference": "a",
			 "BookingDateTime": "2019-09-24T15:33:54.413Z",
			 "ValueDateTime": "2019-09-24T15:33:54.413Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "3.7",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  },
		  {
			 "AccountId": "700004000000000000000002",
			 "TransactionId": "e11bc082-f680-4101-8509-70a0254e811a",
			 "TransactionReference": "5",
			 "BookingDateTime": "2019-09-24T22:43:48.901Z",
			 "ValueDateTime": "2019-09-24T22:43:48.901Z",
			 "ProprietaryBankTransactionCode": {
				"Code": "PMT"
			 },
			 "Amount": {
				"Amount": "3.5",
				"Currency": "GBP"
			 },
			 "CreditDebitIndicator": "Debit",
			 "Status": "Booked"
		  }
	   ]
	},
	"Links": {
	   "Self": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp/transactions"
	},
	"Meta": {
	   "TotalPages": 1,
	   "FirstAvailableDateTime": "2016-01-01T08:40:00.000Z",
	   "LastAvailableDateTime": "2025-12-31T08:40:00.000Z"
	}
 }`)
)
