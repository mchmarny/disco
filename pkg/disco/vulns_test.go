package disco

import (
	"context"
	"strings"
	"testing"

	"github.com/mchmarny/disco/pkg/metric"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	in := &types.VulnsQuery{
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

func TestCounter(t *testing.T) {
	c := &metric.ConsoleCounter{}
	ctx := context.TODO()
	err := MeterVulns(ctx, c, "../../etc/data/report-vuln.json")
	assert.NoError(t, err)
}
