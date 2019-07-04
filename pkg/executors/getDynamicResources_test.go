package executors

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountFromBody(t *testing.T) {
	body := string(goodaccount)
	logger := logrus.WithFields(logrus.Fields{"module": "test"})
	accountid, err := getAccountIDFromJSONResponse(body, logger)
	assert.Nil(t, err)
	assert.Equal(t, "500000000000000000000001", accountid)
}

func TestGetStatementIdFromBody(t *testing.T) {
	body := string(goodstatement)
	logger := logrus.WithFields(logrus.Fields{"module": "test"})

	statementid, err := getStatementIDFromJSONResponse(body, logger)
	assert.Nil(t, err)
	assert.Equal(t, "140000000000000000000001", statementid)

}

var (
	goodaccount = []byte(`{
	"Data": {
	   "Account": [
		  {
			 "AccountId": "500000000000000000000001",
			 "Currency": "GBP",
			 "Nickname": "xxxx0101",
			 "AccountType": "Personal",
			 "AccountSubType": "CurrentAccount",
			 "Account": [
				{
				   "SchemeName": "UK.OBIE.SortCodeAccountNumber",
				   "Identification": "10000119820101",
				   "SecondaryIdentification": "Roll No. 001"
				}
			 ]
		  }
	   ]
	},
	"Links": {
	   "Self": "http://modelobank2018.o3bank.co.uk/open-banking/v3.1/aisp/accounts"
	},
	"Meta": {
	   "TotalPages": 1
	}
 }`)

	goodstatement = []byte(`
	{
	"Data": {
	   "Statement": [
		  {
			 "AccountId": "500000000000000000000001",
			 "StatementId": "140000000000000000000001",
			 "StatementReference": "001",
			 "Type": "RegularPeriodic",
			 "StartDateTime": "2017-03-01T00:00:00+00:00",
			 "EndDateTime": "2017-03-31T23:59:59+00:00",
			 "CreationDateTime": "2017-04-01T00:00:00+00:00",
			 "StatementDescription": [
				"March 2017 Statement",
				"One Free Uber Ride"
			 ]
		  }
	   ]
	},
	"Links": {
	   "Self": "http://modelobank2018.o3bank.co.uk/open-banking/v3.1/aisp/accounts/500000000000000000000001/statements"
	},
	"Meta": {
	   "TotalPages": 1,
	   "FirstAvailableDateTime": "2016-01-01T08:40:00.000Z",
	   "LastAvailableDateTime": "2025-12-31T08:40:00.000Z"
	}
  }
 `)
)
