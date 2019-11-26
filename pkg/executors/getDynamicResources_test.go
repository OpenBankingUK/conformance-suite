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
)
