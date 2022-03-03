package executors

import (
	"strings"

	"github.com/OpenBankingUK/conformance-suite/pkg/manifest"
	"github.com/OpenBankingUK/conformance-suite/pkg/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	resty "gopkg.in/resty.v1"
)

// GetDynamicResourceIds retrieves the accounts and statements resource ids for the current token
func GetDynamicResourceIds(tokenName, token string, ctx *model.Context, requiredTokens []manifest.RequiredTokens, interactionId string) error {
	logger := logrus.WithFields(logrus.Fields{
		"module":    "GetDynamicResourceIds",
		"tokenName": tokenName,
		"token":     token,
	})

	err := getDynamicResourceIds(tokenName, token, ctx, logger, requiredTokens, interactionId)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("exchangeCodeForToken failed")
		return err
	}
	return nil
}

func getDynamicResourceIds(tokenName, token string, ctx *model.Context, logger *logrus.Entry, requiredTokens []manifest.RequiredTokens, interactionId string) error {

	if !strings.HasPrefix(tokenName, "account") {
		return nil
	}

	resourceBaseURL, err := ctx.GetString("resource_server")
	if err != nil {
		return errors.New("cannot get resource_base_url for dynamic_resource_id call")
	}
	apiVersion, err := ctx.GetString("api-version")
	if err != nil {
		return errors.New("cannot get api_version for code for dynamic_resource_id call")
	}
	xFapiFinancialID, err := ctx.GetString("x-fapi-financial-id")
	if err != nil {
		return errors.New("cannot get X-Fapi_Financial for code for dynamic_resource_id call")
	}

	accountsEndpoint := resourceBaseURL + "/open-banking/" + apiVersion + "/aisp/accounts"
	var resp *resty.Response
	resp, err = resty.R().
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("X-Fapi-Financial-Id", xFapiFinancialID).
		SetHeader("X-Fapi-Interaction-Id", interactionId).
		SetHeader("X-Fcs-Testcase-Id", "GetDynamicResourceIdsAccounts").
		Get(accountsEndpoint)

	if err != nil {
		logger.Errorln("error calling /accounts for account number dynamic resource", err)
		return err
	}

	logger.Tracef("response code: %d ", resp.StatusCode())
	body := string(resp.Body())
	accountID, err := getAccountIDFromJSONResponse(body, logger)
	if err != nil { // out of band fix
		return err
	}

	for k, v := range requiredTokens {
		if v.Name == tokenName {
			requiredTokens[k].AccountID = accountID // put dynamic account number into permissions struct
		}
	}

	return nil
}

func getAccountIDFromJSONResponse(body string, logger *logrus.Entry) (string, error) {
	account := gjson.Get(body, "Data.Account.0.AccountId")
	accountString := account.String()
	if len(accountString) == 0 {
		return "", errors.New("DynamicResourceId, zero length account number")
	}
	logger.Infof("DynamicResource account number: %s", accountString)
	return accountString, nil
}
