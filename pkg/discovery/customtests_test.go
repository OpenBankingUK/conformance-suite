package discovery

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomTestsWithReplacementParameters(t *testing.T) {
	disc := readDiscoveryWithCustomReplacementTests(t)
	require.NotNil(t, disc)
	customTests := disc.DiscoveryModel.CustomTests
	assert.Equal(t, "#ct0004", customTests[0].Sequence[3].ID)
}

func readDiscoveryWithCustomReplacementTests(t *testing.T) *Model {
	initialOzone, err := ioutil.ReadFile("./testdata/ozone-parameters-test.json")
	require.Nil(t, err)
	require.NotNil(t, initialOzone)

	disco, err := UnmarshalDiscoveryJSONBytes(initialOzone)
	assert.NoError(t, err)
	assert.NotNil(t, disco.DiscoveryModel)
	return disco
}

func UnmarshalDiscoveryJSONBytes(discoveryJSON []byte) (*Model, error) {
	discovery := &Model{}
	err := json.Unmarshal(discoveryJSON, &discovery)
	return discovery, err
}
