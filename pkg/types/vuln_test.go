package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVulns(t *testing.T) {
	r := &VulnerabilityReport{
		Vulnerabilities: []*Vulnerability{
			{ID: "CVE-2019-0001", Severity: "HIGH", Package: "test"},
			{ID: "CVE-2019-0002", Severity: "LOW", Package: "test"},
			{ID: "CVE-2019-0003", Severity: "MEDIUM", Package: "test"},
			{ID: "CVE-2019-0004", Severity: "HIGH", Package: "test"},
			{ID: "CVE-2019-0005", Severity: "LOW", Package: "test"},
		},
	}

	for _, o := range r.Vulnerabilities {
		assert.NotEmpty(t, o.ID)
		assert.NotEmpty(t, o.Package)
		assert.NotEmpty(t, o.Package)
		assert.NotEmpty(t, o.Hash())
		assert.NotEmpty(t, o.String())
	}
}
