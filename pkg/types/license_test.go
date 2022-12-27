package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLicenses(t *testing.T) {
	report := &LicenseReport{
		Licenses: []*License{
			{Name: "MIT", Source: "test"},
			{Name: "Apache", Source: "test"},
			{Name: "GPL", Source: "test"},
			{Name: "MIT", Source: "test"},
			{Name: "Apache", Source: "test"},
		},
	}
	assert.NotEmpty(t, report.Hash())
	assert.NotEmpty(t, report.Licenses)
	for _, l := range report.Licenses {
		assert.NotEmpty(t, l.Name)
		assert.NotEmpty(t, l.Source)
		assert.NotEmpty(t, l.String())
		assert.NotEmpty(t, l.Hash())
	}
}
