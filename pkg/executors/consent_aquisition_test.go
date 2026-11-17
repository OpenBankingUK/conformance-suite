package executors

import (
	"fmt"
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildParameters(t *testing.T) {

	permissions := map[string][]string{
		"tok0001": {"readbasic", "writebasic", "updatebasic"},
		"tok0002": {"updatebasic"},
		"tok0003": {},
	}

	buildstr := buildPermissionString(permissions["tok0001"])
	fmt.Println(buildstr)
	assert.Equal(t, `"readbasic","writebasic","updatebasic"`, buildstr)

	buildstr = buildPermissionString(permissions["tok0002"])
	fmt.Println(buildstr)
	assert.Equal(t, `"updatebasic"`, buildstr)

	buildstr = buildPermissionString(permissions["tok0003"])
	assert.Equal(t, ``, buildstr)
	fmt.Println(buildstr)

}

func TestStoreConsentItemsCopiesTokenAndClientAccessToken(t *testing.T) {
	ctx := &model.Context{}
	items := TokenConsentIDs{
		{TokenName: "accountsToken", ConsentID: "consent-123", AccessToken: "access-token"},
	}

	storeConsentItems(ctx, items, logrus.NewEntry(logrus.New()))

	consentID, err := ctx.GetString("accountsToken")
	require.NoError(t, err)
	assert.Equal(t, "consent-123", consentID)
	accessToken, err := ctx.GetString("client_access_token")
	require.NoError(t, err)
	assert.Equal(t, "access-token", accessToken)
}
