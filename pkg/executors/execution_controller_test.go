package executors

import (
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/model"
	"github.com/OpenBankingUK/conformance-suite/pkg/test"

	"github.com/OpenBankingUK/conformance-suite/pkg/executors/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTestCaseRunner(t *testing.T) {
	controller := &mocks.DaemonController{}
	definition := RunDefinition{}
	runner := NewTestCaseRunner(test.NullLogger(), definition, controller)

	assert.NotNil(t, runner.runningLock)
	assert.Equal(t, definition, runner.definition)
	assert.Equal(t, controller, runner.daemonController)
	assert.False(t, runner.running)
}

func TestMakeConsentAcquisitionCtxDropsStaleClientAccessToken(t *testing.T) {
	for _, componentFile := range []string{"PSUConsentProviderComponentV3.json", "PSUConsentProviderComponentV4.json"} {
		t.Run(componentFile, func(t *testing.T) {
			runner := NewTestCaseRunner(test.NullLogger(), RunDefinition{}, &mocks.DaemonController{})
			sharedCtx := &model.Context{"client_access_token": "stale-payments-token"}

			ruleCtx := runner.makeConsentAcquisitionCtx(sharedCtx, TokenConsentIDItem{TokenName: "Token001", Permissions: `"ReadAccountsBasic"`})

			_, err := ruleCtx.GetString("client_access_token")
			assert.Error(t, err)
			sharedToken, err := sharedCtx.GetString("client_access_token")
			require.NoError(t, err)
			assert.Equal(t, "stale-payments-token", sharedToken, "shared context must not be mutated")

			comp, err := model.LoadComponent(componentFile)
			require.NoError(t, err)
			tests := comp.GetTests()
			for i := range tests {
				tests[i].ProcessReplacementFields(ruleCtx, false)
			}

			var accountConsent *model.TestCase
			for i := range tests {
				if tests[i].ID == "#compPsuConsent02" {
					accountConsent = &tests[i]
				}
			}
			require.NotNil(t, accountConsent)
			assert.NotContains(t, accountConsent.Input.Headers["authorization"], "stale-payments-token")

			ruleCtx.PutString("client_access_token", "fresh-accounts-token")
			accountConsent.ProcessReplacementFields(ruleCtx, false)
			assert.Equal(t, "Bearer fresh-accounts-token", accountConsent.Input.Headers["authorization"])
		})
	}
}
