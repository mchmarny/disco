package disco

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	in := &VulnsQuery{
		CVE: "CVE-2019-0001",
	}

	list := []string{"CVE-2019-0001", "CVE-2019-0002", "CVE-2019-0003"}

	filter := func(v string) bool {
		if in.CVE == "" {
			return false
		}
		return strings.EqualFold(in.CVE, v)
	}

	found := 0
	for _, v := range list {
		if filter(v) {
			found++
		}
	}
	assert.Equal(t, 1, found)
}
