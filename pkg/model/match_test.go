package model

import (
	"net/http"
	"testing"
)

var (
	responseBody = []byte(`
	{
	    "Data": {
	        "Account": [{
	            "AccountId": "500000000000000000000001",
	            "Currency": "GBP",
	            "Nickname": "xxxx0101",
	            "AccountType": "Personal",
	            "AccountSubType": "CurrentAccount",
	            "Account": [{
	                "SchemeName": "SortCodeAccountNumber",
	                "Identification": "10000119820101",
	                "Name": "Mr. Roberto Rastapopoulos & Mr. Mitsuhirato",
	                "SecondaryIdentification": "Roll No. 001"
	            }]
	        }, {
	            "AccountId": "500000000000000000000007",
	            "Currency": "GBP",
	            "Nickname": "xxxx0001",
	            "AccountType": "Business",
	            "AccountSubType": "CurrentAccount",
	            "Account": [{
	                "SchemeName": "SortCodeAccountNumber",
	                "Identification": "10000190210001",
	                "Name": "Marios Amazing Carpentry Supplies Limited"
	            }]
	        }]
	    },
	    "Links": {
	        "Self": "http://modelobank2018.o3bank.co.uk/open-banking/v2.0/accounts/"
	    },
	    "Meta": {
	        "TotalPages": 1
	    }
    }    
	`)
)

func TestContextPutFromExpectsMatch(t *testing.T) {
	ctx := Context{}
	// look for variable in context - should not be there
	m := Match{Description: "Test Context Put", ContextName: "AccountId", JSON: "Data.Account.1.AccountId"}
	contextPut := ContextAccessor{Context: ctx, Matches: []Match{m}}

	response := http.Response{}
	contextPut.PutValues(&response)

}

/*
	body := string(bodyBytes)
	for _, match := range t.Expect.Matches {
		jsonMatch := match.JSON // check if there is a JSON match to be satisifed
		if len(jsonMatch) > 0 {
			matched := gjson.Get(body, jsonMatch)
			if matched.String() != match.Value { // check the value of the JSON match - is equal to the 'value' parameter of the match section within the testcase 'Expects' area
				return false, fmt.Errorf("(%s):%s: Json Match: expected %s got %s", t.ID, t.Name, match.Value, matched.String())
			}
		}
	}

*/
