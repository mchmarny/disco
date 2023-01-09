package trivy

import (
	"strings"
	"testing"

	"github.com/mchmarny/disco/pkg/types"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestCVEFilterParsing(t *testing.T) {
	cve := "CVE-2020-8911"
	expectedResults := 1

	filter := func(v interface{}) bool {
		vul := v.(*types.Vulnerability)
		exclude := !strings.EqualFold(cve, vul.ID)
		log.Debug().Msgf("filter on cve (want: %s, got: %s, filter out: %t",
			cve, vul.ID, exclude)
		return exclude
	}

	rep, err := testVulnParsing(t, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.Equal(t, expectedResults, len(rep.Vulnerabilities))
}

func TestVulnFilterParsing(t *testing.T) {
	severity := types.VulnSevHigh
	expectedResults := 5

	filter := func(v interface{}) bool {
		vul := v.(*types.Vulnerability)
		vs := types.ParseMinVulnSeverityOrDefault(vul.Severity)
		exclude := !(vs >= severity)
		log.Debug().Msgf("filter on severity (want: %s, got: %s, filter out: %t",
			severity, vul.Severity, exclude)
		return exclude
	}

	rep, err := testVulnParsing(t, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	assert.Equal(t, expectedResults, len(rep.Vulnerabilities))
}

func testVulnParsing(t *testing.T, filter types.ItemFilter) (*types.VulnerabilityReport, error) {
	src := "../../../data/test-vuln.json"
	rep, err := ParseVulnerabilities(src, filter)
	assert.NoError(t, err)
	assert.NotNil(t, rep)
	return rep, err
}
