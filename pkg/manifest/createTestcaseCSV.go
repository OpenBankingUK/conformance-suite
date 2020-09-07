package manifest

import (
	"fmt"
	"strings"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

const accountsPath = "file://manifests/ob_3.1_accounts_transactions_fca.json"
const paymentsPath = "file://manifests/ob_3.1_payment_fca.json"
const cbpiiPath = "file://manifests/ob_3.1_cbpii_fca.json"

type apiTests struct {
	ApiType     string
	PathtoTests string
}

func GenerateTestCaseListCSV() {
	apis := []apiTests{{"accounts", accountsPath}, {"payments", paymentsPath}, {"cbpii", cbpiiPath}}

	var values []interface{}
	values = append(values, "accounts_v0.0.0", "payments_v0.0.0", "cbpii_v0.0.0")
	context := &model.Context{"apiversions": values}

	fmt.Println("Resource,TestCase Id,Method,Path,Condition,Version,Schema,Sig,Description")
	for _, api := range apis {
		scripts, _, err := LoadGenerationResources(api.ApiType, api.PathtoTests, context)
		if err != nil {
			fmt.Printf("Error on loadGenerationResources %v", err)
			return
		}

		for _, v := range scripts.Scripts {
			description := strings.Replace(v.Description, ",", "", -1)
			fmt.Printf("%s,%s,%s,%s,%s,%s,%5t,%5t,%s\n", v.Resource, v.ID, strings.ToUpper(v.Method),
				v.URI, v.URIImplemenation, v.APIVersion, v.SchemaCheck, v.ValidateSignature, description)
		}
	}
}
