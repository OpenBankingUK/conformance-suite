package executors

import (
	"errors"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
)

// TokenConsentIDs captures the token/consentIds awaiting authorisation
type TokenConsentIDs []TokenConsentIDItem

// TokenConsentIDItem is a single consentId mapping to token name
type TokenConsentIDItem struct {
	TokenName   string
	ConsentID   string
	Permissions string
}

var (
	defaultAccountPermissions = []string{"ReadAccountsBasic", "ReadAccountsDetail", "ReadBalances", "ReadBeneficiariesBasic", "ReadBeneficiariesDetail",
		"ReadDirectDebits", "ReadTransactionsBasic", "ReadTransactionsCredits", "ReadTransactionsDebits", "ReadTransactionsDetail", "ReadProducts",
		"ReadStandingOrdersDetail", "ReadProducts", "ReadStandingOrdersDetail"}
	consentChannelTimeout = 30
)

// InitiationConsentAcquisition - get required tokens
func InitiationConsentAcquisition(consentRequirements []model.SpecConsentRequirements, definition RunDefinition, ctx *model.Context) (TokenConsentIDs, error) {
	consentIDChannel := make(chan TokenConsentIDItem, 100)
	tokenParameters := getConsentTokensAndPermissions(consentRequirements)

	for tokenName, permissionList := range tokenParameters {
		runner := NewConsentAcquisitionRunner(definition, NewBufferedDaemonController())
		tokenAcquisitionType := definition.DiscoModel.DiscoveryModel.TokenAcquisition
		permissionString := buildPermissionString(permissionList)
		consentInfo := TokenConsentIDItem{TokenName: tokenName, Permissions: permissionString}
		runner.RunConsentAcquisition(consentInfo, ctx, tokenAcquisitionType, consentIDChannel)
	}

	consentItems, err := waitForConsentIDs(consentIDChannel, tokenParameters)
	for _, v := range consentItems {
		logrus.Debugf("Setting Token: %s, ConsentId: %s", v.TokenName, v.ConsentID)
		ctx.PutString(v.TokenName, v.ConsentID)
	}
	return consentItems, err
}

func waitForConsentIDs(consentIDChannel chan TokenConsentIDItem, tokenParameters map[string][]string) (TokenConsentIDs, error) {
	consentItems := TokenConsentIDs{}
	consentIDsRequired := len(tokenParameters)
	consentIDsReceived := 0
	logrus.Debugf("waiting for consentids items ...")
	for {
		select {
		case item := <-consentIDChannel:
			logrus.Infof("recieved consent channel item item %#v", item)
			consentIDsReceived++
			consentItems = append(consentItems, item)
			if consentIDsReceived == consentIDsRequired {
				logrus.Infof("Got %d required tokens - progressiing..", consentIDsReceived)
				for _, v := range consentItems {
					logrus.Infof("item %s: %s", v.TokenName, v.ConsentID)
				}
				return consentItems, nil
			}
		case <-time.After(time.Duration(consentChannelTimeout) * time.Second):
			logrus.Warnf("consent channel timeout after %d seconds", consentChannelTimeout)
			return consentItems, errors.New("ConsentChannel Timeout")
		}
	}
}

func getConsentTokensAndPermissions(consentRequirements []model.SpecConsentRequirements) map[string][]string {
	tokenParameters := make(map[string][]string)
	for _, v := range consentRequirements {
		for _, namedPermission := range v.NamedPermissions {
			codeset := namedPermission.CodeSet
			for _, b := range codeset.CodeSet {
				mystring := string(b)
				set := tokenParameters[namedPermission.Name]
				set = append(set, mystring)
				tokenParameters[namedPermission.Name] = set
			}
		}
	}
	tokenParameters["DefaultAccountToken"] = defaultAccountPermissions
	logrus.Debugf("required tokens: %#v", tokenParameters)

	return tokenParameters
}

func buildPermissionString(permissionSlice []string) string {
	var permissions string
	first := true
	for _, perms := range permissionSlice {
		if !first {
			permissions += ","
		} else {
			first = !first
		}
		permissions += "\"" + perms + "\""
	}
	return permissions
}
