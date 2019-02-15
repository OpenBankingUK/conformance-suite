package executors

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildParameters(t *testing.T) {

	permissions := map[string][]string{
		"tok0001": {"readbasic", "writebasic", "updatebasic"},
		"tok0002": {"updatebasic"},
		"tok0003": {},
	}

	buildstr := strings.Join(permissions["tok0001"], ",")
	fmt.Println(buildstr)
	assert.Equal(t, `"readbasic","writebasic","updatebasic"`, buildstr)

	buildstr = strings.Join(permissions["tok0002"], ",")
	fmt.Println(buildstr)
	assert.Equal(t, `"updatebasic"`, buildstr)

	buildstr = strings.Join(permissions["tok0003"], ",")
	assert.Equal(t, ``, buildstr)
	fmt.Println(buildstr)

}
