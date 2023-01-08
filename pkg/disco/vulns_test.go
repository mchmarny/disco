package disco

import (
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestFulnFilter(t *testing.T) {
	filter := makeVulnerabilityFilter(&types.VulnsQuery{
		CVE:        "CVE-2019-0001",
		MinVulnSev: types.VulnSevHigh,
	})

	list := []*types.Vulnerability{
		{
			ID:       "CVE-2019-0001",
			Severity: "low",
		},
		{
			ID:       "CVE-2019-0002",
			Severity: "medium",
		},
		{
			ID:       "CVE-2019-0003",
			Severity: "high",
		},
	}

	found := testFilterOnList(filter, list)
	assert.Equal(t, 1, found)
}
